package main

import (
    "log"
    "os"
    "encoding/xml"
    "net/http"
    "io/ioutil"
)

func main() {
    bridge := NewBridge("192.168.1.128")
    log.Println(bridge.IPAddress)
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
    if response.StatusCode != 200 {
        log.Println("Bridge status error: %d", response.StatusCode)
        os.Exit(1)
    } else if error != nil {
        log.Println("Error: ", error)
        os.Exit(1)
    }

    body, error := ioutil.ReadAll(response.Body)
    defer response.Body.Close()
    if error != nil {
        log.Println("Error parsing bridge description XML")
        os.Exit(1)
    }

    data := new(BridgeInfo)
    error = xml.Unmarshal(body, &data)
    if error != nil {
        log.Println("GetBridgeInfo error. Cannot parse data from description.")
        log.Println(error)
    }
    self.Info = *data
}
