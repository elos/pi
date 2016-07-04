package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/elos/pi/grovepi"
	"github.com/elos/pi/grovepi/sensor"
)

// i.e.,
// {
//		"light": "A0",
// }
type plan map[string]string

type Plan map[grovepi.Sensor]grovepi.Pin

func Parse(path string) (Plan, error) {
	p := make(plan)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &p); err != nil {
		return nil, err
	}

	plan := make(Plan)

	for s, p := range plan {
		sensor, ok := Sensors[s]
		if !ok {
			return nil, fmt.Errorf("unrecognized sensor: %q", s)
		}

		pin, ok := Pins[p]
		if !ok {
			return nil, fmt.Errorf("unrecognized pin: %q", p)
		}

		plan[sensor] = pin
	}

	return plan, nil
}

var Sensors = map[string]grovepi.Sensor{
	"light": grovepi.Light,
	"sound": grovepi.Sound,
}

var ExtractorFactories = map[grovepi.Sensor]func(grovepi.Pin) sensor.Extractor{
	Light: LightExtractor,
	Sound: SoundExtractor,
}

var Pins = map[string]grovepi.Pin{
	"A0": grovepi.A0,
	"A1": grovepi.A1,
	"A2": grovepi.A2,

	"D2": grovepi.D2,
	"D3": grovepi.D3,
	"D4": grovepi.D4,
	"D5": grovepi.D5,
	"D6": grovepi.D6,
	"D7": grovepi.D7,
}

func LightExtractor(p grovepi.Pin) sensor.Extractor {
	return func(g grovepi.Interface) (sensor.Findings, error) {
		light, err := g.ReadAnalog(p)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"light": light,
		}, nil
	}
}

func SoundExtractor(p grovepi.Pin) sensor.Extractor {
	return func(g grovepi.Interface) (sensor.Findings, error) {
		sound, err := g.ReadAnalogy(p)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"sound": sound,
		}, nil
	}
}
