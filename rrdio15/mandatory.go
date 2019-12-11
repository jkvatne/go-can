package main

import (
	"go-can/node"
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

func VerifyMandatoryObjects(node *node.Node, id int, vendor int, deviceType int, hw int, sw int, rev int) {
	if node.SkipTest("Verify mandatory objects from 0x1000 to 0x1018") {
		return
	}
	node.Bus.SetPreoperational(0)
	node.VerifyEqual(DEVICE_TYPE, 0, 4, deviceType, "Device type")
	//TODO node.VerifyAbort(DEVICE_TYPE, 1, 4, 3300, "Device type")
	node.VerifyEqual(ERROR_REGISTER, 0, 1, 0, "Error register")
	node.VerifyEqual(PREDEF_ERROR_FIELD, 0, 1, 1, "PREDEF_ERROR_FIELD count")
	node.VerifyEqual(PREDEF_ERROR_FIELD, 1, 4, 0xFF00, "PREDEF_ERROR_FIELD")
	node.VerifyEqual(SYNC_COBID, 0, 4, 0x80000080, "SYNC_COBID")
	node.VerifyEqual(CYCLE_PERIODE, 0, 4, 0, "CYCLE_PERIODE")
	node.VerifyEqual(SYNC_WINDOW, 0, 4, 0, "SYNC_WINDOW")
	node.VerifyEqual(HARDWARE_VERSION, 0, 3, hw, "HARDWARE_VERSION")
	node.VerifyEqual(SOFTWARE_VERSION, 0, 3, sw, "SOFTWARE_VERSION")
	node.VerifyEqual(GUARD_TIME, 0, 2, 0, "GUARD_TIME")
	node.VerifyEqual(GUARD_TIME_FACT, 0, 1, 0, "GUARD_TIME_FACT")
	node.VerifyEqual(COB_ID_TIME, 0, 4, 0x80000100, "COB_ID_TIME")
	node.VerifyEqual(COB_ID_EMCY, 0, 4, 0x80+nodeId, "COB_ID_EMCY")
	node.VerifyRange(HEARTBEAT_TIME, 0, 2, 0, 1000, "HEARTBEAT_TIME")
	node.VerifyEqual(IDENTITY, 0, 1, 4, "Identity size")
	node.VerifyEqual(IDENTITY, 1, 4, vendor, "Vendor id")
	node.VerifyEqual(IDENTITY, 2, 4, deviceType, "Product code")
	node.VerifyEqual(IDENTITY, 3, 4, rev, "Revision number")
	node.VerifyEqual(IDENTITY, 4, 4, 0, "Serial number")
}


