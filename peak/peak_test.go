package peak_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-can/can"
	"go-can/peak"
	"testing"
	"time"
)

var expectedMsg = can.Msg{Id: 0x58B, Type: 0, Len:8, Data:[8]uint8{67,5,16,0,128,0,0,128}}
var response *can.Msg

func handler(m *can.Msg) {
	fmt.Println(m.ToString())
	response = m
}

// TestPeak assumes a CAN-Open node with id=11 is conneted. It polls a SDO at 1547 = 0x1003
func TestPeak(t *testing.T) {
	c,_ := peak.New(peak.PCAN_USBBUS1, peak.PCAN_BAUD_125K, handler)
	msg := can.Msg{Id:1547, Type:peak.PCAN_MESSAGE_STANDARD, Len:8, Data:[8]uint8{64,5,16,0,0,0,0,0} }
	c.Write(msg)
	time.Sleep(100*time.Millisecond)
	assert.Equal(t, response, &expectedMsg, "Error")
	m:=c.Poll(msg, 0x58B)
	if m!=nil {
		fmt.Printf("Poll response is : %s\n", m.ToString())
	}
	for i := 0; i<2; i++ {
		m := c.Poll(msg, 0x58B)
		assert.Equal(t, response, m, "Error")
		fmt.Printf("Poll response is : %s\n", m.ToString())
		time.Sleep(time.Second/2)
	}
	c.Uninitialize()
}
