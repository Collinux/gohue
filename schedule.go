/*
* schedule.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/

package hue

import (
    "fmt"
    "encoding/json"
)

// Schedule struct defines attributes of Alarms and Timers
type Schedule struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
	Command struct {
		Address string   `json:"address"`
		Body    struct {
		    Scene string `json:"scene"`
		} `json:"body"`
		Method string    `json:"method"`
	} `json:"command"`
    Localtime   string   `json:"localtime"`
    Time        string   `json:"time"`
	Created     string   `json:"created"`
	Status      string   `json:"status"`
	Autodelete  bool     `json:"autodelete"`
    ID          string
}

// Bridge.GetSchedules will get Alarms and Timers in a Schedule struct.
func (bridge *Bridge) GetSchedules() ([]Schedule, error) {
    uri := fmt.Sprintf("/api/%s/schedules", bridge.Username)
    body, _, err := bridge.Get(uri)
    if err != nil {
        return []Schedule{}, err
    }

    schedules := map[string]Schedule{}
    err = json.Unmarshal(body, &schedules)
	if err != nil {
		return []Schedule{}, err
	}

    scheduleList := []Schedule{}
    for key, value := range schedules {
        schedule := Schedule{}
        schedule = value
        schedule.ID = key
        scheduleList = append(scheduleList, schedule)
    }

    // for sched := range scheduleList {
    //     fmt.Println("\n\nScheduleoutput: ", scheduleList[sched])
    // }

    return scheduleList, nil
}

// Bridge.GetSchedule will get the attributes for an individual schedule.
// This is used as to optimize time when updating the state of a schedule item.
func (bridge *Bridge) GetSchedule(id string) (Schedule, error) {
    uri := fmt.Sprintf("/api/%s/schedules/%s", bridge.Username, id)
    body, _, err := bridge.Get(uri)
    if err != nil {
        return Schedule{}, err
    }
    return Schedule{}, nil
}

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
