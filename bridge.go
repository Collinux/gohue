package main

import (
    "log"
    "os"
    "encoding/xml"
    "net/http"
    "io/ioutil"
)

func main() {
    log.Println("Test")
}

type Bridge struct {
    IPAddress   string
    Username    string
    Info        BridgeInfo
}

type BridgeInfo struct {
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
        IPAddress: ip,
        // TODO: other defaults here
    }

    bridge.Info = GetBridgeInfo("192.168.1.128") // TODO: IP var
    return &bridge
}

// Go to http://<bridge_ip>/description.xml and parse the info
func GetBridgeInfo(ip string) BridgeInfo {
    response, error := http.Get(ip + "/description.xml")
    if response.StatusCode != 200 {
        log.Println("Bridge status error: %d", response.StatusCode)
        os.Exit(1)
    }

    body, error := ioutil.ReadAll(response.Body)
    if error != nil {
        log.Println("Error parsing bridge description XML")
        os.Exit(1)
    }

    data := make(map[string]string)
    error = xml.Unmarshal(body, data)
    if error != nil {
        log.Println("GetBridgeInfo error. Cannot parse data from description.")
    }

    log.Println(data["FriendlyName"])

    return BridgeInfo{} // TODO: fill in
}
