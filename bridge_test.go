package hue

import (
    "log"
    "testing"
)

func TestNewBridge(t *testing.T) {
    bridge := NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207")
    log.Println(bridge)
}

func TestCreateUser(t *testing.T) {
    user, _ := CreateUser(NewBridge("192.168.1.128", "319b36233bd2328f3e40731b23479207"), "linux")
    log.Println(user)
}