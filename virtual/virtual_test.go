package virtual_test

import (
	"github.com/stretchr/testify/assert"
	"go-can/can"
	"go-can/peak"
	"go-can/virtual"
	"testing"
	"time"
)

var expectedMsg = can.Msg{Id: 0x58B, Type: 0, Len:8, Data:[8]uint8{0x43,0x00,0x10,0,0xE4,0x0C,0}}
var response *can.Msg

// TestPeak assumes a CAN-Open node with id=11 is conneted. It polls a SDO at 1547 = 0x060B
func TestVirtual(t *testing.T) {
	c := virtual.New(peak.PCAN_USBBUS1,125000)
	msg := can.Msg{Id:1547, Type:can.MESSAGE_STANDARD, Len:8, Data:[8]uint8{0x40,0,16,0,0,0,0,0} }
	c.Write(msg)
	time.Sleep(100*time.Millisecond)
	response:=c.Read()
	assert.Equal(t, expectedMsg, *response,  "Error")
	c.Close()
}
