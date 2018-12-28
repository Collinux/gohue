/*
* sensor.go
* GoHue library for Philips Hue
* Copyright (C) 2016 Collin Guarino (Collinux) collin.guarino@gmail.com
* License: GPL version 2 or higher http://www.gnu.org/licenses/gpl.html
 */
// https://developers.meethue.com/documentation/sensors-api

package hue

import (
	"time"
	"strings"
)

// special time type for unmarshal of lastupdated
type UpdateTime struct {
	*time.Time
}

// implement Unmarshal interface
// required for "none" as lastupdated in unused sensor
func (u *UpdateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "none" {
		*u = UpdateTime{&time.Time{}}
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	*u = UpdateTime{&t}
	return nil
}

// Sensor struct defines attributes of a sensor.
type Sensor struct {
	State struct {
		Daylight     bool       `json:"daylight"`    // True if day & false if night
		LastUpdated  UpdateTime `json:"lastupdated"` // Time of last update
		ButtonEvent  uint16     `json:"buttonevent"` // ID of button event
	} `json:"state"`

	Config struct {
		On         bool   `json:"on"`        // Turns the sensor on/off. When off, state changes of the sensor are not reflected in the sensor resource.
		Reachable  bool   `json:"reachable"` // Indicates whether communication with devices is possible
		Battery    uint8  `json:"battery"`   // The current battery state in percent, only for battery powered devices
	} `json:"config"`

	Type             string  `json:"type"`
	Name             string  `json:"name"`
	ModelID          string  `json:"modelid"`
	ManufacturerName string  `json:"manufacturername"`
	UniqueID         string  `json:"uniqueid"`
	SWVersion        string  `json:"swversion"`
	Index            int     // Set by index of sensor array response
	Bridge           *Bridge // Set by the bridge when the sensor is found
}

/// Refresh sensor attributes
func (s *Sensor) Refresh() error {
	sensor, err := s.Bridge.GetSensorByIndex(s.Index)
	if err != nil {
		return err
	}

	s.State = sensor.State
	s.Config = sensor.Config
	s.SWVersion = sensor.SWVersion
	s.Name = sensor.Name
	return nil
}
