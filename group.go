/*
* group.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* Copyright (C) 2018 Niels de Vos (nixpanic) niels@nixpanic.net
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */
// http://www.developers.meethue.com/documentation/groups-api

package hue

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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

// State of a group, allows for detecting if (some) lighs are on.
type State struct {
	AllOn bool `json:"all_on"`
	AnyOn bool `json:"any_on"`
}

// Group struct defines the attributes for a group of lights.
type Group struct {
	Action    Action  ` json:"action"`
	// LightsInt is the string (Index) representation of a `Light`. Use the
	// `Lights` array to access details of the `Light`s
	LightsInt []string `json:"lights"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	// Class of the group, `Room` or `LightGroup` (unknown if others exist).
	Class     string   `json:"class"`
	Recycle   bool     `json:"recycle"`

	// Status of the lights in the group.
	State     State    `json:"state"`

	// Index of this Groups in the list as returned by the Hue bridge.
	Index     int

	// The bridge itsef
	Bridge    *Bridge
	// The Light instances in this Group
	Lights    []Light
}

// GetAllGroups gets the attributes for each group of lights.
func (bridge *Bridge) GetAllGroups() ([]Group, error) {
	uri := fmt.Sprintf("/api/%s/groups", bridge.Username)
	body, _, err := bridge.Get(uri)
	if err != nil {
		return nil, err
	}

	groupMap := map[string]Group{}
	err = json.Unmarshal(body, &groupMap)
	if err != nil {
		return nil, err
	}

	// Ideally Group.Lights should be filled with actual light objects
	lights, err := bridge.GetAllLights()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to get lights for groups: %s", err))
	}

	// Parse the index, add the light to the list, and return the array
	groups := []Group{}
	for index, group := range groupMap {
		group.Index, err = strconv.Atoi(index)
		if err != nil {
			return []Group{}, errors.New("Unable to convert group index to integer. ")
		}

		group.Bridge = bridge

		// this is not optimal, but we'll assume there won't be hundreds of lights
		for _, lightStr := range(group.LightsInt) {
			lightInt, err := strconv.Atoi(lightStr)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("failed to convert group.light.index (%s) to int", lightStr))
			}

			for _, light := range(lights) {
				if light.Index == lightInt {
					group.Lights = append(group.Lights, light)
					break
				}
			}
		}

		groups = append(groups, group)
	}

	return groups, nil
}

// GetGroupByName() uses GetAllGroups() to fetch all the groups, and returns
// the group with the given name if it exists.
func (bridge *Bridge) GetGroupByName(name string) (*Group, error) {
	groups, err := bridge.GetAllGroups()
	if err != nil {
		return nil, errors.New("failed to get groups")
	}

	for _, group := range(groups) {
		if group.Name == name {
			return &group, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("group \"%s\" not found", name))
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

// The request to send to the Hue bridge to create a new group.
type newGroupRequest struct {
	Lights []string `json:"lights,omitempty"`
	Name   string   `json:"name,omitempty"`
	Type   string   `json:"type,omitempty"`
	Class  string   `json:"class,omitempty"`
}

// Create a new group with `name`. The `type` of the group will be `Room` with
// `class` defaulting to `Other`.
func (bridge *Bridge) NewGroup(name string, class string, lights []Light) (*Group, error) {
	uri := fmt.Sprintf("/api/%s/groups", bridge.Username)

	lightsStr := []string{}
	for _, light := range(lights) {
		lightsStr = append(lightsStr, fmt.Sprintf("%d", light.Index))
	}

	req := newGroupRequest{
		Name: name,
		Type: "Room", // TODO: by default the new group is a room (LightGroup exists too)
		Lights: lightsStr,
	}

	if class != "" {
		req.Class = class
	}

	_, _, err := bridge.Post(uri, req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("group %s has been created with %d lights\n", name, len(lights))

	// TODO: return the new group
	return nil, nil
}

// Delete the group from the Hue bridge.
func (group *Group) Delete() error {
	uri := fmt.Sprintf("/api/%s/groups/%d", group.Bridge.Username, group.Index)
	err := group.Bridge.Delete(uri)
	if err != nil {
		return err
	}
	return nil
}

// Turn the lights in the group on.
func (group *Group) On() error {
	on := true
	onAction := Action{
		On: &on,
	}

	return group.Bridge.SetGroupState(group.Index, &onAction)
}

// Turn the lights in the group off.
func (group *Group) Off() error {
	on := false
	offAction := Action{
		On: &on,
	}

	return group.Bridge.SetGroupState(group.Index, &offAction)
}
