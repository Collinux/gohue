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
    "strings"
    "bytes"
    "io"
    "errors"
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

func (self *Bridge) Get(path string) ([]byte, io.Reader, error) {
    resp, err := http.Get("http://" + self.IPAddress + path)
    if self.Error(resp, err) {
        return []byte{}, nil, err
    }
    return handleResponse(resp)
}

// bridge.Post will send an http POST to the bridge with
// a body formatted with parameters.
func (self *Bridge) Post(path string, params map[string]string) ([]byte, io.Reader, error) {
    // Add the params to the request
    request, err := json.Marshal(params)
    if err != nil {
        trace("", err)
        return []byte{}, nil, nil
    }

    // Send the request and handle the response
    uri := fmt.Sprintf("http://" + self.IPAddress + path)
    resp, err := http.Post(uri, "text/json", bytes.NewReader(request))
    if self.Error(resp, err) {
        return []byte{}, nil, nil
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

    return body, reader, nil
}

// bridge.Error handles all bridge response status errors
func (self *Bridge) Error(resp *http.Response, err error) (bool) {
    if err != nil {
        trace("", err)
        return true
    } else if resp.StatusCode != 200 {
        // TODO: handle other status codes
        trace(fmt.Sprintf("Bridge status error: %d", resp.StatusCode), nil)
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

// Error Return Values
// http://www.developers.meethue.com/documentation/error-messages
var (
    // Not from Hue documentation
    NoErr = Error{}
    ErrResponse = Error{0, "Could not read or parse response from bridge",
        "Data structure for return type may be invalid."}


    // Generic Errors from Hue SDK
    ErrAuth         = Error{1, "Unauthorized User",
        `This will be returned if an invalid username is used in the request,
        or if the username does not have the rights to modify the resource.`}
    ErrJson         = Error{2, "Body contains invalid JSON.",
        "This will be returned if the body of the message contains invalid JSON."}
    ErrResource     = Error{3, "Resource, <resource>, not available.",
        `This will be returned if the addressed resource does not exist.
         E.g. the user specifies a light ID that does not exist.`}
    ErrMethod       = Error{4, "Method, <method_name>, not available for resource, <resource>",
        `This will be returned if the method (GET/POST/PUT/DELETE)
         used is not supported by the URL e.g. DELETE is not
         supported on the /config resource`}
    ErrParamMissing = Error{5, "Missing parameters in body.", `Will be returned if
         required parameters are not present in the message body. The presence
         of invalid parameters should not trigger this error as long as all
         required parameters are present.`}
    ErrParamNA      = Error{6, "Parameter, <parameter>, not available.",
        `This will be returned if a parameter sent in the message body does
         not exist. This error is specific to PUT commands; invalid parameters
         in other commands are simply ignored.`}
    ErrParamInvalid = Error{7, "Invalid value, <value>, for parameter, <parameter>",
        `This will be returned if the value set for a parameter is of the
         incorrect format or is out of range.`}
    ErrParamStatic  = Error{8, "Parameter, <parameter>, is not modifiable",
        `This will be returned if an attempt to modify
         a read only parameter is made.`}
    ErrItemOverflow = Error{11, "Too many items to list.",
        "List in request contains too many items"}
    ErrPortalConn   = Error{12, "Portal connection required.",
        `Command requires portal connection.
        Returned if portalservices is “false“ or the portal connection is down`}
    ErrorInternal   = Error{901, "Internal error, <error code>",
        `This will be returned if there is an internal error in the
         processing of the command. This indicates an error in the
         bridge, not in the message being sent.`}

    // Command Specific Errors from Hue SDK
    ErrLink   = Error{101, "Link button not pressed.",
        `/config/linkbutton is false. Link button has
         not been pressed in last 30 seconds.`}
    ErrDHCP   = Error{110, "DHCP cannot be disabled.",
        "DHCP can only be disabled if there is a valid static IP configuration"}
    ErrUpdate = Error{111, "Invalid updatestate.",
        "Checkforupdate can only be set in updatestate 0 and 1."}
    // TODO: Need to add 201, 301, 305, 306, 402, 403, 501, 502, 601...
)

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
