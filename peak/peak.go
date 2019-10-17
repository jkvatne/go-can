package peak

import (
	"fmt"
	"go-can/can"
	"sync"
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
	mutex sync.Mutex
	wg sync.WaitGroup
	pollTime time.Time
	responseId can.Id
	response *can.Msg
	timeout time.Duration
    responseChan chan *can.Msg
    requestChan chan can.Id
	sleepInterval time.Duration
}

// Check if peak interface satisfies can.Device interface
var _ can.Device = &peakChannel{}

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

	PCAN_ERROR_OK            = uintptr(0x00000)  // No error
	PCAN_ERROR_XMTFULL       = uintptr(0x00001)  // Transmit buffer in CAN controller is full
	PCAN_ERROR_OVERRUN       = uintptr(0x00002)  // CAN controller was read too late
	PCAN_ERROR_BUSLIGHT      = uintptr(0x00004)  // Bus error: an error counter reached the 'light' limit
	PCAN_ERROR_BUSHEAVY      = uintptr(0x00008)  // Bus error: an error counter reached the 'heavy' limit
	PCAN_ERROR_BUSOFF        = uintptr(0x00010)  // Bus error: the CAN controller is in bus-off state
	PCAN_ERROR_ANYBUSERR     = uintptr(PCAN_ERROR_BUSLIGHT | PCAN_ERROR_BUSHEAVY | PCAN_ERROR_BUSOFF)  // Mask for all bus errors
	PCAN_ERROR_QRCVEMPTY     = uintptr(0x00020)  // Receive queue is empty
	PCAN_ERROR_QOVERRUN      = uintptr(0x00040)  // Receive queue was read too late
	PCAN_ERROR_QXMTFULL      = uintptr(0x00080)  // Transmit queue is full
	PCAN_ERROR_REGTEST       = uintptr(0x00100)  // Test of the CAN controller hardware registers failed (no hardware found)
	PCAN_ERROR_NODRIVER      = uintptr(0x00200)  // Driver not loaded
	PCAN_ERROR_HWINUSE       = uintptr(0x00400)  // Hardware already in use by a Net
	PCAN_ERROR_NETINUSE      = uintptr(0x00800)  // A Client is already connected to the Net
	PCAN_ERROR_ILLHW         = uintptr(0x01400)  // Hardware handle is invalid
	PCAN_ERROR_ILLNET        = uintptr(0x01800)  // Net handle is invalid
	PCAN_ERROR_ILLCLIENT     = uintptr(0x01C00)  // Client handle is invalid
	PCAN_ERROR_ILLHANDLE     = uintptr(PCAN_ERROR_ILLHW | PCAN_ERROR_ILLNET | PCAN_ERROR_ILLCLIENT)  // Mask for all handle errors
	PCAN_ERROR_RESOURCE      = uintptr(0x02000)  // Resource (FIFO, Client, timeout) cannot be created
	PCAN_ERROR_ILLPARAMTYPE  = uintptr(0x04000)  // Invalid parameter
	PCAN_ERROR_ILLPARAMVAL   = uintptr(0x08000)  // Invalid parameter value
	PCAN_ERROR_UNKNOWN       = uintptr(0x10000)  // Unknow error
	PCAN_ERROR_ILLDATA       = uintptr(0x20000)  // Invalid data, function, or action
	PCAN_ERROR_INITIALIZE    = uintptr(0x40000)  // peakChannel is not initialized
	PCAN_ERROR_ILLOPERATION  = uintptr(0x80000)  // Invalid operation
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


func convertError(e uintptr) can.Status {
	if e & 0xDF == 0 {
		return can.Ok
	} else if (e & PCAN_ERROR_BUSOFF)!=0 {
		return can.BusOff
	} else if (e & PCAN_ERROR_BUSHEAVY)!=0 {
		return can.ErrorActive
	} else if (e & PCAN_ERROR_BUSLIGHT)!=0 {
		return can.ErrorPassive
	} else {
		return can.Error
	}
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
	msg := can.Msg{}
	timestamp := CanTimeStamp{}
	ret, _, callErr := syscall.Syscall6(uintptr(CAN_Read), 3, uintptr(p.channel), uintptr(unsafe.Pointer(&msg)), uintptr(unsafe.Pointer(&timestamp)),0,0,0)
	p.status = convertError(ret)
	if callErr!=0 {
		p.status = can.Error
		return nil
	}
	if msg.Type > 2 {
		p.status = convertError(uintptr(msg.Data[3]))
		fmt.Printf("Got message with status data, status=%s\n",p.StatusString())
		return nil
	}
	if ret==0x20 {
		return nil
	}
	return &msg
}

func (p *peakChannel) Write(msg can.Msg) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	_, _, _ = syscall.Syscall6(uintptr(CAN_Write), 2, uintptr(p.channel), uintptr(unsafe.Pointer(&msg)),0,0,0,0)
}

func New(channel uint8, Btr0Btr1 uint16, callback can.Callback) (*peakChannel, error) {
	bus := &peakChannel{channel: channel, btr: Btr0Btr1, callback:callback}
	err := bus.Initialize(Btr0Btr1)
	bus.timeout = 100*time.Millisecond
	bus.requestChan = make(chan can.Id, 1)
	bus.responseChan = make(chan *can.Msg, 1)
	go Handler(bus,callback)
	return bus, err
}

// Handler is a goroutine that will check for incoming messages and call a callback for each
func Handler(p *peakChannel, callback can.Callback) {
    var responseId can.Id = 0xFFFFFFFF
	for (!p.terminated) {
		select {
		case responseId = <-p.requestChan:
			p.pollTime = time.Now()
		default:
		}

		m := p.Read()
		if m!=nil {
			if m.Id == responseId {
				responseId = 0xFFFFFFFF
				p.responseChan <- m
			} else {
				callback(m)
			}
		} else {
			// No messages found, sleep for some time
			time.Sleep(p.sleepInterval * time.Millisecond)
		}
		// Check for poll timeout
		if  responseId!=0xFFFFFFFF && time.Since(p.pollTime)>p.timeout {
			responseId = 0xFFFFFFFF
			p.responseChan <- nil
		}
	}
}

func (p *peakChannel) Poll(msg can.Msg, responseId can.Id) *can.Msg {
	// Polling on a terminated channel is not allowed
	if p.terminated {
		return nil
	}
	// Allow only one thread to poll at the same time
	p.mutex.Lock()
	defer p.mutex.Unlock()
    p.requestChan<-responseId
	_, _, _ = syscall.Syscall6(uintptr(CAN_Write), 2, uintptr(p.channel), uintptr(unsafe.Pointer(&msg)),0,0,0,0)
	resp :=<-p.responseChan
    return resp
}