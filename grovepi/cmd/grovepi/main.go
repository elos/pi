package main

import (
	"crypto/x509"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/elos/pi/grovepi"
	"github.com/elos/pi/grovepi/config"
	"github.com/elos/pi/grovepi/sensor"
	"github.com/elos/x/auth"
	"github.com/elos/x/data"
	"github.com/elos/x/models"
	"golang.org/x/net/context"
)

var (
	configPath = flag.String("config", "/tmp/grovepi/config", "the configuration path for the sensors to load")
	dbAddr     = flag.String("dbaddr", "", "the address of the database server to connect to")
)

func main() {
	log.Print("--- Starting GrovePi ---")
	flag.Parse()
	if *configPath == "" {
		log.Fatal("*configPath = \"\", must specify config path")
	}
	if *dbAddr == "" {
		log.Fatal("*dbAddr == \"\", must specify database address")
	}

	p, err := config.Parse(*configPath)
	if err != nil {
		log.Fatal("config.Parase(*configPath) error: %v", err)
	}
	log.Print("Parsed Configuration:\n%v", p)

	g := grovepi.InitGrovePi(0x04)
	e := p.Extractor()

	out := make(chan []*models.Quantity)
	r := sensor.NewRecorder(g, 500*time.Millisecond, out)
	go r.Record(context.Background(), e)

	events := make(chan *models.Event)

	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		log.Fatal("failed to parse rootPEM")
	}

	conn, err := grpc.Dial(
		*dbAddr,
		/*
			grpc.WithTransportCredentials(
				credentials.NewClientTLSFromCert(pool, *dbAddr),
			),
		*/
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(
			auth.RawCredentials("public", "private"),
		),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	dbclient := data.NewDBClient(conn)

	go func() {
		var prior *models.Event
		for f := range out {
			log.Print(f)
			e := new(models.Event)
			e.OwnerId = "5698366ed6e4fe081b06570a"
			e.Name = "GrovePi Sensor Update"
			e.Time = models.TimestampFrom(time.Now())
			if prior != nil {
				e.PriorId = prior.Id
			}
			e.Quantities = append(e.Quantities, f...)
			prior = e
			events <- e
		}

		close(events)
	}()

	for e := range events {
		_, err := dbclient.Mutate(
			context.Background(),
			&data.Mutation{
				Op: data.Mutation_CREATE,
				Record: &data.Record{
					Kind:  models.Kind_EVENT,
					Event: e,
				},
			},
		)
		if err != nil {
			log.Printf("dbclient.Mutate error: %v", err)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatalf("recorder.Close() error: %v", err)
	}
}

const rootPEM = `
-----BEGIN CERTIFICATE-----
MIIE/zCCA+egAwIBAgISAyWO1mzQY7ckQXgEN7sdCmA5MA0GCSqGSIb3DQEBCwUA
MEoxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MSMwIQYDVQQD
ExpMZXQncyBFbmNyeXB0IEF1dGhvcml0eSBYMzAeFw0xNjA3MzEyMjM3MDBaFw0x
NjEwMjkyMjM3MDBaMBIxEDAOBgNVBAMTB2Vsb3MucHcwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQC9eUxrKjs3tseQcLjs3y0zP7lZjFmHtpWRhSdLkBtG
LfvWRle3eNOiCB23iBIi52uA41HYXjFnOxZJqkSSC0L28iyQ0AejuSUz7EMGEl3X
ZditiLMTLMx49+1YE2cUl79XSr42XtJu8KssSUXCo6kuiMo6YO6lVB4/FKrEIt0n
Cyclk5KAurpopytDaojckgiC+V22TV8/KBjGJdMbvBnrcnyhME7UJGNrWlccndR7
F81eLpbUY6mkJvo9zmmPUvD1IkxDlYnECzcDRkc/HLLXOpFIQrwjzejxyNv4+thR
aQIhKL53Um+3M7wZg+FlZnKvk/CbM/JPXe/bZPUb9xbVAgMBAAGjggIVMIICETAO
BgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwG
A1UdEwEB/wQCMAAwHQYDVR0OBBYEFAs23FhAEuVRpdhwpcdHFiscrQeNMB8GA1Ud
IwQYMBaAFKhKamMEfd265tE5t6ZFZe/zqOyhMHAGCCsGAQUFBwEBBGQwYjAvBggr
BgEFBQcwAYYjaHR0cDovL29jc3AuaW50LXgzLmxldHNlbmNyeXB0Lm9yZy8wLwYI
KwYBBQUHMAKGI2h0dHA6Ly9jZXJ0LmludC14My5sZXRzZW5jcnlwdC5vcmcvMB8G
A1UdEQQYMBaCB2Vsb3MucHeCC3d3dy5lbG9zLnB3MIH+BgNVHSAEgfYwgfMwCAYG
Z4EMAQIBMIHmBgsrBgEEAYLfEwEBATCB1jAmBggrBgEFBQcCARYaaHR0cDovL2Nw
cy5sZXRzZW5jcnlwdC5vcmcwgasGCCsGAQUFBwICMIGeDIGbVGhpcyBDZXJ0aWZp
Y2F0ZSBtYXkgb25seSBiZSByZWxpZWQgdXBvbiBieSBSZWx5aW5nIFBhcnRpZXMg
YW5kIG9ubHkgaW4gYWNjb3JkYW5jZSB3aXRoIHRoZSBDZXJ0aWZpY2F0ZSBQb2xp
Y3kgZm91bmQgYXQgaHR0cHM6Ly9sZXRzZW5jcnlwdC5vcmcvcmVwb3NpdG9yeS8w
DQYJKoZIhvcNAQELBQADggEBAFQKygv4TEnci3vMwoHHW9bTY6tEozqrd6X1aHsG
1kFOnivo56zPEcyl2KDvQZPrTGG1dY26s5vcexvM3xtqKokyTvHf7G4vFmtnhE9+
G9B8lDwyhA22u2XbGYJu0snQ+b3xUvz+6X6yTGgPMiW8YqPBNTVBoKA4liLnSRGK
vBK3K4osN3FPNQdHqRr/CDGN2k1DkYdhGl9tgnstyFQ5eirIUu1wwA0k8OuBWw4+
+JWch9QNfia0mLDDZLeVYY2RDJTq2wE9Dp2HvLxcCV0BPQVK/5VsIMXVkglOUftf
2PskvS0wWjNi2Nd6Ev2mM26UHdIZs9KJuy1/w9BGrBNEF1s=
-----END CERTIFICATE-----`
