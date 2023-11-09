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
	layout := ui.GoVFlexBoxLayout(parent)
	header := ui.H3Label(layout, label)
	hIconView := &GoIconViewObj{
		GioObject: object,
		GioWidget: widget,
		Layout: layout,
		Label: header,
	}
		
	return hIconView
}

type GoIconViewObj struct {
	ui.GioObject
	ui.GioWidget
	Layout *ui.GoLayoutObj
	//theme *GoThemeObj
	//font text_gio.Font
	//fontSize unit_gio.Sp
	
	Label *ui.GoLabelObj
}