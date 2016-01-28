package hue

import (
    "testing"
    "log"
)

func TestGetAllLights(t *testing.T) {
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    GetAllLights(bridge)
}

func TestGetLight(t *testing.T) {
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    light , err := GetLight(bridge, "Bathroom Light")
    if err != nil {
        trace("", err)
    }
    log.Println("LIGHT: ", light)
}
