/*
* light_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collinux[-at-]users.noreply.github.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"github.com/collinux/gohue"
	"testing"
	"time"
	"os"
)

func TestSetLightState(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	bridge.Login(os.Getenv("HUE_USER_TOKEN"))
	//nameTest, err := bridge.GetLightByName("Desk Light") // Also tests GetAllLights
	//if err != nil {
	//	t.Fatal(err)
	//}
	//_ = nameTest
	selectedLight, err := bridge.GetLightByIndex(1)
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

	// Skip validation of colors
	//selectedLight.SetColor(hue.RED)
	//time.Sleep(time.Second)
	//selectedLight.SetColor(hue.YELLOW)
	//time.Sleep(time.Second)
	//selectedLight.SetColor(hue.GREEN)
	//time.Sleep(time.Second)
	//selectedLight.SetColor(hue.WHITE)
	//time.Sleep(time.Second)
	//selectedLight.Off()

	// TODO
	// Skip validation of deleting and re-adding light
	// _ := selectedLight.Delete()
}

func TestFindNewLights(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	if os.Getenv("HUE_USER_TOKEN") == "" {
		t.Fatal("The environment variable HUE_USER_TOKEN must be set to the value from bridge.CreateUser")
	}
	bridge.Login(os.Getenv("HUE_USER_TOKEN"))
	err = bridge.FindNewLights()
	if err != nil {
		t.Fatal(err)
	}
}
