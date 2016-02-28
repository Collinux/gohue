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

// Group struct defines the attributes for a group of lights.
type Group struct {
	Action struct {
		Alert     string    `json:"alert"`
		Bri       int       `json:"bri"`
		Colormode string    `json:"colormode"`
		Ct        int       `json:"ct"`
		Effect    string    `json:"effect"`
		Hue       int       `json:"hue"`
		On        bool      `json:"on"`
		Sat       int       `json:"sat"`
		XY        []float64 `json:"xy"`
	} `json:"action"`
	Lights []string `json:"lights"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
}

// Bridge.GetGroups gets the attributes for each group of lights.
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
