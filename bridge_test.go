/*
* bridge_test.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
*/

package hue

import (
    "testing"
    "fmt"
)

func TestCreateUser(t *testing.T) {
    bridge, _ := NewBridge("192.168.1.128")
    bridge.CreateUser("test")
    //bridge.DeleteUser(bridge.Username)
}

func TestFindBridge(t *testing.T) {
    bridge, _ := FindBridge()
    fmt.Println(bridge.IPAddress)
}
