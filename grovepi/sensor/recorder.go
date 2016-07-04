package sensor

import (
	"log"
	"time"

	"github.com/elos/pi/grovepi"
	"golang.org/x/net/context"
)

type Findings map[string]interface{}

type Extractor func(g grovepi.Interface) (Findings, error)

func Merge(extractors ...Extractor) Extractor {
	return func(g grovepi.Interface) (Findings, error) {
		findings := make(Findings)
		for _, extractor := range extractors {
			f, err := extractor(g)
			if err != nil {
				return nil, err
			}
			for k, v := range f {
				findings[k] = v
			}
		}

		return findings, nil
	}

}

type Recorder interface {
	Record(ctx context.Context, extractors ...Extractor)
	Close() error
}

func NewRecorder(g grovepi.Interface, i time.Duration, out chan<- Findings) Recorder {
	return &recorder{
		g:        g,
		interval: i,
		out:      out,
	}
}

type recorder struct {
	g        grovepi.Interface
	interval time.Duration
	out      chan<- Findings
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
