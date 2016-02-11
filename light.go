/*
* light.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/

// http://www.developers.meethue.com/documentation/lights-api
package hue

import (
    "fmt"
    "encoding/json"
    "strings"
    "errors"
)

// Light struct defines attributes of a light.
type Light struct {
    State struct {
        On          bool      `json:"on"`     // On or Off state of the light ("true" or "false")
        Bri         int       `json:"bri"`    // Brightness value 1-254
        Hue         int       `json:"hue"`    // Hue value 1-65535
        Saturation  int       `json:"sat"`    // Saturation value 0-254
        Effect      string    `json:"effect"` // "None" or "Colorloop"
        XY          [2]float32 `json:"xy"`    // Coordinates of color in CIE color space
        CT          int       `json:"ct"`     // Mired Color Temperature (google it)
        Alert       string    `json:"alert"`
        ColorMode   string    `json:"colormode"`
        Reachable   bool      `json:"reachable"`
    } `json:"state"`
    Type             string     `json:"type"`
    Name             string     `json:"name"`
    ModelID          string     `json:"modelid"`
    ManufacturerName string     `json:"manufacturername"`
    UniqueID         string     `json:"uniqueid"`
    SWVersion        string     `json:"swversion"`
    Index            int        // Set by index of light array response // TODO: change to smaller int
    Bridge          *Bridge
}

// LightState used in SetLightState to amend light attributes.
type LightState struct {
    On                   bool           `json:"on"`
    Bri                  uint8          `json:"bri,omitempty"`
    Hue                  uint16         `json:"hue,omitempty"`
    Sat                  uint8          `json:"sat,omitempty"`
    XY                   *[2]float32     `json:"xy,omitempty"`
    CT                   uint16         `json:"ct,omitempty"`
    Effect               string         `json:"effect,omitempty"`
    Alert                string         `json:"alert,omitempty"`
    TransitionTime       string         `json:"transitiontime,omitempty"`
    SaturationIncrement  int            `json:"sat_inc,omitempty"` // TODO: -254 to 254
    HueIncrement         int            `json:"hue_inc,omitempty"` // TODO: -65534 to 65534
    BrightnessIncrement  int            `json:"bri_inc,omitempty"` // TODO: -254 to 254
    CTIncrement          int            `json:"ct_inc,omitempty"` // TODO: -65534 to 65534
    XYIncrement          *[2]float32     `json:"xy_inc,omitempty"`
}

func (self *Light) SetName(name string) error {
    uri := fmt.Sprintf("/api/%s/lights/%s", self.Bridge.Username, self.Index)
    body := make(map[string]string)
    body["name"] = name
    _, _, err := self.Bridge.Put(uri, body)
    if err != nil {
        return err
    }
    return nil
}

// Light.Off will turn the light source off
func (self *Light) Off() error {
    return SetLightState(self, LightState{On: false})
}

// Light.Off will turn the light source on
func (self *Light) On() error {
    return SetLightState(self, LightState{On: true})
}

// Light.Toggle will toggle the light source on and off
func (self *Light) Toggle() error {
    if self.State.On {
        return self.Off()
    } else {
        return self.On()
    }
}

// Light.ColorLoop will set the light state to 'colorloop' if `active`
// is true or it will set the light state to "none" if `activate` is false.
func (self *Light) ColorLoop(activate bool) error {
    var state = "none"
    if activate {
        state = "colorloop"
    }
    return SetLightState(self, LightState{On: true, Effect: state})
}

// SetLightState will modify light attributes such as on/off, saturation,
// brightness, and more. See `SetLightState` struct.
func SetLightState(light *Light, newState LightState) error {
    uri := fmt.Sprintf("/api/%s/lights/%d/state", light.Bridge.Username, light.Index)
    _, _, err := light.Bridge.Put(uri, newState)
    return err
}

// GetAllLights retreives the state of all lights that the bridge is aware of.
func GetAllLights(bridge *Bridge) ([]Light, error) {
    // Loop through all light indicies to see if they exist
    // and parse their values. Supports 100 lights.
    var lights []Light
    for index := 1; index < 101; index++ {

        // Send an http GET and inspect the response
        uri := fmt.Sprintf("/api/%s/lights/%d", bridge.Username, index)
        body, _, err := bridge.Get(uri)
        if err != nil {
            return lights, err
        }
        if strings.Contains(string(body), "not available") {
            // Handle end of searchable lights
            //fmt.Printf("\n\n%d lights found.\n\n", index)
            break
        }

        // Parse and load the response into the light array
        data := Light{}
        err = json.Unmarshal(body, &data)
        if err != nil {
            trace("", err)
        }
        data.Index = index
        data.Bridge = bridge
        lights = append(lights, data)
    }
    return lights, nil
}

// GetLight will return a light struct containing data on a given name.
func GetLight(bridge *Bridge, name string) (Light, error) {
    lights, _ := GetAllLights(bridge)
    for index := 0; index < len(lights); index++ {
        if lights[index].Name == name {
            return lights[index], nil
        }
    }
    return Light{}, errors.New("Light not found.")
}
