package peak_test

import (
	"github.com/stretchr/testify/assert"
	"go-can/can"
	"go-can/peak"
	"testing"
	"time"
)

var expectedMsg = can.Msg{Id: 0x58B, Type: 0, Len:8, Data:[8]uint8{0x43,0x00,0x10,0,0xE4,0x0C,0}}
var response *can.Msg

// TestPeak assumes a CAN-Open node with id=11 is conneted. It polls a SDO at 1547 = 0x060B
func TestPeak(t *testing.T) {
	c, err := peak.New(peak.PCAN_USBBUS1,125000)
	assert.NoError(t, err, "Peak driver")
	msg := can.Msg{Id:1547, Type:can.Standard, Len:8, Data:[8]uint8{0x40,0,16,0,0,0,0,0} }
	c.Write(msg)
	time.Sleep(100*time.Millisecond)
	response:=c.Read()
	assert.NotNil(t, response, "No response from external card with id 11")
	if response!=nil {
		assert.Equal(t, expectedMsg.Data[0], response.Data[0], "Error")
		assert.Equal(t, expectedMsg.Data[1], response.Data[1], "Error")
	}
	c.Close()
}
