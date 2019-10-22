package main

import (
	"fmt"
	"go-can/can"
	"go-can/virtual"
	"time"
)

func handler(m *can.Msg) {
	fmt.Printf("Callback with message: %s\n",m.ToString())
}

func main() {
	fmt.Printf("Setting up adapter and bus\n")
	connection := connection.New(
		//peak.New(peak.PCAN_USBBUS1, 125000),
		virtual.New(0,125000),
		100*time.Millisecond,
		handler)

	fmt.Println("Sending a message (sdo read)")
	msg := can.Msg{Id:1547, Type:can.Standard, Len:8, Data:[8]uint8{64,3,16,0,0,0,0,0} }
	connection.Dev.Write(msg)
	time.Sleep(100*time.Millisecond)

	fmt.Println("Polling on bus")
	m:=connection.Poll(msg, 0x58B)
	if m==nil {
		fmt.Printf("No response from peak canbus poll ")
	} else {
		fmt.Printf("Poll response is : %s\n", m.ToString())
	}
	fmt.Println("Testing SdoRead()")
	devId, _ := connection.SdoRead(11, 0x1000, 0, 4)
	fmt.Printf("SdoRead from index 0x1000 from node 11, result is %d\n", devId)




	connection.Close()
	fmt.Printf("Done\n")


}
