package mistergui

import _ "embed"

//go:embed embed/Roboto-Regular.ttf
var robotoRegular []byte

//go:embed embed/Roboto-Bold.ttf
var robotoBold []byte

//go:embed embed/Roboto-Black.ttf
var robotoBlack []byte

type EmbededData struct {
	RobotoRegular *[]byte
	RobotoBold    *[]byte
	RobotoBlack   *[]byte
}

var Embeds *EmbededData = &EmbededData{
	RobotoRegular: &robotoRegular,
	RobotoBold:    &robotoBold,
	RobotoBlack:   &robotoBlack,
}
