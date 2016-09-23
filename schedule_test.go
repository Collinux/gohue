/*
* schedule_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"github.com/collinux/GoHue"
	"testing"
)

func TestGetAllSchedules(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	bridge.Login("427de8bd6d49f149c8398e4fc08f")
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
	bridge.Login("427de8bd6d49f149c8398e4fc08f")
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

// func TestCreateSchedule(t *testing.T) {
//     bridge, _ := NewBridge("192.168.1.128", "427de8bd6d49f149c8398e4fc08f")
//     schedule := Schedule{
//         Command: interface{
//             Address: "/api/fffffffff6338294fffffffff585692d/groups/0/action",
//             Body: {
//                 Scene: "e5e33fdf2-off-0" // scene id
//             }
//             Method: "PUT",
//         }
//     }
//     bridge.CreateSchedule()
// }
