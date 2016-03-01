/*
* bridge_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */

package hue

import (
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	bridge, _ := NewBridge("192.168.1.128")
	username, _ := bridge.CreateUser("test")
	bridge.Login(username)
	//bridge.DeleteUser(bridge.Username)
}

func TestFindBridges(t *testing.T) {
	bridges, _ := FindBridges()
	fmt.Println(bridges[0].IPAddress)
}

func TestBridgeLogin(t *testing.T) {
	bridges, err := FindBridges()
	if err != nil {
		fmt.Println("Error on TestBridgeLogin")
	}
	bridges[0].Login("427de8bd6d49f149c8398e4fc08f")

}
