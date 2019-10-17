package main

import (
	"fmt"
	"go-can/can"
	"go-can/peak"
	"time"
)

func handler(m *can.Msg) {
	fmt.Println(m.ToString())
}

func main() {
	fmt.Printf("Start\n")
	canbus,_ := peak.New(peak.PCAN_USBBUS1, peak.PCAN_BAUD_125K, handler)
	msg := can.Msg{Id:1547, Type:peak.PCAN_MESSAGE_STANDARD, Len:8, Data:[8]uint8{64,3,16,0,0,0,0,0} }
	canbus.Write(msg)
	time.Sleep(100*time.Millisecond)
	m:=canbus.Poll(msg, 0x58B)
	if m!=nil {
		fmt.Printf("Poll response is : %s\n", m.ToString())
	}
	for i := 0; i<100; i++ {
		m := canbus.Poll(msg, 0x58B)
		fmt.Printf("Poll response is : %s\n", m.ToString())
		time.Sleep(time.Second/2)
	}
	canbus.Uninitialize()
	fmt.Printf("Done\n")
}
