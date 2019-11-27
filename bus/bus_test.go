package bus_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-can/bus"
	"go-can/can"
	"go-can/peak"
	"testing"
	"time"
)

var expectedMsg = can.Msg{Id: 0x58B, Type: 0, Len:8, Data:[8]uint8{0x43, 0x00, 0x10, 0,0xE4, 0x0C,0}}
var response can.Msg

func handler(m *can.Msg) {
	if m.Id==0 {
		fmt.Println("Error")
	}
	fmt.Printf("Callback with message: %s\n",m.ToString())
	response = *m
}

func TestConnectin(t *testing.T) {
	adapter, err := peak.New(peak.PCAN_USBBUS1, 125000)
	assert.NoError(t, err, "Peak adapter not found")
	if err!=nil {
		return
	}
	fmt.Printf("Setting up adapter and State\n")
	bus := bus.New(adapter, 100*time.Millisecond)
	assert.NotNil(t, bus, "Bus error")
	if err!=nil {
		return
	}
	defer bus.Close()

	fmt.Println("Sending a message (sdo read)")
	msg := can.Msg{Id:0x60B, Type:can.Standard, Len:8, Data:[8]uint8{0x40,0,0x10,0,0,0,0,0} }
	fmt.Println("Polling on State")
	m:=bus.Poll(msg, 0x58B)
	assert.NotNil(t, m, "No response")
	if m!=nil {
		fmt.Printf("Poll response is : %s\n", m.ToString())
		assert.Equal(t, expectedMsg, *m,  "Error")
	}
}
