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
    fmt.Println("\nTESTING LIGHT STATE:\n\n")
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    lights, _ := GetAllLights(bridge)
    selectedLight := lights[0]
    newState := LightState{On: true,} //On: false, *XY: [2]float32{5.0, 5.0},
    //fmt.Println("\n\nSTATE: ", newState)
    SetLightState(bridge, selectedLight.Index, newState)
}
