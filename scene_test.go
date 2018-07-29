/*
* scene_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"testing"
)

func TestGetAllScenes(t *testing.T) {
	bridges, err := FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	bridge.Login("427de8bd6d49f149c8398e4fc08f")
	scenes, err := bridge.GetAllScenes()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(scenes)
}

// TODO not functional
// func TestCreateScene(t *testing.T) {
// 	bridges, err := FindBridges()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	bridge := bridges[0]
// 	bridge.Login("427de8bd6d49f149c8398e4fc08f")
// 	scene := Scene{Name: "Testing", Lights: []string{"1", "2"}}
// 	err = bridge.CreateScene(scene)
// 	if err != nil {
// t.Fatal(err)
// }
// }
