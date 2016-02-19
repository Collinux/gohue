/*
* scene_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/

package hue

import (
    "testing"
    //"fmt"
)

func TestGetAllScenes(t *testing.T) {
    bridge, _ := NewBridge("192.168.1.128")
    bridge.Login("427de8bd6d49f149c8398e4fc08f")
    scenes, _ := bridge.GetAllScenes()
    // for scene := range scenes {
    //     fmt.Println("SCENE: ", scenes[scene])
    // }

    individual, _ := bridge.GetScene(scenes[0].ID)
    _ = individual
    //fmt.Println("Individual scene: ", individual)
}

func TestCreateScene(t *testing.T) {
    bridge, _ := NewBridge("192.168.1.128")
    bridge.Login("427de8bd6d49f149c8398e4fc08f")
    scene := Scene{Lights: []string{"1", "2"}}
    _ = bridge.CreateScene(scene)
}
