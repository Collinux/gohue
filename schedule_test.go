/*
* schedule_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"testing"
)

func TestGetAllSchedules(t *testing.T) {
	bridge, _ := NewBridge("192.168.1.128")
	bridge.Login("427de8bd6d49f149c8398e4fc08f")
	_, _ = bridge.GetAllSchedules()
}

func TestGetSchedule(t *testing.T) {
	bridge, _ := NewBridge("192.168.1.128")
	bridge.Login("427de8bd6d49f149c8398e4fc08f")
	_, _ = bridge.GetSchedule("4673980164949558")
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
