package groovymister

import "encoding/binary"

const (
	bitpos_Right uint16 = (1 << 0)
	bitpos_Left  uint16 = (1 << 1)
	bitpos_Down  uint16 = (1 << 2)
	bitpos_Up    uint16 = (1 << 3)
	bitpos_B1    uint16 = (1 << 4)
	bitpos_B2    uint16 = (1 << 5)
	bitpos_B3    uint16 = (1 << 6)
	bitpos_B4    uint16 = (1 << 7)
	bitpos_B5    uint16 = (1 << 8)
	bitpos_B6    uint16 = (1 << 9)
	bitpos_B7    uint16 = (1 << 10)
	bitpos_B8    uint16 = (1 << 11)
	bitpos_B9    uint16 = (1 << 12)
)

type GroovyInputPacket struct {
	JoyFrame uint32
	JoyOrder uint8
	Joy1Mask uint16 // bitmask
	Joy2Mask uint16 // bitmask
	//Joy1     GroovyJoy
	//Joy2     GroovyJoy
}

func InputPacketFromBuffer(buf []uint8) GroovyInputPacket {
	m1 := binary.LittleEndian.Uint16(buf[5:7])
	m2 := binary.LittleEndian.Uint16(buf[7:9])
	return GroovyInputPacket{
		JoyFrame: binary.LittleEndian.Uint32(buf[0:4]),
		JoyOrder: buf[4],
		Joy1Mask: m1,
		Joy2Mask: m2,
		//Joy1:     GroovyJoyFromMask(m1),
		//Joy2:     GroovyJoyFromMask(m2),
	}
}

type GroovyJoy struct {
	Right bool
	Left  bool
	Down  bool
	Up    bool
	B1    bool
	B2    bool
	B3    bool
	B4    bool
	B5    bool
	B6    bool
	B7    bool
	B8    bool
	B9    bool
}

func GroovyJoyFromMask(m uint16) GroovyJoy {
	return GroovyJoy{
		Right: m&bitpos_Right != 0,
		Left:  m&bitpos_Left != 0,
		Down:  m&bitpos_Down != 0,
		Up:    m&bitpos_Up != 0,
		B1:    m&bitpos_B1 != 0,
		B2:    m&bitpos_B2 != 0,
		B3:    m&bitpos_B3 != 0,
		B4:    m&bitpos_B4 != 0,
		B5:    m&bitpos_B5 != 0,
		B6:    m&bitpos_B6 != 0,
		B7:    m&bitpos_B7 != 0,
		B8:    m&bitpos_B8 != 0,
		B9:    m&bitpos_B9 != 0,
	}
}

type InputKey uint8

const (
	InputRight InputKey = iota
	InputLeft
	InputDown
	InputUp
	InputB1
	InputB2
	InputB3
	InputB4
	InputB5
	InputB6
	InputB7
	InputB8
	InputB9
)

func isJoyPressed(joy GroovyJoy, key InputKey) bool {
	switch key {
	case InputRight:
		return joy.Right
	case InputLeft:
		return joy.Left
	case InputDown:
		return joy.Down
	case InputUp:
		return joy.Up
	case InputB1:
		return joy.B1
	case InputB2:
		return joy.B2
	case InputB3:
		return joy.B3
	case InputB4:
		return joy.B4
	case InputB5:
		return joy.B5
	case InputB6:
		return joy.B6
	case InputB7:
		return joy.B7
	case InputB8:
		return joy.B8
	case InputB9:
		return joy.B9
	}
	return false
}

type GroovyInput struct {
	PrevJoy1 GroovyJoy
	Joy1     GroovyJoy
	PrevJoy2 GroovyJoy
	Joy2     GroovyJoy
}

func (input *GroovyInput) AddInputPacket(packet GroovyInputPacket) {
	input.PrevJoy1 = input.Joy1
	input.Joy1 = GroovyJoyFromMask(packet.Joy1Mask)
	input.PrevJoy2 = input.Joy2
	input.Joy2 = GroovyJoyFromMask(packet.Joy2Mask)
}

func (input *GroovyInput) IsPressed(joyNum uint8, key InputKey) bool {
	if joyNum == 2 {
		return isJoyPressed(input.Joy2, key)
	}
	return isJoyPressed(input.Joy1, key)
}

func (input *GroovyInput) IsJustPressed(joyNum uint8, key InputKey) bool {
	if joyNum == 2 {
		return isJoyPressed(input.Joy2, key) && !isJoyPressed(input.PrevJoy2, key)
	}
	return isJoyPressed(input.Joy1, key) && !isJoyPressed(input.PrevJoy1, key)
}

func (input *GroovyInput) IsJustReleased(joyNum uint8, key InputKey) bool {
	if joyNum == 2 {
		return !isJoyPressed(input.Joy2, key) && isJoyPressed(input.PrevJoy2, key)
	}
	return !isJoyPressed(input.Joy1, key) && isJoyPressed(input.PrevJoy1, key)
}
