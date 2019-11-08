package peak

import (
	"fmt"
	"go-can/can"
	"os"
	"syscall"
	"unsafe"
)

type peakChannel struct {
	channel uint8
	bitrate int
	status  can.Status
}

// Check if peak interface satisfies can.Device interface
var _ can.Device = &peakChannel{}

const (
	PCAN_USBBUS1            = 0x51          // PCAN-USB interface, channel 1
	PCAN_USBBUS2            = 0x52          // PCAN-USB interface, channel 1
	PCAN_USBBUS3            = 0x53          // PCAN-USB interface, channel 1
	PCAN_USBBUS4            = 0x54          // PCAN-USB interface, channel 1
)

const (
	pPCAN_BAUD_1M             = uint16(0x0014)  // 1 MBit/s
	pPCAN_BAUD_800K           = uint16(0x0016)  // 800 kBit/s
	pPCAN_BAUD_500K           = uint16(0x001C)  // 500 kBit/s
	pPCAN_BAUD_250K           = uint16(0x011C)  // 250 kBit/s
	pPCAN_BAUD_125K           = uint16(0x031C)  // 125 kBit/s
	pPCAN_BAUD_100K           = uint16(0x432F)  // 100 kBit/s
	pPCAN_BAUD_50K            = uint16(0x472F)  // 50 kBit/s
	pPCAN_BAUD_20K            = uint16(0x532F)  // 20 kBit/s
	pPCAN_BAUD_10K            = uint16(0x672F)  // 10 kBit/s
	pPCAN_MESSAGE_STANDARD    = uint8(0x00)  // The PCAN message is a CAN Standard Frame (11-bit identifier)
	pPCAN_MESSAGE_RTR         = uint8(0x01)  // The PCAN message is a CAN Remote-Transfer-Request Frame
	pPCAN_MESSAGE_EXTENDED    = uint8(0x02)  // The PCAN message is a CAN Extended Frame (29-bit identifier)
	pPCAN_MESSAGE_STATUS      = uint8(0x80)  // The PCAN message represents a PCAN status message
	pPCAN_ERROR_OK            = uintptr(0x00000)  // No error
	pPCAN_ERROR_XMTFULL       = uintptr(0x00001)  // Transmit buffer in CAN controller is full
	pPCAN_ERROR_OVERRUN       = uintptr(0x00002)  // CAN controller was read too late
	pPCAN_ERROR_BUSLIGHT      = uintptr(0x00004)  // Bus error: an error counter reached the 'light' limit
	pPCAN_ERROR_BUSHEAVY      = uintptr(0x00008)  // Bus error: an error counter reached the 'heavy' limit
	pPCAN_ERROR_BUSOFF        = uintptr(0x00010)  // Bus error: the CAN controller is in bus-off state
	pPCAN_ERROR_QRCVEMPTY     = uintptr(0x00020)  // Receive queue is empty
	pPCAN_ERROR_QOVERRUN      = uintptr(0x00040)  // Receive queue was read too late
	pPCAN_ERROR_QXMTFULL      = uintptr(0x00080)  // Transmit queue is full
	pPCAN_ERROR_REGTEST       = uintptr(0x00100)  // Test of the CAN controller hardware registers failed (no hardware found)
	pPCAN_ERROR_NODRIVER      = uintptr(0x00200)  // Driver not loaded
	pPCAN_ERROR_HWINUSE       = uintptr(0x00400)  // Hardware already in use by a Net
	pPCAN_ERROR_NETINUSE      = uintptr(0x00800)  // A Client is already connected to the Net
	pPCAN_ERROR_ILLHW         = uintptr(0x01400)  // Hardware handle is invalid
	pPCAN_ERROR_ILLNET        = uintptr(0x01800)  // Net handle is invalid
	pPCAN_ERROR_ILLCLIENT     = uintptr(0x01C00)  // Client handle is invalid
	pPCAN_ERROR_RESOURCE      = uintptr(0x02000)  // Resource (FIFO, Client, timeout) cannot be created
	pPCAN_ERROR_ILLPARAMTYPE  = uintptr(0x04000)  // Invalid parameter
	pPCAN_ERROR_ILLPARAMVAL   = uintptr(0x08000)  // Invalid parameter value
	pPCAN_ERROR_UNKNOWN       = uintptr(0x10000)  // Unknow error
	pPCAN_ERROR_ILLDATA       = uintptr(0x20000)  // Invalid data, function, or action
	pPCAN_ERROR_INITIALIZE    = uintptr(0x40000)  // peakChannel is not initialized
	pPCAN_ERROR_ILLOPERATION  = uintptr(0x80000)  // Invalid operation

	PCAN_BUSOFF_AUTORESET  = 0x07
)

var (
	peak, peakError = syscall.LoadLibrary("pcanbasic.dll")
	CAN_Initialize, _ = syscall.GetProcAddress(peak, "CAN_Initialize")
	CAN_Uninitialize, _ = syscall.GetProcAddress(peak, "CAN_Uninitialize")
	CAN_Reset, _ = syscall.GetProcAddress(peak, "CAN_Reset")
	CAN_GetStatus, _ = syscall.GetProcAddress(peak, "CAN_GetStatus")
	CAN_Read, _ = syscall.GetProcAddress(peak, "CAN_Read")
	CAN_Write, _ = syscall.GetProcAddress(peak, "CAN_Write")
	CAN_FilterMessages, _ = syscall.GetProcAddress(peak, "CAN_FilterMessages")
	CAN_GetValue, _ = syscall.GetProcAddress(peak, "CAN_GetValue")
	CAN_SetValue, _ = syscall.GetProcAddress(peak, "CAN_SetValue")
	CAN_GetErrorText, _ = syscall.GetProcAddress(peak, "CAN_GetErrorText")
)

// Total Microseconds = micros + 1000 * millis + 0x100000000 * 1000 * millis_overflow
type canTimeStamp struct {
	Millisec         uint32 // [("millis", c_ulong),           // Base-value: milliseconds: 0.. 2^32-1
	MillisecOverflow uint16 // ("millis_overflow", c_ushort),  // Roll-arounds of millis
	MicroSec         uint16 //("micros", c_ushort)]            // Microseconds: 0..999
}


func convertError(e uintptr) can.Status {
	if e & 0xFFFFFFDF == 0 {
		return can.Ok
	} else if (e & pPCAN_ERROR_BUSOFF)!=0 {
		return can.BusOff
	} else if (e & pPCAN_ERROR_BUSHEAVY)!=0 {
		return can.ErrorActive
	} else if (e & pPCAN_ERROR_BUSLIGHT)!=0 {
		return can.ErrorPassive
	} else {
		return can.Error
	}
}

func (p *peakChannel) Reset() {
//	_, _, _ = syscall.Syscall6(uintptr(CAN_Reset), 1, uintptr(p.channel),0,0,0,0,0)
	p.Close()
	p.Initialize(p.bitrate)
}

func (p *peakChannel) Close() {
	if p!=nil {
		_, _, _ = syscall.Syscall6(uintptr(CAN_Uninitialize), 1, uintptr(p.channel), 0, 0, 0, 0, 0)
	}
}

func (p *peakChannel) Status() can.Status {
	ret, _, _ := syscall.Syscall6(uintptr(CAN_GetStatus), 1, uintptr(p.channel),0,0,0,0,0)
	p.status = convertError(ret)
	return p.status
}

func (p *peakChannel) StatusString() string {
	return can.StatusToString(p.status)
}

func  (p *peakChannel)ErrorText(err uint32) string {
	var ErrorName = [...]string{"Ok", "ErrorPassive", "ErrorActive", "BusOff","InternalError"}
	return ErrorName[p.status]
}

func (p *peakChannel) Read() *can.Msg {
	if p==nil {
		return nil
	}
	msg := can.Msg{}
	timestamp := canTimeStamp{}
	ret, _, callErr := syscall.Syscall6(uintptr(CAN_Read), 3, uintptr(p.channel), uintptr(unsafe.Pointer(&msg)), uintptr(unsafe.Pointer(&timestamp)),0,0,0)
	p.status = convertError(ret)
	if callErr!=0 {
		p.status = can.Error
		return nil
	}
	if msg.Type > 2 {
		p.status = convertError(uintptr(msg.Data[3]))
		//fmt.Printf("Got message with status data, status=%s, msg=%s\n",p.StatusString(),msg.ToString())
		return nil
	}
	if ret==pPCAN_ERROR_QRCVEMPTY {
		return nil
	}
	if ret> pPCAN_ERROR_QXMTFULL {
		return nil
	}
	return &msg
}

func (p *peakChannel) Write(msg can.Msg) {
	if p!=nil {
		_, _, _ = syscall.Syscall6(uintptr(CAN_Write), 2, uintptr(p.channel), uintptr(unsafe.Pointer(&msg)),0,0,0,0)
	}
}

var value = 1

func (p *peakChannel) Initialize(bitrate int) error {
	var Btr0Btr1 uint16
	p.bitrate = bitrate
	switch bitrate {
	case 1000000:
		Btr0Btr1 = pPCAN_BAUD_1M
	case 800000:
		Btr0Btr1 = pPCAN_BAUD_800K
	case 500000:
		Btr0Btr1 = pPCAN_BAUD_500K
	case 250000:
		Btr0Btr1 = pPCAN_BAUD_250K
	case 125000:
		Btr0Btr1 = pPCAN_BAUD_125K
	case 100000:
		Btr0Btr1 = pPCAN_BAUD_100K
	case 50000:
		Btr0Btr1 = pPCAN_BAUD_50K
	case 20000:
		Btr0Btr1 = pPCAN_BAUD_20K
	case 10000:
		Btr0Btr1 = pPCAN_BAUD_10K
	}
	if peakError!= nil {
		return fmt.Errorf("could not find \"pcanbasic.dll\", please check installation of peak driver")
	}
	if CAN_Initialize==uintptr(0) {
		return fmt.Errorf("error loading \"pcanbasic.dll\", please check installation of peak driver")
	}
	ret, _, callErr := syscall.Syscall6(uintptr(CAN_Initialize), 5, uintptr(p.channel), uintptr(Btr0Btr1), uintptr(0), uintptr(0), uintptr(0), 0)
	if ret&0x400!=0 {
		return fmt.Errorf("adapter not connected")
	} else if ret!=0 {
		return fmt.Errorf("error %d", ret)
	}
	if callErr!=0 {
		return fmt.Errorf("Initialize error %d", callErr)
	}
	if unsafe.Pointer(CAN_SetValue) == nil {
		fmt.Printf("CAN_SetValue==nil")
		os.Exit(255)
	}
	ret, _ , callErr = syscall.Syscall6(uintptr(CAN_SetValue), 4, uintptr(p.channel), uintptr(PCAN_BUSOFF_AUTORESET), uintptr(unsafe.Pointer(&value)), uintptr(4),0,0)
	if callErr!=0 || ret!=0 {
		return fmt.Errorf("Set pcan busoff autoreset error %d", callErr)
	}
	return nil
}


func New(channel uint8, bitrate int) (*peakChannel, error) {
	bus := &peakChannel{channel: channel}
	err := bus.Initialize(bitrate)
	if err!=nil {
		return nil, err
	}
	return bus, nil
}
