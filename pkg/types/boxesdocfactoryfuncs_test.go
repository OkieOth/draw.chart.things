package types_test

import (
	"testing"
)

type DummyDimensionCalculator struct {
	captionWidth  int32
	captionHeight int32
	text1Width    int32
	text1Height   int32
	text2Width    int32
	text2Height   int32
}

func (d *DummyDimensionCalculator) CaptionDimensions(txt string) (width, height int32) {
	return d.captionWidth, d.captionHeight
}

func (d *DummyDimensionCalculator) Text1Dimensions(txt string) (width, height int32) {
	return d.text1Width, d.text1Height
}

func (d *DummyDimensionCalculator) Text2Dimensions(txt string) (width, height int32) {
	return d.text2Width, d.text2Height
}

func NewDummyDimensionCalculator(captionWidth, captionHeight, text1Width, text1Height, text2Width, text2Height int32) *DummyDimensionCalculator {
	return &DummyDimensionCalculator{
		captionWidth:  captionWidth,
		captionHeight: captionHeight,
		text1Width:    text1Width,
		text1Height:   text1Height,
		text2Width:    text2Width,
		text2Height:   text2Height,
	}
}

func TestInitDimensions(t *testing.T) {
	// TODO
}
