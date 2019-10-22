package main

import (
	"go-can/node"
)

func VerifyRxPdoParameters(node *node.Node) {
	if node.SkipTest("Verify pdo parameters at 0x1400..") {
		return
	}
	node.SetPreOperational()
	node.VerifyEqual(0x1400, 0, 1, 3, "RxPdo")
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
}

func VerifyRxPdoMapping(node *node.Node) {
	if node.SkipTest("Verify pdo mapping at 0x1600..") {
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
	if node.SkipTest("Verify tx pdo parameters at 0x1800..") {
		return
	}
	node.SetPreOperational()
	node.VerifyEqual(0x1800, 0, 1, 3, "TxPdo")
	node.VerifyEqual(0x1800, 1, 4, 0x4000018B, "Cob ID")
	node.VerifyEqual(0x1800, 2, 1, 4, "Transmission type")
	node.VerifyEqual(0x1800, 3, 1, 0, "Inhibit time") // Error in port code - should be 2 byte
	node.VerifyReadAbort(0x1800, 4, 1, "Not implemented")
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
}

func VerifyTxPdoMapping(node *node.Node) {
	if node.SkipTest("Verify tx pdo mapping at 0x1A00..") {
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

