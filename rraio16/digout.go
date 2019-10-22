package main

import (
	"go-can/node"
	"time"
)

func VerifyDigOut(n *node.Node) {
	if n.SkipTest("Testing digital outputs, reading voltage back") {
		return
	}
	n.SetOperational()
	time.Sleep(10*time.Millisecond)
	// Turn on both outputs
	n.SetPdoValue(1, 0, 1, 1)
	n.SendPdo(1)
	time.Sleep(100*time.Millisecond)
	n.SendPdo(1)
	time.Sleep(100*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 1, 23.0, 25.0, "First digital output 1 voltage readback should be high" )
	n.VerifyRangeFloat(0x4021, 2, 0.0, 0.5, "First digital output 2 voltage readbck" )
	time.Sleep(100*time.Millisecond)
	n.SetPdoValue(1, 0, 1, 2)
	n.SendPdo(1)
	time.Sleep(100*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 1, 0.0, 0.5, "Second digital output 1 voltage readbck" )
	n.VerifyRangeFloat(0x4021, 2, 23.0, 25.0, "Second digital output 2 voltage readbck  should be high" )
}

