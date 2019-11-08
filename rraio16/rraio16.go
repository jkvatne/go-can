package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"go-can/bus"
	"go-can/node"
	"go-can/peak"
	"go-can/psu"
	"os"
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
var Vsupply = 20.1
var togglePower bool

func VerifyBootupMessages(node *node.Node) {
	node.Verify(node.EmcyCount == 1, "Emergency count should be 1, was %d", node.EmcyCount)
	node.Verify(node.HeartbeatCount == 1, "Emergency count should be 1, was %d", node.HeartbeatCount)
}

func main() {
	flag.IntVar(&node.SubTest, "subtest", 0, "Set to the subtest number that should be executed. Set to zero to run all tests")
	flag.IntVar(&nodeId, "node", 11, "Node number to test. Uses 11 as default.")
	flag.BoolVar(&togglePower, "toggle-power", true, "Set to false to disable power off-on at start of test. This also skips test of bootup messages.")
	info := flag.Bool("info",false,"Print help info")
	testPower := flag.Bool("test-power",false,"Check if a power supply is connected and verify connection")
	pwrPort := flag.String("power-port","","Name of com-port used to control the power supply, defaults to device with highest com-port number")
	flag.Parse()
	if *info {
		fmt.Printf("This is a functional testing of RRAIO16 software\n")
		fmt.Printf("* A TTiCPX4000 power supply connected via USB will be used if present. Fallback is manual settings.\n")
		fmt.Printf("* The can bus is connected via a Peak USB adapter as device 1\n")
		fmt.Printf("* Power supply channel 2 is 20.000V, and must be connected to terminal 1-2 and 9-10\n")
		fmt.Printf("* Channel 9 and 13 is connected\n")
		fmt.Printf("* Channel 10 and 14 is connected\n")
		fmt.Printf("* Channel 11 and 15 is connected\n")
		fmt.Printf("* Channel 12 and 16 is connected\n")
		fmt.Printf("* Channel 1 and 3 is connected\n")
		fmt.Printf("* Channel 2 and 4 is connected\n")
		fmt.Printf("* Channel 1 and 7 is connected via 1 kohm resistor\n")
		fmt.Printf("* Channel 2 and 8 is connected via 1 kohm resistor\n")
	}

	// Setup power supply
	pwr, _ = psu.NewTtiPsu(*pwrPort)
	if pwr!=nil {
		fmt.Printf("Power supply name is \"%s\"\n", pwr.Name())
	}
	if *testPower {
		if pwr==nil {
			fmt.Printf("Could not find power supply\n")
			return
		}
		volt,cur,err := pwr.GetOutput(1)
		if err!=nil {
			fmt.Printf("Error reading volt/current setting, %s\n", err)
		} else {
			fmt.Printf("Actual outputs on channel 1 is %0.3fV, %0.3fA\n", volt, cur)
		}
		volt,cur,err = pwr.GetSetpoint(1)
		if err!=nil {
			fmt.Printf("Error reading volt/current setting, %s\n", err)
		} else {
			fmt.Printf("Setpoint on channel en is %0.3fV, %0.3fA\n", volt, cur)
		}
		return
	}

	if pwr==nil {
		fmt.Printf("Using manual power supply control\n")
		pwr, _ = psu.NewManualPsu("")
	}

	// Setup can bus peak usb adapter
	dev,err := peak.New(peak.PCAN_USBBUS1, 125000)
	if err!=nil {
		fmt.Printf("Peak adapter initialization failed, %s", err)
		os.Exit(1)
	}
	if togglePower {
		_ = pwr.SetOutput(2, 20.0, 0.5, false)
		time.Sleep(time.Millisecond * 1000)
	}
	b := bus.New(dev, 100*time.Millisecond)
	n := node.New(b, nodeId)

	if togglePower {
		_ = pwr.SetOutput(2, 20.0, 0.5, true)
		fmt.Printf("Waiting for boot-up\n")
		time.Sleep(time.Millisecond * 1500)
		VerifyBootupMessages(n)
		n.Bus.Reset()
	}
	time.Sleep(time.Millisecond*100)
	VerifyMandatoryObjects(n, nodeId, 250,3300, 0x622e38, 0x362e32, )
	VerifyTxPdoParameters(n)
	VerifyTxPdoMapping(n)
	VerifyRxPdoParameters(n)
	VerifyRxPdoMapping(n)
	VerifyHeartbeat(n)
	VerifyEmcyOkStartingPdos(n)
	VerifyRxPdo(n)
	VerifyTxPdo(n)
	VerifyDigOut(n)
	VerifyAin(n)
	VerifyFrequency(n)
	VerifyIsolationModeTime(n)
	if n.Failed {
		color.Error.Printf("Test failed!\n")
	} else {
		fmt.Printf("Test ok\n")
	}
	b.Close()
}
