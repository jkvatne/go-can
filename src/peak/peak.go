package peak

import (
	"can"
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

type peakChannel struct {
	channel uint8
	btr     uint16
	status  can.Status
	callback can.Callback
	terminated bool
}

const (
	PCAN_BAUD_1M             = uint16(0x0014)  // 1 MBit/s
	PCAN_BAUD_800K           = uint16(0x0016)  // 800 kBit/s
	PCAN_BAUD_500K           = uint16(0x001C)  // 500 kBit/s
	PCAN_BAUD_250K           = uint16(0x011C)  // 250 kBit/s
	PCAN_BAUD_125K           = uint16(0x031C)  // 125 kBit/s
	PCAN_BAUD_100K           = uint16(0x432F)  // 100 kBit/s
	PCAN_BAUD_95K            = uint16(0xC34E)  // 95,238 kBit/s
	PCAN_BAUD_83K            = uint16(0x852B)  // 83,333 kBit/s
	PCAN_BAUD_50K            = uint16(0x472F)  // 50 kBit/s
	PCAN_BAUD_47K            = uint16(0x1414)  // 47,619 kBit/s
	PCAN_BAUD_33K            = uint16(0x8B2F)  // 33,333 kBit/s
	PCAN_BAUD_20K            = uint16(0x532F)  // 20 kBit/s
	PCAN_BAUD_10K            = uint16(0x672F)  // 10 kBit/s
	PCAN_BAUD_5K             = uint16(0x7F7F)  // 5 kBit/s
	
	PCAN_TYPE_ISA            = uint8(0x01)  // PCAN-ISA 82C200
	PCAN_TYPE_ISA_SJA        = uint8(0x09)  // PCAN-ISA SJA1000
	PCAN_TYPE_ISA_PHYTEC     = uint8(0x04)  // PHYTEC ISA
	PCAN_TYPE_DNG            = uint8(0x02)  // PCAN-Dongle 82C200
	PCAN_TYPE_DNG_EPP        = uint8(0x03)  // PCAN-Dongle EPP 82C200
	PCAN_TYPE_DNG_SJA        = uint8(0x05)  // PCAN-Dongle SJA1000
	PCAN_TYPE_DNG_SJA_EPP    = uint8(0x06)  // PCAN-Dongle EPP SJA1000

	CAN_MAX_STANDARD_ID 	= 0x7ff
	CAN_MAX_EXTENDED_ID 	= 0x1fffffff
	PCAN_USBBUS1            = 0x51          // PCAN-USB interface, channel 1

	PCAN_MESSAGE_STANDARD    = uint8(0x00)  // The PCAN message is a CAN Standard Frame (11-bit identifier)
	PCAN_MESSAGE_RTR         = uint8(0x01)  // The PCAN message is a CAN Remote-Transfer-Request Frame
	PCAN_MESSAGE_EXTENDED    = uint8(0x02)  // The PCAN message is a CAN Extended Frame (29-bit identifier)
	PCAN_MESSAGE_STATUS      = uint8(0x80)  // The PCAN message represents a PCAN status message

	PCAN_ERROR_OK            = uint32(0x00000)  // No error
	PCAN_ERROR_XMTFULL       = uint32(0x00001)  // Transmit buffer in CAN controller is full
	PCAN_ERROR_OVERRUN       = uint32(0x00002)  // CAN controller was read too late
	PCAN_ERROR_BUSLIGHT      = uint32(0x00004)  // Bus error: an error counter reached the 'light' limit
	PCAN_ERROR_BUSHEAVY      = uint32(0x00008)  // Bus error: an error counter reached the 'heavy' limit
	PCAN_ERROR_BUSOFF        = uint32(0x00010)  // Bus error: the CAN controller is in bus-off state
	PCAN_ERROR_ANYBUSERR     = uint32(PCAN_ERROR_BUSLIGHT | PCAN_ERROR_BUSHEAVY | PCAN_ERROR_BUSOFF)  // Mask for all bus errors
	PCAN_ERROR_QRCVEMPTY     = uint32(0x00020)  // Receive queue is empty
	PCAN_ERROR_QOVERRUN      = uint32(0x00040)  // Receive queue was read too late
	PCAN_ERROR_QXMTFULL      = uint32(0x00080)  // Transmit queue is full
	PCAN_ERROR_REGTEST       = uint32(0x00100)  // Test of the CAN controller hardware registers failed (no hardware found)
	PCAN_ERROR_NODRIVER      = uint32(0x00200)  // Driver not loaded
	PCAN_ERROR_HWINUSE       = uint32(0x00400)  // Hardware already in use by a Net
	PCAN_ERROR_NETINUSE      = uint32(0x00800)  // A Client is already connected to the Net
	PCAN_ERROR_ILLHW         = uint32(0x01400)  // Hardware handle is invalid
	PCAN_ERROR_ILLNET        = uint32(0x01800)  // Net handle is invalid
	PCAN_ERROR_ILLCLIENT     = uint32(0x01C00)  // Client handle is invalid
	PCAN_ERROR_ILLHANDLE     = uint32(PCAN_ERROR_ILLHW | PCAN_ERROR_ILLNET | PCAN_ERROR_ILLCLIENT)  // Mask for all handle errors
	PCAN_ERROR_RESOURCE      = uint32(0x02000)  // Resource (FIFO, Client, timeout) cannot be created
	PCAN_ERROR_ILLPARAMTYPE  = uint32(0x04000)  // Invalid parameter
	PCAN_ERROR_ILLPARAMVAL   = uint32(0x08000)  // Invalid parameter value
	PCAN_ERROR_UNKNOWN       = uint32(0x10000)  // Unknow error
	PCAN_ERROR_ILLDATA       = uint32(0x20000)  // Invalid data, function, or action
	PCAN_ERROR_INITIALIZE    = uint32(0x40000)  // peakChannel is not initialized
	PCAN_ERROR_ILLOPERATION  = uint32(0x80000)  // Invalid operation
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}


var (
	peak, _ = syscall.LoadLibrary("pcanbasic.dll")
	CAN_Initialize, _ = syscall.GetProcAddress(peak, "CAN_Initialize")
	CAN_Uninitialize, _ = syscall.GetProcAddress(peak, "CAN_Uninitialize")
	CAN_Reset, _ = syscall.GetProcAddress(peak, "CAN_Reset")
	CAN_GetStatus, _ = syscall.GetProcAddress(peak, "CAN_GetStatus")
	CAN_Read, _ = syscall.GetProcAddress(peak, "CAN_Read")
	CAN_Write, _ = syscall.GetProcAddress(peak, "CAN_Write")
	CAN_FilterMessages, _ = syscall.GetProcAddress(peak, "CAN_FilterMessages")
	CAN_GetValue, _ = syscall.GetProcAddress(peak, "CAN_GetValue")
	CAN_SetValue, _ = syscall.GetProcAddress(peak, "CAN_SetValue	")
	CAN_GetErrorText, _ = syscall.GetProcAddress(peak, "CAN_GetErrorText")
)


type CanTimeStamp struct {
	Millisec         uint32 // [("millis", c_ulong),           // Base-value: milliseconds: 0.. 2^32-1
	MillisecOverflow uint16 // ("millis_overflow", c_ushort),  // Roll-arounds of millis
	MicroSec         uint16 //("micros", c_ushort)]            // Microseconds: 0..999
}

func (p *peakChannel) Initialize(Btr0Btr1 uint16) error {
	ret, _, callErr := syscall.Syscall6(uintptr(CAN_Initialize), 5, uintptr(p.channel), uintptr(Btr0Btr1), uintptr(0), uintptr(0), uintptr(0), 0)
	if ret!=0 {
		return fmt.Errorf("Initialize error %d", ret)
	}
	if callErr!=0 {
		return fmt.Errorf("Initialize error %d", callErr)
	}
	return nil
}

func (p *peakChannel) Uninitialize() {
	p.terminated = true
	_, _, _ = syscall.Syscall6(uintptr(CAN_Uninitialize), 1, uintptr(p.channel),0,0,0,0,0)
	return
}

func (p *peakChannel) convertError(ret uintptr) {
	if uint32(ret) & PCAN_ERROR_BUSOFF != 0 {
		p.status = can.BusOff
	} else if uint32(ret) & PCAN_ERROR_BUSHEAVY != 0 {
		p.status = can.ErrorActive
	} else if uint32(ret) & PCAN_ERROR_BUSLIGHT != 0 {
		p.status = can.ErrorPassive
	} else if uint32(ret) != PCAN_ERROR_QRCVEMPTY {
		p.status = can.Error
	} else {
		p.status = can.Ok
	}
}

func (p *peakChannel) Status() can.Status {
	ret, _, _ := syscall.Syscall6(uintptr(CAN_GetStatus), 1, uintptr(p.channel),0,0,0,0,0)
	p.convertError(ret)
	return p.status
}

func  (p *peakChannel)ErrorText(err uint32) string {
	var ErrorName = [...]string{"Ok", "ErrorPassive", "ErrorActive", "BusOff","InternalError"}
	return ErrorName[p.status]
}

func (p *peakChannel) Read() (*can.Msg, error) {
	msg := can.Msg{}
	timestamp := CanTimeStamp{}
	ret, _, callErr := syscall.Syscall6(uintptr(CAN_Read), 3, uintptr(p.channel), uintptr(unsafe.Pointer(&msg)), uintptr(unsafe.Pointer(&timestamp)),0,0,0)
	p.convertError(ret)
	if ret!=0 {
		return nil, fmt.Errorf("Read error : %s", p.ErrorText(uint32(ret)))
	}
	if callErr!=0 {
		return nil, fmt.Errorf("Initialize error %d", callErr)
	}
	return &msg, nil
}

func (p *peakChannel) Write(msg can.Msg) error {
	ret, _, callErr := syscall.Syscall6(uintptr(CAN_Write), 2, uintptr(p.channel), uintptr(unsafe.Pointer(&msg)),0,0,0,0)
	if ret!=0 {
		return fmt.Errorf("Write error : %s", p.ErrorText(uint32(ret)))
	}
	if callErr!=0 {
		return fmt.Errorf("Initialize error %d", callErr)
	}
	return nil
}

func New(channel uint8, Btr0Btr1 uint16, callback can.Callback) (*peakChannel, error) {
	bus := &peakChannel{channel: channel, btr: Btr0Btr1, callback:callback}
	err := bus.Initialize(Btr0Btr1)
	go Handler(bus,callback)
	return bus, err
}

// Handler is a go routine that will check for incoming messages and call a callback for each
func Handler(p *peakChannel, c can.Callback) {
	for (!p.terminated) {
		m, err := p.Read()
		if err==nil {
			c(m)
		} else {
			time.Sleep(50*time.Millisecond)
		}
	}
}