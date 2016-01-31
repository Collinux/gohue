package hue

import (
    "testing"
    "fmt"
)

func TestGetAllLights(t *testing.T) {
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    GetAllLights(bridge)
}

func TestGetLight(t *testing.T) {
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    GetLight(bridge, "Bathroom Light")
}

func TestSetLightState(t *testing.T) {
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    randomLight := GetAllLights(bridge)[0]
    fmt.Println(randomLight.Name)
    newState := LightState{On: true}
    SetLightState(bridge, randomLight.UniqueID, newState)
}
