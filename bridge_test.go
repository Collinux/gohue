/*
* bridge_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	bridges, err := FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	username, _ := bridge.CreateUser("test")
	bridge.Login(username)
	//bridge.DeleteUser(bridge.Username)
}

func TestFindBridges(t *testing.T) {
	bridges, err := FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bridges)
}

func TestBridgeLogin(t *testing.T) {
	bridges, err := FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridges[0].Login("427de8bd6d49f149c8398e4fc08f")

}
