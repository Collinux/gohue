/*
* scene.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/

package hue

import (
    "fmt"
    "encoding/json"
)

type Scene struct {
    Appdata struct {
    	Data    string   `json:"data"`
    	Version int      `json:"version"`
    } `json:"appdata"`
    Lastupdated string   `json:"lastupdated"`
    Lights      []string `json:"lights"`
    Locked      bool     `json:"locked"`
    Name        string   `json:"name"`
    Owner       string   `json:"owner"`
    Picture     string   `json:"picture"`
    Recycle     bool     `json:"recycle"`
    Version     int      `json:"version"`
    ID          string
}

// Bridge.GetScenes will get attributes for all scenes.
func (bridge *Bridge) GetScenes() ([]Scene, error) {
    uri := fmt.Sprintf("/api/%s/scenes", bridge.Username)
    body, _, err := bridge.Get(uri)
    if err != nil  {
        return []Scene{}, err
    }

    fmt.Sprintf("GET SCENES BODY: ", string(body))

    scenes := map[string]Scene{}
    err = json.Unmarshal(body, &scenes)
    if err != nil  {
        return []Scene{}, err
    }
    scenesList := []Scene{}
    for key, value := range scenes {
        scene := Scene{}
        scene = value
        scene.ID = key
        scenesList = append(scenesList, scene)
    }
    return scenesList, nil
}
