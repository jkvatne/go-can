package main

import (
	"fmt"
	"go-can/can"
	"go-can/peak"
	"time"
)
const (
	DEVICE_TYPE         = 0x1000
	ERROR_REGISTER      = 0x1001
	PREDEF_ERROR_FIELD  = 0x1003
	SYNC_COBID          = 0x1005
	CYCLE_PERIODE       = 0x1006
	SYNC_WINDOW         = 0x1007
	HARDWARE_VERSION    = 0x1009
	SOFTWARE_VERSION    = 0x100A
	GUARD_TIME          = 0x100C
	GUARD_TIME_FACT     = 0x100D
	COB_ID_TIME         = 0x1012
	COB_ID_EMCY         = 0x1014
	HEARTBEAT_TIME      = 0x1017
	IDENTITY            = 0x1018
)

var nodeId = 11
var failed bool

func verify(ok bool, description string) {
	if !ok {
		failed = true
		fmt.Println("Error : " + description)
	}
}

type Node struct {
	connection *can.Connection
	id uint8
}

func (node *Node) VerifyObject(index can.Index, subindex can.Subindex, bytecount int, expected int, description string) {
	value, err := node.connection.SdoRead(node.id, index, subindex, bytecount)
	if err!=nil {
		failed = true
		fmt.Printf("Error reading sdo")
	}
	verify(value==expected, fmt.Sprintf("Expected %x, Actual %x, %s", expected, value, description))
}

func (node *Node) SetPreoperational() {
	node.connection.SetPreoperational(node.id)
}

func (node *Node) VerifyMandatoryObjects() {
	fmt.Println("Verify mandatory objects from 0x1000 to 0x1018")
	node.SetPreoperational()
	node.VerifyObject(DEVICE_TYPE, 0, 4, 3300, "Device type")
	//TODO node.VerifyAbort(node, DEVICE_TYPE, 1, 4, 3300, "Device type")
	node.VerifyObject(SOFTWARE_VERSION, 0, 4, 0x20362e32, "Software version text")
	node.VerifyObject(ERROR_REGISTER, 0, 1, 0, "Error register")
	node.VerifyObject(PREDEF_ERROR_FIELD, 0, 1, 1, "PREDEF_ERROR_FIELD count")
	node.VerifyObject(PREDEF_ERROR_FIELD, 1, 4, 0xFF00, "PREDEF_ERROR_FIELD")
	node.VerifyObject(SYNC_COBID, 0, 4, 0x80000080, "SYNC_COBID")
	node.VerifyObject(CYCLE_PERIODE, 0, 4, 0, "CYCLE_PERIODE")
	node.VerifyObject(SYNC_WINDOW, 0, 4, 0, "SYNC_WINDOW")
	node.VerifyObject(SYNC_WINDOW, 0, 4, 0, "SYNC_WINDOW")
	node.VerifyObject(HARDWARE_VERSION, 0, 4, 0x20622e38, "HARDWARE_VERSION")
	node.VerifyObject(SOFTWARE_VERSION, 0, 4, 0x20362e32, "SOFTWARE_VERSION")
	node.VerifyObject(GUARD_TIME, 0, 2, 0, "GUARD_TIME")
	node.VerifyObject(GUARD_TIME_FACT, 0, 1, 0, "GUARD_TIME_FACT")
	node.VerifyObject(COB_ID_TIME, 0, 4, 0x80000100, "COB_ID_TIME")
	node.VerifyObject(COB_ID_EMCY, 0, 4, 0x8B, "COB_ID_EMCY")
	node.VerifyObject(HEARTBEAT_TIME, 0, 2, 0, "HEARTBEAT_TIME")
	node.VerifyObject(IDENTITY, 0, 1, 4, "Identity size")
	node.VerifyObject(IDENTITY, 1, 4, 250, "Vendor id")
	node.VerifyObject(IDENTITY, 2, 4, 3300, "Product code")
	node.VerifyObject(IDENTITY, 3, 4, 0x1F0101, "Revision number")
	node.VerifyObject(IDENTITY, 4, 4, 0, "Serial number")

	node.VerifyObject(0x1400, 0, 1, 3, "RxPdo")
	node.VerifyObject(0x1400, 1, 4, 0x4000020B, "Cob ID")
	node.VerifyObject(0x1400, 2, 1, 254, "Transmission type")
	node.VerifyObject(0x1400, 3, 2, 0, "Inhibit time")
	node.VerifyObject(0x1400, 4, 1, 0x6020000, "Not implemented")
	node.VerifyObject(0x1401, 0, 1, 3, "RxPdo")
	node.VerifyObject(0x1401, 1, 4, 0x4000030B, "Cob ID")
	node.VerifyObject(0x1401, 2, 1, 254, "Transmission type")
	node.VerifyObject(0x1401, 3, 2, 0, "Inhibit time")
	node.VerifyObject(0x1401, 4, 1, 0x6020000, "Not implemented")
	node.VerifyObject(0x1402, 0, 1, 3, "RxPdo")
	node.VerifyObject(0x1402, 1, 4, 0x4000040B, "Cob ID")
	node.VerifyObject(0x1402, 2, 1, 254, "Transmission type")
	node.VerifyObject(0x1402, 3, 2, 0, "Inhibit time")
	node.VerifyObject(0x1402, 4, 1, 0x6020000, "Not implemented")
	node.VerifyObject(0x1403, 0, 1, 3, "RxPdo")
	node.VerifyObject(0x1403, 1, 4, 0x4000050B, "Cob ID")
	node.VerifyObject(0x1403, 2, 1, 254, "Transmission type")
	node.VerifyObject(0x1403, 3, 2, 0, "Inhibit time")
	node.VerifyObject(0x1403, 4, 1, 0x6020000, "Not implemented")
}

func handler(m *can.Msg) {
	fmt.Printf("Callback with message: %s\n",m.ToString())
}

func main() {
	fmt.Printf("Setting up adapter and connection\n")
	connection := can.NewConnection(
		peak.New(peak.PCAN_USBBUS1, 125000),
		100*time.Millisecond,
		handler)

	fmt.Println("Testing SdoRead()")
	devId, _ := connection.SdoRead(11, 0x1000, 0, 4)
	fmt.Printf("SdoRead from index 0x1000 from node 11, result is %d\n", devId)

	node := Node{connection,11}
	node.VerifyMandatoryObjects()
	connection.Close()
	fmt.Printf("Done\n")
}
