package virtual

import (
	"fmt"
	"go-can/can"
)

type virtualChannel struct {
	channel uint8
	speed   int
	status  can.Status
	open 	bool
	response *can.Msg
}

// Check if peak interface satisfies can.Device interface
var _ can.Device = &virtualChannel{}


func (p *virtualChannel) Initialize(bitrate int) error {
	return nil
}

func (p *virtualChannel) Reset() {
}

func (p *virtualChannel) Close() {
	p.open = false
}

func (p *virtualChannel) Status() can.Status {
	return p.status
}

func (p *virtualChannel) StatusString() string {
	return fmt.Sprintf("Status %x",p.status)
}

func (p *virtualChannel) Read() *can.Msg {
	msg := p.response
	p.response = nil
	return msg
}

func (p *virtualChannel) Write(msg can.Msg) {
	if msg.Id&0xFF80 == 0x600 {
		p.response = &can.Msg{Id: 0x580+(msg.Id&0x7F), Type: 0, Len:8, Data:[8]uint8{0x43, 0x00, 0x10, 0,0xE4, 0x0C,0}}
	}
}

func New(channel uint8, bitrate int) *virtualChannel {
	dev := &virtualChannel{channel: channel, speed: bitrate}
	dev.open = true
	dev.Initialize(bitrate)
	return dev
}
