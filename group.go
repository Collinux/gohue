/*
* group.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */
// http://www.developers.meethue.com/documentation/groups-api

package hue

import (
	"encoding/json"
	"fmt"
)

// Action struct defines the state of a group
type Action struct {
	Alert     string    `json:"alert,omitempty"`
	Bri       int       `json:"bri,omitempty"`
	Colormode string    `json:"colormode,omitempty"`
	Ct        int       `json:"ct,omitempty"`
	Effect    string    `json:"effect,omitempty"`
	Hue       *int      `json:"hue,omitempty"`
	On        *bool     `json:"on,omitempty"`
	Sat       *int      `json:"sat,omitempty"`
	XY        []float64 `json:"xy,omitempty"`
	Scene     string    `json:"scene,omitempty"`
}

// Group struct defines the attributes for a group of lights.
type Group struct {
	Action Action   `json:"action"`
	Lights []string `json:"lights"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
}

// GetGroups gets the attributes for each group of lights.
// TODO: NOT TESTED, NOT FULLY IMPLEMENTED
func (bridge *Bridge) GetGroups() ([]Group, error) {
	uri := fmt.Sprintf("/api/%s/groups", bridge.Username)
	body, _, err := bridge.Get(uri)
	if err != nil {
		return []Group{}, err
	}

	//fmt.Println("GROUP GET: ", string(body))

	groups := map[string]Group{}
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return []Group{}, err
	}
	//fmt.Println("GROUPS: ", groups)

	return []Group{}, nil
}

// SetGroupState sends an action to group
func (bridge *Bridge) SetGroupState(group int, action *Action) error {
	uri := fmt.Sprintf("/api/%s/groups/%d/action", bridge.Username, group)
	_, _, err := bridge.Put(uri, action)
	if err != nil {
		return err
	}
	return nil
}
