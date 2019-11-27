// Assumes that a TTi CPX4000 power supply is connected to a USB port,
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

func TestPsu(t *testing.T) {
	list, err := serial.GetPortsList()
	assert.NoError(t, err, "error fetching com port list")
	assert.True(t, len(list)>0, "no ports found")

	// Use the last port in the list!
	name := list[len(list)-1]
	fmt.Printf("Using port %s\n", name)

	p, err := psu.NewTtiPsu(name)
	assert.NoError(t, err, "Failed to open COM9")
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

	// Start by disabling both outputs
	fmt.Printf("Turn off both outputs\n")
	err = p.DisableOutput(1)
	assert.NoError(t, err, "Disable output 1")
	err = p.DisableOutput(2)
	assert.NoError(t, err, "Disable output 2")
	time.Sleep(200*time.Millisecond)

	// Set both to 20.0V
	fmt.Printf("Set both outputs to 20.0V\n")
	err = p.SetOutput(1, 20.0, 0.2, false)
	assert.NoError(t, err, "set output 1")
	err = p.SetOutput(2, 20.0, 0.15, false)
	assert.NoError(t, err, "set output 1")
	volt, current, err := p.GetOutput(1)
	assert.NoError(t, err, "get output 1")
	assert.InDelta(t, 0.0, volt, 0.1, "voltage 1 setpoint")
	assert.InDelta(t, 0.0, current, 0.1, "current 1 setpoint")

	// Turn on outputs
	fmt.Printf("Turn on both outputs\n")
	err = p.EnableOutput(1)
	assert.NoError(t, err, "Enable output 1")
	err = p.EnableOutput(2)
	assert.NoError(t, err, "Enable output 2")
	time.Sleep(300*time.Millisecond)
	volt, current, err = p.GetOutput(1)
	assert.NoError(t, err, "get output 1")
	assert.InDelta(t, 20.0, volt, 0.1, "voltage 1 setpoint")
	assert.InDelta(t, 0.0, current, 0.1, "voltage 1 setpoint")
	volt, current, err = p.GetOutput(2)
	assert.NoError(t, err, "get output 1")
	assert.InDelta(t, 20.0, volt, 0.1, "voltage 2 setpoint")
	fmt.Printf("Shutdown\n")
	p.Shutdown()
}

