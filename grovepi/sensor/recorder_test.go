package sensor_test

import (
	"testing"
	"time"

	"github.com/elos/pi/grovepi"
	"github.com/elos/pi/sensor"
	"golang.org/x/net/context"
)

type mockGrovePi struct {
	analog int
	err    error
}

func (gp *mockGrovePi) ReadAnalog(pin byte) (int, error)                { return gp.analog, gp.err }
func (gp *mockGrovePi) SetPinMode(pin byte, mode grovepi.PinMode) error { return gp.err }
func (gp *mockGrovePi) Close() error                                    { return gp.err }

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
	out := make(chan sensor.Findings, 6)
	r := sensor.NewRecorder(g, 50*time.Millisecond, out)

	r.Record(ctx, func(g grovepi.Interface) (sensor.Findings, error) {
		a, err := g.ReadAnalog(grovepi.A0)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"light": a,
		}, nil
	})

	for f := range out {
		if _, ok := f["light"]; !ok {
			t.Errorf("f[\"light\"]: got %t, want %t", ok, true)
		}
	}

	if err := r.Close(); err != nil {
		t.Fatalf("r.Close() error: %v", err)
	}
}
