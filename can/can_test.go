package can_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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


func TestPeak(t *testing.T) {
	fmt.Printf("Setting up adapter and Connection\n")
	connection := can.NewConnection(
		peak.New(peak.PCAN_USBBUS1, 125000),
		100*time.Millisecond,
		handler)

	fmt.Println("Sending a message (sdo read)")
	msg := can.Msg{Id:1547, Type:can.MESSAGE_STANDARD, Len:8, Data:[8]uint8{0x40,0,0x10,0,0,0,0,0} }
	connection.Dev.Write(msg)
	time.Sleep(100*time.Millisecond)
	assert.Equal(t, expectedMsg, response,  "Error")

	fmt.Println("Polling on Connection")
	m:=connection.Poll(msg, 0x58B)
	if m==nil {
		fmt.Printf("No response from peak canbus poll ")
	} else {
		fmt.Printf("Poll response is : %s\n", m.ToString())
	}
	assert.Equal(t, expectedMsg, *m,  "Error")

	fmt.Println("Testing SdoRead()")
	devId, _ := connection.SdoRead(11, 0x1000, 0, 4)
	fmt.Printf("SdoRead from index 0x1000 from node 11, result is %d\n", devId)
	time.Sleep(time.Second)
	connection.Close()
	fmt.Printf("Done\n")
}
