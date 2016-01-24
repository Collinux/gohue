package main

import (
    "log"
    "fmt"
    "os"
    "net/http"
    "io/ioutil"
    "encoding/json"
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
    }
    Type             string     `json:"type"`
    Name             string     `json:"name"`
    ModelID          string     `json:"modelid"`
    ManufacturerName string     `json:"manufacturername"`
    UniqueID         string     `json:"uniqueid"`
    SWVersion        string     `json:"swversion"`
}

//http://192.168.1.128/api/319b36233bd2328f3e40731b23479207/lights/
// http://<bridge_ip>/api/<username>/lights/
func GetAllLights(bridge *Bridge) {
    response, error := http.Get(
        fmt.Sprintf("http://%s/api/%s/lights/", bridge.IPAddress, bridge.Username)) 
    if error != nil {
        trace("", error)
        os.Exit(1)
    } else if response.StatusCode != 200 {
        trace(fmt.Sprintf("Bridge status error %d", response.StatusCode), nil)
        os.Exit(1)
    }
    
    body, error := ioutil.ReadAll(response.Body)
    defer response.Body.Close()
    if error != nil {
        trace("", error)
        os.Exit(1)
    }
    
    data := []Light{}
    error = json.Unmarshal(body, &data)
    if error != nil {
        trace("", error)
        os.Exit(1)
    }
    
    log.Println(data)
    
}