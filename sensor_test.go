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
	"fmt"
)

func TestGetAllSensors(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	bridge.Login("427de8bd6d49f149c8398e4fc08f")

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

