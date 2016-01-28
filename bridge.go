package hue

// username: 319b36233bd2328f3e40731b23479207
import (
    "log"
    "os"
    "encoding/xml"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "runtime"
    "fmt"
    "bytes"
)

type Bridge struct {
    IPAddress   string
    Username    string
    Info        BridgeInfo
}

type BridgeInfo struct {
    XMLName	    xml.Name	`xml:"root"`
    Device      Device      `xml:"device"`
}

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
// Go to http://<bridge_ip>/description.xml
func GetBridgeInfo(self *Bridge) {
    response, error := http.Get("http://" + self.IPAddress + "/description.xml")
    if error != nil {
        trace("", error)
    } else if response.StatusCode != 200 {
        trace(fmt.Sprintf("Bridge status error: %d", response.StatusCode), nil)
        os.Exit(1)
    }
    defer response.Body.Close()

    body, error := ioutil.ReadAll(response.Body)
    if error != nil {
        trace("Error parsing bridge description xml.", nil)
        os.Exit(1)
    }

    data := new(BridgeInfo)
    error = xml.Unmarshal(body, &data)
    if error != nil {
        trace("Error using unmarshal to split xml.", nil)
        os.Exit(1)
    }
    self.Info = *data
}

// Error Struct - REST Error Response Format
// http://www.developers.meethue.com/documentation/error-messages
type Error struct {
    response struct {
        Type        string  `xml:"type"`
        Address     string  `xml:"address"`
        Description string  `xml:"description"`
    } `xml:"error"`
}

// CreateUser posts to ./api on the bridge to create a new whitelisted user.
// If the button on the bridge was not pressed then _____todo_____
func CreateUser(bridge *Bridge, deviceType string) (string, error) {
    // Construct the http POST
    params := map[string]string{"devicetype": deviceType}
    request, err := json.Marshal(params)
    if err != nil {
        return "", err
    }

    // Send the request to create the user and read the response
    uri := fmt.Sprintf("http://%s/api", bridge.IPAddress)
    response, err := http.Post(uri, "text/json", bytes.NewReader(request))
    if err != nil {
        return "", err
    }
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    fmt.Println(string(body))

    // TODO: decode and return
    // TODO: handle errors. http://www.developers.meethue.com/documentation/error-messages

    return "", err
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
