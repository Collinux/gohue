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
    "errors"
)

type Bridge struct {
    IPAddress   string
    Username    string
    Info        BridgeInfo
}

type BridgeInfo struct {
    Info struct {
        Info struct {
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
    } `xml:"root"`
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
    // Generic Errors
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

    // TODO: Command specific error numbers and descriptions
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
    if err != nil {
        trace("", err)
    }

    result := string(body)
    errFound := strings.Contains(result, "error")
    noLink := strings.Contains(result, "link button not pressed")
    if errFound && noLink {
        return "", errors.New("Link button not pressed.")
    }

    // TODO: decode and return
    // TODO: handle errors. http://www.developers.meethue.com/documentation/error-messages

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
