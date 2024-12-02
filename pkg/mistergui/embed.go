package mistergui

import _ "embed"

//go:embed embed/PTSans-Regular.ttf
var ptSansRegular []byte

//go:embed embed/listing-bg.png
var listingBg []byte

/*
//go:embed embed/confirm-bg.png
var confirmBg []byte

//go:embed embed/images-bg.png
var imagesBg []byte

//go:embed embed/info-bg.png
var infoBg []byte
*/

type EmbededData struct {
	PtSansRegular *[]byte
	ListingBg     *[]byte
	//ConfirmBg     *[]byte
	//ImagesBg      *[]byte
	//InfoBg        *[]byte
}

var Embeds *EmbededData = &EmbededData{
	PtSansRegular: &ptSansRegular,
	ListingBg:     &listingBg,
	//ConfirmBg:     &confirmBg,
	//ImagesBg:      &imagesBg,
	//InfoBg:        &infoBg,
}
