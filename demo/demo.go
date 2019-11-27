package main

import (
	"fmt"
	"go-can/bus"
	"go-can/can"
	"go-can/node"
	"go-can/peak"
	"time"
)

func handler(m *can.Msg) {
	fmt.Printf("Callback with message: %s\n",m.ToString())
}

func main() {
	fmt.Printf("Setting up adapter and bus\n")
	adapter, err := peak.New(peak.PCAN_USBBUS1, 125000)
	if err!=nil {
		fmt.Println("Peak adapter not found")
		return
	}
	bus := bus.New(adapter, 100*time.Millisecond)
	fmt.Println("Sending a message (sdo read)")
	msg := can.Msg{Id:1547, Type:can.Standard, Len:8, Data:[8]uint8{64,3,16,0,0,0,0,0} }
	bus.Write(msg)
	time.Sleep(100*time.Millisecond)

	fmt.Println("Polling on bus")
	m:=bus.Poll(msg, 0x58B)
	if m==nil {
		fmt.Printf("No response from peak canbus poll ")
	} else {
		fmt.Printf("Poll response is : %s\n", m.ToString())
	}

	myNode := node.New(bus, 11)
	fmt.Println("Testing SdoRead()")
	devId, _ := myNode.ReadObject( 0x1000, 0, 4)
	fmt.Printf("SdoRead from index 0x1000 from node 11, result is %d\n", devId)

	bus.Close()
	fmt.Printf("Done\n")


}
