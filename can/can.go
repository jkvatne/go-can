package can

import "fmt"

type Id uint32
type MsgType uint8
type Status uint32
type Callback func(msg *Msg)

type Msg struct {
	Id Id           // 11/29-bit message identifier
	Type uint8      // MsgType of the message
	Len uint8       // Data Length of the message (0..8)
	Data [8]uint8   // Data of the message, 8 bytes
}

type Device interface {
	Read() *Msg
	Write(Msg)
	Status() Status
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


