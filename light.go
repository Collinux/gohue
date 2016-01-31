// http://www.developers.meethue.com/documentation/lights-api

package hue

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strings"
    "errors"
    "bytes"
)

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
}

// LightState used in SetLightState to ammend light attributes.
type LightState struct {
    On                   bool
    Bri                  uint8
    Hue                  uint16
    Sat                  uint8
    XY                   [2]float32
    CT                   uint16
    Alert                string
    Effect               string
    TransitionTime       string
    BrightnessIncrement  int // TODO: -254 to 254
    SaturationIncrement  int // TODO: -254 to 254
    HueIncrement         int // TODO: -65534 to 65534
    CTIncrement          int // TODO: -65534 to 65534
    XYIncrement          [2]float32
}

// SetLightState will modify light attributes such as on/off, saturation,
// brightness, and more. See `SetLightState` struct.
func SetLightState(bridge *Bridge, lightID string, newState LightState) error {
    // Construct the http POST
    req, err := json.Marshal(newState)
    if err != nil {
        trace("", err)
        return err
    }

    // Send the request and read the response
    uri := fmt.Sprintf("http://%s/api/%s/lights/%s/state",
        bridge.IPAddress, bridge.Username, lightID)
    resp, err := http.Post(uri, "text/json", bytes.NewReader(req))
    if err != nil {
        trace("", err)
        return err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        trace("", err)
        return err
    }

    // TODO: Parse the response and return any error
    fmt.Println(string(body))
    return nil
}

// GetAllLights retreives the state of all lights that the bridge is aware of.
func GetAllLights(bridge *Bridge) ([]Light, error) {
    // Loop through all light indicies to see if they exist
    // and parse their values. Supports 100 lights.
    var lights []Light
    for index := 1; index < 101; index++ {

        // Send an http GET and inspect the response
        resp, err := http.Get(
            fmt.Sprintf("http://%s/api/%s/lights/%d", bridge.IPAddress, bridge.Username, index))
        if err != nil {
            trace("", err)
            return lights, err
        } else if resp.StatusCode != 200 {
            trace(fmt.Sprintf("Bridge status error %d", resp.StatusCode), nil)
        }

        // Read and inspect the response content
        body, err := ioutil.ReadAll(resp.Body)
        defer resp.Body.Close()
        if err != nil {
            trace("", err)
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
