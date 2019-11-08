package main

import (
	"go-can/node"
	"time"
)

func VerifyFrequency(n *node.Node) {
	if n.SkipTest("Frequency measurement") {
		return
	}
	n.SetPreOperational()
	_ = n.WriteObject(0x4010, 1, 1, SENS_DIG_OUT)
	_ = n.WriteObject(0x4010, 2, 1, SENS_DIG_OUT)
	_ = n.WriteObject(0x4010, 3, 1, SENS_FREQ)
	_ = n.WriteObject(0x4010, 4, 1, SENS_FREQ)
	n.SetOperational()
	time.Sleep(10*time.Millisecond)
	// Send 20 pulses at about 5Hz.
	out := 0;
	t:= time.Now()
	for i:=1; i<20; i++ {
		out = out ^ 3
		n.SetPdoValue(1, 0, 1, out)
		n.SendPdo(1)
		for time.Since(t)<time.Duration(i)*50*time.Millisecond {
			time.Sleep(time.Millisecond)
		}
	}
	// This does not work on rev 2.6
	//n.VerifyRangeFloat(0x4021, 7, 1, 1000, "Frequency chan 3")
	//n.VerifyRangeFloat(0x4021, 8, 1, 1000, "Frequency chan 8")
}



