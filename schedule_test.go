/*
* schedule_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collinux[-at-]users.noreply.github.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"github.com/collinux/gohue"
	"os"
	"testing"
)

func TestGetAllSchedules(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	if os.Getenv("HUE_USER_TOKEN") == "" {
		t.Fatal("The environment variable HUE_USER_TOKEN must be set to the value from bridge.CreateUser")
	}
	bridge.Login(os.Getenv("HUE_USER_TOKEN"))
	schedules, err := bridge.GetAllSchedules()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(schedules)
}

func TestGetSchedule(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	if os.Getenv("HUE_USER_TOKEN") == "" {
		t.Fatal("The environment variable HUE_USER_TOKEN must be set to the value from bridge.CreateUser")
	}
	bridge.Login(os.Getenv("HUE_USER_TOKEN"))
	schedules, err := bridge.GetAllSchedules()
	if err != nil {
		t.Fatal(err)
	}
	schedule, err := bridge.GetSchedule(schedules[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(schedule)
}

// TODO
// func TestCreateSchedule(t *testing.T) { }
