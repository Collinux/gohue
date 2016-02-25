# GoHue
Package hue interfaces Philips Hue devices to control lights, scenes, schedules, and groups.

[![GoDoc](https://camo.githubusercontent.com/b3b2a2b7fad4e76052830945cd839a3bba5be723/687474703a2f2f696d672e736869656c64732e696f2f62616467652f676f646f632d7265666572656e63652d3532373242342e706e67)](https://godoc.org/github.com/Collinux/GoHue)

## Installation
```
go get github.com/collinux/gohue
```

## Usage
```
package main

import (
    "github.com/collinux/gohue"
)

func main() {
    bridge, _ := hue.NewBridge("192.168.1.128")
    bridge.Login("new_user")

    lights, _ := bridge.GetAllLights()
    for _, light := range lights {
        light.On()
        light.SetBrightness(100)
        light.ColorLoop(true)
    }

    nightstandLight, _ := bridge.GetLightByName("Nightstand")
    nightstandLight.Blink(5)
    nightstandLight.SetName("Bedroom Lamp")

    lights[0].SetColor(hue.RED)
    lights[1].SetColor(hue.BLUE)
    lights[2].SetColor(hue.GREEN)

    for _, light := range lights {
        light.Off()
    }
}
```

## Features
##### Lights
- [x] Get all lights
- [x] Get light by name
- [x] Get light by index on bridge
- [x] Get lights attributes and state
- [x] Set lights attributes (rename)
- [x] Set light state (color, effects, brightness, etc)
- [x] Delete light
- [x] Turn On, Off, Toggle
- [x] Blink
- [x] Colorloop On/Off

##### Bridge
- [x] Create user
- [x] Delete user
- [x] Get configuration
- [ ] Modify configuration
- [ ] Get full state (datastore)
- [ ] Search for bridges
- [ ] Search for new lights
- [ ] Get all timezones

##### Schedules
- [x] Get all schedules
- [x] Get schedule by ID
- [x] Get schedule attributes
- [ ] Create schedules
- [ ] Set schedule attributes
- [ ] Delete schedule

##### Scenes
- [x] Get all scenes
- [x] Get scene by ID
- [x] Create scene
- [ ] Modify scene
- [ ] Recall scene
- [ ] Delete scene

##### Groups
- [ ] Get all groups
- [ ] Create group
- [ ] Get group attributes
- [ ] Set group attributes
- [ ] Set group state
- [ ] Delete Group

##### Sensors
- [ ] Get all sensors
- [ ] Create sensor
- [ ] Find new sensors
- [ ] Get new sensors
- [ ] Get sensor
- [ ] Update sensor
- [ ] Delete sensor
- [ ] Change sensor configuration

##### Rules
- [ ] Get all rules
- [ ] Get rule
- [ ] Create rule
- [ ] Update rule
- [ ] Delete rule

## API Documentation
For official Philips Hue documentation check out the [Philips Hue website](http://www.developers.meethue.com/philips-hue-api)

## License
Copyright (C) 2016 Collin Guarino (Collinux)  
GPL version 2 or higher http://www.gnu.org/licenses/gpl.html  
GoHue project maintained by Collin Guarino, collin.guarino@gmail.com

## Contributing  
Pull requests happily accepted on GitHub
