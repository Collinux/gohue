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

// Scene struct defines attributes for Scene items
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
func (bridge *Bridge) GetAllScenes() ([]Scene, error) {
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

// Bridge.GetScene will get the attributes for an individual scene.
// This is used to optimize time when updating the state of the scene.
// Note: The ID is not an index, it's a unique key generated for each scene.
func (bridge *Bridge) GetScene(id string) (Scene, error) {
    uri := fmt.Sprintf("/api/%s/scenes/%s", bridge.Username, id)
    body, _, err := bridge.Get(uri)
    if err != nil {
        return Scene{}, err
    }

    scene := Scene{}
    err = json.Unmarshal(body, &scene)
    if err != nil {
        return Scene{}, err
    }
    return scene, nil
}
