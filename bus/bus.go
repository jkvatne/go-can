package bus

import (
	"go-can/can"
	"sync"
	"time"
)

type State struct {
	dev           can.Device
	mutex         sync.Mutex
	timeout       time.Duration
	terminated    bool
	responseChan  chan *can.Msg
	requestChan   chan can.CobId
	sleepInterval time.Duration
	pollTime      time.Time
}

type MsgHandler func(msg *can.Msg)
var Handler MsgHandler

func (c *State) Status() can.Status {
	return c.dev.Status()
}

func (c *State) Reset() {
	c.dev.Reset()
}

func (c *State) Poll(msg can.Msg, responseId can.CobId) *can.Msg {
	// Polling on a terminated channel is not allowed
	if c.terminated {
		return nil
	}
	// Allow only one thread to poll at the same time
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.requestChan<-responseId
	c.dev.Write(msg)
	resp :=<-c.responseChan
	return resp
}


// RxThread is a goroutine that will check for incoming messages and call a callback for each
func RxThread(c *State) {
	var responseId can.CobId = 0xFFFFFFFF
	for !c.terminated {
		select {
		case responseId = <-c.requestChan:
			c.pollTime = time.Now()
		default:
		}

		m := c.dev.Read()
		if m!=nil {
			if m.Id == responseId {
				responseId = 0xFFFFFFFF
				c.responseChan <- m
			} else {
				Handler(m)
			}
		} else {
			// No messages found, sleep for some time
			time.Sleep(c.sleepInterval)
		}
		// Check for poll timeout
		if  responseId!=0xFFFFFFFF && time.Since(c.pollTime)>c.timeout {
			responseId = 0xFFFFFFFF
			c.responseChan <- nil
		}
	}
}

func (c *State) SetOperational() {
	c.dev.Write(can.NewStdMsg(0, []uint8{1,0}))
}

func (c *State) SetPreoperational(nodeId uint8) {
	c.dev.Write(can.NewStdMsg(0, []uint8{128,0}))
}

func (c *State) SendSync() {
	c.dev.Write(can.NewStdMsg(0x680, []uint8{}))
}

func (c *State) Write(msg can.Msg) {
	c.dev.Write(msg)
}


func (c *State) Close() {
	c.terminated = true
	// Wait for go routine to terminate
	time.Sleep(2*c.sleepInterval)
	c.dev.Close()
}


func New(d can.Device, timeout time.Duration) *State {
	c := &State{dev: d, timeout: timeout}
	c.mutex = sync.Mutex{}
	c.sleepInterval = 10*time.Millisecond
	c.requestChan = make(chan can.CobId, 1)
	c.responseChan = make(chan *can.Msg, 1)
	go RxThread(c)
	return c
}
