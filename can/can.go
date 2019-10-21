package can

import (
	"fmt"
	"time"
)

type CobId uint32
type MsgType uint8
type Status uint32
type Callback func(msg *Msg)

type Msg struct {
	Id CobId           // 11/29-bit message identifier
	Type MsgType    // MsgType of the message
	Len uint8       // Data Length of the message (0..8)
	Data [8]uint8   // Data of the message, 8 bytes
	Time time.Time  // Time of message (on received messages only)
}

type Device interface {
	Read() *Msg
	Write(Msg)
	Status() Status
	Close()
	Initialize(speed int) error
	Reset()
}

const (
	Standard  = MsgType(0x00) // Message is a CAN Standard Frame (11-bit identifier)
	Rtr       = MsgType(0x01) // Message is a CAN Remote-Transfer-Request Frame
	Extended  = MsgType(0x02) // Message is a CAN Extended Frame (29-bit identifier)
	StatusMsg = MsgType(0x80) // Message represents a status message
)

const (
	Ok Status = iota
	ErrorPassive
	ErrorActive
	BusOff
	Error
)

var errorNames = [...]string{"Ok", "ErrorPassive", "ErrorActive", "BusOff", "Error"}

func StatusToString(sts Status) string {
	if sts<0 || sts > Error {
		return "Unknown error"
	}
	return errorNames[sts]
}

func (m *Msg) ToString() string {
	return fmt.Sprintf("CobID %x, Data %d,%d,%d,%d,%d,%d,%d,%d", m.Id, m.Data[0], m.Data[1], m.Data[2], m.Data[3],
		m.Data[4], m.Data[5], m.Data[6], m.Data[7])
}


func NewStdMsg(id CobId, data []uint8) Msg {
	var d [8]uint8
	for i:=0; i<len(data); i++ {
		d[i]=data[i]
	}
	return Msg{Id: id, Type: Standard, Len: uint8(len(data)), Data: d}
}

