package groovymister

import (
	"strconv"
	"strings"
)

type Modeline struct {
	PixelClock float64
	HActive    uint16
	HBegin     uint16
	HEnd       uint16
	HTotal     uint16
	VActive    uint16
	VBegin     uint16
	VEnd       uint16
	VTotal     uint16
	FrameRate  float64
	Interlace  bool
}

func safeInt16(intstr string) uint16 {
	v, err := strconv.ParseInt(intstr, 10, 16)
	if err != nil {
		panic("unable to parse modeline uint16")
	}
	return uint16(v)
}

func safeFloat64(floatstr string) float64 {
	v, err := strconv.ParseFloat(floatstr, 64)
	if err != nil {
		panic("unable to parse modeline float64")
	}
	return v
}

func ModelineFromSettings(modeline string, frameRate float64, interlace bool) *Modeline {
	parts := strings.Fields(modeline)
	if len(parts) < 9 {
		panic("unable to parse modeline")
	}
	return &Modeline{
		PixelClock: safeFloat64(parts[0]),
		HActive:    safeInt16(parts[1]),
		HBegin:     safeInt16(parts[2]),
		HEnd:       safeInt16(parts[3]),
		HTotal:     safeInt16(parts[4]),
		VActive:    safeInt16(parts[5]),
		VBegin:     safeInt16(parts[6]),
		VEnd:       safeInt16(parts[7]),
		VTotal:     safeInt16(parts[8]),
		FrameRate:  frameRate,
		Interlace:  interlace,
	}
}
