package grovepi

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/mrmorphic/hwio"
)

// A Pin is the physical location on the board to which
// we read and write voltages.
type Pin byte

const (
	// Analog pins.
	A0 Pin = 0
	A1     = 1
	A2     = 2

	// Digital pins.
	D2 = 2
	D3 = 3
	D4 = 4
	D5 = 5
	D6 = 6
	D7 = 7
	D8 = 8
)

// PinMode specifies the direction of interaction.
// (Output or Input)
type PinMode int

const (
	Output PinMode = iota
	Input
)

// A Sensor is one of the Grove Seeed studio sensor, which
// can be connected to the GrovePi board.
type Sensor int

const (
	Light Sensor = iota
	Sound
)

type command byte

const (
	digitalRead  = 1
	digitalWrite = 2
	analogRead   = 3
	analogWrite  = 4
	pinMode      = 5
	dhtRead      = 40
)

// Interface specifies interaction with the GrovePi.
type Interface interface {
	// SetPinMode sets the pin mode for reading or writing.
	SetPinMode(pin Pin, mode PinMode) error

	// ReadAnalogy reads the analog value at the pin.
	ReadAnalog(pin Pin) (int, error)

	// Close disables the underlying I2C module handle.
	Close() error
}

// grovepi is the implementation of the Interface.
type grovePi struct {
	i2cmodule hwio.I2CModule
	i2cDevice hwio.I2CDevice
}

func InitGrovePi(address int) Interface {
	grovePi := new(grovePi)
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

func (grovePi *grovePi) Close() error {
	return grovePi.i2cmodule.Disable()
}

func (grovePi *grovePi) ReadAnalog(p Pin) (int, error) {
	pin := byte(p)
	b := []byte{analogRead, pin, 0, 0}
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

func (grovePi *grovePi) DigitalRead(pin byte) (byte, error) {
	b := []byte{digitalRead, pin, 0, 0}
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

func (grovePi *grovePi) DigitalWrite(pin byte, val byte) error {
	b := []byte{digitalWrite, pin, val, 0}
	err := grovePi.i2cDevice.Write(1, b)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return err
	}
	return nil
}

func (grovePi *grovePi) SetPinMode(p Pin, mode PinMode) error {
	pin := byte(p)
	var b []byte
	switch mode {
	case Output:
		b = []byte{pinMode, pin, 1, 0}
	case Input:
		b = []byte{pinMode, pin, 0, 0}
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

func (grovePi *grovePi) ReadDHT(pin byte) (float32, float32, error) {
	b := []byte{dhtRead, pin, 1, 0}
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

func (grovePi *grovePi) readDHTRawData(cmd []byte) ([]byte, error) {

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
