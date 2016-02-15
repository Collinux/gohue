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

func TestGetSchedules(t *testing.T) {
    bridge, _ := NewBridge("192.168.1.128", "427de8bd6d49f149c8398e4fc08f")
    _, _ = bridge.GetSchedules()
}

func TestGetSchedule(t *testing.T) {
    bridge, _ := NewBridge("192.168.1.128", "427de8bd6d49f149c8398e4fc08f")
    _, _ = bridge.GetSchedule("4673980164949558")
}
