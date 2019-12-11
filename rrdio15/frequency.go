package main

import (
	"fmt"
	"go-can/node"
	"time"
)

func VerifyFrequency(n *node.Node) {
	if n.SkipTest("Frequency measurement") {
		return
	}
	n.SetPreOperational()
	_ = n.WriteObject(0x4010, 1, 1, IO_FREQ_HZ)
	_ = n.WriteObject(0x4010, 2, 1, IO_FREQ_HZx10)
	_ = n.WriteObject(0x4010, 3, 1, IO_FREQ_HZx100)
	_ = n.WriteObject(0x4010, 4, 1, IO_FREQ_HZ)
	n.SetOperational()
	time.Sleep(10*time.Millisecond)
	// Send 20 pulses at about 5Hz.
	out := 0;
	for i:=1; i<10; i++ {
		time.Sleep(100*time.Millisecond)
		out = out ^ 0x0F00
		n.SetPdoValue(1, 0, 2, out)
		n.SendPdo(1)
	}
	// This does not work on rev 2.6
	freq := n.ReadFloat(0x4021, 3)
	fmt.Printf("Frequency was %0.3f\n", freq)
	n.VerifyRangeFloat(0x4021, 1, 4.9, 5.1, "Frequency chan 1")
	n.VerifyRangeFloat(0x4021, 2, 4.9, 5.1, "Frequency chan 2")
	n.VerifyRangeFloat(0x4021, 3, 4.9, 5.1, "Frequency chan 3")
	n.VerifyRangeFloat(0x4021, 4, 4.9, 5.1, "Frequency chan 4")
}



