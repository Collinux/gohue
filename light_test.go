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
    //"time"
)

func TestSetLightState(t *testing.T) {
    bridge, _ := NewBridge("192.168.1.128")
    bridge.Login("427de8bd6d49f149c8398e4fc08f")
    nameTest, _ := bridge.GetLightByName("Desk Light")  // Also tests GetAllLights
    _ = nameTest
    selectedLight, _ := bridge.GetLightByIndex(7)

    selectedLight.On()
    selectedLight.SetBrightness(100)
    time.Sleep(time.Second)
    selectedLight.Off()
    time.Sleep(time.Second)
    selectedLight.Toggle()
    time.Sleep(time.Second)

    selectedLight.ColorLoop(false)

    selectedLight.SetName(selectedLight.Name)

    selectedLight.Blink(3)

    selectedLight.Dim(20)

    // selectedLight.SetColor(RED)
    // time.Sleep(time.Second)
    // selectedLight.SetColor(YELLOW)
    // time.Sleep(time.Second)
    // selectedLight.SetColor(ORANGE)
    // time.Sleep(time.Second)
    // selectedLight.SetColor(GREEN)
    // time.Sleep(time.Second)
    // selectedLight.SetColor(CYAN)
    // time.Sleep(time.Second)
    // selectedLight.SetColor(BLUE)
    // time.Sleep(time.Second)
    // selectedLight.SetColor(PURPLE)
    // time.Sleep(time.Second)
    // selectedLight.SetColor(PINK)
    // time.Sleep(time.Second)
    // selectedLight.SetColor(WHITE)

    // _ := selectedLight.Delete()
}
