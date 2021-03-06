package psu

import (
	"fmt"
	"go-can/serial"
	"strconv"
	"strings"
	"time"
)

// Check if tti interface satisfies Psu interface
var _ Psu = &TtiUsbPsu{}

// TtiUsbPsu stores setup for a TTI CPX4000 power supply
type TtiUsbPsu struct {
	timeout time.Duration
	port *serial.Port
}


// NewTtiPsu returns a PSU instance for the tti supply
func NewTtiPsu(port string) (*TtiUsbPsu, error) {
	var err error
	psu := &TtiUsbPsu{}
	port, err = CheckComPort(port)
	if err!=nil {
		return nil, err
	}
	c := &serial.Config{Name: port, Baud: 115200, ReadTimeout: time.Second/5, IntervalTimeout: time.Millisecond * 30}
	psu.port, err = serial.OpenPort(c)
	if err!=nil {
		return nil, fmt.Errorf("Error opening port, %s", err)
	}
	if psu.VerifyType() {
		return psu, nil
	}
	return nil, fmt.Errorf("Com port %s has not a TTi supply connected", port)
}

func (psu *TtiUsbPsu) PortCount() int {
	return 2
}

func (psu *TtiUsbPsu) VerifyType() bool {
	psu.Write("\n")
	name := psu.Name()
	return (name[0:15]=="THURLBY THANDAR")
}

func (psu *TtiUsbPsu) Name() string {
	name,err := psu.Ask("*IDN?")
	if err!=nil {
		return fmt.Sprintf("Error, %s",err)
	}
	return name
}

// Write will send a commend to the supply
func (psu *TtiUsbPsu) Write(data string) error {
	if data[len(data)-1] != '\n' {
		data = data + "\n"
	}
	b := []byte(data)
	n, err := psu.port.Write(b)
	if n!= len(b) {
		return fmt.Errorf("did not send all characters")
	}
	return err
}

// Ask will query the instrument for a string response
func  (psu *TtiUsbPsu) Ask(query string) (string, error) {
	var err error
	if psu==nil {
		return "",fmt.Errorf("No power supply defined")
	}
	buf := make([]byte, 64)
	err = psu.port.Flush()
	time.Sleep(time.Millisecond*100)
	err = psu.Write(query)
	if err!=nil {
		return "", err
	}
	time.Sleep(time.Millisecond*100)
	n, err := psu.port.Read(buf)
	if err!=nil {
		return "", err
	}
	if n==0 {
		return "", fmt.Errorf("no response")
	}
	response := string(buf)
	response = strings.TrimRight(response, "\n\r\000")
	return response, nil
}

// SetOutput will set output voltage and current limit for a given channel
func (psu *TtiUsbPsu) SetOutput(channel int, voltage float64, current float64) error {
	// Set output voltage
	err := psu.Write(fmt.Sprintf("V%d %0.3f", channel, voltage))
	if err!=nil {
		return err
	}
	// Set current limit
	err = psu.Write(fmt.Sprintf("I%d %0.2f", channel, current))
	if err!=nil {
		return err
	}
	return  psu.Write(fmt.Sprintf("OP%d 1", channel))
}

// DisableOutput will turn off the given output channel
func (psu *TtiUsbPsu) Disable(channel int) error {
	return psu.Write(fmt.Sprintf("OP%d 0", channel))
}

// GetOutput will return the actual output voltage and current from the channel
func (psu *TtiUsbPsu) GetOutput(channel int) (float64, float64, error) {
	// Read back output voltage
	voltageString, err1 := psu.Ask(fmt.Sprintf("V%dO?", channel))
	voltageString = strings.TrimRight(voltageString, "V\n")
	// Read back output current
	currentString, err2 := psu.Ask(fmt.Sprintf("I%dO?", channel))
	currentString = strings.TrimRight(currentString, "A\n")
	volt, err3 := strconv.ParseFloat(voltageString, 64)
	if err1!=nil || err2!=nil || err3!=nil {
		return 0,0,fmt.Errorf("error reding voltage, %s", err1)
	}
	curr, err := strconv.ParseFloat(currentString,64)
	if err!=nil {
		return volt, 0, fmt.Errorf("error reding current, %s", err)
	}
	return volt, curr, nil
}

// GetOutput will return the actual output voltage and current from the channel
func (psu *TtiUsbPsu) GetSetpoint(channel int) (float64, float64, error) {
	// Read back output voltage setpoint
	voltageString, err1 := psu.Ask(fmt.Sprintf("V%d?", channel))
	voltageString = strings.TrimPrefix(voltageString,fmt.Sprintf("V%d ",channel))
	voltageString = strings.TrimRight(voltageString, "V\n")
	// Read back output current setpoint
	currentString, err2 := psu.Ask(fmt.Sprintf("I%d?", channel))
	currentString = strings.TrimPrefix(currentString,fmt.Sprintf("I%d ",channel))
	currentString = strings.TrimRight(currentString, "A\n")

	volt, err3 := strconv.ParseFloat(voltageString, 64)
	if err1!=nil || err2!=nil || err3!=nil {
		return 0,0,fmt.Errorf("error reding voltage, %s", err1)
	}
	curr, err := strconv.ParseFloat(currentString,64)
	if err!=nil {
		return volt, 0, fmt.Errorf("error reding current, %s", err)
	}
	return volt, curr, nil
}

// Shutdown will turn off all outputs and close the communication
func (psu *TtiUsbPsu) Shutdown() {
	psu.SetOutput(1, 0, 0)
	psu.SetOutput(2, 0, 0)
	_ = psu.port.Close()
}

