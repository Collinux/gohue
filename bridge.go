/*
* bridge.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux)
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/

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

// bridge.Get will send an http GET to the bridge
func (self *Bridge) Get(path string) ([]byte, io.Reader, error) {
    resp, err := http.Get("http://" + self.IPAddress + path)
    if self.Error(resp, err) {
        return []byte{}, nil, err
    }
    return handleResponse(resp)
}

// bridge.Post will send an http POST to the bridge with
// a body formatted with parameters (in a generic interface)
func (self *Bridge) Post(path string, params interface{}) ([]byte, io.Reader, error) {
    // Add the params to the request
    request, err := json.Marshal(params)
    if err != nil {
        trace("", err)
        return []byte{}, nil, nil
    }
    log.Println("\nSending POST body: ", string(request))

    // Send the request and handle the response
    uri := fmt.Sprintf("http://" + self.IPAddress + path)
    resp, err := http.Post(uri, "text/json", bytes.NewReader(request))
    if self.Error(resp, err) {
        return []byte{}, nil, nil
    }
    return handleResponse(resp)
}

func (self *Bridge) Put(path string, params interface{}) ([]byte, io.Reader, error) {
    uri := fmt.Sprintf("http://" + self.IPAddress + path)
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
    return handleResponse(resp)
}

// HandleResponse manages the http.Response content from a
// bridge Get/Put/Post/Delete by checking it for errors
// and invalid return types.
func handleResponse(resp *http.Response) ([]byte, io.Reader, error) {
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        trace("Error parsing bridge description xml.", nil)
        return []byte{}, nil, err
    }
    reader := bytes.NewReader(body)
    log.Println("Handled response:\n--------------------\n", string(body) +
        "\n--------------------\n")
    return body, reader, nil
}

// bridge.Error handles all bridge response status errors
func (self *Bridge) Error(resp *http.Response, err error) (bool) {
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

// Error Struct
// http://www.developers.meethue.com/documentation/error-messages
type Error struct {
    ID          int
    Description string
    Details     string
}

// NewBridge defines hardware that is compatible with Hue.
func NewBridge(ip string, username string) *Bridge {
    bridge := Bridge {
        IPAddress: ip,
        Username: username,
    }
    GetBridgeInfo(&bridge)
    return &bridge
}

// GetBridgeInfo retreives the description.xml file from the bridge.
func GetBridgeInfo(self *Bridge) error {
    _, reader, err := self.Get("/description.xml")
    if err != nil {
        return err
    }
    data := BridgeInfo{}
    err = xml.NewDecoder(reader).Decode(&data)
    if err != nil {
        return err
    }
    self.Info = data
    return nil
}

// CreateUser posts to ./api on the bridge to create a new whitelisted user.
func CreateUser(bridge *Bridge, deviceType string) (string, error) {
    // Send an HTTP POST with the body content
    params := map[string]string{"devicetype": deviceType}
    body, _, err := bridge.Post("/api", params)
    if err != nil {
        return "", err
    }

    // Parse the result and return it
    result := string(body)
    errFound := strings.Contains(result, "error")
    noLink := strings.Contains(result, "link button not pressed")
    if errFound && noLink {
        return "", errors.New("Bridge link button not pressed.")
    }
    return "", nil
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
