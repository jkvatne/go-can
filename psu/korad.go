// korad is an interface to the KD3005 series power supplies.
// An example is the Elfa RND320 supply
// It is speical in that it does not use CR/LF as command endings, but depends on timeouts.

package psu

import (
	"fmt"
	"go-can/serial"
	"math"
	"strconv"
	"strings"
	"time"
)

// Check if KoradUsbPsu interface satisfies Psu interface
var _ Psu = &KoradUsbPsu{}

// KoradUsbPsu stores setup for a TTI CPX4000 power supply
type KoradUsbPsu struct {
	timeout time.Duration
	port *serial.Port
	voltage float64
	current float64
}

// NewTtiPsu returns a PSU instance for the tti supply
func NewKoradPsu(port string) (*KoradUsbPsu, error) {
	var err error
	if port == "" {
		list, err := serial.GetPortsList()
		if err!=nil || len(list)==0 {
			return nil, err
		}
		port = list[len(list)-1]
	}
	psu := &KoradUsbPsu{}
	c := &serial.Config{Name: port, Baud: 9600, ReadTimeout: time.Millisecond * 100, IntervalTimeout: time.Millisecond * 30}
	psu.port, err = serial.OpenPort(c)
	if err!=nil {
		return nil, fmt.Errorf("Error opening port, %s", err)
	}
	return psu, nil
}

func (psu *KoradUsbPsu) VerifyType() bool {
	name := psu.Name()
	return strings.Contains(name, "KD3005P")
}

func (psu *KoradUsbPsu) PortCount() int {
	return 2
}

func (psu *KoradUsbPsu) Name() string {
	name,err := psu.Ask("*IDN?")
	if err!=nil {
		return fmt.Sprintf("Error, %s",err)
	}
	return name
}

// Write will send a commend to the supply
func (psu *KoradUsbPsu) Write(data string) error {
	b := []byte(data)
	n, err := psu.port.Write(b)
	time.Sleep(50*time.Millisecond)
	if n!= len(b) {
		return fmt.Errorf("did not send all characters")
	}
	return err
}

// Ask will query the instrument for a string response
func  (psu *KoradUsbPsu) Ask(query string) (string, error) {
	var err error
	if psu==nil {
		return "",fmt.Errorf("No power supply defined")
	}
	buf := make([]byte, 64)
	err = psu.port.Flush()
	time.Sleep(time.Millisecond*10)
	err = psu.Write(query)
	if err!=nil {
		return "", err
	}
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
func (psu *KoradUsbPsu) SetOutput(channel int, voltage float64, current float64) error {
	// Korad has no enable command, so we save the setpoints and set outputs to zero if not enabled
	psu.voltage = voltage
	psu.current = current
	// The output voltage rate of change is ca 10V/sec
	var wait time.Duration
	if voltage>psu.voltage {
		wait = 100*time.Millisecond
	} else {
		wait = 50*time.Millisecond + time.Duration(math.Round(math.Abs(voltage-psu.voltage)*30))*time.Millisecond
	}
	// Set output voltage
	err := psu.Write(fmt.Sprintf("VSET%d:%0.2f", channel, voltage))
	if err!=nil {
		return err
	}
	// Set current limit
	err = psu.Write(fmt.Sprintf("ISET%d:%0.3f", channel, current))
	if err!=nil {
		return err
	}
	time.Sleep(wait)
	return nil
}

// DisableOutput will turn off the given output channel
func (psu *KoradUsbPsu) Disable(channel int) error {
	return psu.SetOutput(channel, 0, 0)
}

// GetOutput will return the actual output voltage and current from the channel
func (psu *KoradUsbPsu) GetOutput(channel int) (float64, float64, error) {
	// Read back output voltage
	voltageString, err1 := psu.Ask(fmt.Sprintf("VOUT%d?", channel))
	voltageString = strings.TrimRight(voltageString, "V\n")
	// Read back output current
	currentString, err2 := psu.Ask(fmt.Sprintf("IOUT%d?", channel))
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
func (psu *KoradUsbPsu) GetSetpoint(channel int) (float64, float64, error) {
	// Read back output voltage setpoint
	voltageString, err1 := psu.Ask(fmt.Sprintf("VSET%d?", channel))
	voltageString = strings.TrimPrefix(voltageString,fmt.Sprintf("V%d ",channel))
	voltageString = strings.TrimRight(voltageString, "V\n")
	// Read back output current setpoint
	currentString, err2 := psu.Ask(fmt.Sprintf("ISET%d?", channel))
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
func (psu *KoradUsbPsu) Shutdown() {
	_ = psu.port.Close()
}

