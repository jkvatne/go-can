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
	node.VerifyRange(0x1400, 0, 1, 2,3, "RxPdo")
	node.VerifyEqual(0x1400, 1, 4, 0x4000020B, "Cob ID")
	node.VerifyEqual(0x1400, 2, 1, 254, "Transmission type")
	//node.VerifyEqual(0x1400, 3, 2, 0, "Inhibit time")
	node.VerifyReadAbort(0x1400, 4, 1, "Not implemented")
	node.VerifyRange(0x1401, 0, 1, 2, 3, "RxPdo")
	node.VerifyEqual(0x1401, 1, 4, 0x4000030B, "Cob ID")
	node.VerifyEqual(0x1401, 2, 1, 254, "Transmission type")
	//node.VerifyEqual(0x1401, 3, 2, 0, "Inhibit time")
	//node.VerifyEqual(0x1401, 4, 1, 0x6020000, "Not implemented")
	node.VerifyRange(0x1402, 0, 1, 2, 3, "RxPdo")
	node.VerifyEqual(0x1402, 1, 4, 0x4000040B, "Cob ID")
	node.VerifyEqual(0x1402, 2, 1, 254, "Transmission type")
	//node.VerifyEqual(0x1402, 3, 2, 0, "Inhibit time")
	//node.VerifyEqual(0x1402, 4, 1, 0x6020000, "Not implemented")
	node.VerifyRange(0x1403, 0, 1, 2, 3, "RxPdo")
	node.VerifyEqual(0x1403, 1, 4, 0x4000050B, "Cob ID")
	node.VerifyEqual(0x1403, 2, 1, 254, "Transmission type")
	//node.VerifyEqual(0x1403, 3, 2, 0, "Inhibit time")
	//node.VerifyEqual(0x1403, 4, 1, 0x6020000, "Not implemented")
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

	node.VerifyEqual(0x1601, 0, 1, 4, "Rx pdo mapping count")
	node.VerifyEqual(0x1601, 1, 4, 0x24110510, "Rx pdo mapping 1")
	node.VerifyEqual(0x1601, 2, 4, 0x24110610, "Rx pdo mapping 2")
	node.VerifyEqual(0x1601, 3, 4, 0x24110710, "Rx pdo mapping 3")
	node.VerifyEqual(0x1601, 4, 4, 0x24110810, "Rx pdo mapping 4")

	node.VerifyEqual(0x1602, 0, 1, 4, "Rx pdo mapping count")
	node.VerifyEqual(0x1602, 1, 4, 0x24110910, "Rx pdo mapping 1")
	node.VerifyEqual(0x1602, 2, 4, 0x24110A10, "Rx pdo mapping 2")
	node.VerifyEqual(0x1602, 3, 4, 0x24110B10, "Rx pdo mapping 3")
	node.VerifyEqual(0x1602, 4, 4, 0x24110C10, "Rx pdo mapping 4")

	node.VerifyEqual(0x1603, 0, 1, 4, "Rx pdo mapping count")
	node.VerifyEqual(0x1603, 1, 4, 0x24110D10, "Rx pdo mapping 1")
	node.VerifyEqual(0x1603, 2, 4, 0x24110E10, "Rx pdo mapping 2")
	node.VerifyEqual(0x1603, 3, 4, 0x24110F10, "Rx pdo mapping 3")
	node.VerifyEqual(0x1603, 4, 4, 0x24111010, "Rx pdo mapping 4")
}

func VerifyTxPdoParameters(node *node.Node) {
	if node.SkipTest("Verify tx pdo parameters at 0x1800") {
		return
	}
	node.SetPreOperational()
	fmt.Printf("Version 2.6 will return 1 byte for inhibit time, not 2 as required by CAN-Open specification\n")
	node.VerifyEqual(0x1800, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1800, 1, 4, 0x4000018B, "Cob ID")
	node.VerifyEqual(0x1800, 2, 1, 4, "Transmission type")
	node.VerifyEqual(0x1800, 3, 2, 0, "Inhibit time") // Error in port code - should be 2 byte
	node.VerifyReadAbort(0x1800, 4, 1, "Not implemented")
	node.VerifyEqual(0x1801, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1801, 1, 4, 0x4000028B, "Cob ID")
	node.VerifyEqual(0x1801, 2, 1, 2, "Transmission type")
	node.VerifyEqual(0x1801, 3, 2, 0, "Inhibit time")
	//node.VerifyEqual(0x1801, 4, 1, 0x6020000, "Not implemented")
	node.VerifyEqual(0x1802, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1802, 1, 4, 0x4000038B, "Cob ID")
	node.VerifyEqual(0x1802, 2, 1, 2, "Transmission type")
	node.VerifyEqual(0x1802, 3, 2, 0, "Inhibit time")
	//node.VerifyEqual(0x1802, 4, 1, 0x6020000, "Not implemented")
	node.VerifyEqual(0x1803, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1803, 1, 4, 0x4000048B, "Cob ID")
	node.VerifyEqual(0x1803, 2, 1, 3, "Transmission type")
	node.VerifyEqual(0x1803, 3, 2, 0, "Inhibit time")
	//node.VerifyEqual(0x1803, 4, 1, 0x6020000, "Not implemented")
}

func VerifyTxPdoMapping(node *node.Node) {
	if node.SkipTest("Verify tx pdo mapping at 0x1A00") {
		return
	}
	node.SetPreOperational()
	node.VerifyEqual(0x1A00, 0, 1, 4, "Tx pdo mapping count")
	node.VerifyEqual(0x1A00, 1, 4, 0x24010110, "Tx pdo mapping 1")
	node.VerifyEqual(0x1A00, 2, 4, 0x24010210, "Tx pdo mapping 2")
	node.VerifyEqual(0x1A00, 3, 4, 0x24010310, "Tx pdo mapping 3")
	node.VerifyEqual(0x1A00, 4, 4, 0x24010410, "Tx pdo mapping 4")

	node.VerifyEqual(0x1A01, 0, 1, 4, "Tx pdo mapping count")
	node.VerifyEqual(0x1A01, 1, 4, 0x24010510, "Tx pdo mapping 1")
	node.VerifyEqual(0x1A01, 2, 4, 0x24010610, "Tx pdo mapping 2")
	node.VerifyEqual(0x1A01, 3, 4, 0x24010710, "Tx pdo mapping 3")
	node.VerifyEqual(0x1A01, 4, 4, 0x24010810, "Tx pdo mapping 4")

	node.VerifyEqual(0x1A02, 0, 1, 4, "Tx pdo mapping count")
	node.VerifyEqual(0x1A02, 1, 4, 0x24010910, "Tx pdo mapping 1")
	node.VerifyEqual(0x1A02, 2, 4, 0x24010A10, "Tx pdo mapping 2")
	node.VerifyEqual(0x1A02, 3, 4, 0x24010B10, "Tx pdo mapping 3")
	node.VerifyEqual(0x1A02, 4, 4, 0x24010C10, "Tx pdo mapping 4")

	node.VerifyEqual(0x1A03, 0, 1, 4, "Tx pdo mapping count")
	node.VerifyEqual(0x1A03, 1, 4, 0x24010D10, "Tx pdo mapping 1")
	node.VerifyEqual(0x1A03, 2, 4, 0x24010E10, "Tx pdo mapping 2")
	node.VerifyEqual(0x1A03, 3, 4, 0x24010F10, "Tx pdo mapping 3")
	node.VerifyEqual(0x1A03, 4, 4, 0x24011010, "Tx pdo mapping 4")
}


func VerifyRxPdo(n *node.Node) {
	if n.SkipTest("Testing rx pdo operation - transfer output values to card") {
		return
	}
	// Set input 3-8 to analog in
	n.SetPreOperational()
	for i:=3; i<9; i++ {
		n.WriteObject(0x4010, node.SubIndex(i), 1, SENS_VOLT)
		time.Sleep(50*time.Millisecond)
	}

	n.SetOperational()
	time.Sleep(10*time.Millisecond)
	n.SetPdoValue(4, 0, 2, 15000)
	for i:=0; i<4; i++ {
		n.Bus.SendSync()
		n.SendPdo(4)
		time.Sleep(100 * time.Millisecond)
	}
	n.ResetPdoCount()
	// Send 12 sync pulses
	for i:=0; i<12; i++ {
		n.Bus.SendSync()
		n.SendPdo(1)
		n.SendPdo(2)
		n.SendPdo(3)
		n.SendPdo(4)
		time.Sleep(100*time.Millisecond)
	}
	// The tx types are 4,2,2,3, giving 12/4=3, 12/2=6 and 12/3=4
	n.VerifyPdoCount(3, 6, 6, 4)


}

func SetIsolationMode(n *node.Node, timeoutMs int) {
	n.SetPreOperational()
	n.Check(n.WriteObject(ISOLATION_PDO_NUM, 0, 1, 4))
	if timeoutMs==0 {
		n.Check(n.WriteObject(ISOLATION_PDO_RATE, 0, 1, 0))
	} else {
		n.Check(n.WriteObject(ISOLATION_PDO_RATE, 0, 1, 100))
	}
	n.Check(n.WriteObject(ISOLATION_PDO_COUNT, 0, 1, (timeoutMs+99)/100))

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
	n.SetPreOperational()
	SetIsolationMode(n, 500)
	n.SetOperational()
	n.SetPdoValue(4, 0, 2, 15000)
	SendPdos(n,3, time.Millisecond*100)
	n.VerifyRange(0x2401, 9, 2, 14500, 15500, "Output ok")
	time.Sleep(800 * time.Millisecond)
	n.VerifyRange(0x2401, 9, 2, 0, 500, "Should have Isolation mode after 0.8sec")
	SetIsolationMode(n, 2000)

}

func VerifyTxPdo(n *node.Node) {
	if n.SkipTest("Verify tx pdos") {
		return
	}
	// Turn off isolation mode
	n.SetOperational()
	// Send pdo setting output 13 to +10mA (15000)
	n.SetPdoValue(4, 0, 2, 15000)
	SendPdos(n,4, time.Millisecond*100)
	// Verify that the value on channel 9 is also 15000 (because 9 and 13 is connected)
	n.VerifyRange(0x2401, 9, 2, 14500, 15500, "Sdo read of chan 9 at 10mA")
	// Verify that the content of pdo 3 from the node matches value read by sdo (ca 15000)
	n.VerifyPdoInt16(3, 0, 14500, 15501, "Analog channel 9 should be ca 15000 (10mA)")
	// Verify analog current at channel 13- should be 10mA.
	n.VerifyRangeFloat(0x4021, 13, 0.009, 0.011, "Float input chan 13")
	// Verify analog current at channel 13- should be -10mA.
	n.VerifyRangeFloat(0x4021, 9, -0.011, -0.009, "Float input chan 9")
}

