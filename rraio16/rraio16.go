package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"go-can/bus"
	"go-can/node"
	"go-can/peak"
	"go-can/psu"
	"time"
)
var pwr psu.Psu

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

var nodeId int
var failed bool


func SetOutput16bit(n *node.Node, channel int, value int) {
	_ = n.WriteObject(OUTPUT_16BIT, node.SubIndex(channel), 2, value)
}

func VerifyAin(n *node.Node) {
	if n.SkipTest("Testing analog inputs 3-8") {
		return
	}
	n.WriteObject(0x4010, node.SubIndex(5), 1, SENS_VOLT)
	_ = pwr.SetOutput(1, 0.0, 0.1)
	time.Sleep(500*time.Millisecond)
	_ = pwr.SetOutput(1, 0.1, 0.1)
	time.Sleep(500*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 5, 0.05, 0.15, "supply is 0.1V" )
	_ = pwr.SetOutput(1, 1.0, 0.1)
	time.Sleep(200*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 5, 0.9, 1.1, "supply is 1V" )
	_ = pwr.SetOutput(1, 20.0, 0.1)
	time.Sleep(200*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 5, 19.95, 20.1, "supply is 20V" )
	time.Sleep(100*time.Millisecond)
	n.VerifyRangeFloat(0x4021, 5, 19.95, 20.1, "supply is 20V" )
	n.WriteObject(0x4010, node.SubIndex(5), 1, SENS_DIG_IN)
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
	//fmt.Printf("Number of pdos: %d, %d, %d, %d\n", n.PdoCount[1], n.PdoCount[2],n.PdoCount[3],n.PdoCount[4])
	n.VerifyPdoCount(2, 5, 5, 3)
	n.VerifyPdoInt16(2, 0, 14500, 15500, "Analog channel 5 should be ca 10V")
	n.VerifyPdoInt16(3, 0, 14500, 15500, "Analog channel 9 should be ca 10mA")
	// Wait 1000mS and check for timeout while we send sync messages
	time.Sleep(3000 * time.Millisecond)
	n.VerifyRange(0x2401, 9, 2, 0, 800, "After pdo timeout of 2.5 sec")
	n.Bus.SendSync()
}

func  VerifyHeartbeat(node *node.Node) {
	if node.SkipTest("Verify heartbeat operation") {
		return
	}
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

func VerifyBootupMessages(node *node.Node) {
	if node.SkipTest("Verify bootup messages") {
		return
	}
	node.Verify(node.EmcyCount == 1, "Emergency count should be 1, was %d", node.EmcyCount)
	node.Verify(node.HeartbeatCount == 1, "Emergency count should be 1, was %d", node.HeartbeatCount)
}

func main() {
	flag.IntVar(&node.SubTest, "no", 0, "Set to the subtest number that should be executed. Set to zero to run all tests")
	flag.IntVar(&nodeId, "node", 11, "Node number to test. Uses 11 as default.")
	flag.Parse()

	fmt.Printf("This is a functional testing of RRAIO16 software\n")
	fmt.Printf("A TTi power supply connected via USB will be used if present. Fallback is manual settings.\n")
	fmt.Printf("The can bus is connected via a Peak USB adapter as device 1\n")
	fmt.Printf("Power supply channel 1 is 10.000V, and must be connected to AI13 (treminal 23)\n")
	fmt.Printf("Power supply channel 2 is 24.000V, and must be connected to terminal 1 and 9\n")
	fmt.Printf("Power supply ground for cahnnel 1 and 2 must be connected to terminal 2 and 10\n\n")
	fmt.Printf("Setting up adapter and bus\n")

	// Setup power supply
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

	n := node.New(b, nodeId)

	fmt.Printf("Turn on power, please wait...\n")
	_ = pwr.SetOutput(1, 10.0, 0.1)
	_ = pwr.SetOutput(2, 24.0, 0.5)
	_ = pwr.EnableOutput(2)
	_ = pwr.EnableOutput(1)
	time.Sleep(time.Millisecond*1500)
	VerifyBootupMessages(n)
	n.Bus.Reset()
	time.Sleep(time.Millisecond*100)
	VerifyMandatoryObjects(n, nodeId, 250,3300, 0x622e38, 0x362e32, )
	VerifyTxPdoParameters(n)
	VerifyTxPdoMapping(n)
	VerifyRxPdoParameters(n)
	VerifyRxPdoMapping(n)
	VerifyHeartbeat(n)
	VerifyDigOut(n)
	VerifyRxPdo(n)
	VerifyAin(n)
	if n.Failed {
		color.Error.Printf("Test failed!\n")
	} else {
		fmt.Printf("Test ok\n")
	}
	b.Close()
}
