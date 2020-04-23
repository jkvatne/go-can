package main

import (
	"go-can/node"
	"time"
)

func VerifyAin(n *node.Node) {
	if n.SkipTest("Testing analog inputs") {
		return
	}
	n.SetPreOperational()
	_ = n.WriteObject(0x4010, 1, 1, SENS_DIG_OUT)
	_ = n.WriteObject(0x4010, 2, 1, SENS_DIG_OUT)
	_ = n.WriteObject(0x4010, 7, 1, SENS_VOLT)
	_ = n.WriteObject(0x4010, 8, 1, SENS_VOLT)
	n.SetOperational()
	time.Sleep(10*time.Millisecond)
	// Turn on output 1 and off output 2
	n.SetPdoValue(1, 0, 1, 3)
	SendPdos(n,4, 100*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 1, Vsupply-1.0, Vsupply+0.5, "DO1 high" )
	n.VerifyRangeFloat(0x4021, 2,  Vsupply-1.0, Vsupply+0.5, "DO2 high" )
	n.VerifyRangeFloat(0x4021, 7, Vsupply-1.0, Vsupply+0.5, "DO7 high")
	n.VerifyRangeFloat(0x4021, 8, Vsupply-1.0,  Vsupply+0.5, "DO8 high")
	// Set as current in
	_ = n.WriteObject(0x4010, 7, 1, SENS_MA)
	_ = n.WriteObject(0x4010, 8, 1, SENS_MA)
	SendPdos(n,4, 100*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 7, 0.015, 0.019, "Chan 7 current high" )
	n.VerifyRangeFloat(0x4021, 8, 0.015,  0.019, "Chan 8 current high" )

}


