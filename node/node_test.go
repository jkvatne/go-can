package node_test

import (
	"fmt"
	"go-can/bus"
	"go-can/node"
	"go-can/peak"
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	adapter := peak.New(peak.PCAN_USBBUS1, 125000)
	bus := bus.New(adapter, 100*time.Millisecond, node.MsgHandler)
	n11 := node.New(bus, 11)
	o1, err := n11.ReadObject(0x1000, 0, 4)
	if err!=nil {
		fmt.Printf("Error reading 0x1000:0, %s", err)
	}
	fmt.Printf("0x1000:0 = %d\n", o1)
}

