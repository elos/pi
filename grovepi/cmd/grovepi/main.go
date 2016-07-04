package main

import (
	"flag"
	"log"
	"time"

	"github.com/elos/models"
	"github.com/elos/pi/grovepi"
	"github.com/elos/pi/grovepi/config"
	"github.com/elos/pi/grovepi/sensor"
)

var (
	configPath = flag.String("config", "", "the configuration path for the sensors to load")
)

func main() {
	if *configPath == "" {
		log.Fatal("*configPath = \"\", must specify config path")
	}

	p, err := config.Parse(*configPath)
	if err != nil {
		log.Fatal("config.Parase(*configPath) error: %v", err)
	}

	g := grovepi.InitGrovePi(0x04)
	Execute(p, g)

	out := make(chan sensor.Findings)
	r := sensor.NewRecorder(g, 500*time.Millisecond, out)

	events := make(chan *models.Event)

	go func() {
		var prior *models.Event
		for f := range out {
			e := new(models.Event)
			e.SetID(db.NewID())
			e.OwnerId = "5698366ed6e4fe081b06570a"
			e.Name = "GrovePi Sensor Update"
			e.CreatedAt = time.Now()
			e.UpdatedAt = time.Now()
			e.Time = time.Now()
			if prior != nil {
				e.SetPrior(prior)
			}
			e.Data = f
			prior = e
		}

		close(events)
	}()

	for e := range events {
		if err := db.Save(e); err != nil {
			log.Fatalf("db.Save error: %v", err)
		}
	}

	if err := recorder.Close(); err != nil {
		log.Fatalf("recorder.Close() error: %v", err)
	}
}

func Execute(p Plan, g grovepi.Interface) {
	extractors := make([]Extractor, 0, len(p))
	for sensor, pin := range p {
		g.SetPinMode(pin, grovepi.Input)
		extractors = append(extractors, config.ExtractorFactories[sensor](pin))
	}

	extractor := sensor.Merge(extractors)
}
