package main

import (
    "log"
)

type Light struct {
    State struct {
        On        bool
        Bri       int
        hue       int
        sat       int
        effect    string     // TODO: what is this?
        xy        []string   // TODO: what is this?
        ct        int        // TODO: what is this?
        alert     string
        colormode string
        reachable bool
    }
    Type             string
    Name             string
    ModelID          string
    ManufacturerName string
    UniqueID         string
    SWVersion        string  // TODO: what is this?
}

func GetAllLights() {
    
}