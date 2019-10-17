package can

type Msg struct {
	Id uint32  // _fields_ = [("ID", c_ulong),                // 11/29-bit message identifier
	Type uint8 //("MSGTYPE", TPCANMessageType),  // Type of the message
	Len uint8  // ("LEN", c_ubyte),              // Data Length Code of the message (0..8)
	Data [8]uint8 //("DATA", c_ubyte * 8)]       // Data of the message (DATA[0]..DATA[7])
}

type Status uint32
type Callback func(msg *Msg)

const (
	Ok Status = iota
	ErrorPassive
	ErrorActive
	BusOff
	Error
)

type Device interface {
	Read() (*Msg, error)
	Write(Msg) (error)
	Status() Status
	HasData() bool
}

