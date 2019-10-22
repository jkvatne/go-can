package main

import (
	"fmt"
	"github.com/gookit/color"
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
const (
SENS_NONE        = 0
SENS_MA          = 1
SENS_20MA        = 2
SENS_VOLT        = 3
SENS_VOLT_OUT    = 4
SENS_MA_OUT      = 5
SENS_DIG_OUT     = 6
SENS_FREQ        = 7
SENS_DIG_IN      = 8
SENS_QUADRATURE  = 10
SENS_FREQX10     = 11
SENS_FREQX100    = 12
SENS_MA_SCALED   = 21
SENS_VOLT_SCALED = 23
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

func VerifyDigOut(n *node.Node) {
	fmt.Printf("Testing digital outputs, reading voltage back\n")
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

func VerifyRxPdo(n *node.Node) {
	fmt.Printf("Testing rx pdo operation - transfer output values to card\n")
	// Set input 3-8 to analog in
	n.SetPreOperational()
	for i:=3; i<9; i++ {
		n.WriteObject(0x4010, node.SubIndex(i), 1, SENS_VOLT)
		time.Sleep(50*time.Millisecond)
	}

	n.SetOperational()
	time.Sleep(50*time.Millisecond)
	n.SetPdoValue(4, 0, 2, 15000)
	n.SendPdo(4)
	time.Sleep(50*time.Millisecond)
	// Send 12 sync pulses
	for i:=0; i<12; i++ {
		time.Sleep(100*time.Millisecond)
		n.Bus.SendSync()
		n.SendPdo(1)
		n.SendPdo(2)
		n.SendPdo(3)
		n.SendPdo(4)
	}
	fmt.Printf("Number of pdos: %d, %d, %d, %d\n", n.PdoCount[1], n.PdoCount[2],n.PdoCount[3],n.PdoCount[4])
	n.VerifyPdoCount(2, 5, 5, 3)
	n.VerifyPdoInt16(2, 0, 14500, 15500, "Analog channel 5 should be ca 10V")
	n.VerifyPdoInt16(3, 0, 14500, 15500, "Analog channel 9 should be ca 10mA")
	// Wait 1000mS and check for timeout while we send sync messages
	time.Sleep(3000 * time.Millisecond)
	n.VerifyRange(0x2401, 9, 2, 0, 800, "After pdo timeout of 2.5 sec")
	n.Bus.SendSync()
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
	fmt.Printf("This is a functional testing of RRAIO16 software\n")
	fmt.Printf("A TTi power supply connected via USB will be used if present. Fallback is manual settings.\n")
	fmt.Printf("The can bus is connected via a Peak USB adapter as device 1\n")
	fmt.Printf("Power supply channel 1 is 10.000V, and must be connected to AI13 (treminal 23)\n")
	fmt.Printf("Power supply channel 2 is 24.000V, and must be connected to terminal 1 and 9\n")
	fmt.Printf("Power supply ground for cahnnel 1 and 2 must be connected to terminal 2 and 10\n\n")
	fmt.Printf("Setting up adapter and bus\n")

	// Setup power supply
	var pwr psu.Psu
	pwr, _ = psu.NewTtiPsu("")
	if pwr==nil {
		pwr, _ = psu.NewManualPsu("")
	}
	_ = pwr.DisableOutput(1)
	_ = pwr.DisableOutput(2)
	time.Sleep(time.Millisecond*100)

	// Setup can bus peak usb adapter
	b := bus.New(
		peak.New(peak.PCAN_USBBUS1, 125000),
		100*time.Millisecond)

	n := node.New(b, 11)

	fmt.Printf("Turn on power, please wait...\n")
	_ = pwr.SetOutput(1, 10.0, 0.1)
	_ = pwr.SetOutput(2, 24.0, 0.5)
	_ = pwr.EnableOutput(2)
	_ = pwr.EnableOutput(1)
	time.Sleep(time.Millisecond*1500)
	fmt.Printf("Verify bootup messages\n")
	n.Verify(n.EmcyCount==1, "Emergency count should be 1, was %d", n.EmcyCount)
	n.Verify(n.HeartbeatCount==1, "Emergency count should be 1, was %d", n.EmcyCount)

	n.Reset()
	// NB: The peak adapter needs 1.2sec delay after a reset!
	time.Sleep(time.Millisecond*1200)

	VerifyMandatoryObjects(n, 11, 250,3300, 0x622e38, 0x362e32, )
	VerifyPdoParameters(n)
	VerifyHeartbeat(n)
	VerifyDigOut(n)
	VerifyRxPdo(n)
	if n.Failed {
		color.Error.Printf("Test failed!\n")
	} else {
		fmt.Printf("Test ok\n")
	}
	b.Close()
}
