package hue

import (
    "log"
    "fmt"
    "os"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strings"
)

type Light struct {
    State struct {
        On          bool      `json:"on"`
        Bri         int       `json:"bri"`
        hue         int       `json:"hue"`
        sat         int       `json:"sat"`
        effect      string    `json:"effect"`
        xy          []string  `json:"xy"`     // TODO: what is this?
        ct          int       `json:"ct"`     // TODO: what is this?
        alert       string    `json:"alert"`
        colormode   string    `json:"colormode"`
        reachable   bool      `json:"reachable"`
    } `json:"state"`
    Type             string     `json:"type"`
    Name             string     `json:"name"`
    ModelID          string     `json:"modelid"`
    ManufacturerName string     `json:"manufacturername"`
    UniqueID         string     `json:"uniqueid"`
    SWVersion        string     `json:"swversion"`
}



//http://192.168.1.128/api/319b36233bd2328f3e40731b23479207/lights/

// GetAllLights retreives the state of all lights that the bridge is aware of.
func GetAllLights(bridge *Bridge) {
    // Loop through all light indicies to see if they exist
    // and parse their values. Supports 100 lights.
    var lights []Light
    for index := 1; index < 101; index++ {
        response, error := http.Get(
            fmt.Sprintf("http://%s/api/%s/lights/%d", bridge.IPAddress, bridge.Username, index))
        if error != nil {
            trace("", error)
            os.Exit(1)
        } else if response.StatusCode != 200 {
            trace(fmt.Sprintf("Bridge status error %d", response.StatusCode), nil)
            os.Exit(1)
        }

        // Read the response
        body, error := ioutil.ReadAll(response.Body)
        defer response.Body.Close()
        if error != nil {
            if strings.Contains(string(body), "not available") {
                // Handle end of searchable lights
                fmt.Printf("\n\n%d lights found.\n\n", index)
                break
            } else {
                // Other error found while parsing
                trace("", error)
                os.Exit(1)
            }
        }

        // Parse and load the response into the light array
        data := Light{}
        error = json.Unmarshal(body, &data)
        if error != nil {
            trace("", error)
            os.Exit(1)
        }
        lights = append(lights, data)
    }
    return lights
}
