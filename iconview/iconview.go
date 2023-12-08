/* iconview.go */

package iconview

import (
	//"log"
	//"encoding/hex"
	//"strconv"
	//"strings"
	ui "github.com/utopiagio/utopia"
	//layout_gio "github.com/utopiagio/gio/layout"
	//pointer_gio "github.com/utopiagio/gio/io/pointer"
	//text_gio "github.com/utopiagio/gio/text"
	//unit_gio "github.com/utopiagio/gio/unit"
)



// GoIconView creates an icon viewer.
func GoIconView(parent ui.GoObject, label string) (hObj *GoIconViewObj) {
	//theme := ui.GoApp.Theme()
	//theme := GoTheme(gofont_gio.Collection())
	object := ui.GioObject{parent, parent.ParentWindow(), []ui.GoObject{}, ui.GetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)}
	widget := ui.GioWidget{
		GoBorder: ui.GoBorder{ui.BorderNone, ui.Color_Black, 0, 0},
		GoMargin: ui.GoMargin{0,0,0,0},
		GoPadding: ui.GoPadding{0,0,0,0},
		FocusPolicy: ui.NoFocus,
		Visible: true,
	}

	hIconView := &GoIconViewObj{
		GioObject: object,
		GioWidget: widget,
		layout: nil,
		hdrLabel: nil,
	}
	hIconView.layout = ui.GoVFlexBoxLayout(hIconView)
	hIconView.hdrLabel = ui.H3Label(hIconView.layout, hdrLabel)
	return hIconView
}

type GoIconViewObj struct {
	ui.GioObject
	ui.GioWidget
	layout *ui.GoLayoutObj
	//theme *GoThemeObj
	//font text_gio.Font
	//fontSize unit_gio.Sp
	
	hdrLabel *ui.GoLabelObj
}

func (ob *GoIconViewObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	dims = ui.layout_gio.Dimensions {Size: gtx.Constraints.Max,}
	//log.Println("gtx.Constraints.Max: ", dims)
	if ob.Visible {
	//margin := layout_gio.Inset(ob.margin.Left)
		dims = ob.GoMargin.Layout(gtx, func(gtx ui.C) ui.D {
			borderDims := ob.GoBorder.Layout(gtx, func(gtx ui.C) ui.D {
				paddingDims := ob.GoPadding.Layout(gtx, func(gtx ui.C) ui.D {
					return ob.Layout(gtx)
				})
				//log.Println("PaddingDims: ", paddingDims)
				return paddingDims
			})
			//log.Println("BorderDims: ", borderDims)
			return borderDims
		})
		ob.dims = dims
		//log.Println("ButtonDims: ", dims)
		ob.Width = (int(float32(dims.Size.X) / GoDpr))
		ob.Height = (int(float32(dims.Size.Y) / GoDpr))
	}
	return dims
}

func (ob *GoIconViewObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.ReceiveEvents(gtx)
	width := gtx.Dp(unit_gio.Dp(ob.Width))
	height := gtx.Dp(unit_gio.Dp(ob.Height))
	if ob.SizePolicy().HFlex {
		width = gtx.Constraints.Max.X
	}
	if ob.SizePolicy().VFlex {
		height = gtx.Constraints.Max.Y
	}
	dims := image.Point{X: width, Y: height}
	rr := gtx.Dp(ob.cornerRadius)
	defer clip_gio.UniformRRect(image.Rectangle{Max: dims}, rr).Push(gtx.Ops).Pop()
	// paint background
	background := ob.background.NRGBA()
	paint_gio.Fill(gtx.Ops, background)
	// add the events handler to receive widget pointer events
	ob.SignalEvents(gtx)
	ob.layout.Draw(gtx)
	return layout_gio.Dimensions{Size: dims}
}