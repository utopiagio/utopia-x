/* messagebox.go */

package messagebox

import (
	"fmt"
	"image"
	"log"
	"os"
	//"encoding/hex"
	//"strconv"
	"strings"
	ui "github.com/utopiagio/utopia"

	clip_gio "github.com/utopiagio/gio/op/clip"
	layout_gio "github.com/utopiagio/gio/layout"
	paint_gio "github.com/utopiagio/gio/op/paint"
	text_gio "github.com/utopiagio/gio/text"
	unit_gio "github.com/utopiagio/gio/unit"

	icons "golang.org/x/exp/shiny/materialdesign/icons"
)

var modalDialog *ui.GoWindowObj

const (
	NoIcon byte[] = []byte{}
	Question byte[] = icons.FileFolder
	ShowInformation byte[]
	Warning byte[]
	Critical byte[]
)

func ShowInformation(parent ui.GoObject, caption string, text string) (action int, selectedPath string) {
	modalDialog = ui.GoModalWindow("GoMessageDialog", "Information")
	modalDialog.SetSize(600, 400)
		
	GoMessageBox(modalDialog.Layout(), caption, text, "")
	ui.GoApp.AddWindow(modalDialog)
	action, _ = modalDialog.ShowModal()
	return
}





// GoMessageBox creates a messagebox dialog.
func GoMessageBox(parent ui.GoObject, caption string, text string) (hObj *GoMessageBoxObj) {
	//theme := ui.GoApp.Theme()
	//theme := GoTheme(gofont_gio.Collection())
	object := ui.GioObject{parent, parent.ParentWindow(), []ui.GoObject{}, ui.GetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)}
	widget := ui.GioWidget{
		GoBorder: ui.GoBorder{ui.BorderNone, ui.Color_Black, 0, 0, 0},
		GoMargin: ui.GoMargin{0,0,0,0},
		GoPadding: ui.GoPadding{0,0,0,0},
		GoSize: ui.GoSize{100, 100, 200, 200, 1000, 1000},
		FocusPolicy: ui.NoFocus,
		Visible: true,
	}
	
	hMessageBox := &GoMessageBoxObj{
		GioObject: object,
		GioWidget: widget,
		layout: nil,
		ActionBar: nil,
		MessageLayout: nil,
		CommandBar: nil,
		
		
		
	}
	hMessageBox.layout = ui.GoVFlexBoxLayout(hMessageBox)
	hMessageBox.ActionBar = ui.GoHFlexBoxLayout(hMessageBox.layout)
	hMessageBox.MessageLayout = ui.GoVFlexBoxLayout(hMessageBox.layout)
	hMessageBox.CommandBar = ui.GoVFlexBoxLayout(hMessageBox.layout)
	
	// ***************** Action Bar ********************* //

	hMessageBox.ActionBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	hMessageBox.ActionBar.SetPadding(4, 4, 4, 4)

	ui.GoSpacer(selectionLayout, 50)
	/*cmdInfo := */ui.GoIconButton(hFileDialog.ActionBar, ui.GoIcon(icons.NavigationArrowBack))
	//cmdBack.SetOnClick(hFileDialog.ActionSelectionBack)

	ui.GoSpacer(selectionLayout, 50)
	lblCaption = ui.GoLabel(selectionLayout, caption)		

	// ***************** Command Bar ********************* //

	hMessageBox.CommandBar.SetBackgroundColor(ui.Color_WhiteSmoke)
	hMessageBox.CommandBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	hMessageBox.CommandBar.SetPadding(10, 10, 10, 10)
	
		selectionLayout := ui.GoHFlexBoxLayout(hFileDialog.CommandBar)
		selectionLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)

			spacer1 := ui.GoSpacer(selectionLayout, 50)
			spacer1.SetSizePolicy(ui.FixedWidth, ui.FixedHeight)

			/*lblFileName := */ui.GoLabel(selectionLayout, "File name:")

			

		actionLayout := ui.GoHFlexBoxLayout(hFileDialog.CommandBar)
		actionLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
		actionLayout.SetPadding(0, 10, 0, 0)

			spacer2 := ui.GoSpacer(actionLayout, 0)
			spacer2.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
			cancel := ui.GoButton(actionLayout, "OK")
			cancel.SetOnClick(hFileDialog.ActionCancelDialog)
			cancel.SetMargin(10,0,0,0)

	// ***************** Center Layout ********************* //

	//fileBar := ui.GoHFlexBoxLayout(hFileDialog.FileLayout)
	//fileBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	//fileBar.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	fileBar := ui.GoMenuBar(hFileDialog.FileLayout)
	//fileBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	fileBar.SetBackgroundColor(ui.Color_WhiteSmoke)
	fileBar.SetPadding(10,10,10,10)
	fileBar.Show()
	//fileBar.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	mnuView := fileBar.AddMenu("View")

	//mnuView.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	mnuView.AddItem("List")
	mnuView.AddItem("Details")

	/*button := ui.GoButton(fileBar, "Help")*/

	fileViewLayout := ui.GoHFlexBoxLayout(hFileDialog.FileLayout)
	hFileDialog.dirView = ui.GoListView(fileViewLayout)
	hFileDialog.dirView.SetLayoutMode(ui.Vertical)
	hFileDialog.dirView.SetSizePolicy(ui.FixedWidth, ui.ExpandingHeight)
	hFileDialog.dirView.SetWidth(260)
	hFileDialog.dirView.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	hFileDialog.dirView.SetOnItemClicked(hFileDialog.DirViewItemClicked)
	hFileDialog.dirView.SetOnItemDoubleClicked(hFileDialog.DirViewItemDoubleClicked)
	hFileDialog.fileView = ui.GoListView(fileViewLayout)
	hFileDialog.fileView.SetLayoutMode(ui.Vertical)
	hFileDialog.fileView.SetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)
	hFileDialog.fileView.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	hFileDialog.fileView.SetOnItemClicked(hFileDialog.FileViewItemClicked)
	//icon := ui.GoIcon(icons.NavigationArrowBack)
	
	hFileDialog.populateDirView()
	hFileDialog.history = ui.GoFileHistory(hFileDialog.dirPath)
	parent.AddControl(hFileDialog)
	return hFileDialog
}

type GoMessageBoxObj struct {
	ui.GioObject
	ui.GioWidget
	layout *ui.GoLayoutObj
	ActionBar *ui.GoLayoutObj
	FileLayout *ui.GoLayoutObj
	CommandBar *ui.GoLayoutObj
	FilePath *ui.GoLayoutObj
	dirView *ui.GoListViewObj
	fileView *ui.GoListViewObj

	lblFilePath *ui.GoLabelObj
	lblSelectionFilter *ui.GoLabelObj
	lblSelectedName *ui.GoLabelObj
	//theme *GoThemeObj
	//font text_gio.Font
	//fontSize unit_gio.Sp
	detailedText string
	font text_gio.Font
	fontSize unit_gio.Sp
	icon string
	iconButton *GoIconButtonObj
	informativeText string
	text string
	color ui.GoColor

	background ui.GoColor
	cornerRadius unit_gio.Dp
	
	dirPath string 		// dirPath is current directory path in FilePath(lblFilePath)
	filePath string 	// filePath is path of selected file in FileView
	rootPath string		// rootPath is root directory in DirView
	selectedFile string // selectedFile is name of selected file in FileView
	//currentDir string

	history *ui.GoHistoryObj
	//Label *ui.GoLabelObj
	showHiddenFiles bool
}


func (ob *GoMessageBoxObj) ObjectType() (string) {
	return "GoMessageBoxObj"
}

func (ob *GoMessageBoxObj) ShowHiddenFiles(show bool) {
	ob.showHiddenFiles = show
}

func (ob *GoMessageBoxObj) Widget() (*ui.GioWidget) {
	return &ob.GioWidget
}

func (ob *GoMessageBoxObj) ActionCancel() {
	log.Println("ActionCancel_Clicked().......")
	ob.ParentWindow().ModalAction = -1
	ob.ParentWindow().ModalInfo = ""
	ob.ParentWindow().Close()

}

func (ob *GoMessageBoxObj)ActionOK() {
	log.Println("ActionOK_Clicked().......")
	var path string = ob.filePath
	//if path == "/" {path = ""}
	ob.ParentWindow().ModalAction = 1
	ob.ParentWindow().ModalInfo = ""
	ob.ParentWindow().Close()
}