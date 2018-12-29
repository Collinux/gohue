/*
* sensor_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"github.com/collinux/GoHue"

	"testing"
	"os"
	"fmt"
)

func TestGetAllSensors(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	if os.Getenv("HUE_USER_TOKEN") == "" {
		t.Fatal("The environment variable HUE_USER_TOKEN must be set to the value from bridge.CreateUser")
	}
	bridge.Login(os.Getenv("HUE_USER_TOKEN"))

	sensors, err := bridge.GetAllSensors()
	if err != nil {
		t.Fatal(err)
	}

	for _, sensor := range sensors {
		fmt.Println(sensor.Name)
		fmt.Println(sensor.Config.Battery)
		fmt.Println(sensor.State.LastUpdated)
		fmt.Println("------")
	}
}

