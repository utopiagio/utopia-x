/*
Package colorpicker provides simple widgets for selecting an RGBA color
and for choosing one of a set of colors.

The PickerStyle type can be used to render a colorpicker (the state will be
stored in a State). Colorpickers allow choosing specific RGBA values with
sliders or providing an RGB hex code.

The MuxStyle type can be used to render a color multiplexer (the state will
be stored in a MuxState). Color multiplexers provide a choice from among a
set of colors.
*/
package colorpicker

import (
	//"log"
	"encoding/hex"
	"image"
	"image/color"
	"strconv"
	"strings"

	ui "github.com/utopiagio/utopia"

	gofont_gio "github.com/utopiagio/gio/font/gofont"
	layout_gio "github.com/utopiagio/gio/layout"
	op_gio "github.com/utopiagio/gio/op"
	clip_gio "github.com/utopiagio/gio/op/clip"
	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"
	widget_gio "github.com/utopiagio/gio/widget"
	material_gio "github.com/utopiagio/gio/widget/material"
)

// MuxState holds the state of a color multiplexer. A color multiplexer allows
// choosing from among a set of colors.
type MuxState struct {
	widget_gio.Enum
	Options        map[string]*color.NRGBA
	OrderedOptions []string
}

// NewMuxState creates a MuxState that will provide choices between
// the MuxOptions given as parameters.
func NewMuxState(options ...MuxOption) MuxState {
	keys := make([]string, 0, len(options))
	mapped := make(map[string]*color.NRGBA)
	for _, opt := range options {
		keys = append(keys, opt.Label)
		mapped[opt.Label] = opt.Value
	}
	state := MuxState{
		Options:        mapped,
		OrderedOptions: keys,
	}
	if len(keys) > 0 {
		state.Enum.Value = keys[0]
	}
	return state
}

// MuxOption is one choice for the value of a color multiplexer.
type MuxOption struct {
	Label string
	Value *color.NRGBA
}

// Color returns the currently-selected color.
func (m MuxState) Color() *color.NRGBA {
	return m.Options[m.Enum.Value]
}

// MuxStyle renders a MuxState as a material design themed widget.
type MuxStyle struct {
	*MuxState
	Theme *material_gio.Theme
	Label string
}

// Mux creates a MuxStyle from a theme and a state.
func Mux(theme *material_gio.Theme, state *MuxState, label string) MuxStyle {
	return MuxStyle{
		Theme:    theme,
		MuxState: state,
		Label:    label,
	}
}

// Layout renders the MuxStyle into the provided context.
func (m MuxStyle) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	gtx.Constraints.Min.Y = 0
	var children []layout_gio.FlexChild
	inset := layout_gio.UniformInset(unit_gio.Dp(2))
	children = append(children, layout_gio.Rigid(func(gtx C) D {
		return inset.Layout(gtx, func(gtx C) D {
			return material_gio.Body1(m.Theme, m.Label).Layout(gtx)
		})
	}))
	for i := range m.OrderedOptions {
		opt := m.OrderedOptions[i]
		children = append(children, layout_gio.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				return m.layoutOption(gtx, opt)
			})
		}))
	}
	return layout_gio.Flex{Axis: layout_gio.Vertical}.Layout(gtx, children...)
}

func (m MuxStyle) layoutOption(gtx C, option string) D {
	return layout_gio.Flex{Alignment: layout_gio.Middle}.Layout(gtx,
		layout_gio.Rigid(func(gtx C) D {
			return material_gio.RadioButton(m.Theme, &m.Enum, option, option).Layout(gtx)
		}),
		layout_gio.Flexed(1, func(gtx C) D {
			return layout_gio.Inset{Left: unit_gio.Dp(8)}.Layout(gtx, func(gtx C) D {
				color := m.Options[option]
				if color == nil {
					return D{}
				}
				return borderedSquare(gtx, *color)
			})
		}),
	)
}

func borderedSquare(gtx C, c color.NRGBA) D {
	dims := square(gtx, unit_gio.Dp(20), color.NRGBA{A: 255})

	off := gtx.Dp(unit_gio.Dp(1))
	defer op_gio.Offset(image.Pt(off, off)).Push(gtx.Ops).Pop()
	square(gtx, unit_gio.Dp(18), c)
	return dims
}

func square(gtx C, sizeDp unit_gio.Dp, color color.NRGBA) D {
	return rect(gtx, sizeDp, sizeDp, color)
}

func rect(gtx C, width, height unit_gio.Dp, color color.NRGBA) D {
	w, h := gtx.Dp(width), gtx.Dp(height)
	return rectAbs(gtx, w, h, color)
}

func rectAbs(gtx C, w, h int, color color.NRGBA) D {
	size := image.Point{X: w, Y: h}
	bounds := image.Rectangle{Max: size}
	paint_gio.FillShape(gtx.Ops, color, clip_gio.Rect(bounds).Op())
	return D{Size: image.Pt(w, h)}
}

// State is the state of a colorpicker.
type State struct {
	R, G, B, A widget_gio.Float
	widget_gio.Editor

	changed bool
}

// SetColor changes the color represented by the colorpicker.
func (s *State) SetColor(c color.NRGBA) {
	s.R.Value = float32(c.R) / 255.0
	s.G.Value = float32(c.G) / 255.0
	s.B.Value = float32(c.B) / 255.0
	s.A.Value = float32(c.A) / 255.0
	s.updateEditor()
}

// Color returns the currently selected color.
func (s State) Color() color.NRGBA {
	return color.NRGBA{
		R: s.Red(),
		G: s.Green(),
		B: s.Blue(),
		A: s.Alpha(),
	}
}

// Red returns the red value of the currently selected color.
func (s State) Red() uint8 {
	return uint8(s.R.Value * 255)
}

// Green returns the green value of the currently selected color.
func (s State) Green() uint8 {
	return uint8(s.G.Value * 255)
}

// Blue returns the blue value of the currently selected color.
func (s State) Blue() uint8 {
	return uint8(s.B.Value * 255)
}

// Alpha returns the alpha value of the currently selected color.
func (s State) Alpha() uint8 {
	return uint8(s.A.Value * 255)
}

// Changed returns whether the color has changed since last frame.
func (s State) Changed() bool {
	return s.changed
}

// Layout handles all state updates from the underlying widgets.
func (s *State) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	s.changed = false
	if s.R.Changed() || s.G.Changed() || s.B.Changed() || s.A.Changed() {
		s.updateEditor()
	}
	if events := s.Editor.Events(); len(events) != 0 {
		out, err := hex.DecodeString(s.Editor.Text())
		if err == nil && len(out) == 3 {
			s.R.Value = (float32(out[0]) / 255.0)
			s.G.Value = (float32(out[1]) / 255.0)
			s.B.Value = (float32(out[2]) / 255.0)
			s.changed = true
		}
	}
	return layout_gio.Dimensions{}
}

func (s *State) updateEditor() {
	s.Editor.SetText(hex.EncodeToString([]byte{s.Red(), s.Green(), s.Blue()}))
	s.changed = true
}

// PickerStyle renders a color picker using material widgets.
type GoPickerObj struct {
	ui.GioObject
	ui.GioWidget
	*State
	Theme *material_gio.Theme
	Label string
}

type (
	C = layout_gio.Context
	D = layout_gio.Dimensions
)

// Picker creates a pickerstyle from a theme and a state.
func GoPicker(parent ui.GoObject, label string) (hObj *GoPickerObj) {
	//var theme *ui.GoThemeObj = ui.GoApp.Theme()
	th := material_gio.NewTheme(gofont_gio.Collection())
	state := &State{}
	object := ui.GioObject{parent, parent.ParentWindow(), []ui.GoObject{}, ui.GetSizePolicy(ui.FixedWidth, ui.FixedHeight)}
	widget := ui.GioWidget{
		GoBorder: ui.GoBorder{ui.BorderNone, ui.Color_Black, 0, 0},
		GoMargin: ui.GoMargin{0,0,0,0},
		GoPadding: ui.GoPadding{0,0,0,0},
		Visible: true,
	}
	hPicker := &GoPickerObj{
		object,
		widget,
		state,
		th,
		label,
	}
	parent.AddControl(hPicker)
	return hPicker
}

func (ob *GoPickerObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx C) D {
			return ob.GoBorder.Layout(gtx, func(gtx C) D {
				return ob.GoPadding.Layout(gtx, func(gtx C) D {
					return ob.layout(gtx)
				})
			})
		})
	}
	return dims
}

// Layout renders the PickerStyle into the provided context.
func (p *GoPickerObj) layout(gtx layout_gio.Context) layout_gio.Dimensions {
	p.State.Layout(gtx)

	// lay out the label and editor to compute their width
	leftSide := op_gio.Record(gtx.Ops)
	leftSideDims := p.layoutLeftPane(gtx)
	layoutLeft := leftSide.Stop()

	// lay out the sliders in the remaining horizontal space
	rgtx := gtx
	rgtx.Constraints.Max.X -= leftSideDims.Size.X
	rightSide := op_gio.Record(gtx.Ops)
	rightSideDims := p.layoutSliders(rgtx)
	layoutRight := rightSide.Stop()

	// compute the space beneath the editor that will not extend
	// past the sliders vertically
	margin := gtx.Dp(unit_gio.Dp(4))
	sampleWidth, sampleHeight := leftSideDims.Size.X, rightSideDims.Size.Y-leftSideDims.Size.Y

	// lay everything out for real, starting with the editor/label
	layoutLeft.Add(gtx.Ops)

	// offset downwards and lay out the color sample
	var stack op_gio.TransformStack
	stack = op_gio.Offset(image.Pt(margin, leftSideDims.Size.Y)).Push(gtx.Ops)
	rectAbs(gtx, sampleWidth-(2*margin), sampleHeight-(2*margin), p.State.Color())
	stack.Pop()

	// offset to the right to lay out the sliders
	defer op_gio.Offset(image.Pt(leftSideDims.Size.X, 0)).Push(gtx.Ops).Pop()
	layoutRight.Add(gtx.Ops)

	return layout_gio.Dimensions{
		Size: image.Point{
			X: gtx.Constraints.Max.X,
			Y: rightSideDims.Size.Y,
		},
	}
}

func (p *GoPickerObj) layoutLeftPane(gtx C) D {
	gtx.Constraints.Min.X = 0
	inset := layout_gio.UniformInset(unit_gio.Dp(4))
	dims := layout_gio.Flex{Axis: layout_gio.Vertical}.Layout(gtx,
		layout_gio.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				return material_gio.Body1(p.Theme, p.Label).Layout(gtx)
			})
		}),
		layout_gio.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				return layout_gio.Stack{}.Layout(gtx,
					layout_gio.Expanded(func(gtx C) D {
						return rectAbs(gtx, gtx.Constraints.Min.X, gtx.Constraints.Min.Y, color.NRGBA{R: 230, G: 230, B: 230, A: 255})
					}),
					layout_gio.Stacked(func(gtx C) D {
						return layout_gio.UniformInset(unit_gio.Dp(2)).Layout(gtx, func(gtx C) D {
							return layout_gio.Flex{Alignment: layout_gio.Baseline}.Layout(gtx,
								layout_gio.Rigid(func(gtx C) D {
									label := material_gio.Body1(p.Theme, "#")
									label.Font.Variant = "Mono"
									return label.Layout(gtx)
								}),
								layout_gio.Rigid(func(gtx C) D {
									editor := material_gio.Editor(p.Theme, &p.Editor, "rrggbb")
									editor.Font.Variant = "Mono"
									return editor.Layout(gtx)
								}),
							)
						})
					}),
				)
			})
		}),
	)
	return dims
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (p *GoPickerObj) layoutSliders(gtx C) D {
	return layout_gio.Flex{Axis: layout_gio.Vertical}.Layout(gtx,
		layout_gio.Rigid(func(gtx C) D {
			return p.layoutSlider(gtx, &p.R, "R:", valueString(p.Red()))
		}),
		layout_gio.Rigid(func(gtx C) D {
			return p.layoutSlider(gtx, &p.G, "G:", valueString(p.Green()))
		}),
		layout_gio.Rigid(func(gtx C) D {
			return p.layoutSlider(gtx, &p.B, "B:", valueString(p.Blue()))
		}),
		layout_gio.Rigid(func(gtx C) D {
			return p.layoutSlider(gtx, &p.A, "A:", valueString(p.Alpha()))
		}),
	)
}

func valueString(in uint8) string {
	s := strconv.Itoa(int(in))
	delta := 3 - len(s)
	if delta > 0 {
		s = strings.Repeat(" ", delta) + s
	}
	return s
}

func (p *GoPickerObj) layoutSlider(gtx C, value *widget_gio.Float, label, valueStr string) D {
	inset := layout_gio.UniformInset(unit_gio.Dp(2))
	layoutDims := layout_gio.Flex{}.Layout(gtx,
		layout_gio.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				label := material_gio.Body1(p.Theme, label)
				label.Font.Variant = "Mono"
				return label.Layout(gtx)
			})
		}),
		layout_gio.Flexed(1, func(gtx C) D {
			sliderDims := inset.Layout(gtx, material_gio.Slider(p.Theme, value, 0, 1).Layout)
			return sliderDims
		}),
		layout_gio.Rigid(func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				label := material_gio.Body1(p.Theme, valueStr)
				label.Font.Variant = "Mono"
				return label.Layout(gtx)
			})
		}),
	)
	return layoutDims
}

func (ob *GoPickerObj) ObjectType() (string) {
	return "GoPickerObj"
}
