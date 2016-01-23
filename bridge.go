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
    DeviceType          string  `xml:"deviceType"`
    FriendlyName        string  `xml:"friendlyName"`
    Manufacturer        string  `xml:"manufacturer"`
    ManufacturerURL     string  `xml:"manufacturerURL"`
    ModelDescription    string  `xml:"modelDescription"`
    ModelName           string  `xml:"modelName"`
    ModelNumber         string  `xml:"modelNumber"`
    ModelURL            string  `xml:"modelURL"`
    SerialNumber        string  `xml:"serialNumber"`
    UDN                 string  `xml:"UDN"`
}

func NewBridge(ip string) *Bridge {
    // TODO: if yaml file exists then return
    bridge := Bridge {
        IPAddress: ip,
        // TODO: other defaults here
    }

    bridge.Info = GetBridgeInfo(bridge.IPAddress) // TODO: IP var
    return &bridge
}

// Go to http://<bridge_ip>/description.xml and parse the info
func GetBridgeInfo(ip string) BridgeInfo {
    response, error := http.Get("http://" + ip + "/description.xml")
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

    data := BridgeInfo{}
    error = xml.Unmarshal(body, &data)
    if error != nil {
        log.Println("GetBridgeInfo error. Cannot parse data from description.")
        log.Println(error)
    }

    log.Println(data)

    return data
}
