package node

import (
	"fmt"
	"go-can/can"
	"go-can/bus"
)

type NodeId uint8
type Index uint16
type SubIndex uint8
type Callback func(msg *can.Msg)


type Node struct {
	Bus            *bus.State
	Id             NodeId
	HeartbeatCount int
	State          int
	Failed         bool
	LastEmcyMsg    [8]byte
	EmcyCount      int
}

var Nodes[127] *Node

func AddNode(node *Node) {
	Nodes[node.Id] = node
}

const (
	FUNC_NMT    = 0x00
	FUNC_EMCY   = 0x01
	FUNC_SYNC   = 0x01
	FUNC_TIME   = 0x02
	FUNC_TXPDO1 = 0x03
	FUNC_RXPDO1 = 0x04
	FUNC_TXPDO2 = 0x05
	FUNC_RXPDO2 = 0x06
	FUNC_TXPDO3 = 0x07
	FUNC_RXPDO3 = 0x08
	FUNC_TXPDO4 = 0x09
	FUNC_RXPDO4 = 0x0A
	FUNC_TXSDO  = 0x0B
	FUNC_RXSDO  = 0x0C
	FUNC_HEARTBEAT = 0x0E
	FUNC_MGMT   = 0x0F
)

func DefaultId(base int, nodeId NodeId) can.CobId {
	return can.CobId(base)+can.CobId(nodeId)
}

func MsgHandler(msg *can.Msg) {
	funcCode := (int(msg.Id) & 0x780) >> 7
	nodeId := int(msg.Id) & 0x07F

	if funcCode == FUNC_TXSDO {
		// Handle received sdo request from external node
		// If the node id is own id, return a value from own object dictionary
		// Skip if the request was to another slave
	} else if funcCode==FUNC_HEARTBEAT {
		// Incoming heartbeat messages
		if Nodes[nodeId]!=nil {
			Nodes[nodeId].HeartbeatCount++
			Nodes[nodeId].State = int(msg.Data[0])
		}
	} else if funcCode == FUNC_TXPDO1 {
		Nodes[nodeId].HandlePdo((funcCode+1)/2-1, msg.Data)

	} else if funcCode == FUNC_EMCY {
		Nodes[nodeId].HandleEmcy(msg.Data)
		Nodes[nodeId].EmcyCount++
	} else {
		fmt.Printf("Unknown msg with func=%d\n", funcCode)
	}
}

func (node *Node) Reset() {
	node.Bus.Write(can.NewStdMsg(0, []uint8{129,uint8(node.Id)}))
}

func (node *Node) SetOperational() {
	node.Bus.Write(can.NewStdMsg(0, []uint8{1,uint8(node.Id)}))
}

func (node *Node) SetPreOperational() {
	node.Bus.Write(can.NewStdMsg(0, []uint8{128,uint8(node.Id)}))
}

func (node *Node) HandleEmcy(msg [8]uint8) {
	node.LastEmcyMsg = msg
}

func (node *Node) HandlePdo(no int, data [8]uint8) {
	fmt.Printf("PDO %d data=%d, %d, %d, %d, %d, %d, %d %d\n", no, data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7])
}

func NewMuxMsg(base int, nodeId NodeId, op uint8, index Index, subIndex SubIndex, Value int) can.Msg {
	data := [8]uint8{op, uint8(index&0xFF), uint8(index>>8), uint8(subIndex), uint8(Value&0xFF) ,uint8((Value>>8)&0xFF),
		uint8((Value>>16)&0xFF) ,uint8((Value>>24)&0xFF)}
	return can.Msg{Id: DefaultId(base,nodeId), Type: can.Standard, Len: 8, Data: data}
}

func (node *Node) ReadObject(index Index, subIndex SubIndex, byteCount uint8) (int, error) {
	var sdoReadOpcode = [5]uint8 {0x40, 0x4F, 0x4B, 0x47, 0x43}
	var mask = [5]int{0, 0xFF, 0xFFFF, 0xFFFFFF, 0xFFFFFFFF}
	if byteCount<1 || byteCount>4 {
		return 0, fmt.Errorf("Byte count was %d, must be 1..4", byteCount)
	}
	msg := NewMuxMsg(0x600, node.Id, 0x40 /*sdoReadOpcode[byteCount]*/, index, subIndex, 0)

	resp := node.Bus.Poll(msg, DefaultId(0x580, node.Id))
	if resp==nil {
		return 0, fmt.Errorf("no response")
	}
	if resp.Data[0]!=sdoReadOpcode[byteCount] {
		return 0, fmt.Errorf("Read Object size mismatch for %x:%d, got opcode=%x, expected %x",index,subIndex, resp.Data[0], sdoReadOpcode[byteCount])
	}
	return mask[byteCount]&(int(resp.Data[4])+int(resp.Data[5])*256 + int(resp.Data[6])*256*256+int(resp.Data[7])*256*256*256) , nil
}

func (node *Node) WriteObject(index Index, subIndex SubIndex, byteCount uint8, value int) error {
	var sdoWriteOpcode = [5]uint8 {0x23, 0x2F, 0x2B, 0x27, 0x23}
	if byteCount<1 || byteCount>4 {
		return fmt.Errorf("byte count was %d, must be 1..4", byteCount)
	}
	msg := NewMuxMsg(0x600, node.Id, sdoWriteOpcode[byteCount], index, subIndex, value)
	resp := node.Bus.Poll(msg, DefaultId(0x580, node.Id))
	if resp==nil {
		return fmt.Errorf("sdo write to node %d timed out", node.Id)
	}
	if resp.Data[0]==0x80 {
		return fmt.Errorf("write aborted by node, %s", resp.ToString() )
	}
	return nil
}

func (node *Node) VerifyReadAbort(index Index, subIndex SubIndex, byteCount uint8, description string) {
	msg := NewMuxMsg(0x600, node.Id, 0x40, index, subIndex, 0)
	resp := node.Bus.Poll(msg, DefaultId(0x580, node.Id))
	if resp==nil {
		node.Failed = true
		fmt.Printf("Polling returned no data for %x:%d (%s)", index, subIndex, description)
		return
	}
	if resp.Data[0]!=0x80 {
		node.Failed = true
		fmt.Printf("Expected abort code, got data for %x:%d (%s)", index, subIndex, description)
	}
}

func (node *Node) VerifyEqual(index Index, subIndex SubIndex, byteCount uint8, expected int, description string) {
	value, err := node.ReadObject(index, subIndex, byteCount)
	if err!=nil {
		node.Failed = true
		fmt.Printf("Error reading Object %x:%d (%s), error %s\n", index, subIndex, description, err)
	}
	if  value!=expected {
		node.Failed = true
		fmt.Printf("Expected %x, Actual %x, %s\n", expected, value, description)
	}
}

func (node *Node) VerifyRange(index Index, subIndex SubIndex, byteCount uint8, min int, max int, description string) {
	value, err := node.ReadObject(index, subIndex, byteCount)
	if err!=nil {
		node.Failed = true
		fmt.Printf("Error reading Object %x:%d (%s)\n", index, subIndex, description)
	}
	if  value<min || value>max {
		node.Failed = true
		fmt.Printf("Expected %x..%x, Actual %x, %s\n", min, max, value, description)
	}
}

func (node *Node) Verify(cond bool, msg string, par...interface{}) {
	if !cond {
		node.Failed = true
		fmt.Printf(msg+"\n", par)
	}
}

func New(con *bus.State, id NodeId) *Node {
	n:= &Node{Bus: con, Id: id}
	Nodes[id] = n
	return n
}

func init() {
	bus.Handler = MsgHandler
}