package main

import (
	"go-can/node"
	"time"
)

func SendPdos(n *node.Node, count int, delay time.Duration) {
	if delay<time.Millisecond {
		panic("Time delay in SendPdos() is too small!")
	}
	for i := 0; i < count; i++ {
		n.Bus.SendSync()
		for pdoNo:=1; pdoNo<=node.PdoCount; pdoNo++ {
			n.SendPdo(pdoNo)
		}
		time.Sleep(delay)
	}
}

func VerifyDigOut(n *node.Node) {
	if n.SkipTest("Readback digital outputs") {
		return
	}
	SetIsolationMode(n, 1000)
	n.SetOperational()
	// Turn on output 1
	n.SetPdoValue(1, 0, 1, 1)
	SendPdos(n,2, 100*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 1, Vsupply-1.0, Vsupply+0.5, "DO1 high readback" )
	n.VerifyRangeFloat(0x4021, 2, -0.5, 0.5, "DO2 low readback" )
	n.SetPdoValue(1, 0, 1, 2)
	SendPdos(n,2, 100*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 1, -0.5, 0.5, "DO1 low readback" )
	n.VerifyRangeFloat(0x4021, 2, Vsupply-1.0,  Vsupply+0.5, "DO2 high readback" )
	SetIsolationMode(n, 2000)
}

