// Assumes that a RND KD3005 power supply is connected to a USB port,
// and that this port has the highest com port number

package psu_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-can/psu"
	"go-can/serial"
	"os"
	"testing"
	"time"
)

func TestKorad(t *testing.T) {
	list, err := serial.GetPortsList()
	assert.NoError(t, err, "error fetching com port list")
	assert.True(t, len(list)>0, "no ports found")

	// Use the last port in the list!
	name := list[len(list)-1]
	fmt.Printf("Using port %s\n", name)

	p, err := psu.NewKoradPsu(name)
	assert.NoError(t, err, "Failed to open COM port")
	if err!=nil {
		fmt.Printf("Error opening port %s, %s\n", "COM9", err)
		os.Exit(1)
	}
	id, err := p.Ask("*IDN?")
	assert.NoError(t, err, "error fetching *IDN?")
	fmt.Printf("Found power supply \"%s\"\n",id)
	if err!=nil {
		return
	}

	// Start by disabling output
	fmt.Printf("Turn off output\n")
	err = p.Disable(1)
	assert.NoError(t, err, "Disable output 1")
	time.Sleep(800*time.Millisecond)
	volt, current, err := p.GetOutput(1)
	assert.NoError(t, err, "get output 1")
	assert.InDelta(t, 0.0, volt, 0.1, "voltage 1 setpoint")
	assert.InDelta(t, 0.0, current, 0.1, "voltage 1 setpoint")

	// Set output to 20.0V
	fmt.Printf("Set output to 24.0V\n")
	err = p.SetOutput(1, 24.0, 0.2)
	assert.NoError(t, err, "set output 1")

	time.Sleep(400*time.Millisecond)
	volt, current, err = p.GetOutput(1)
	assert.NoError(t, err, "get output 1")
	assert.InDelta(t, 24.0, volt, 0.1, "voltage 1 setpoint")
	assert.InDelta(t, 0.0, current, 0.1, "voltage 1 setpoint")

	volt, current, err = p.GetSetpoint(1)
	assert.InDelta(t, 24.0, volt, 0.1, "voltage 1 setpoint")
	assert.InDelta(t, 0.2, current, 0.1, "voltage 1 setpoint")

	// Turn off output
	err = p.Disable(1)
	time.Sleep(400*time.Millisecond)
	assert.NoError(t, err, "set output 1")
	volt, current, err = p.GetOutput(1)
	assert.NoError(t, err, "get output 1")
	assert.InDelta(t, 0.0, volt, 1.0, "voltage 1 setpoint")
	assert.InDelta(t, 0.0, current, 0.1, "current 1 setpoint")

	fmt.Printf("Shutdown\n")
	p.Shutdown()
}


