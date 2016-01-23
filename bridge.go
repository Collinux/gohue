package main

import (
    "log"
)

func main() {
    log.Println("Test")
}

type Bridge struct {
    IPAddress   string
    Username    string
    Debug       bool
    Info        Description
}

type Description struct {
    DeviceType          string
    FriendlyName        string
    ManufacturerURL     string
    About               string
    ModelName           string
    ModelNumber         string
    ModelURL            string
    SerialNumber        string
    UDN                 string
}
