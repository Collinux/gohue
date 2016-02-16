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
    "time"
)

// Light struct defines attributes of a light.
type Light struct {
    State struct {
        On          bool       `json:"on"`     // On or Off state of the light ("true" or "false")
        Bri         int        `json:"bri"`    // Brightness value 1-254
        Hue         int        `json:"hue"`    // Hue value 1-65535
        Saturation  int        `json:"sat"`    // Saturation value 0-254
        Effect      string     `json:"effect"` // "None" or "Colorloop"
        XY          [2]float32 `json:"xy"`     // Coordinates of color in CIE color space
        CT          int        `json:"ct"`     // Mired Color Temperature (google it)
        Alert       string     `json:"alert"`
        ColorMode   string     `json:"colormode"`
        Reachable   bool       `json:"reachable"`
    } `json:"state"`
    Type             string    `json:"type"`
    Name             string    `json:"name"`
    ModelID          string    `json:"modelid"`
    ManufacturerName string    `json:"manufacturername"`
    UniqueID         string    `json:"uniqueid"`
    SWVersion        string    `json:"swversion"`
    Index            int        // Set by index of light array response // TODO: change to smaller int
    Bridge          *Bridge
}

// LightState used in Light.SetState to amend light attributes.
type LightState struct {
    On                   bool           `json:"on"`
    Bri                  uint8          `json:"bri,omitempty"`
    Hue                  uint16         `json:"hue,omitempty"`
    Sat                  uint8          `json:"sat,omitempty"`
    XY                   *[2]float32    `json:"xy,omitempty"`
    CT                   uint16         `json:"ct,omitempty"`
    Effect               string         `json:"effect,omitempty"`
    Alert                string         `json:"alert,omitempty"`
    TransitionTime       string         `json:"transitiontime,omitempty"`
    SaturationIncrement  int16          `json:"sat_inc,omitempty"`
    HueIncrement         int32          `json:"hue_inc,omitempty"`
    BrightnessIncrement  int16          `json:"bri_inc,omitempty"`
    CTIncrement          int32          `json:"ct_inc,omitempty"`
    XYIncrement          *[2]float32    `json:"xy_inc,omitempty"`
    Name                 string         `json:"name,omitempty"`
}

// Light.SetName assigns a new name in the light's
// attributes as recognized by the bridge.
func (self *Light) SetName(name string) error {
    uri := fmt.Sprintf("/api/%s/lights/%d", self.Bridge.Username, self.Index)
    body := make(map[string]string)
    body["name"] = name
    _, _, err := self.Bridge.Put(uri, body)
    if err != nil {
        return err
    }
    return nil
}

// Light.Off turns the light source off
func (self *Light) Off() error {
    return self.SetState(LightState{On: false})
}

// Light.Off turns the light source on
func (self *Light) On() error {
    return self.SetState(LightState{On: true})
}

// Light.Toggle switches the light source on and off
func (self *Light) Toggle() error {
    if self.State.On {
        return self.Off()
    } else {
        return self.On()
    }
    return nil
}

// Light.Delete removes the light from the
// list of lights available on the bridge.
func (self *Light) Delete() error {
    uri := fmt.Sprintf("/api/%s/lights/%d", self.Bridge.Username, self.Index)
    err := self.Bridge.Delete(uri)
    if err != nil {
        return err
    }
    return nil
}

// Light.Blink increases and decrease the brightness
// repeatedly for a given seconds interval and return the
// light back to its off  or on state afterwards.
// Note: time will vary based on connection speed and algorithm speed.
func (self *Light) Blink(seconds int) error {
    originalPosition := self.State.On
    originalBrightness := self.State.Bri
    blinkMax := LightState{On: true, Bri: uint8(200)}
    blinkMin := LightState{On: true, Bri: uint8(50)}

    // Start with near maximum brightness and toggle between that and
    // a lesser brightness to create a blinking effect.
    err := self.SetState(blinkMax)
    if err != nil {
        return err
    }
    for i := 0; i <= seconds*2; i++ {
        if i % 2 == 0 {
            err = self.SetState(blinkMax)
            if err != nil {
                return err
            }
        } else {
            err = self.SetState(blinkMin)
            if err != nil {
                return err
            }
        }
        time.Sleep(time.Second/2)
    }

    // Return the light to its original on or off state and brightness
    if self.State.Bri != originalBrightness || self.State.On != originalPosition {
        self.SetState(LightState{On: originalPosition, Bri: uint8(originalBrightness)})
    }
    return nil
}

// Light.ColorLoop sets the light state to 'colorloop' if `active`
// is true or it sets the light state to "none" if `activate` is false.
func (self *Light) ColorLoop(activate bool) error {
    var state = "none"
    if activate {
        state = "colorloop"
    }
    return self.SetState(LightState{On: true, Effect: state})
}

// Light.SetState modifyies light attributes. See `LightState` struct for attributes.
// Brightness must be between 1 and 254 (inclusive)
// Hue must be between 0 and 65535 (inclusive)
// Sat must be between 0 and 254 (inclusive)
// See http://www.developers.meethue.com/documentation/lights-api for more info
func (self *Light) SetState(newState LightState) error {
    uri := fmt.Sprintf("/api/%s/lights/%d/state", self.Bridge.Username, self.Index)
    _, _, err := self.Bridge.Put(uri, newState)
    if err != nil {
        return err
    }

    // Get the new light state and update the current Light struct
    *self, err = GetLightByIndex(self.Bridge, self.Index)
    if err != nil {
        return err
    }
    return nil
}

// GetAllLights retreives the state of all lights that the bridge is aware of.
func GetAllLights(bridge *Bridge) ([]Light, error) {
    // Loop through all light indicies to see if they exist
    // and parse their values. Supports 100 lights.
    var lights []Light
    for index := 1; index < 101; index++ {
        light, err := GetLightByIndex(bridge, index)
        if err != nil {
            break
        }
        lights = append(lights, light)
    }
    return lights, nil
}

// GetLightByIndex returns a light struct containing data on
// a light given its index stored on the bridge. This is used for
// quickly updating an individual light.
func GetLightByIndex(bridge *Bridge, index int) (Light, error) {
    // Send an http GET and inspect the response
    uri := fmt.Sprintf("/api/%s/lights/%d", bridge.Username, index)
    body, _, err := bridge.Get(uri)
    if err != nil {
        return Light{}, err
    }
    if strings.Contains(string(body), "not available") {
        return Light{}, errors.New("Index Error")
    }

    // Parse and load the response into the light array
    light := Light{}
    err = json.Unmarshal(body, &light)
    if err != nil {
        trace("", err)
    }
    light.Index = index
    light.Bridge = bridge
    return light, nil
}

// GetLight returns a light struct containing data on a given name.
func GetLightByName(bridge *Bridge, name string) (Light, error) {
    lights, _ := GetAllLights(bridge)
    for _, light := range lights {
        if light.Name == name {
            return light, nil
        }
    }
    return Light{}, errors.New("Light not found.")
}
