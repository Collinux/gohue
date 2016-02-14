/*
* schedule.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/

package hue

import (
    //"fmt"
)

type Timer struct {
    Index       int
    Name        string      `json:"name,omitempty"`
    Description string      `json:"description,omitempty"`
    Command     CommandInfo `json:"command,omitempty"`
    Time        string      `json:"time,omitempty"`
    Created     string      `json:"created,omitempty"`
    Status      string      `json:"status,omitempty"`
    AutoDelete  bool        `json:"autodelete,omitempty"`
    StartTime   string      `json:"starttime,omitempty"`
}

type Alarm struct {
    Index       string
    Name        string      `json:"name,omitempty"`
    Description string      `json:"description,omitempty"`
    Command     CommandInfo `json:"command,omitempty"`
    LocalTime   string      `json:"localtime,omitempty"`
    Time        string      `json:"time,omitempty"`
    Created     string      `json:"created,omitempty"`
    Status      string      `json:"status,omitempty"`
    AutoDelete  bool        `json:"autodelete,omitempty"`
}

type CommandInfo struct {
    Address     string      `json:"address,omitempty"`
    Body        string      `json:"body,omitempty"` // TODO: may be diff type
    Method      string      `json:"method,omitempty"`
}

// func (bridge *Bridge) GetSchedules() ([]interface{}, error) {
//     return []interface{}, nil
// }
//
// func (bridge *Bridge) CreateSchedule(schedule interface{}) error {
//     return nil
// }
//
// func (bridge *Bridge) GetSchedule(index int) (interface{}, error) {
//     return []interface{}, nil
// }
//
// func (bridge *Bridge) SetSchedule(index int, schedule interface{}) error {
//     return nil
// }
//
// func (bridge *Bridge) DeleteSchedule(index int) error {
//     return nil
// }
