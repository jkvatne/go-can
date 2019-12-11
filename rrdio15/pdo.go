package main

import (
	"fmt"
	"go-can/node"
	"time"
)

func VerifyRxPdoParameters(node *node.Node) {
	if node.SkipTest("Verify pdo parameters at 0x1400") {
		return
	}
	node.SetPreOperational()
	node.VerifyEqual(0x1400, 0, 1, 2, "RxPdo")
	node.VerifyEqual(0x1400, 1, 4, 0x4000020B, "Cob ID")
	node.VerifyEqual(0x1400, 2, 1, 254, "Transmission type")
	//node.VerifyEqual(0x1400, 3, 2, 0x11, "Inhibit time")
	node.VerifyReadAbort(0x1400, 4, 1, "Not implemented")
}

func VerifyRxPdoMapping(node *node.Node) {
	if node.SkipTest("Verify pdo mapping at 0x1600") {
		return
	}
	node.SetPreOperational()
	node.VerifyEqual(0x1600, 0, 1, 8, "Rx pdo mapping count")
	node.VerifyEqual(0x1600, 1, 4, 0x22000108, "Rx pdo mapping 1")
	node.VerifyEqual(0x1600, 2, 4, 0x22000208, "Rx pdo mapping 2")
	node.VerifyEqual(0x1600, 3, 4, 0x22010108, "Rx pdo mapping 3")
	node.VerifyEqual(0x1600, 4, 4, 0x22010208, "Rx pdo mapping 4")
	node.VerifyEqual(0x1600, 5, 4, 0x22000308, "Rx pdo mapping 5")
	node.VerifyEqual(0x1600, 6, 4, 0x22000408, "Rx pdo mapping 6")
	node.VerifyEqual(0x1600, 7, 4, 0x22010308, "Rx pdo mapping 7")
	node.VerifyEqual(0x1600, 8, 4, 0x22010408, "Rx pdo mapping 8")
}

func VerifyTxPdoParameters(node *node.Node) {
	if node.SkipTest("Verify tx pdo parameters at 0x1800") {
		return
	}
	node.SetPreOperational()
	node.VerifyEqual(0x1800, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1800, 1, 4, 0x4000018b, "Cob ID")
	node.VerifyEqual(0x1800, 2, 1, 4, "Transmission type")
	node.VerifyEqual(0x1800, 3, 2, 0, "Inhibit time") // Error in port code - should be 2 byte
	node.VerifyReadAbort(0x1800, 4, 1, "Not implemented")
}

func VerifyTxPdoMapping(node *node.Node) {
	if node.SkipTest("Verify tx pdo mapping at 0x1A00") {
		return
	}
	node.SetPreOperational()
	node.VerifyEqual(0x1A00, 0, 1, 6, "Tx pdo mapping count")
	node.VerifyEqual(0x1A00, 1, 4, 0x21000108, "Tx pdo mapping 1")
	node.VerifyEqual(0x1A00, 2, 4, 0x21000208, "Tx pdo mapping 2")
	node.VerifyEqual(0x1A00, 3, 4, 0x40290108, "Tx pdo mapping 3")
	node.VerifyEqual(0x1A00, 4, 4, 0x40290208, "Tx pdo mapping 4")
	node.VerifyEqual(0x1A00, 5, 4, 0x40280110, "Tx pdo mapping 5")
	node.VerifyEqual(0x1A00, 6, 4, 0x40280210, "Tx pdo mapping 6")
}


func VerifyRxPdo(n *node.Node) {
	if n.SkipTest("Testing rx pdo operation - transfer output values to card") {
		return
	}
	n.SetPreOperational()
	// Set input 1-8 to digital in
	for i:=1; i<9; i++ {
		_ = n.WriteObject(0x4010, node.SubIndex(i), 1, IO_PNP)
	}
	// Set channel 9-15 to digital out
	for i:=9; i<16; i++ {
		_ = n.WriteObject(0x4010, node.SubIndex(i), 1, IO_OUTPUT)
	}
	// Assume channel 1-8 is connected to channel 9-16
	n.SetOperational()
	time.Sleep(10*time.Millisecond)
	// Turn on outputs 9-15 via pdo
	n.SetPdoValue(1, 0, 2, 0xEF00)
	for i:=0; i<4; i++ {
		n.Bus.SendSync()
		n.SendPdo(1)
		time.Sleep(100 * time.Millisecond)
	}
	n.ResetPdoCount()
	// Send 12 sync pulses
	for i:=0; i<12; i++ {
		n.Bus.SendSync()
		n.SendPdo(1)
		time.Sleep(100*time.Millisecond)
	}
	n.VerifyPdoCount(3, 0, 0, 0)
}

func SetIsolationMode(n *node.Node, timeoutMs int) {
	n.SetPreOperational()
	n.Check(n.WriteObject(ISOLATION_PDO_NUM, 0, 1, 4))
	if timeoutMs==0 {
		n.Check(n.WriteObject(ISOLATION_PDO_RATE, 0, 1, 0))
	} else {
		n.Check(n.WriteObject(ISOLATION_PDO_RATE, 0, 1, 100))
	}
	n.Check(n.WriteObject(ISOLATION_PDO_COUNT, 0, 1, timeoutMs/100))

}

func VerifyEmcyOkStartingPdos(n *node.Node) {
	if n.SkipTest("Verify emergency message when starting pdos") {
		return
	}
	n.SetOperational()
	n.EmcyCount = 0
	SendPdos(n,2, time.Millisecond*100)
	// V2.6 sends one emergency message. This is actualy not correct!
	if n.EmcyCount==1 {
		fmt.Printf("One emergency message found at start of pdo trafic. This is standard v2.6 behavior.\n")
	}
}

func VerifyIsolationModeTime(n *node.Node) {
	if n.SkipTest("Verify isolation mode timeout delay") {
		return
	}
	// Set input 1-8 to digital in
	for i:=1; i<=8; i++ {
		_ = n.WriteObject(SENSOR_TYPE, node.SubIndex(i), 1, IO_PNP)
	}
	// Set channel 9-15 to digital out
	for i:=9; i<16; i++ {
		_ = n.WriteObject(SENSOR_TYPE, node.SubIndex(i), 1, IO_OUTPUT)
	}
	n.SetPreOperational()
	SetIsolationMode(n, 400)
	n.SetOperational()
	n.SetPdoValue(1, 0, 2, 0x5555)
	SendPdos(n,3, time.Millisecond*100)
	n.VerifyRange(0x4020, 1, 2, 1, 1, "Channel 1 digital output")
	n.VerifyRange(0x4020, 2, 2, 0, 0, "Channel 0 digital output")
	time.Sleep(2800 * time.Millisecond)
	n.VerifyRange(0x4020, 1, 2, 0, 0, "Channel 1 digital output")
	n.VerifyRange(0x4020, 2, 2, 0, 0, "Channel 0 digital output")
	SetIsolationMode(n, 2000)

}


