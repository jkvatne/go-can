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
	// Set input 1-8 to digital in
	for i:=1; i<=8; i++ {
		_ = n.WriteObject(SENSOR_TYPE, node.SubIndex(i), 1, IO_OUTPUT)
	}
	// Set channel 9-15 to digital out
	for i:=9; i<16; i++ {
		_ = n.WriteObject(SENSOR_TYPE, node.SubIndex(i), 1, IO_OUTPUT)
	}
	//SetIsolationMode(n, 0)
	n.SetOperational()
	// Turn on output 9,11,13,15
	n.SetPdoValue(1, 0, 2, 0xFFFF)
	SendPdos(n,5, 100*time.Millisecond)
	n.VerifyRangeFloat(INPUT_VALUE, 1, Vsupply-1.0, Vsupply+2.0, "DO1 high readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 2, 0.5, 4.5, "DO2 low readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 3, Vsupply-1.0, Vsupply+2.0, "DO3 high readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 4, 0.5, 4.5, "DO4 low readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 5, Vsupply-1.0, Vsupply+2.0, "DO5 high readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 6, 0.5, 4.5, "DO6 low readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 7, Vsupply-1.0, Vsupply+2.0, "DO7 high readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 8, Vsupply-1.0, Vsupply+2.0, "DO8 low readback" )
	n.SetPdoValue(1, 0, 2, 0xAA00)
	SendPdos(n,5, 100*time.Millisecond)
	n.VerifyRangeFloat(INPUT_VALUE, 1, 0.5, 4.5, "DO1 low readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 2, Vsupply-1.0,  Vsupply+2.0, "DO2 high readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 3, 0.5, 4.5, "DO3 low readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 4, Vsupply-1.0,  Vsupply+2.0, "DO4 high readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 5, 0.5, 4.5, "DO5 low readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 6, Vsupply-1.0,  Vsupply+2.0, "DO6 high readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 7, 0.5, 4.5, "DO7 low readback" )
	n.VerifyRangeFloat(INPUT_VALUE, 8, Vsupply-1.0,  Vsupply+2.0, "DO8 high readback" )
	SetIsolationMode(n, 2000)
}

