/*
* light_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/

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

func TestGetLightByName(t *testing.T) {
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    GetLightByName(bridge, "Bathroom Light")
}

func TestSetLightState(t *testing.T) {
    fmt.Println("\nTESTING LIGHT STATE:\n\n")
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    lights, _ := GetAllLights(bridge)
    selectedLight := lights[5]

    selectedLight.On()
    time.Sleep(time.Second)
    selectedLight.Off()
    time.Sleep(time.Second)
    selectedLight.Toggle()
    time.Sleep(time.Second)

    selectedLight.ColorLoop(false)

    selectedLight.SetName(selectedLight.Name)

    selectedLight.Blink(2)

    // err := selectedLight.Delete()
    // if err != nil {
    //     fmt.Println("error on delete")
    // }
}
