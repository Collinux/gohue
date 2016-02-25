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
    "strconv"
)

// Bridge struct defines hardware that is used to communicate with the lights.
type Bridge struct {
    IPAddress   string
    Username    string
    Info        BridgeInfo
}

// BridgeInfo struct is the format for parsing xml from a bridge.
type BridgeInfo struct {
    XMLName	    xml.Name	`xml:"root"`
    Device      struct {
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
    } `xml:"device"`
}

// bridge.Get sends an http GET to the bridge
func (bridge *Bridge) Get(path string) ([]byte, io.Reader, error) {
    resp, err := http.Get("http://" + bridge.IPAddress + path)
    if err != nil {
        err = errors.New("Unable to access bridge.")
        log.Println(err)
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
        err = errors.New("Error: Unable marshal PUT request interface.")
        log.Println(err)
        return []byte{}, nil, err
    }
    //fmt.Println("\n\nPARAMS: ", params)
    //log.Println("\nSending PUT body: ", string(data))

	request, err := http.NewRequest("PUT", uri, bytes.NewReader(data))
    resp, err := client.Do(request)
    if err != nil {
        err = errors.New("Error: Unable to access bridge.")
        log.Println(err)
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
        err = errors.New("Error: Unable to marshal request from bridge http POST")
        log.Println(err)
        return []byte{}, nil, err
    }

    // Send the request and handle the response
    uri := fmt.Sprintf("http://" + bridge.IPAddress + path)
    resp, err := http.Post(uri, "text/json", bytes.NewReader(request))
    if err != nil {
        err = errors.New("Error: Unable to access bridge.")
        log.Println(err)
        return []byte{}, nil, err
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
        err = errors.New("Error: Unable to access bridge.")
        log.Println(err)
        return err
    }
    _, _, err = HandleResponse(resp)
    return err
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
    if strings.Contains(string(body), "error") {
        errString := string(body)
        errNum := errString[strings.Index(errString, "type\":")+6 :
            strings.Index(errString, ",\"address")]
        errDesc := errString[strings.Index(errString, "description\":\"")+14 :
            strings.Index(errString, "\"}}")]
        errOut := fmt.Sprintf("Error type %s: %s.", errNum, errDesc)
        err = errors.New(errOut)
        log.Println(err)
        return []byte{}, nil, err
    }
    return body, reader, nil
}

// NewBridge defines hardware that is compatible with Hue.
// The function is the core of all functionality, it's necessary
// to call `NewBridge` and `Login` or `CreateUser` to access any
// lights, scenes, groups, etc.
func NewBridge(ip string) (*Bridge, error) {
    bridge := Bridge {
        IPAddress: ip,
    }
    // Test the connection by attempting to get the bridge info.
    err := bridge.GetInfo()
    if err != nil {
        log.Fatal("Error: Unable to access bridge. ", err)
        return &Bridge{}, err
    }
    return &bridge, nil
}

// GetBridgeInfo retreives the description.xml file from the bridge.
func (bridge *Bridge) GetInfo() error {
    _, reader, err := bridge.Get("/description.xml")
    if err != nil {
        return err
    }
    data := BridgeInfo{}
    err = xml.NewDecoder(reader).Decode(&data)
    if err != nil {
        err = errors.New("Error: Unable to decode XML response from bridge. ")
        log.Println(err)
        return err
    }
    bridge.Info = data
    return nil
}

// Bridge.Login assigns a username to access the bridge with and
// will create the username key if it does not exist.
func (bridge *Bridge) Login(username string) error {
    bridge.Username = username
    // err := bridge.CreateUser(username)
    // if err != nil {
    //     log.Fatal("Error: Unable to login as user ", username, " ", err)
    //     return err
    // }
    return nil
}

// Bridge.CreateUser posts to ./api on the bridge to create a new whitelisted user.
func (bridge *Bridge) CreateUser(deviceType string) error {
    params := map[string]string{"devicetype": deviceType}
    body, _, err := bridge.Post("/api", params)
    if err != nil {
        log.Fatal("Error: Unable to create user. ", err)
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

// GetAllLights retreives the state of all lights that the bridge is aware of.
func (bridge *Bridge) GetAllLights() ([]Light, error) {
    uri := fmt.Sprintf("/api/%s/lights", bridge.Username)
    body, _, err := bridge.Get(uri)
    if err != nil {
        return []Light{}, err
    }

    // An index is at the top of every Light in the array
    lightMap := map[string]Light{}
    err = json.Unmarshal(body, &lightMap)
    if err != nil {
        return []Light{}, errors.New("Unable to marshal GetAllLights response. ")
    }

    // Parse the index, add the light to the list, and return the array
    lights := []Light{}
    for index, light := range lightMap {
        light.Index, err = strconv.Atoi(index)
        if err != nil {
            return []Light{}, errors.New("Unable to convert light index to integer. ")
        }
        light.Bridge = bridge
        lights = append(lights, light)
    }
    return lights, nil
}

// GetLightByIndex returns a light struct containing data on
// a light given its index stored on the bridge. This is used for
// quickly updating an individual light.
func (bridge *Bridge) GetLightByIndex(index int) (Light, error) {
    // Send an http GET and inspect the response
    uri := fmt.Sprintf("/api/%s/lights/%d", bridge.Username, index)
    body, _, err := bridge.Get(uri)
    if err != nil {
        return Light{}, err
    }
    if strings.Contains(string(body), "not available") {
        return Light{}, errors.New("Error: Light selection index out of bounds. ")
    }

    // Parse and load the response into the light array
    light := Light{}
    err = json.Unmarshal(body, &light)
    if err != nil {
        trace("", err)
    }
    light.Index = index
    light.Bridge = bridge
    return light, nil
}

// GetLight returns a light struct containing data on a given name.
func (bridge *Bridge) GetLightByName(name string) (Light, error) {
    lights, _ := bridge.GetAllLights()
    for _, light := range lights {
        if light.Name == name {
            return light, nil
        }
    }
    errOut := fmt.Sprintf("Error: Light name '%s' not found. ", name)
    return Light{}, errors.New(errOut)
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
