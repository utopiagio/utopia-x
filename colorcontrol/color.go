/* color.go */

package colorcontrol

import (
	"log"
	"encoding/hex"
	"strconv"
	"strings"
	ui "github.com/utopiagio/utopia"
	pointer_gio "github.com/utopiagio/gio/io/pointer"
)

// GoColorState is the state of a colorcontrol.
type GoColorState struct {
	R, G, B, A uint8
}

// GoColorControlObj lays out a color control using utopia widgets.
type GoColorControlObj struct {
	Layout *ui.GoLayoutObj
	GoColorState
	ColorView *ui.GoCanvasObj
	TagGoColor *ui.GoTextEditObj
	TagImColor *ui.GoTextEditObj
	//Theme *material_gio.Theme
	Label *ui.GoLabelObj
	RBar *ui.GoSliderObj
	GBar *ui.GoSliderObj
	BBar *ui.GoSliderObj
	ABar *ui.GoSliderObj
	RBarValue *ui.GoTextEditObj
	GBarValue *ui.GoTextEditObj
	BBarValue *ui.GoTextEditObj
	ABarValue *ui.GoTextEditObj
}

// GoColorControl creates a color control with a color state.
func GoColorControl(parent ui.GoObject, label string) (hObj *GoColorControlObj) {
	layout := ui.GoVFlexBoxLayout(parent)
	header := ui.H3Label(layout, label)
	hCCtrl := &GoColorControlObj{
		Layout: layout,
		Label: header,
		GoColorState: GoColorState{0,0,0,255},
	}

	viewLayout := ui.GoHFlexBoxLayout(layout)
	viewLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	//viewLayout.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Blue)

	colorLayout := ui.GoVFlexBoxLayout(viewLayout)
	colorLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	colorLayout.SetMargin(4,4,4,4)
	//colorLayout.SetBorder(ui.BorderSingleLine, 1, 5, ui.Color_Blue)
	colorLayout.SetPadding(10,10,10,10)

	hCCtrl.ColorView = ui.GoCanvas(colorLayout)
	hCCtrl.ColorView.SetSizePolicy(ui.FixedWidth, ui.FixedHeight)
	hCCtrl.ColorView.SetOnPointerClick(hCCtrl.CopyToClipBoard)
	
	shape := ui.GoCanvasPath(hCCtrl.ColorView)
	shape.AddPoint(50,0)
	shape.AddPoint(100,0)
	shape.AddPoint(100,50)
	shape.AddPoint(50,50)
	shape.AddPoint(50,0)
	shape.SetFillColor(ui.Color_Orange)
	shape.SetLineColor(ui.Color_Gray)
	shape.SetLineWidth(4)

	circle := ui.GoCanvasEllipse(hCCtrl.ColorView)
	circle.SetSize(20, 40)
	circle.SetFillColor(ui.Color_Orange)
	circle.SetLineColor(ui.Color_Gray)
	circle.SetLineWidth(4)
	circle.Move(20, 20)
	//hCCtrl.ColorView.SetFocus(true)

	colorTag := ui.GoLabel(colorLayout, "ColorViewer")
	colorTag.SetSizePolicy(ui.FixedWidth, ui.FixedHeight)
	//colorTag.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Blue)
	colorTag.SetMargin(4,4,4,4)

	
	goTagLayout := ui.GoHFlexBoxLayout(colorLayout)
	goTagLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)

	imTagLayout := ui.GoHFlexBoxLayout(colorLayout)
	imTagLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)

	lblGoColor := ui.GoLabel(goTagLayout, "GoColor ARGB:")
	lblGoColor.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	lblGoColor.SetMargin(4,4,4,4)
	lblGoColor.SetSelectable(true)

	hCCtrl.TagGoColor = ui.GoTextEdit(goTagLayout, "0x00000000")
	hCCtrl.TagGoColor.SetFontVariant("Mono")
	hCCtrl.TagGoColor.SetMargin(2,2,2,2)
	//hCCtrl.TagGoColor.SetBorder(ui.BorderSingleLine, 2, 2, ui.Color_Blue)
	hCCtrl.TagGoColor.SetPadding(6,2,6,2)

	lblImColor := ui.GoLabel(imTagLayout, "ImageColor RGBA:")
	lblImColor.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	lblImColor.SetMargin(4,4,4,4)
	lblImColor.SetSelectable(true)

	hCCtrl.TagImColor = ui.GoTextEdit(imTagLayout, "0x00000000")
	hCCtrl.TagImColor.SetFontVariant("Mono")
	hCCtrl.TagImColor.SetMargin(2,2,2,2)
	//hCCtrl.TagImColor.SetBorder(ui.BorderSingleLine, 2, 2, ui.Color_Blue)
	hCCtrl.TagImColor.SetPadding(6,2,6,2)

	// color selector bars
	barLayout := ui.GoVFlexBoxLayout(viewLayout)
	barLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	barLayout.SetMargin(4,4,4,4)
	barLayout.SetBorder(ui.BorderSingleLine, 1, 5, ui.Color_Blue)
	barLayout.SetPadding(10,10,10,10)

	rbarLayout := ui.GoHFlexBoxLayout(barLayout)
	rbarLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	gbarLayout := ui.GoHFlexBoxLayout(barLayout)
	gbarLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	bbarLayout := ui.GoHFlexBoxLayout(barLayout)
	bbarLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	abarLayout := ui.GoHFlexBoxLayout(barLayout)
	abarLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)

	rLabel := ui.H5Label(rbarLayout, "R: ")
	rLabel.SetFontVariant("Mono")
	
	hCCtrl.RBar = ui.GoSlider(rbarLayout, 0, 1)
	hCCtrl.RBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	hCCtrl.RBar.SetOnChange(hCCtrl.RBar_Drag)
	hCCtrl.RBarValue = ui.GoTextEdit(rbarLayout, "")
	hCCtrl.RBarValue.SetFontVariant("Mono")
	hCCtrl.RBarValue.SetMargin(8,0,0,0)
	
	gLabel := ui.H5Label(gbarLayout, "G: ")
	gLabel.SetFontVariant("Mono")

	hCCtrl.GBar = ui.GoSlider(gbarLayout, 0, 1)
	hCCtrl.GBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	hCCtrl.GBar.SetOnChange(hCCtrl.GBar_Drag)
	hCCtrl.GBarValue = ui.GoTextEdit(gbarLayout, "")
	hCCtrl.GBarValue.SetFontVariant("Mono")
	hCCtrl.GBarValue.SetMargin(8,0,0,0)

	bLabel := ui.H5Label(bbarLayout, "B: ")
	bLabel.SetFontVariant("Mono")

	hCCtrl.BBar = ui.GoSlider(bbarLayout, 0, 1)
	hCCtrl.BBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	hCCtrl.BBar.SetOnChange(hCCtrl.BBar_Drag)
	hCCtrl.BBarValue = ui.GoTextEdit(bbarLayout, "")
	hCCtrl.BBarValue.SetFontVariant("Mono")
	hCCtrl.BBarValue.SetMargin(8,0,0,0)

	aLabel := ui.H5Label(abarLayout, "A: ")
	aLabel.SetFontVariant("Mono")

	hCCtrl.ABar = ui.GoSlider(abarLayout, 0, 1)
	hCCtrl.ABar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	hCCtrl.ABar.SetOnChange(hCCtrl.ABar_Drag)
	hCCtrl.ABarValue = ui.GoTextEdit(abarLayout, "")
	hCCtrl.ABarValue.SetFontVariant("Mono")
	hCCtrl.ABarValue.SetMargin(8,0,0,0)

	hCCtrl.SetColor(ui.Color_Blue)

	return hCCtrl
}

func (ob *GoColorControlObj) CopyToClipBoard(e pointer_gio.Event) {
	log.Println("CopyToClipBoard:", e.Type)
	/*log.Println("Type:", e.Type)
	log.Println("Source:", e.Source)
	log.Println("PointerID:", e.PointerID)
	log.Println("Priority:", e.Priority)
	log.Println("Time:", e.Time)
	log.Println("Buttons:", e.Buttons)
	log.Println("Position:", e.Position)
	log.Println("Scroll:", e.Scroll)
	log.Println("Modifiers:", e.Modifiers)*/
}

func (ob *GoColorControlObj) GetColor() (color ui.GoColor) {
	return ui.RGBAColor(ob.R, ob.G, ob.B, ob.A)
}

func (ob *GoColorControlObj) Refresh() {
	ob.Layout.ParentWindow().Refresh()
}

func (ob *GoColorControlObj) SetBorder(style ui.GoBorderStyle, width int, radius int, color ui.GoColor) {
	ob.Layout.SetBorder(style, width, radius, color)
}

// SetColor function updates the color bar components
func (ob *GoColorControlObj) SetColor(color ui.GoColor) {
	r, g, b, a := color.RGBA()
	ob.GoColorState = GoColorState{r, g, b, a}
	ob.RBar.SetValue(float32(r) / 255)
	ob.RBar_Drag(ob.RBar.Value())
	ob.GBar.SetValue(float32(g) / 255)
	ob.GBar_Drag(ob.GBar.Value())
	ob.BBar.SetValue(float32(b) / 255)
	ob.BBar_Drag(ob.BBar.Value())
	ob.ABar.SetValue(float32(a / 255))
	ob.ABar_Drag(ob.ABar.Value())
	ob.updateColorView()
	log.Println("ob.ObjectType:", ob.Layout.ParentWindow().ObjectType())
	ob.Refresh()
}

// SetMargin function sets the margin of the controls within its parent layout
func (ob *GoColorControlObj) SetMargin(left int, top int, right int, bottom int) {
	ob.Layout.GoMargin = ui.GoMargin{left, top, right, bottom}
}

// SetPadding function sets the padding of the control interface within the control layout
func (ob *GoColorControlObj) SetPadding(left int, top int, right int, bottom int) {
	ob.Layout.GoPadding = ui.GoPadding{left, top, right, bottom}
}

// SetSizePolicy function sets how the control resizes to fill its parent layout
func (ob *GoColorControlObj) SetSizePolicy(horiz ui.GoSizeType, vert ui.GoSizeType) {	// widget sizing policy - GoSizePolicy{horiz, vert, fixed}
	ob.Layout.SetSizePolicy(horiz, vert)
}

// RBar_Drag function updates the ui with the new R component value
func (ob *GoColorControlObj) RBar_Drag(value float32) {
	val := uint8(value * 255)
	ob.RBarValue.SetText(valueString(val))
	ob.R = val
	ob.updateColorView()
}

// GBar_Drag function updates the ui with the new G component value
func (ob *GoColorControlObj) GBar_Drag(value float32) {
	val := uint8(value * 255)
	ob.GBarValue.SetText(valueString(val))
	ob.G = val
	ob.updateColorView()
}

// BRBar_Drag function updates the ui with the new B component value
func (ob *GoColorControlObj) BBar_Drag(value float32) {
	val := uint8(value * 255)
	ob.BBarValue.SetText(valueString(val))
	ob.B = val
	ob.updateColorView()
}

// ABar_Drag function updates the ui with the new A component value
func (ob *GoColorControlObj) ABar_Drag(value float32) {
	val := uint8(value * 255)
	ob.ABarValue.SetText(valueString(val))
	ob.A = val
	ob.updateColorView()
}

// Internal function: updateColorView updates the ui when the color control state changes
func (ob *GoColorControlObj) updateColorView() {
	color := ob.GetColor()
	ob.ColorView.SetBackgroundColor(color)
	ob.TagGoColor.SetText("0x" + hex.EncodeToString([]byte{ob.A, ob.R, ob.G, ob.B}))
	ob.TagImColor.SetText("0x" + hex.EncodeToString([]byte{ob.R, ob.G, ob.B, ob.A}))
}

// Internal function: valueString returns a string representative of the number(uint8)
func valueString(in uint8) string {
	s := strconv.Itoa(int(in))
	delta := 3 - len(s)
	if delta > 0 {
		s = strings.Repeat(" ", delta) + s
	}
	return s
}