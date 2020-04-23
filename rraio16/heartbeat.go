package main

import (
	"fmt"
	"go-can/node"
	"time"
)

func  VerifyHeartbeat(node *node.Node) {
	if node.SkipTest("Verify heartbeat operation") {
		return
	}
	node.HeartbeatCount = 0
	// Set heartbeet at 100mS
	err := node.WriteObject(HEARTBEAT_TIME, 0, 2, 100)
	if err!=nil {
		node.Failed = true
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
		node.Failed = true
		fmt.Printf("Did not get correct number of hearbeat messages, expected ca 10, got %d\n", n)
	}
	_ = node.WriteObject(HEARTBEAT_TIME, 0, 2, 0)

}


