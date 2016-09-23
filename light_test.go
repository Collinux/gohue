/*
* light_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"github.com/collinux/GoHue"
	"testing"
	"time"
)

func TestSetLightState(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	bridge.Login("427de8bd6d49f149c8398e4fc08f")
	nameTest, err := bridge.GetLightByName("Desk Light") // Also tests GetAllLights
	if err != nil {
		t.Fatal(err)
	}
	_ = nameTest
	selectedLight, err := bridge.GetLightByIndex(7)
	if err != nil {
		t.Fatal(err)
	}

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
	selectedLight.Brighten(20)

	selectedLight.SetColor(hue.RED)
	time.Sleep(time.Second)
	selectedLight.SetColor(hue.YELLOW)
	time.Sleep(time.Second)
	selectedLight.SetColor(hue.ORANGE)
	time.Sleep(time.Second)
	selectedLight.SetColor(hue.GREEN)
	time.Sleep(time.Second)
	selectedLight.SetColor(hue.CYAN)
	time.Sleep(time.Second)
	selectedLight.SetColor(hue.BLUE)
	time.Sleep(time.Second)
	selectedLight.SetColor(hue.PURPLE)
	time.Sleep(time.Second)
	selectedLight.SetColor(hue.PINK)
	time.Sleep(time.Second)
	selectedLight.SetColor(hue.WHITE)

	// _ := selectedLight.Delete()
}

func TestFindNewLights(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	bridge.Login("427de8bd6d49f149c8398e4fc08f")
	err = bridge.FindNewLights()
	if err != nil {
		t.Fatal(err)
	}
}
