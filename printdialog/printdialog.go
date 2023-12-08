/* printdialog.go */

package printdialog

import (
	filepath_go "path/filepath"
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
	Dir int = iota
	File
)

func Print(parent ui.GoObject, openPath string, label string) (action int, selectedPath string) {
	modalDialog = ui.GoModalWindow("GoFileDialog", "Open")
	modalDialog.SetSize(600, 400)
	if openPath == "" {
		openPath = "/"
	}
	
	GoOpenFile(modalDialog.Layout(), openPath, "Open File Dialog", "")
	ui.GoApp.AddWindow(modalDialog)
	action, selectedPath = modalDialog.ShowModal()
	return
}

func GetSaveFileName(parent ui.GoObject, savePath string, label string) (action int, selectedPath string) {
	modalDialog = ui.GoModalWindow("GoFileDialog", "Save As")
	modalDialog.SetSize(600, 400)
	if savePath == "" {
		savePath = "/"
	}
	GoSaveFile(modalDialog.Layout(), savePath, "Save File Dialog", "")
	ui.GoApp.AddWindow(modalDialog)
	action, selectedPath = modalDialog.ShowModal()
	return
}



// GoOpenFile creates an open file dialog.
func GoPrint(parent ui.GoObject, label string) (hObj *GoFileDialogObj) {
	//theme := ui.GoApp.Theme()
	//theme := GoTheme(gofont_gio.Collection())
	openDir, openFile := filepath_go.Split(openPath)
	object := ui.GioObject{parent, parent.ParentWindow(), []ui.GoObject{}, ui.GetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)}
	widget := ui.GioWidget{
		GoBorder: ui.GoBorder{ui.BorderNone, ui.Color_Black, 0, 0, 0},
		GoMargin: ui.GoMargin{0,0,0,0},
		GoPadding: ui.GoPadding{0,0,0,0},
		GoSize: ui.GoSize{100, 100, 200, 200, 1000, 1000},
		FocusPolicy: ui.NoFocus,
		Visible: true,
	}
	
	hFileDialog := &GoFileDialogObj{
		GioObject: object,
		GioWidget: widget,
		layout: nil,
		ActionBar: nil,
		FileLayout: nil,
		CommandBar: nil,
		FilePath: nil,
		dirPath: openDir,
		filePath: openDir,
		rootPath: rootPath,
		selectedFile: openFile,
		showHiddenFiles: false,
	}
	hFileDialog.layout = ui.GoVFlexBoxLayout(hFileDialog)
	hFileDialog.ActionBar = ui.GoHFlexBoxLayout(hFileDialog.layout)
	hFileDialog.FileLayout = ui.GoVFlexBoxLayout(hFileDialog.layout)
	hFileDialog.CommandBar = ui.GoVFlexBoxLayout(hFileDialog.layout)
	
	// ***************** Action Bar ********************* //

	hFileDialog.ActionBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	hFileDialog.ActionBar.SetPadding(4, 4, 4, 4)

		cmdBack := ui.GoIconButton(hFileDialog.ActionBar, ui.GoIcon(icons.NavigationArrowBack))
		cmdBack.SetOnClick(hFileDialog.ActionSelectionBack)
		cmdForward := ui.GoIconButton(hFileDialog.ActionBar, ui.GoIcon(icons.NavigationArrowForward))
		cmdForward.SetOnClick(hFileDialog.ActionSelectionForward)
		cmdUp := ui.GoIconButton(hFileDialog.ActionBar, ui.GoIcon(icons.NavigationArrowUpward))
		cmdUp.SetOnClick(hFileDialog.ActionSelectionUp)

		hFileDialog.FilePath = ui.GoHFlexBoxLayout(hFileDialog.ActionBar)
		hFileDialog.FilePath.SetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)
		hFileDialog.FilePath.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)

			ui.GoIconButton(hFileDialog.FilePath, ui.GoIcon(icons.FileFolder))

			hFileDialog.lblFilePath = ui.GoLabel(hFileDialog.FilePath, hFileDialog.dirPath + " > ")
			hFileDialog.lblFilePath.SetSelectable(true)
			hFileDialog.lblFilePath.SetMargin(8,8,8,8)
			hFileDialog.lblFilePath.SetSizePolicy(ui.FixedWidth, ui.ExpandingHeight)
			spacer3 := ui.GoSpacer(hFileDialog.FilePath, 20)
			spacer3.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
			ui.GoIconButton(hFileDialog.FilePath, ui.GoIcon(icons.NavigationRefresh))

	// ***************** Command Bar ********************* //

	hFileDialog.CommandBar.SetBackgroundColor(ui.Color_WhiteSmoke)
	hFileDialog.CommandBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	hFileDialog.CommandBar.SetPadding(10, 10, 10, 10)
	
		selectionLayout := ui.GoHFlexBoxLayout(hFileDialog.CommandBar)
		selectionLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)

			spacer1 := ui.GoSpacer(selectionLayout, 50)
			spacer1.SetSizePolicy(ui.FixedWidth, ui.FixedHeight)

			/*lblFileName := */ui.GoLabel(selectionLayout, "File name:")

			hFileDialog.lblSelectedName = ui.GoLabel(selectionLayout, "")
			hFileDialog.lblSelectedName.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
			hFileDialog.lblSelectedName.SetBackgroundColor(ui.Color_White)
			hFileDialog.lblSelectedName.SetMargin(10, 0, 0, 0)
			hFileDialog.lblSelectedName.SetPadding(4, 4, 4, 4)
			hFileDialog.lblSelectedName.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)

			hFileDialog.lblSelectionFilter = ui.GoLabel(selectionLayout, "All Files (*.*)")
			hFileDialog.lblSelectionFilter.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
			hFileDialog.lblSelectionFilter.SetBackgroundColor(ui.Color_White)
			hFileDialog.lblSelectionFilter.SetMargin(10, 0, 0, 0)
			hFileDialog.lblSelectionFilter.SetPadding(4, 4, 4, 4)
			hFileDialog.lblSelectionFilter.SetMinWidth(200)
			hFileDialog.lblSelectionFilter.SetSizePolicy(ui.FixedWidth, ui.FixedHeight)

		actionLayout := ui.GoHFlexBoxLayout(hFileDialog.CommandBar)
		actionLayout.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
		actionLayout.SetPadding(0, 10, 0, 0)

			spacer2 := ui.GoSpacer(actionLayout, 0)
			spacer2.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
			action := ui.GoButton(actionLayout, "Open")
			action.SetOnClick(hFileDialog.ActionOpenDialog)
			action.SetMargin(10,0,0,0)
			cancel := ui.GoButton(actionLayout, "Cancel")
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