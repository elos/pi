package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"./grovepi"
	"github.com/elos/gaia"
	"github.com/elos/models"
)

var prior *models.Event

func main() {
	gdb := &gaia.DB{
		URL:      "http://elos.pw",
		Username: "public",
		Password: "private",
		Client:   http.DefaultClient,
	}

	var g grovepi.GrovePi
	g = *grovepi.InitGrovePi(0x04)
	err := g.PinMode(grovepi.A0, "input")
	if err != nil {
		log.Fatal(err)
	}
	err = g.PinMode(grovepi.A1, "input")
	if err != nil {
		log.Fatal(err)
	}
	for {
		e := new(models.Event)
		e.OwnerId = "5698366ed6e4fe081b06570a"
		e.Name = "GrovePi Sensor Update"
		e.CreatedAt = time.Now()
		e.UpdatedAt = time.Now()
		e.Time = time.Now()
		light, err := g.AnalogRead(grovepi.A0)
		if err != nil {
			log.Fatal(err)
		}
		sound, err := g.AnalogRead(grovepi.A1)
		if err != nil {
			log.Fatal(err)
		}
		e.Data = map[string]interface{}{
			"light": light, "sound": sound,
		}
		fmt.Printf("Light: %d, Sound: %d\n", light, sound)
		if prior != nil {
			e.SetPrior(prior)
		}
		if err := gdb.Save(e); err != nil {
			log.Fatal(err)
		}
		prior = e
		time.Sleep(500 * time.Millisecond)
	}
	g.CloseDevice()
}
