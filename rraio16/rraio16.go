package main

import (
	"fmt"
	"go-can/bus"
	"go-can/node"
	"go-can/peak"
	"go-can/psu"
	"time"
)

const (
	INPUT_16BIT         = 0x2401    // iaAnalogInput
	OUTPUT_16BIT        = 0x2411    // iaAnalogOutput
	ISOLATION_PDO_NUM   = 0x5020
	ISOLATION_PDO_RATE  = 0x5021
	ISOLATION_PDO_COUNT = 0x5022
	DIGITAL_ISO_OUT     = 0x5501    // baDouIniStat
	DIGITAL_PASSIVE_OUT = 0x5504    // baDouPassive
)

var nodeId = 11
var failed bool


func SetOutput16bit(n *node.Node, channel int, value int) {
	_ = n.WriteObject(OUTPUT_16BIT, node.SubIndex(channel), 2, value)
}

func VerifyPdoParameters(node *node.Node) {	node.VerifyEqual(0x1400, 0, 1, 3, "RxPdo")

	fmt.Println("Verify pdo parameters at 0x1400..")
	node.VerifyEqual(0x1400, 1, 4, 0x4000020B, "Cob ID")
	node.VerifyEqual(0x1400, 2, 1, 254, "Transmission type")
	node.VerifyEqual(0x1400, 3, 2, 0, "Inhibit time")
	node.VerifyReadAbort(0x1400, 4, 1, "Not implemented")
	node.VerifyEqual(0x1401, 0, 1, 3, "RxPdo")
	node.VerifyEqual(0x1401, 1, 4, 0x4000030B, "Cob ID")
	node.VerifyEqual(0x1401, 2, 1, 254, "Transmission type")
	node.VerifyEqual(0x1401, 3, 2, 0, "Inhibit time")
	//node.VerifyEqual(0x1401, 4, 1, 0x6020000, "Not implemented")
	node.VerifyEqual(0x1402, 0, 1, 3, "RxPdo")
	node.VerifyEqual(0x1402, 1, 4, 0x4000040B, "Cob ID")
	node.VerifyEqual(0x1402, 2, 1, 254, "Transmission type")
	node.VerifyEqual(0x1402, 3, 2, 0, "Inhibit time")
	//node.VerifyEqual(0x1402, 4, 1, 0x6020000, "Not implemented")
	node.VerifyEqual(0x1403, 0, 1, 3, "RxPdo")
	node.VerifyEqual(0x1403, 1, 4, 0x4000050B, "Cob ID")
	node.VerifyEqual(0x1403, 2, 1, 254, "Transmission type")
	node.VerifyEqual(0x1403, 3, 2, 0, "Inhibit time")
	//node.VerifyEqual(0x1403, 4, 1, 0x6020000, "Not implemented")

	fmt.Println("Verify rx pdo mapping at 0x1600..")
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

	fmt.Println("Verify tx pdo parameters at 0x1800..")
	node.VerifyEqual(0x1800, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1800, 1, 4, 0x4000018B, "Cob ID")
	node.VerifyEqual(0x1800, 2, 1, 4, "Transmission type")
	node.VerifyEqual(0x1800, 3, 1, 0, "Inhibit time")  // Error in port code - should be 2 byte
	node.VerifyReadAbort(0x1800, 4, 1,  "Not implemented")
	node.VerifyEqual(0x1801, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1801, 1, 4, 0x4000028B, "Cob ID")
	node.VerifyEqual(0x1801, 2, 1, 2, "Transmission type")
	node.VerifyEqual(0x1801, 3, 1, 0, "Inhibit time")
	//node.VerifyEqual(0x1801, 4, 1, 0x6020000, "Not implemented")
	node.VerifyEqual(0x1802, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1802, 1, 4, 0x4000038B, "Cob ID")
	node.VerifyEqual(0x1802, 2, 1, 2, "Transmission type")
	node.VerifyEqual(0x1802, 3, 1, 0, "Inhibit time")
	//node.VerifyEqual(0x1802, 4, 1, 0x6020000, "Not implemented")
	node.VerifyEqual(0x1803, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1803, 1, 4, 0x4000048B, "Cob ID")
	node.VerifyEqual(0x1803, 2, 1, 3, "Transmission type")
	node.VerifyEqual(0x1803, 3, 1, 0, "Inhibit time")
	//node.VerifyEqual(0x1803, 4, 1, 0x6020000, "Not implemented")

	fmt.Println("Verify tx pdo mapping at 0x1A00...")
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

func VerifyRxPdo(node *node.Node) {
	node.SetOperational()
	time.Sleep(50*time.Millisecond)
	// Send 12 sync pulses
	for i:=0; i<13; i++ {
		node.Bus.SendSync()
		time.Sleep(50*time.Millisecond)
	}
}

func  VerifyHeartbeat(node *node.Node) {
	fmt.Println("Verify heartbeat operation")
	//node.Reset()
	//time.Sleep(800*time.Millisecond)
	node.HeartbeatCount = 0
	// Set heartbeet at 100mS
	err := node.WriteObject(HEARTBEAT_TIME, 0, 2, 100)
	if err!=nil {
		fmt.Printf("Error writing heartbeat time, %s\n", err)
		return
	}
	// and set operational
	node.SetOperational()
	// and wait 1 second
	time.Sleep(time.Second)
	// and set preoperational
	node.SetPreOperational()
	n := node.HeartbeatCount
	if n < 9 || n > 11 {
		failed = true
		fmt.Printf("Did not get correct number of hearbeat messages, expected ca 10, got %d\n", n)
	}
}

func main() {
	fmt.Printf("Setting up adapter and bus\n")
	pwr, _ := psu.NewTtiPsu("")
	_ = pwr.DisableOutput(1)
	_ = pwr.DisableOutput(2)
	time.Sleep(time.Millisecond*1000)

	b := bus.New(
		peak.New(peak.PCAN_USBBUS1, 125000),
		100*time.Millisecond)

	n := node.New(b, 11)

	fmt.Printf("Turn on power, please wait...\n")
	_ = pwr.SetOutput(1, 10.0, 0.1)
	_ = pwr.SetOutput(2, 24.0, 0.3)
	_ = pwr.EnableOutput(2)
	_ = pwr.EnableOutput(1)
	time.Sleep(time.Second*2)
	fmt.Printf("Verify bootup messages\n")
	n.Verify(n.EmcyCount==1, "Emergency count should be 1, was %d", n.EmcyCount)
	n.Verify(n.HeartbeatCount==1, "Emergency count should be 1, was %d", n.EmcyCount)
	// We must reset here to get out of BusOff
	b.Reset()

	VerifyMandatoryObjects(n, 11, 250,3300, 0x622e38, 0x362e32, )
	VerifyPdoParameters(n)
	VerifyHeartbeat(n)
	VerifyRxPdo(n)
	if n.Failed {
		fmt.Printf("******* Test failed!\n")
	} else {
		fmt.Printf("Test ok\n")
	}
	b.Close()
	//p.Shutdown()
}
