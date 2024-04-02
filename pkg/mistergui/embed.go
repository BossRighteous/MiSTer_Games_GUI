package mistergui

import _ "embed"

//go:embed embed/Roboto-Regular.ttf
var robotoRegular []byte

//go:embed embed/Roboto-Bold.ttf
var robotoBold []byte

//go:embed embed/Roboto-Black.ttf
var robotoBlack []byte

//go:embed embed/powerstone.png
var powerstone []byte

type EmbededData struct {
	RobotoRegular *[]byte
	RobotoBold    *[]byte
	RobotoBlack   *[]byte
	Powerstone    *[]byte
}

var Embeds *EmbededData = &EmbededData{
	RobotoRegular: &robotoRegular,
	RobotoBold:    &robotoBold,
	RobotoBlack:   &robotoBlack,
	Powerstone:    &powerstone,
}
