package hue

import (
    "testing"
    "fmt"
    "time"
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
    fmt.Println("\nTESTING LIGHT STATE:\n\n")
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    lights, _ := GetAllLights(bridge)
    selectedLight := lights[0]

    selectedLight.TurnOn()
    time.Sleep(time.Second)
    selectedLight.TurnOff()
    time.Sleep(time.Second)
    selectedLight.TurnOn()
}
