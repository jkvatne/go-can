package main

import (
	"can"
	"fmt"
	"peak"
	"time"
)

func handler(m *can.Msg) {
	fmt.Printf("ID %x, Data %d,%d,%d,%d,%d,%d,%d,%d\n", m.Id, m.Data[0], m.Data[1], m.Data[2], m.Data[3],
		m.Data[4], m.Data[5], m.Data[6], m.Data[7])
}

func main() {
	fmt.Printf("Start\n")
	canbus,_ := peak.New(peak.PCAN_USBBUS1, peak.PCAN_BAUD_125K, handler)
	msg := can.Msg{Id:1547, Type:peak.PCAN_MESSAGE_STANDARD, Len:8, Data:[8]uint8{64,3,16,0,0,0,0,0} }
	_ = canbus.Write(msg)
	time.Sleep(100*time.Millisecond)
	canbus.Uninitialize()
}
