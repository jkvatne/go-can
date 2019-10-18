package can

import (
	"fmt"
	"sync"
	"time"
)

type Id uint32
type MsgType uint8
type Status uint32
type Callback func(msg *Msg)

type Msg struct {
	Id Id           // 11/29-bit message identifier
	Type MsgType    // MsgType of the message
	Len uint8       // Data Length of the message (0..8)
	Data [8]uint8   // Data of the message, 8 bytes
	Time time.Time  // Time of message (on received messages only)
}

type Device interface {
	Read() *Msg
	Write(Msg)
	Status() Status
	Uninitialize()
	Reset()
}

const (
	MESSAGE_STANDARD    = MsgType(0x00) // The PCAN message is a CAN Standard Frame (11-bit identifier)
	MESSAGE_RTR         = MsgType(0x01) // The PCAN message is a CAN Remote-Transfer-Request Frame
	MESSAGE_EXTENDED    = MsgType(0x02) // The PCAN message is a CAN Extended Frame (29-bit identifier)
	MESSAGE_STATUS      = MsgType(0x80) // The PCAN message represents a PCAN status message
)

const (
	Ok Status = iota
	ErrorPassive
	ErrorActive
	BusOff
	Error
)

type connection struct {
	Dev           Device
	mutex         sync.Mutex
	timeout       time.Duration
	terminated    bool
	responseChan  chan *Msg
	requestChan   chan Id
	sleepInterval time.Duration
	pollTime      time.Time
}

var errorNames = [...]string{"Ok", "ErrorPassive", "ErrorActive", "BusOff", "Error"}

func StatusToString(sts Status) string {
	if sts<0 || sts > Error {
		return "Unknown error"
	}
	return errorNames[sts]
}

func (m *Msg) ToString() string {
	return fmt.Sprintf("ID %x, Data %d,%d,%d,%d,%d,%d,%d,%d", m.Id, m.Data[0], m.Data[1], m.Data[2], m.Data[3],
		m.Data[4], m.Data[5], m.Data[6], m.Data[7])
}

func (c connection) Poll(msg Msg, responseId Id) *Msg {
	// Polling on a terminated channel is not allowed
	if c.terminated {
		return nil
	}
	// Allow only one thread to poll at the same time
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.requestChan<-responseId
	c.Dev.Write(msg)
	//_, _, _ = syscall.Syscall6(uintptr(CAN_Write), 2, uintptr(c.channel), uintptr(unsafe.Pointer(&msg)),0,0,0,0)
	resp :=<-c.responseChan
	return resp
}


// Handler is a goroutine that will check for incoming messages and call a callback for each
func Handler(c connection, callback Callback) {
	var responseId Id = 0xFFFFFFFF
	for !c.terminated {
		select {
		case responseId = <-c.requestChan:
			c.pollTime = time.Now()
		default:
		}

		m := c.Dev.Read()
		if m!=nil {
			if m.Id == responseId {
				responseId = 0xFFFFFFFF
				c.responseChan <- m
			} else {
				callback(m)
			}
		} else {
			// No messages found, sleep for some time
			time.Sleep(c.sleepInterval * time.Millisecond)
		}
		// Check for poll timeout
		if  responseId!=0xFFFFFFFF && time.Since(c.pollTime)>c.timeout {
			responseId = 0xFFFFFFFF
			c.responseChan <- nil
		}
	}
}

func (c connection) SdoRead(nodeId uint16, index uint16, subIndex uint8, byteCount int) (int, error) {
	var sdoReadOpcode = [5]uint8 {0x40, 0x4C, 0x48, 0x44, 0x40}
	if byteCount<0 || byteCount>4 {
		return 0, fmt.Errorf("Byte count ws %d, must be 1..4", byteCount)
	}
	msg := Msg{
		Id: Id(0x600+nodeId),
		Type: MESSAGE_STANDARD,
		Len:8,
		Data:[8]uint8{sdoReadOpcode[byteCount],uint8(index&0xFF),uint8(index>>8),subIndex,0,0,0,0} }
	resp := c.Poll(msg, Id(0x580+nodeId))
	if resp==nil {
		return 0, fmt.Errorf("No response")
	}
	return int(resp.Data[4])+int(resp.Data[5])*256 + int(resp.Data[6])*256*256+int(resp.Data[7])*256*256*256 , nil
}

func (c connection) SdoWrite(nodeId uint16, index uint16, subIndex uint8, byteCount int, value int) error {
	var sdoReadOpcode = [5]uint8 {0x40, 0x4C, 0x48, 0x44, 0x40}
	if byteCount<0 || byteCount>4 {
		return fmt.Errorf("byte count was %d, must be 1..4", byteCount)
	}
	msg := Msg{
		Id: Id(0x600+nodeId),
		Type: MESSAGE_STANDARD,
		Len:8,
		Data:[8]uint8{sdoReadOpcode[byteCount],uint8(index&0xFF),uint8(index>>8),subIndex,
			uint8(value&0xFF),uint8((value>>8)&0xFF),uint8((value>>16)&0xFF),uint8((value>>24)&0xFF)} }
	resp := c.Poll(msg, Id(0x580+nodeId))
	if resp==nil {
		return fmt.Errorf("sdo write to node %d timed out", nodeId)
	}
	if resp.Data[0]==0x80 {
		return fmt.Errorf("write aborted by node, %s", resp.ToString() )
	}
	return nil
}

func (c connection) Close() {
	c.terminated = true
	// Wait for go routine to terminate
	time.Sleep(c.sleepInterval)
	c.Dev.Uninitialize()
}

func NewConnection(d Device, timeout time.Duration, callback Callback) *connection {
	c := connection{Dev: d, timeout: timeout}
	c.mutex = sync.Mutex{}
	c.requestChan = make(chan Id, 1)
	c.responseChan = make(chan *Msg, 1)
	go Handler(c, callback)
	return &c
}
