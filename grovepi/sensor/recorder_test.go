package sensor_test

import (
	"testing"
	"time"

	"github.com/elos/pi/grovepi"
	"github.com/elos/pi/grovepi/sensor"
	"github.com/elos/x/models"
	"golang.org/x/net/context"
)

type mockGrovePi struct {
	analog int
	err    error
}

func (gp *mockGrovePi) ReadAnalog(pin grovepi.Pin) (int, error)                { return gp.analog, gp.err }
func (gp *mockGrovePi) SetPinMode(pin grovepi.Pin, mode grovepi.PinMode) error { return gp.err }
func (gp *mockGrovePi) Close() error                                           { return gp.err }

func TestRecorder(t *testing.T) {
	ctx, _ := context.WithTimeout(
		context.Background(),
		350*time.Millisecond,
	)

	// Never error, mock a return of 100 for AnalogReads.
	g := grovepi.Interface(&mockGrovePi{
		analog: 100,
		err:    nil,
	})

	// 6 reads
	out := make(chan []*models.Quantity, 6)
	r := sensor.NewRecorder(g, 50*time.Millisecond, out)

	r.Record(ctx, func(g grovepi.Interface) ([]*models.Quantity, error) {
		a, err := g.ReadAnalog(grovepi.A0)
		if err != nil {
			return nil, err
		}

		return []*models.Quantity{
			{
				Unit:      models.Quantity_GROVEPI_LIGHT,
				Magnitude: 0,
				Value:     float64(a),
			},
		}, nil
	})

	for f := range out {
		if got, want := len(f), 1; got != want {
			t.Fatalf("len(f): got %d, want %d", got, want)
		}

		if got, want := f[0].Unit, models.Quantity_GROVEPI_LIGHT; got != want {
			t.Errorf("f[0].Unit: got %s, want %s", got, want)
		}
	}

	if err := r.Close(); err != nil {
		t.Fatalf("r.Close() error: %v", err)
	}
}
