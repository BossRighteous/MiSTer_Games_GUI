package groovymister

import "encoding/binary"

/*
memcpy(&inputs.joyFrame, &m_bufferInputsReceive[0], 4);
memcpy(&inputs.joyOrder, &m_bufferInputsReceive[4], 1);
memcpy(&inputs.joy1, &m_bufferInputsReceive[5], 2);
memcpy(&inputs.joy2, &m_bufferInputsReceive[7], 2);

#define GM_JOY_RIGHT (1 << 0)
#define GM_JOY_LEFT  (1 << 1)
#define GM_JOY_DOWN  (1 << 2)
#define GM_JOY_UP    (1 << 3)
#define GM_JOY_B1    (1 << 4)
#define GM_JOY_B2    (1 << 5)
#define GM_JOY_B3    (1 << 6)
#define GM_JOY_B4    (1 << 7)
#define GM_JOY_B5    (1 << 8)
#define GM_JOY_B6    (1 << 9)
#define GM_JOY_B7    (1 << 10)
#define GM_JOY_B8    (1 << 11)
#define GM_JOY_B9    (1 << 12)
*/

type GroovyInput struct {
	JoyFrame uint32
	JoyOrder uint8
	Joy1Mask uint16 // bitmask
	Joy2Mask uint16 // bitmask
	Joy1     GroovyJoy
	Joy2     GroovyJoy
}

func maskToBool(v uint16, pos uint16) bool {
	return v&(1<<pos) != 0
}

func InputFromBuffer(buf []uint8) GroovyInput {
	m1 := binary.LittleEndian.Uint16(buf[5:7])
	m2 := binary.LittleEndian.Uint16(buf[7:9])
	return GroovyInput{
		JoyFrame: binary.LittleEndian.Uint32(buf[0:4]),
		JoyOrder: buf[4],
		Joy1Mask: m1,
		Joy2Mask: m2,
		Joy1:     GroovyJoyFromMask(m1),
		Joy2:     GroovyJoyFromMask(m2),
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
		Right: maskToBool(m, 0),
		Left:  maskToBool(m, 1),
		Down:  maskToBool(m, 2),
		Up:    maskToBool(m, 3),
		B1:    maskToBool(m, 4),
		B2:    maskToBool(m, 5),
		B3:    maskToBool(m, 6),
		B4:    maskToBool(m, 7),
		B5:    maskToBool(m, 8),
		B6:    maskToBool(m, 9),
		B7:    maskToBool(m, 10),
		B8:    maskToBool(m, 11),
		B9:    maskToBool(m, 12),
	}
}
