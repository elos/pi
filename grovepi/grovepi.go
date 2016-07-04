package grovepi

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/elos/pi/grovepi"
	"github.com/mrmorphic/hwio"
)

type Sensor int

const (
	Light Sensor = iota
	Sound
)

type Pin byte

const (
	//Pins
	A0 Pin = 0
	A1     = 1
	A2     = 2

	D2 = 2
	D3 = 3
	D4 = 4
	D5 = 5
	D6 = 6
	D7 = 7
	D8 = 8
)

const (
	//Cmd format
	DIGITAL_READ  = 1
	DIGITAL_WRITE = 2
	ANALOG_READ   = 3
	ANALOG_WRITE  = 4
	PIN_MODE      = 5
	DHT_READ      = 40
)

type PinMode int

const (
	Output PinMode = iota
	Input
)

type Interface interface {
	ReadAnalog(pin byte) (int, error)
	SetPinMode(pin byte, mode PinMode) error
	Close() error
}

type GrovePi struct {
	i2cmodule hwio.I2CModule
	i2cDevice hwio.I2CDevice
}

func InitGrovePi(address int) *GrovePi {
	grovePi := new(GrovePi)
	m, err := hwio.GetModule("i2c")
	if err != nil {
		fmt.Printf("could not get i2c module: %s\n", err)
		return nil
	}
	grovePi.i2cmodule = m.(hwio.I2CModule)
	grovePi.i2cmodule.Enable()

	grovePi.i2cDevice = grovePi.i2cmodule.GetDevice(address)
	return grovePi
}

func (grovePi *GrovePi) Close() error {
	return grovePi.i2cmodule.Disable()
}

func (grovePi *GrovePi) ReadAnalog(p grovepi.Pin) (int, error) {
	pin = byte(p)
	b := []byte{ANALOG_READ, pin, 0, 0}
	err := grovePi.i2cDevice.Write(1, b)
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	grovePi.i2cDevice.ReadByte(1)
	val, err2 := grovePi.i2cDevice.Read(1, 4)
	if err2 != nil {
		return 0, err
	}
	var v1 int = int(val[1])
	var v2 int = int(val[2])
	return ((v1 * 256) + v2), nil
}

func (grovePi *GrovePi) DigitalRead(pin byte) (byte, error) {
	b := []byte{DIGITAL_READ, pin, 0, 0}
	err := grovePi.i2cDevice.Write(1, b)
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	val, err2 := grovePi.i2cDevice.ReadByte(1)
	if err2 != nil {
		return 0, err2
	}
	return val, nil
}

func (grovePi *GrovePi) DigitalWrite(pin byte, val byte) error {
	b := []byte{DIGITAL_WRITE, pin, val, 0}
	err := grovePi.i2cDevice.Write(1, b)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return err
	}
	return nil
}

func (grovePi *GrovePi) SetPinMode(pin byte, mode PinMode) error {
	var b []byte
	switch mode {
	case Output:
		b = []byte{PIN_MODE, pin, 1, 0}
	case Input:
		b = []byte{PIN_MODE, pin, 0, 0}
	default:
		return fmt.Errorf("unknown pin mode: %d", mode)
	}
	err := grovePi.i2cDevice.Write(1, b)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return err
	}
	return nil
}

func (grovePi *GrovePi) ReadDHT(pin byte) (float32, float32, error) {
	b := []byte{DHT_READ, pin, 1, 0}
	rawdata, err := grovePi.readDHTRawData(b)
	if err != nil {
		return 0, 0, err
	}
	temperatureData := rawdata[1:5]

	tInt := int32(temperatureData[0]) | int32(temperatureData[1])<<8 | int32(temperatureData[2])<<16 | int32(temperatureData[3])<<24
	t := (*(*float32)(unsafe.Pointer(&tInt)))

	humidityData := rawdata[5:9]
	humInt := int32(humidityData[0]) | int32(humidityData[1])<<8 | int32(humidityData[2])<<16 | int32(humidityData[3])<<24
	h := (*(*float32)(unsafe.Pointer(&humInt)))
	return t, h, nil
}

func (grovePi *GrovePi) readDHTRawData(cmd []byte) ([]byte, error) {

	err := grovePi.i2cDevice.Write(1, cmd)
	if err != nil {
		return nil, err
	}
	time.Sleep(600 * time.Millisecond)
	grovePi.i2cDevice.ReadByte(1)
	time.Sleep(100 * time.Millisecond)
	raw, err := grovePi.i2cDevice.Read(1, 9)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
