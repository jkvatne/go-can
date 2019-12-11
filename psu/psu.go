package psu

import (
	"fmt"
	"go-can/serial"
	"time"
)

type Psu interface {
	SetOutput(channel int, voltage float64, current float64) error
	GetOutput(channel int) (float64, float64, error)
	GetSetpoint(channel int) (float64, float64, error)
	Disable(channel int) error
	Name() string
	Shutdown()
	PortCount() int
}

var channels []struct {
	psu Psu
	psuChanel int
	name string
}

func CheckComPort(port string) (string, error) {
	list, _ := serial.GetPortsList()
	if port!="" {
		ok := false;
		// Check if port exists in list, if given
		for i := 0; i < len(list); i++ {
			if list[i] == port {
				ok = true
			}
		}
		if !ok {
			return "", fmt.Errorf("port %s does not exist", port)
		}
	}
	// Verify that given port is not used
	if port!="" {
		c := &serial.Config{Name: port, Baud: 115200, ReadTimeout: time.Second / 5, IntervalTimeout: time.Millisecond * 30}
		p, e1 := serial.OpenPort(c)
		if p!=nil {
			p.Close()
		}
		if e1!=nil || p==nil {
			return "", fmt.Errorf("port %s returns error (might be in use), %s", port, e1)
		}
		return port, nil
	}
	// Search for last unused port
	for j:=len(list)-1; j>=0; j-- {
		c := &serial.Config{Name: list[j], Baud: 115200, ReadTimeout: time.Second / 5, IntervalTimeout: time.Millisecond * 30}
		p, err := serial.OpenPort(c)
		if p!=nil {
			p.Close()
		}
		if err==nil && p!=nil {
			return list[j], nil
		}
	}
	return "", fmt.Errorf("no free ports found")
}

// NewPsu will return a psu on the given port
// if port is empty, use highest com-port found that is working
// Use channel given if the psu has more than one channel.

func NewPsu(port string) Psu {
	port, _ = CheckComPort(port)
	if port=="" {

	}
	var psu Psu
	var err error
	psu, err = NewTtiPsu(port)
	if psu!=nil && err==nil {
		return psu
	}
	psu, err = NewKoradPsu(port)
	if psu!=nil && err==nil {
		return psu
	}
	psu, err = NewManualPsu(port)
	return psu
}
