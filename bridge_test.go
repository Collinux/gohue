/*
* bridge_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collinux[-at-]users.noreply.github.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"github.com/collinux/GoHue"
	"testing"
	"os"
)

func TestCreateUser(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	bridge := bridges[0]
	username, _ := bridge.CreateUser("test")
	bridge.Login(username)
	//bridge.DeleteUser(bridge.Username)
}

func TestFindBridges(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bridges)
}

func TestBridgeLogin(t *testing.T) {
	bridges, err := hue.FindBridges()
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("HUE_USER_TOKEN") == "" {
		t.Fatal("The environment variable HUE_USER_TOKEN must be set to the value from bridge.CreateUser")
	}
	bridges[0].Login(os.Getenv("HUE_USER_TOKEN"))

}
