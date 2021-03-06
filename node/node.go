package node

import (
	"encoding/binary"
	"fmt"
	"github.com/gookit/color"
	"go-can/bus"
	"go-can/can"
	"math"
)

type NodeId uint8
type Index uint16
type SubIndex uint8
type Callback func(msg *can.Msg)


const (
	PdoCount = 4
	CanMessageSize = 8
    MaxNodeId = 127
)

type Node struct {
	Bus            *bus.State
	Id             NodeId
	HeartbeatCount int
	State          int
	Failed         bool
	LastEmcyMsg    [8]byte
	EmcyCount      int
	RxPdosReceived [PdoCount]int
	rxPdo          [PdoCount][CanMessageSize]byte  // Data received from node
	txPdo          [PdoCount][CanMessageSize]byte  // Data to send to node
	testNo         int
}

var Nodes[MaxNodeId] *Node
var SubTest int

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
	FUNC_HEARTBEAT = 0x0E
	FUNC_MGMT   = 0x0F
	SYNC_COB_ID = 0x80
)

func (node *Node) SkipTest(description string) bool {
	node.testNo++
	if SubTest>0 && SubTest!=node.testNo {
		fmt.Printf("%d : Skipping test \"%s\"\n", node.testNo, description)
		return true
	}
	fmt.Printf("%d : %s\n", node.testNo, description)
	return false
}

// SetPdoValue will set a value in the node's txPdo[] data structure.
// The offset and count defines where in the 8-byte array the data is stored
// Actual transmission of data to the remote node is done with node.SendPdo(n)
func (node *Node) SetPdoValue(pdoNo int, offset int, count int, value int) {
	// Verify that the given pdoNo is valid (1..4)
	CheckPdoNo(pdoNo)
	for i:=offset; i<count+offset; i++ {
		node.txPdo[pdoNo-1][i] = uint8(value & 0xFF)
		value = value >> 8
	}
}

// SendPdo will send a number of messages to the node.
func (node *Node) SendPdo(pdoNo int) {
	// Verify that the given pdoNo is valid (1..4)
	CheckPdoNo(pdoNo)
	// Make the message from data in node.txPdo[]
	msg := can.Msg{Id: can.CobId((pdoNo+1)*0x100+int(node.Id)), Len: 8}
	for i:=0; i<CanMessageSize; i++ {
		msg.Data[i] = node.txPdo[pdoNo-1][i]
	}
	// Send message
	node.Bus.Write(msg)
}

func (node *Node) GetPdoInt16(pdoNo int, ofs int) int {
	CheckPdoNo(pdoNo)
	var value int
	value = int(node.rxPdo[pdoNo-1][ofs])
	value += int(node.rxPdo[pdoNo-1][ofs+1])*256
	return value
}

func CheckPdoNo(pdoNo int) {
	if pdoNo<1 || pdoNo>4 {
		panic("Pdo number error")
	}
}

func (node *Node) VerifyPdoInt16(pdoNo int, ofs, min int, max int, msg string) {
	CheckPdoNo(pdoNo)
	value := node.GetPdoInt16(pdoNo, ofs)
	if value<min || value>max {
		node.Failed = true
		color.Error.Printf("Pdo %d value at ofs=%d was %d, should be %d..%d, %s\n", pdoNo, ofs, value, min, max, msg)
	}
}

func (node *Node) VerifyPdoCount(n1 int, n2 int, n3 int, n4 int) {
	if n1!= node.RxPdosReceived[0] || n2!= node.RxPdosReceived[1] || n3!= node.RxPdosReceived[2] || n4!= node.RxPdosReceived[3] {
		node.Failed = true
		color.Error.Printf("Expected %d, %d, %d, %d pdos, got %d, %d, %d, %d\n", n1,n2,n3,n4,
			node.RxPdosReceived[0], node.RxPdosReceived[1], node.RxPdosReceived[2], node.RxPdosReceived[3])
	}
}

func DefaultId(base int, nodeId NodeId) can.CobId {
	return can.CobId(base)+can.CobId(nodeId)
}

func MsgHandler(msg *can.Msg) {
	funcCode := (int(msg.Id) & 0x780) >> 7
	nodeId := int(msg.Id) & 0x07F
	if Nodes[nodeId]==nil {
		return
	}
	if msg.Id == SYNC_COB_ID {
		// Got sync message from another node
	} else if funcCode==FUNC_HEARTBEAT {
		// Incoming heartbeat messages
		if Nodes[nodeId]!=nil {
			Nodes[nodeId].HeartbeatCount++
			Nodes[nodeId].State = int(msg.Data[0])
		}
	} else if funcCode >= FUNC_TXPDO1 && funcCode <= FUNC_RXPDO4 {
		if funcCode==FUNC_TXPDO1 {
			Nodes[nodeId].HandlePdo(1, msg.Data)
		} else if funcCode==FUNC_TXPDO2 {
			Nodes[nodeId].HandlePdo(2, msg.Data)
		} else if funcCode==FUNC_TXPDO3 {
			Nodes[nodeId].HandlePdo(3, msg.Data)
		} else if funcCode==FUNC_TXPDO4 {
			Nodes[nodeId].HandlePdo(4, msg.Data)
		}

	} else if funcCode == FUNC_EMCY {
		Nodes[nodeId].HandleEmcy(msg.Data)
		Nodes[nodeId].EmcyCount++

	} else {
		//fmt.Printf("Unknown msg with func=%d\n", funcCode)
	}
}

func (node *Node) ResetNode() {
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

func (node *Node) HandlePdo(pdoNo int, data [8]uint8) {
	CheckPdoNo(pdoNo)
	for i:=0; i<8; i++ {
		node.rxPdo[pdoNo-1][i] = data[i]
	}
	node.RxPdosReceived[pdoNo-1]++
}

func (node *Node) ResetPdoCount() {
	for i:=0; i<4; i++ {
		node.RxPdosReceived[i] = 0
	}
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
	msg := NewMuxMsg(0x600, node.Id, 0x40, index, subIndex, 0)
	resp := node.Bus.Poll(msg, DefaultId(0x580, node.Id))
	if resp==nil {
		return 0, fmt.Errorf("no response")
	}
	if resp.Data[0]!=sdoReadOpcode[byteCount] {
		n := (0x53-resp.Data[0])>>2
		fmt.Printf("Warning: Read object size mismatch for 0x%x:%d, expected %d bytes, got %d\n", index, subIndex,byteCount, n)
		// Clear bytes not used - using received length info
		for i:=4+n; i<8; i++ {
			resp.Data[i] = 0
		}
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
		color.Error.Printf("Polling returned no data for %x:%d (%s)\n", index, subIndex, description)
		return
	}
	if resp.Data[0]!=0x80 {
		node.Failed = true
		color.Error.Printf("Expected abort code, got data for %x:%d (%s)\n", index, subIndex, description)
	}
}

func (node *Node) VerifyEqual(index Index, subIndex SubIndex, byteCount uint8, expected int, description string) {
	value, err := node.ReadObject(index, subIndex, byteCount)
	if err!=nil {
		node.Failed = true
		color.Error.Printf("Error reading Object %x:%d (%s), %s\n", index, subIndex, description, err)
	}
	if  value!=expected {
		node.Failed = true
		color.Error.Printf("Expected 0x%x, Actual 0x%x, Object 0x%x:%d (%s)\n", expected, value, index, subIndex, description)
	}
}

func (node *Node) VerifyRange(index Index, subIndex SubIndex, byteCount uint8, min int, max int, description string) {
	value, err := node.ReadObject(index, subIndex, byteCount)
	if err!=nil {
		node.Failed = true
		color.Error.Printf("Error reading Object %x:%d (%s)\n", index, subIndex, description)
	}
	if  value<min || value>max {
		node.Failed = true
		color.Error.Printf("Expected %d..%d, Actual %d(0x%x), %s\n", min, max, value, value, description)
	}
}

func (node *Node) Check(err error) {
	if err!=nil {
		node.Failed = true
		color.Error.Printf("Error %s", err)
	}
}

func Float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func (node *Node) ReadFloat(index Index, subIndex SubIndex) float64 {
	value, err := node.ReadObject(index, subIndex, 4)
	if err!=nil {
		node.Failed = true
		color.Error.Printf("Error reading float value in object %x:%d %s\n", index, subIndex, err)
		return 0.0
	}
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint32(bs, uint32(value))
	return float64(Float32frombytes(bs))
}

func (node *Node) VerifyRangeFloat(index Index, subIndex SubIndex, min float64, max float64, description string) {
	floatValue := node.ReadFloat(index, subIndex)
	if floatValue<min || floatValue>max {
		node.Failed = true
		color.Error.Printf("Error reading float value in %x:%d, expected %0.3f..%0.3f, was %0.3f, %s\n",
			index, subIndex, min, max,floatValue, description)
	}
}

func (node *Node) Verify(cond bool, msg string, par...interface{}) {
	if !cond {
		node.Failed = true
		if par==nil {
			color.Error.Printf(msg+"\n")
		} else {
			color.Error.Printf(msg+"\n", par)
		}
	}
}

func New(con *bus.State, id int) *Node {
	n:= &Node{Bus: con, Id: NodeId(id)}
	Nodes[id] = n
	return n
}

func init() {
	bus.Handler = MsgHandler
}