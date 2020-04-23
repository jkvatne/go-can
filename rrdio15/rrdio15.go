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
	INPUT_DIG           = 0x2100    // Two bytes with digital inputs
	OUTPUT_DIG          = 0x2200    // Five bytes with digital outputs (including relays)
	OUTPUT_STATE        = 0x2201    // FIve bytes with digital output state
	INPUT_16BIT         = 0x2401    // iaAnalogInputs, 16 words
	OUTPUT_16BIT        = 0x2411    // iaAnalogOutput (not used)
	OUTPUT_ANLG_STATE   = 0x2412    // Analog Output states (not used)
	SENSOR_TYPE         = 0x4010    // 16 bytes of sensor types
	ANALOG16BIT         = 0x4020
	INPUT_VALUE         = 0x4021
	SCALED_VALUE        = 0x4022
	OUTPUT_VALUE        = 0x4023
	BIT_COUNT           = 0x4025
	SENSOR_SCAN         = 0x4026
	LOOP_SCAN           = 0x4027
	FREQUENCY           = 0x4028
	ERROR_STATUS        = 0x4029
	NODE_STATE          = 0x5001
	SYNC_PAR_1          = 0x5010
	SYNC_PAR_2          = 0x5011
	SYNC_PAR_3          = 0x5012
	SYNC_PAR_4          = 0x5013
	ISOLATION_PDO_NUM   = 0x5020
	ISOLATION_PDO_RATE  = 0x5021
	ISOLATION_PDO_COUNT = 0x5022
	DIN_LOGIC           = 0x5400
	DOU_LOGIC           = 0x5500
	DOU_INI             = 0x5501
	DOU_OUTEN           = 0x5502
	DOU_PASSIVE         = 0x5504
	AIN_EUNIT           = 0x5602  // Not used
	AIN_DEC             = 0x5603  // Not used
	AIN_LOWRANGE        = 0x5604  // Not used (no rangecheck enabled)
	AIN_HIGHRANGE       = 0x5605  // Not used             "


	DIGITAL_ISO_OUT     = 0x5501    // baDouIniStat
	DIGITAL_PASSIVE_OUT = 0x5504    // baDouPassive
)

const (
IO_NONE        = 0
IO_OUTPUT      = 1
IO_PNP         = 2
IO_NPN         = 3
IO_TTL         = 4
IO_NAMUR       = 5
IO_ANALOG      = 6
IO_FREQ_HZ     = 7
IO_FREQ_HZx10  = 8
IO_FREQ_HZx100 = 9
IO_PWM_OUTPUT  = 10
)

var nodeId int
var Vsupply = 20.0   // 21.8 when using old card, Is actualy 20.0
var togglePower bool

func VerifyBootupMessages(node *node.Node) {
	node.Verify(node.EmcyCount == 1, "Emergency count should be 1, was %d", node.EmcyCount)
	fmt.Printf("LastEmcyMsg was %x, %x, %x, %x, %x, %x, %x, %x\n",
		node.LastEmcyMsg[0], node.LastEmcyMsg[1], node.LastEmcyMsg[2],node.LastEmcyMsg[3],node.LastEmcyMsg[4],node.LastEmcyMsg[5],node.LastEmcyMsg[6],node.LastEmcyMsg[7])
	if node.LastEmcyMsg[2]!=0 {
		fmt.Printf("Error register (third byte) should be 1, was actualy 0x%x\n", node.LastEmcyMsg[2])
	}
	node.Verify(node.HeartbeatCount == 1, "Heartbeat count should be 1, was %d", node.HeartbeatCount)
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
		os.Exit(0)
	}

	// Setup power supply
	pwr = psu.NewPsu(*pwrPort)
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
	// Toggle power to give a clean startup
	if togglePower {
		_ = pwr.Disable(1)
		time.Sleep(time.Millisecond * 1000)
	}
	b := bus.New(dev, 100*time.Millisecond)
	n := node.New(b, nodeId)
	if togglePower {
		_ = pwr.SetOutput(1, 20.0, 0.5)
		fmt.Printf("Waiting for boot-up\n")
		time.Sleep(time.Millisecond * 1500)
		VerifyBootupMessages(n)
		n.Bus.Reset()
	}
	time.Sleep(time.Millisecond*100)

	_, err = n.ReadObject(0x1000, 0, 4)
	if err!=nil {
		fmt.Printf("No card found, terminating test\n")
		os.Exit(1)
	}

	VerifyMandatoryObjects(n, nodeId, 250,2300, 0x302e39, 0x352e32,0x1f0102 )
	VerifyTxPdoParameters(n)
	VerifyTxPdoMapping(n)
	VerifyRxPdoParameters(n)
	VerifyRxPdoMapping(n)
	VerifyHeartbeat(n)
	VerifyEmcyOkStartingPdos(n)
	VerifyRxPdo(n)
	VerifyDigOut(n)
	VerifyFrequency(n)
	VerifyIsolationModeTime(n)
	if n.Failed {
		color.Error.Printf("Test failed!\n")
	} else {
		fmt.Printf("Test ok\n")
	}
	b.Close()
}
