package main

import (
    "log"
    "encoding/xml"
)

func main() {
    log.Println("Test")
}

type Bridge struct {
    IPAddress   string
    Username    string
    Info        Info
}

Info struct {
    DeviceType          string
    FriendlyName        string
    ManufacturerURL     string
    ModelDescription    string
    ModelName           string
    ModelNumber         string
    ModelURL            string
    SerialNumber        string
    UDN                 string
}

func NewBridge(ip string) *Bridge {
    // TODO: if yaml file exists then return
    bridge := Bridge {
        IPAddress: ip
        // TODO: other defaults here
    }

    response := GetBridgeInfo()
    bridge.Info = GetBridgeInfo("192.168.1.128") // TODO: IP var
}

// Go to http://<bridge_ip>/description.xml and parse the info
func GetBridgeInfo(ip string) Info {
    data := Info{}

    err := xml.Unmarshal([]byte(data), &data)
    if err != nil {
        fmt.Println("GetBridgeInfo error. Cannot parse data from description.")
    }

    log.Println("FriendlyName")
}
