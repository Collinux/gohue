package main

import (
    "log"
    "os"
    "encoding/xml"
    "net/http"
    "io/ioutil"
    "runtime"
    "fmt"
)

func main() {
    bridge := NewBridge("192.168.1.128")
    log.Println(bridge.IPAddress)

    //trace("this is a trace")
}

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

func NewBridge(ip string) *Bridge {
    bridge := Bridge {
        IPAddress: ip,
        // TODO: other defaults here
    }

    GetBridgeInfo(&bridge)
    return &bridge
}

// Go to http://<bridge_ip>/description.xml set the bridge.Info
func GetBridgeInfo(self *Bridge) {
    response, error := http.Get("http://" + self.IPAddress + "/description.xml")
    if error != nil {
        trace("", error)
    } else if response.StatusCode != 200 {
        trace(fmt.Sprintf("Bridge status error: %d", response.StatusCode), nil)
        os.Exit(1)
    }

    body, error := ioutil.ReadAll(response.Body)
    defer response.Body.Close()
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
