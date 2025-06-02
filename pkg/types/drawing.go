package types

type TextDimensionCalculator interface {
	Dimensions(txt string, format *FontDef) (width, height int)
}

type RectWithTextFormat struct {

	// Padding of the box
	Padding int `yaml:"padding"`

	FontCaption FontDef `yaml:"fontCaption"`

	FontText1 FontDef `yaml:"fontText1"`

	FontText2 FontDef `yaml:"fontText2"`

	Border *LineDef `yaml:"border,omitempty"`

	Fill *FillDef `yaml:"fill,omitempty"`

	// Minimum margin between boxes
	MinBoxMargin int `yaml:"minBoxMargin"`

	// If true, the text will be displayed vertically
	VerticalTxt bool `yaml:"verticalTxt"`
}

type Drawing interface {
	Start(title string, height, width int) error
	DrawRectWithText(id, caption, text1, text2 string, x, y, width, height int, format RectWithTextFormat) error
	DrawLine(x1, y1, x2, y2 int, format LineDef) error
	DrawArrow(x, y, angle int, format LineDef) error
	DrawSolidRect(x, y, width, height int, format LineDef) error
	DrawText(text string, x, y, width int, fontDef *FontDef) int
	DrawVerticalText(text string, currentX, y, height int, fontDef *FontDef) int
	Done() error
}
