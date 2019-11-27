package node_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-can/bus"
	"go-can/node"
	"go-can/peak"
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	adapter, err := peak.New(peak.PCAN_USBBUS1, 125000)
	assert.NoError(t, err, "Peak adapter not found")
	if err!=nil {
		return
	}
	bus := bus.New(adapter, 100*time.Millisecond)
	n11 := node.New(bus, 11)
	deviceType, err := n11.ReadObject(0x1000, 0, 4)
	assert.NoError(t, err, "Node 11 not found")
	if err==nil {
		fmt.Printf("0x1000:0 = %d\n", deviceType)
	}
	assert.NotZero(t, deviceType, "Object 1000h should not be zero")
}

