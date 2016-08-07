package sensor

import (
	"log"
	"time"

	"github.com/elos/pi/grovepi"
	"github.com/elos/x/models"
	"golang.org/x/net/context"
)

type Extractor func(g grovepi.Interface) ([]*models.Quantity, error)

func Merge(extractors ...Extractor) Extractor {
	return func(g grovepi.Interface) ([]*models.Quantity, error) {
		findings := make([]*models.Quantity, 0)
		for _, extractor := range extractors {
			f, err := extractor(g)
			if err != nil {
				return nil, err
			}
			findings = append(findings, f...)
		}

		return findings, nil
	}

}

type Recorder interface {
	Record(ctx context.Context, extractors ...Extractor)
	Close() error
}

func NewRecorder(g grovepi.Interface, i time.Duration, out chan<- []*models.Quantity) Recorder {
	return &recorder{
		g:        g,
		interval: i,
		out:      out,
	}
}

type recorder struct {
	g        grovepi.Interface
	interval time.Duration
	out      chan<- []*models.Quantity
	err      error
}

func (r *recorder) Record(ctx context.Context, extractors ...Extractor) {
read:
	for {
		select {
		case <-time.After(r.interval):
			for _, e := range extractors {
				f, err := e(r.g)
				if err != nil {
					log.Printf("TRANSIENT ERROR: %v", err)
					r.err = err
					continue read
				}
				r.out <- f
			}
		case <-ctx.Done():
			break read
		}
	}
	close(r.out)
}

func (r *recorder) Close() error {
	return r.err
}
