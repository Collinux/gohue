/*
* bridge.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/
// All things start with the bridge. You will find many Bridge.Func() items
// to use once a bridge has been created and identified.
// See the getting started guide on the Philips hue website:
// http://www.developers.meethue.com/documentation/getting-started

package hue

import (
    "log"
    "encoding/xml"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "runtime"
    "fmt"
    "strings"
    "bytes"
    "io"
    "errors"
)

// Bridge struct defines hardware that is used to communicate with the lights.
type Bridge struct {
    IPAddress   string
    Username    string
    Info        BridgeInfo
}

// BridgeInfo struct is the outermost (root) structure
// for parsing xml from a bridge.
type BridgeInfo struct {
    XMLName	    xml.Name	`xml:"root"`
    Device      Device      `xml:"device"`
}

// Device struct is the innermost (base) structure
// for parsing device info xml from a bridge.
type Device struct {
    XMLName	            xml.Name    `xml:"device"`
    DeviceType          string      `xml:"deviceType"`
    FriendlyName        string      `xml:"friendlyName"`
    Manufacturer        string      `xml:"manufacturer"`
    ManufacturerURL     string      `xml:"manufacturerURL"`
    ModelDescription    string      `xml:"modelDescription"`
    ModelName           string      `xml:"modelName"`
    ModelNumber         string      `xml:"modelNumber"`
    ModelURL            string      `xml:"modelURL"`
    SerialNumber        string      `xml:"serialNumber"`
    UDN                 string      `xml:"UDN"`
}

// bridge.Get sends an http GET to the bridge
func (bridge *Bridge) Get(path string) ([]byte, io.Reader, error) {
    resp, err := http.Get("http://" + bridge.IPAddress + path)
    if bridge.Error(resp, err) {
        return []byte{}, nil, err
    }
    return HandleResponse(resp)
}

// Bridge.Put sends an http PUT to the bridge with
// a body formatted with parameters (in a generic interface)
func (bridge *Bridge) Put(path string, params interface{}) ([]byte, io.Reader, error) {
    uri := fmt.Sprintf("http://" + bridge.IPAddress + path)
    client := &http.Client{}

    data, err := json.Marshal(params)
    if err != nil {
        return []byte{}, nil, err
    }
    //fmt.Println("\n\nPARAMS: ", params)
    log.Println("\nSending PUT body: ", string(data))

	request, err := http.NewRequest("PUT", uri, bytes.NewReader(data))
    resp, err := client.Do(request)
	if err != nil {
		return []byte{}, nil, err
    }
    return HandleResponse(resp)
}

// bridge.Post sends an http POST to the bridge with
// a body formatted with parameters (in a generic interface)
func (bridge *Bridge) Post(path string, params interface{}) ([]byte, io.Reader, error) {
    // Add the params to the request
    request, err := json.Marshal(params)
    if err != nil {
        trace("", err)
        return []byte{}, nil, nil
    }
    log.Println("\nSending POST body: ", string(request))

    // Send the request and handle the response
    uri := fmt.Sprintf("http://" + bridge.IPAddress + path)
    resp, err := http.Post(uri, "text/json", bytes.NewReader(request))
    if bridge.Error(resp, err) {
        return []byte{}, nil, nil
    }
    return HandleResponse(resp)
}

// Bridge.Delete sends an http DELETE to the bridge
func (bridge *Bridge) Delete(path string) error {
    uri := fmt.Sprintf("http://" + bridge.IPAddress + path)
    client := &http.Client{}
    req, err := http.NewRequest("DELETE", uri, nil)
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    _, _, err = HandleResponse(resp)
    if err != nil {
        return err
    }
    return nil
}

// HandleResponse manages the http.Response content from a
// bridge Get/Put/Post/Delete by checking it for errors
// and invalid return types.
func HandleResponse(resp *http.Response) ([]byte, io.Reader, error) {
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        trace("Error parsing bridge description xml.", nil)
        return []byte{}, nil, err
    }
    reader := bytes.NewReader(body)
    log.Println("Handled response:\n--------------------\n", string(body) +
        "\n--------------------\n")
    if strings.Contains(string(body), "error") {
        return []byte{}, nil, errors.New(string(body))
    }
    return body, reader, nil
}

// Bridge.Error handles all bridge response status errors
func (bridge *Bridge) Error(resp *http.Response, err error) (bool) {
    if err != nil {
        trace("", err)
        return true
    } else if resp.StatusCode != 200 {
        // TODO: handle other status codes
        log.Println(fmt.Sprintf("Bridge status error: %d", resp.StatusCode))
        return true
    }
    return false
}

// NewBridge defines hardware that is compatible with Hue.
func NewBridge(ip string, username string) (*Bridge, error) {
    bridge := Bridge {
        IPAddress: ip,
        Username: username,
    }
    info, err := bridge.GetInfo()
    if err != nil {
        return &Bridge{}, err
    }
    bridge.Info = info
    return &bridge, nil
}

// GetBridgeInfo retreives the description.xml file from the bridge.
func (bridge *Bridge) GetInfo() (BridgeInfo, error) {
    _, reader, err := bridge.Get("/description.xml")
    if err != nil {
        return BridgeInfo{}, err
    }
    data := BridgeInfo{}
    err = xml.NewDecoder(reader).Decode(&data)
    if err != nil {
        return BridgeInfo{}, err
    }
    bridge.Info = data
    return data, nil
}

// Bridge.CreateUser posts to ./api on the bridge to create a new whitelisted user.
func (bridge *Bridge) CreateUser(deviceType string) error {
    params := map[string]string{"devicetype": deviceType}
    body, _, err := bridge.Post("/api", params)
    if err != nil {
        return err
    }
    content := string(body)
    username := content[strings.LastIndex(content, ":\"")+2 :
        strings.LastIndex(content, "\"")]
    bridge.Username = username
    return nil
}

// Bridge.DeleteUser deletes a user given its USER KEY, not the string name.
// See http://www.developers.meethue.com/documentation/configuration-api
// for description on `username` deprecation in place of the devicetype key.
func (bridge *Bridge) DeleteUser(username string) error {
    uri := fmt.Sprintf("/api/%s/config/whitelist/%s", bridge.Username, username)
    err := bridge.Delete(uri)
    if err != nil {
        return err
    }
    return nil
}

// Log the date, time, file location, line number, and function.
// Message can be "" or Err can be nil (not both)
func trace(message string, err error) {
    pc := make([]uintptr, 10)
    runtime.Callers(2, pc)
    f := runtime.FuncForPC(pc[0])
    file, line := f.FileLine(pc[0])
    if err != nil {
        log.Printf("%s:%d %s: %s\n", file, line, f.Name(), err)
    } else {
        log.Printf("%s:%d %s: %s\n", file, line, f.Name(), message)
    }
}
