package psu

import (
	"bufio"
	"fmt"
	"os"
)

// Check if tti interface satisfies Psu interface
var _ Psu = &ManualPsu{}

// TtiUsbPsu stores setup for a TTI CPX4000 power supply
type ManualPsu struct {
	voltage float64
	current float64
}

// NewTtiPsu returns a PSU instance for the tti supply
func NewManualPsu(port string) (*ManualPsu, error) {
	psu := &ManualPsu{}
	psu.voltage = 0
	return psu, nil
}

func (psu *ManualPsu) Name() string {
	return "Manual power supply control"
}

func (psu *ManualPsu) PortCount() int {
	return 100
}

// SetOutput will set output voltage and current limit for a given channel
func (p *ManualPsu) SetOutput(channel int, voltage float64, current float64) error {
	fmt.Printf("Set output voltage cahnnel %d to %0.3fV, and current limit to %0.3fA\n", channel, voltage, current)
	fmt.Printf("Press <enter> key to continue...")
	p.voltage = voltage
	p.current = current
	_,_ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

func (p *ManualPsu) GetSetpoint(channel int) (float64, float64, error) {
	return 0.0,0.0,nil
}

// EnableOutput will turn on the given output channel.
// Voltage and current limit should be set first
func (p *ManualPsu) EnableOutput(channel int) error {
	fmt.Printf("Turn on output channel %d\n", channel)
	fmt.Printf("Press <enter> key to continue...")
	_,_ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

// DisableOutput will turn off the given output channel
func (p *ManualPsu) Disable(channel int) error {
	fmt.Printf("Turn off output channel %d\n", channel)
	fmt.Printf("Press <enter> key to continue...")
	_,_ = bufio.NewReader(os.Stdin).ReadBytes('\n')
	return nil
}

// GetOutput will return the actual output voltage and current from the channel
func (p *ManualPsu) GetOutput(channel int) (float64, float64, error) {
	return 0.0, 0.0, nil
}

// Shutdown will turn off all outputs and close the communication
func (p *ManualPsu) Shutdown() {
	fmt.Printf("Turn off power supply")
}

