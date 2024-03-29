/* filedialog.go */

package filedialog

import (
	filepath_go "path/filepath"
	//"fmt"
	"image"
	"log"
	"os"
	"runtime"
	//"encoding/hex"
	//"strconv"
	"strings"
	ui "github.com/utopiagio/utopia"
	"github.com/utopiagio/utopia/metrics"
	"github.com/utopiagio/utopia/history"
	

	clip_gio "github.com/utopiagio/gio/op/clip"
	font_gio "github.com/utopiagio/gio/font"
	layout_gio "github.com/utopiagio/gio/layout"

	paint_gio "github.com/utopiagio/gio/op/paint"
	unit_gio "github.com/utopiagio/gio/unit"


	icons "golang.org/x/exp/shiny/materialdesign/icons"
)

const (
	Dir int = iota
	File
)

func GetOpenFileName(parent ui.GoObject, openPath string, label string) (action int, selectedPath string) {
	var modalDialog *ui.GoWindowObj
	modalDialog = ui.GoModalWindow("GoFileDialog", "Open")
	modalDialog.SetSize(600, 400)
	if openPath == "" {
		openPath = "/"
	}
	
	GoOpenFile(modalDialog.Layout(), openPath, "Open File Dialog", "")
	action, selectedPath = modalDialog.ShowModal()
	return
}

func GetSaveFileName(parent ui.GoObject, savePath string, label string) (action int, selectedPath string) {
	var modalDialog *ui.GoWindowObj
	modalDialog = ui.GoModalWindow("GoFileDialog", "Save As")
	modalDialog.SetSize(600, 400)
	if savePath == "" {
		savePath = "/"
	}
	GoSaveFile(modalDialog.Layout(), savePath, "Save File Dialog", "")
	action, selectedPath = modalDialog.ShowModal()
	return
}



// GoOpenFile creates an open file dialog.
func GoOpenFile(parent ui.GoObject, openPath string, label string, rootPath string) (hObj *GoFileDialogObj) {
	//theme := ui.GoApp.Theme()
	//theme := GoTheme(gofont_gio.Collection())
	openDir, openFile := filepath_go.Split(openPath)
	object := ui.GioObject{parent, parent.ParentWindow(), []ui.GoObject{}, ui.GetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)}
	widget := ui.GioWidget{
		GoBorder: ui.GoBorder{ui.BorderNone, ui.Color_Black, 0, 0, 0},
		GoMargin: ui.GoMargin{0,0,0,0},
		GoPadding: ui.GoPadding{0,0,0,0},
		GoSize: ui.GoSize{100, 100, 200, 200, 1000, 1000, 200, 200},
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

	hFileDialog.ActionBar.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
	hFileDialog.ActionBar.SetPadding(4, 4, 4, 4)

		cmdBack := ui.GoIconVGButton(hFileDialog.ActionBar, ui.GoIconVG(icons.NavigationArrowBack))
		cmdBack.SetOnClick(hFileDialog.ActionSelectionBack)
		cmdForward := ui.GoIconVGButton(hFileDialog.ActionBar, ui.GoIconVG(icons.NavigationArrowForward))
		cmdForward.SetOnClick(hFileDialog.ActionSelectionForward)
		cmdUp := ui.GoIconVGButton(hFileDialog.ActionBar, ui.GoIconVG(icons.NavigationArrowUpward))
		cmdUp.SetOnClick(hFileDialog.ActionSelectionUp)

		hFileDialog.FilePath = ui.GoHFlexBoxLayout(hFileDialog.ActionBar)
		hFileDialog.FilePath.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
		hFileDialog.FilePath.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)

			ui.GoIconVGButton(hFileDialog.FilePath, ui.GoIconVG(icons.FileFolder))
			path := hFileDialog.dirPath
			if runtime.GOOS == "windows" {
				path = "C:" + hFileDialog.dirPath
			}
			hFileDialog.lblFilePath = ui.GoLabel(hFileDialog.FilePath, path + " > ")
			hFileDialog.lblFilePath.SetSelectable(true)
			hFileDialog.lblFilePath.SetMargin(8,4,8,2)
			hFileDialog.lblFilePath.SetSizePolicy(ui.FixedWidth, ui.PreferredHeight)
			spacer3 := ui.GoSpacer(hFileDialog.FilePath, 20)
			spacer3.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
			ui.GoIconVGButton(hFileDialog.FilePath, ui.GoIconVG(icons.NavigationRefresh))

	// ***************** Command Bar ********************* //

	hFileDialog.CommandBar.SetBackgroundColor(ui.Color_WhiteSmoke)
	hFileDialog.CommandBar.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
	hFileDialog.CommandBar.SetPadding(10, 10, 10, 10)
	
		selectionLayout := ui.GoHFlexBoxLayout(hFileDialog.CommandBar)
		selectionLayout.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)

			ui.GoSpacer(selectionLayout, 50)
			/*lblFileName := */ui.GoLabel(selectionLayout, "File name:")

			hFileDialog.lblSelectedName = ui.GoLabel(selectionLayout, "")
			hFileDialog.lblSelectedName.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
			hFileDialog.lblSelectedName.SetBackgroundColor(ui.Color_White)
			hFileDialog.lblSelectedName.SetMargin(10, 0, 0, 0)
			hFileDialog.lblSelectedName.SetPadding(4, 4, 4, 4)
			hFileDialog.lblSelectedName.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)

			hFileDialog.lblSelectionFilter = ui.GoLabel(selectionLayout, "All Files (*.*)")
			hFileDialog.lblSelectionFilter.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
			hFileDialog.lblSelectionFilter.SetBackgroundColor(ui.Color_White)
			hFileDialog.lblSelectionFilter.SetMargin(10, 0, 0, 0)
			hFileDialog.lblSelectionFilter.SetPadding(4, 4, 4, 4)
			hFileDialog.lblSelectionFilter.SetWidth(160)
			hFileDialog.lblSelectionFilter.SetSizePolicy(ui.FixedWidth, ui.PreferredHeight)

		actionLayout := ui.GoHFlexBoxLayout(hFileDialog.CommandBar)
		actionLayout.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
		actionLayout.SetPadding(0, 10, 0, 0)

			spacer2 := ui.GoSpacer(actionLayout, 0)
			spacer2.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
			action := ui.GoButton(actionLayout, "Open")
			action.SetOnClick(hFileDialog.ActionOpenDialog)
			action.SetMargin(10,0,0,0)
			cancel := ui.GoButton(actionLayout, "Cancel")
			cancel.SetOnClick(hFileDialog.ActionCancelDialog)
			cancel.SetMargin(10,0,0,0)

	// ***************** MenuBar Layout ********************* //

	//fileBar := ui.GoHFlexBoxLayout(hFileDialog.FileLayout)
	//fileBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	//fileBar.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	fileBar := ui.GoMenuBar(hFileDialog.FileLayout)
	//fileBar.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
	fileBar.SetBackgroundColor(ui.Color_WhiteSmoke)
	fileBar.SetPadding(0,0,0,0)
	fileBar.Show()
	//fileBar.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	mnuView := fileBar.AddMenu("View")

	//mnuView.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	mnuView.AddItem("List")
	mnuView.AddItem("Details")

	/*button := ui.GoButton(fileBar, "Help")*/

	// ***************** Center Layout ********************* //

	fileViewLayout := ui.GoHFlexBoxLayout(hFileDialog.FileLayout)
	hFileDialog.dirView = ui.GoListView(fileViewLayout)
	hFileDialog.dirView.SetLayoutMode(ui.Vertical)
	hFileDialog.dirView.SetSizePolicy(ui.FixedWidth, ui.PreferredHeight)
	hFileDialog.dirView.SetWidth(260)
	hFileDialog.dirView.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	hFileDialog.dirView.SetOnItemClicked(hFileDialog.DirViewItem_Clicked)
	hFileDialog.dirView.SetOnItemDoubleClicked(hFileDialog.DirViewItem_DoubleClicked)
	hFileDialog.fileView = ui.GoListView(fileViewLayout)
	hFileDialog.fileView.SetLayoutMode(ui.Vertical)
	hFileDialog.fileView.SetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)
	hFileDialog.fileView.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	hFileDialog.fileView.SetOnItemClicked(hFileDialog.FileViewItem_Clicked)
	hFileDialog.fileView.SetOnItemDoubleClicked(hFileDialog.FileViewItem_DoubleClicked)
	//icon := ui.GoIconVG(icons.NavigationArrowBack)
	
	hFileDialog.populateDirView()
	hFileDialog.history = history.GoFileHistory(hFileDialog.dirPath)
	parent.AddControl(hFileDialog)
	return hFileDialog
}

// GoOpenFile creates an open file dialog.
func GoSaveFile(parent ui.GoObject, savePath string, label string, rootPath string) (hObj *GoFileDialogObj) {
	//theme := ui.GoApp.Theme()
	//theme := GoTheme(gofont_gio.Collection())
	saveDir, saveFile := filepath_go.Split(savePath)
	object := ui.GioObject{parent, parent.ParentWindow(), []ui.GoObject{}, ui.GetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)}
	widget := ui.GioWidget{
		GoBorder: ui.GoBorder{ui.BorderNone, ui.Color_Black, 0, 0, 0},
		GoMargin: ui.GoMargin{0,0,0,0},
		GoPadding: ui.GoPadding{0,0,0,0},
		GoSize: ui.GoSize{100, 100, 200, 200, 1000, 1000, 200, 200},
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
		dirPath: saveDir,
		filePath: saveDir,
		rootPath: rootPath,
		selectedFile: saveFile,
		showHiddenFiles: false,
	}
	hFileDialog.layout = ui.GoVFlexBoxLayout(hFileDialog)
	hFileDialog.ActionBar = ui.GoHFlexBoxLayout(hFileDialog.layout)
	hFileDialog.FileLayout = ui.GoVFlexBoxLayout(hFileDialog.layout)
	hFileDialog.CommandBar = ui.GoVFlexBoxLayout(hFileDialog.layout)
	
	// ***************** Action Bar ********************* //

	hFileDialog.ActionBar.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
	hFileDialog.ActionBar.SetPadding(4, 4, 4, 4)

		cmdBack := ui.GoIconVGButton(hFileDialog.ActionBar, ui.GoIconVG(icons.NavigationArrowBack))
		cmdBack.SetOnClick(hFileDialog.ActionSelectionBack)
		cmdForward := ui.GoIconVGButton(hFileDialog.ActionBar, ui.GoIconVG(icons.NavigationArrowForward))
		cmdForward.SetOnClick(hFileDialog.ActionSelectionForward)
		cmdUp := ui.GoIconVGButton(hFileDialog.ActionBar, ui.GoIconVG(icons.NavigationArrowUpward))
		cmdUp.SetOnClick(hFileDialog.ActionSelectionUp)
		
		hFileDialog.FilePath = ui.GoHFlexBoxLayout(hFileDialog.ActionBar)
		hFileDialog.FilePath.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
		hFileDialog.FilePath.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)

			ui.GoIconVGButton(hFileDialog.FilePath, ui.GoIconVG(icons.FileFolder))
			path := hFileDialog.dirPath
			if runtime.GOOS == "windows" {
				path = "C:" + hFileDialog.dirPath
			}
			hFileDialog.lblFilePath = ui.GoLabel(hFileDialog.FilePath, path + " > ")
			hFileDialog.lblFilePath.SetSelectable(true)
			hFileDialog.lblFilePath.SetMargin(8,4,8,4)
			hFileDialog.lblFilePath.SetSizePolicy(ui.FixedWidth, ui.PreferredHeight)
			spacer3 := ui.GoSpacer(hFileDialog.FilePath, 20)
			spacer3.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
			ui.GoIconVGButton(hFileDialog.FilePath, ui.GoIconVG(icons.NavigationRefresh))

	// ***************** Command Bar ********************* //

	hFileDialog.CommandBar.SetBackgroundColor(ui.Color_WhiteSmoke)
	hFileDialog.CommandBar.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
	hFileDialog.CommandBar.SetPadding(10, 10, 10, 10)
	
		selectionLayout := ui.GoHFlexBoxLayout(hFileDialog.CommandBar)
		selectionLayout.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)

			spacer1 := ui.GoSpacer(selectionLayout, 50)
			spacer1.SetSizePolicy(ui.FixedWidth, ui.FixedHeight)

			/*lblFileName := */ui.GoLabel(selectionLayout, "File name:")

			hFileDialog.lblSelectedName = ui.GoLabel(selectionLayout, "")
			hFileDialog.lblSelectedName.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
			hFileDialog.lblSelectedName.SetBackgroundColor(ui.Color_White)
			hFileDialog.lblSelectedName.SetMargin(10, 0, 0, 0)
			hFileDialog.lblSelectedName.SetPadding(4, 4, 4, 4)
			hFileDialog.lblSelectedName.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)

			hFileDialog.lblSelectionFilter = ui.GoLabel(selectionLayout, "All Files (*.*)")
			hFileDialog.lblSelectionFilter.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
			hFileDialog.lblSelectionFilter.SetBackgroundColor(ui.Color_White)
			hFileDialog.lblSelectionFilter.SetMargin(10, 0, 0, 0)
			hFileDialog.lblSelectionFilter.SetPadding(4, 4, 4, 4)
			hFileDialog.lblSelectionFilter.SetMinWidth(200)
			hFileDialog.lblSelectionFilter.SetSizePolicy(ui.FixedWidth, ui.PreferredHeight)

		actionLayout := ui.GoHFlexBoxLayout(hFileDialog.CommandBar)
		actionLayout.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
		actionLayout.SetPadding(0, 10, 0, 0)

			spacer2 := ui.GoSpacer(actionLayout, 0)
			spacer2.SetSizePolicy(ui.ExpandingWidth, ui.FixedHeight)
			action := ui.GoButton(actionLayout, "Save")
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
	//fileBar.SetSizePolicy(ui.ExpandingWidth, ui.FixededHeight)
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
	hFileDialog.dirView.SetSizePolicy(ui.FixedWidth, ui.ExpandingHeight)
	hFileDialog.dirView.SetLayoutMode(ui.Vertical)
	hFileDialog.dirView.SetWidth(260)
	hFileDialog.dirView.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	hFileDialog.dirView.SetOnItemClicked(hFileDialog.DirViewItem_Clicked)
	hFileDialog.dirView.SetOnItemDoubleClicked(hFileDialog.DirViewItem_DoubleClicked)
	hFileDialog.fileView = ui.GoListView(fileViewLayout)
	hFileDialog.fileView.SetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)
	hFileDialog.fileView.SetLayoutMode(ui.Vertical)
	hFileDialog.fileView.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
	hFileDialog.fileView.SetOnItemClicked(hFileDialog.FileViewItem_Clicked)
	
	//icon := ui.GoIconVG(icons.NavigationArrowBack)

	hFileDialog.populateDirView()
	hFileDialog.populateFileView(hFileDialog.dirPath)
 	hFileDialog.history = history.GoFileHistory(hFileDialog.dirPath)
	parent.AddControl(hFileDialog)
	return hFileDialog
}


type GoFileDialogObj struct {
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
	font font_gio.Font
	fontSize unit_gio.Sp
	text string
	color ui.GoColor

	background ui.GoColor
	cornerRadius unit_gio.Dp
	
	dirPath string 		// dirPath is current directory path in FilePath(lblFilePath)
	filePath string 	// filePath is path of selected file in FileView
	rootPath string		// rootPath is root directory in DirView
	selectedFile string // selectedFile is name of selected file in FileView
	currentDirItem *ui.GoListViewItemObj

	history *history.GoHistoryObj
	//Label *ui.GoLabelObj
	showHiddenFiles bool
}

func (ob *GoFileDialogObj) Draw(gtx layout_gio.Context) (dims layout_gio.Dimensions) {
	cs := gtx.Constraints
	width := metrics.DpToPx(ui.GoDpr, ob.Width)
	height := metrics.DpToPx(ui.GoDpr, ob.Height)
	minWidth := metrics.DpToPx(ui.GoDpr, ob.MinWidth)
	minHeight := metrics.DpToPx(ui.GoDpr, ob.MinHeight)
	maxWidth := metrics.DpToPx(ui.GoDpr, ob.MaxWidth)
	maxHeight := metrics.DpToPx(ui.GoDpr, ob.MaxHeight)
	
	switch ob.SizePolicy().Horiz {
	case ui.FixedWidth:			// SizeHint is Fixed
		w := min(maxWidth, width)			// constrain to ob.MaxWidth
		cs.Min.X = max(minWidth, w)				// constrain to ob.MinWidth 
		cs.Max.X = cs.Min.X						// set to cs.Min.X
	case ui.MinimumWidth:			// SizeHint is Minimum
		cs.Min.X = minWidth						// set to ob.MinWidth
		cs.Max.X = minWidth						// set to ob.MinWidth
	case ui.PreferredWidth:		// SizeHint is Preferred
		cs.Min.X = minWidth						// constrain to ob.MinWidth
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
	case ui.MaximumWidth:			// SizeHint is Maximum
		cs.Max.X = maxWidth						// set to ob.MaxWidth
		cs.Min.X = maxWidth						// set to ob.MaxWidth
	case ui.ExpandingWidth:
		cs.Max.X = min(cs.Max.X, maxWidth)		// constrain to ob.MaxWidth
		cs.Min.X = cs.Max.X						// set to cs.Max.X
	}

	switch ob.SizePolicy().Vert {
	case ui.FixedHeight:			// SizeHint is Fixed 
		w := min(maxHeight, height)				// constrain to ob.MaxHeight
		cs.Min.Y = max(minHeight, w)			// constrain to ob.MinHeight 
		cs.Max.Y = cs.Min.Y						// set to cs.Min.Y
	case ui.MinimumHeight:			// SizeHint is Minimum
		cs.Min.Y = minHeight					// set to ob.MinHeight
		cs.Max.Y = minHeight					// set to ob.MinHeight
	case ui.PreferredHeight:		// SizeHint is Preferred
		cs.Min.Y = max(0, minHeight)			// constrain to ob.MinHeight
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
	case ui.MaximumHeight:			// SizeHint is Maximum
		cs.Max.Y = maxHeight					// set to ob.MaxHeight
		cs.Min.Y = maxHeight					// set to ob.MaxHeight
	case ui.ExpandingHeight:
		cs.Max.Y = min(cs.Max.Y, maxHeight)		// constrain to ob.MaxHeight
		cs.Min.Y = cs.Max.Y						// set to cs.Max.Y
	}

	gtx.Constraints = cs
	dims = layout_gio.Dimensions {Size: gtx.Constraints.Min,}
	if ob.Visible {
		dims = ob.GoMargin.Layout(gtx, func(gtx ui.C) ui.D {
			borderDims := ob.GoBorder.Layout(gtx, func(gtx ui.C) ui.D {
				paddingDims := ob.GoPadding.Layout(gtx, func(gtx ui.C) ui.D {
					return ob.Layout(gtx)
				})
				return paddingDims
			})
			return borderDims
		})
		//ob.Dims = dims
		ob.AbsWidth = (int(float32(dims.Size.X) / ui.GoDpr))
		ob.AbsHeight = (int(float32(dims.Size.Y) / ui.GoDpr))
	}
	return dims
}

func (ob *GoFileDialogObj) Layout(gtx layout_gio.Context) layout_gio.Dimensions {
	ob.ReceiveEvents(gtx, nil)
	width := gtx.Constraints.Max.X
	height := gtx.Constraints.Max.Y
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

func (ob *GoFileDialogObj) DirViewItem_Clicked(nodeId []int) {
	var path string = "/" + ob.rootPath 	// ob.rootPath can = ("")!
	if len(path) > 1 {
		if path[0:2] == "//" {path = path[1:]}
	}
	_, path = ob.getDirIndex(nodeId, path)
	ob.dirPath = path
	ob.filePath = path
	if runtime.GOOS == "windows" {
		ob.lblFilePath.SetText("C:" + path)
	} else {
		ob.lblFilePath.SetText(path)
	}
	ob.lblSelectedName.SetText("")
	ob.populateFileView(path)
}

func (ob *GoFileDialogObj) DirViewItem_DoubleClicked(nodeId []int) {
	ob.expandDirView(nodeId)
}

func (ob *GoFileDialogObj) FileViewItem_Clicked(nodeId []int) {
	var listViewItem  *ui.GoListViewItemObj
	var path string
	path = ob.dirPath
	if path == "/" {path = ""}
	listViewItem = ob.fileView.Objects()[nodeId[0]].(*ui.GoListViewItemObj)
	path = ob.filePath + listViewItem.Text()
	for level := 1; level <= len(nodeId) - 1; level++ {
		listViewItem = listViewItem.Objects()[nodeId[level]].(*ui.GoListViewItemObj)
		path = path + "/" + listViewItem.Text()
	}
	info, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	if info.Mode().IsRegular() {
		fileName := filepath_go.Base(path)
		ob.lblSelectedName.SetText(fileName)
	}
	ob.ParentWindow().Refresh()
}

func (ob *GoFileDialogObj) FileViewItem_DoubleClicked(nodeId []int) {
	if !ob.currentDirItem.IsSelected() {
		ob.currentDirItem.Expand()
	}
	idx, _ := ob.getFileIndex(nodeId)
	fileViewItem := ob.fileView.Objects()[idx].(*ui.GoListViewItemObj)
	for id := 0; id < ob.currentDirItem.ItemCount(); id++ {
		dirViewItem := ob.currentDirItem.Objects()[id].(*ui.GoListViewItemObj)
		if dirViewItem.Text() == fileViewItem.Text() {
			dirViewItem.Trigger()
		}
	}
}

func (ob *GoFileDialogObj) ObjectType() (string) {
	return "GoFileDialogObj"
}

func (ob *GoFileDialogObj) ShowHiddenFiles(show bool) {
	ob.showHiddenFiles = show
}

func (ob *GoFileDialogObj) Widget() (*ui.GioWidget) {
	return &ob.GioWidget
}

func (ob *GoFileDialogObj)ActionCancelDialog() {
	//log.Println("ActionCancel_Clicked().......")
	ob.ParentWindow().ModalAction = -1
	ob.ParentWindow().ModalInfo = ob.filePath + ob.lblSelectedName.Text()
	ob.ParentWindow().Close()

}

func (ob *GoFileDialogObj)ActionOpenDialog() {
	//log.Println("ActionOpen_Clicked().......")
	ob.ParentWindow().ModalAction = 4
	ob.ParentWindow().ModalInfo = ob.filePath + ob.lblSelectedName.Text()
	ob.ParentWindow().Close()
}


func (ob *GoFileDialogObj)ActionSelectionBack() {

}

func (ob *GoFileDialogObj)ActionSelectionForward() {
	
}

func (ob *GoFileDialogObj)ActionSelectionUp() {
	
}

/*fmt.Println("DirEntryName:", file.Name())
		fmt.Println("DirEntryIsDir:", file.IsDir())
		fmt.Println("DirEntryType:", strconv.FormatInt(int64(file.Type()), 2))
		fsInfo, _ := file.Info()
		fmt.Println("InfoName:", fsInfo.Name())
		fmt.Println("InfoSize:", fsInfo.Size())
		fmt.Println("InfoModTime:", fsInfo.ModTime())
		fmt.Println("InfoIsDir:", fsInfo.IsDir())
		fmt.Println("InfoSys:", fsInfo.Sys())
		fmt.Println("")
		fmt.Println("ModeIsDir:", fsInfo.Mode().IsDir())
		fmt.Println("ModeIsRegular:", fsInfo.Mode().IsRegular())
		fmt.Println("ModePerm:", fsInfo.Mode().Perm())
		fmt.Println("ModeString:", fsInfo.Mode().String())
		fmt.Println("ModeType:", fsInfo.Mode().Type())
		//fmt.Println(file.Info())
		fmt.Println("")*/

func (ob *GoFileDialogObj)populateDirView() {
	//fmt.Println("GoFileDialogObj::populateDirView()..............")
	var listId int
	var listIdx int
	var parentListIdx int
	var path string
	var parentPath string
	var listItem *ui.GoListViewItemObj
	var nextLevelItem *ui.GoListViewItemObj
	var parentListItem *ui.GoListViewItemObj
	var dirList []string
	//fmt.Println("ob.dirPath =", ob.dirPath)
	dirList = strings.Split(ob.dirPath, "/")
	//fmt.Println("len(dirList) =", len(dirList))
	ob.dirView.ClearList()
	path = ob.rootPath + "/"	// ob.rootPath can = ("")!
	for x := 0; x < len(dirList) - 1; x++ {
		parentListItem = nextLevelItem
		parentListIdx = listIdx
		parentPath = path
		//fmt.Println("x =", x)
		//fmt.Println("listItem =", dirList[x])
		//path = path + "/" + dirList[x]
		//fmt.Println("path =", path)
		files, err := os.ReadDir(parentPath)
		if err != nil {
			log.Fatal(err)
		}
		listId = 0
		for _, file := range files {
			hidden, err := isHidden(parentPath, file.Name())
			if err != nil {
				//log.Print(err)
				continue
			}
			if !hidden || ob.showHiddenFiles {
				if file.IsDir() {
					if x == 0 {
						listItem = ob.dirView.AddListItem(icons.FileFolder, file.Name())
						listItem.SetWidth(listItem.ListView().MaxWidth)
					} else {
						listItem = parentListItem.InsertListItem(icons.FileFolder, file.Name(), parentListIdx + listId)
						listItem.SetWidth(listItem.ListView().MaxWidth)
					}
					if x < len(dirList) - 1 {
						if file.Name() == dirList[x + 1] {
							listIdx += listId + 1
							path = path + file.Name() + "/"
							listItem.SetExpanded(true)
							nextLevelItem = listItem
						}
					}
					listId++
				}
			}
		}
	}
	ob.populateFileView(path)
	ob.ParentWindow().Refresh()
}

func (ob *GoFileDialogObj)populateFileView(filepath string) {
	ob.filePath = filepath
	files, err := os.ReadDir(filepath)
	if err != nil {
		//fmt.Println("error =", err)
		log.Fatal(err)
	}
	ob.fileView.ClearList()
	for _, file := range files {
		hidden, err := isHidden(filepath, file.Name())
		if err != nil {
			//log.Print(err)
			continue
		}
		if !hidden || ob.showHiddenFiles {
			if file.IsDir() {
				//fmt.Println("file.IsDir ", file.Name())
				ob.fileView.AddListItem(icons.FileFolder, file.Name())
			} else if file.Type() == 0 {
				//fmt.Println("file ", file.Name())
				ob.fileView.AddListItem(icons.MapsLocalOffer, file.Name())
			}
		}
	}
	ob.ParentWindow().Refresh()
}

func (ob *GoFileDialogObj)expandDirView(nodeId []int) {
	//fmt.Println("GoFileDialogObj::expandDirView().............")
	//log.Println("nodeId[] =", nodeId)
	var listViewItem  *ui.GoListViewItemObj
	var idx int
	var path string = "/" + ob.rootPath 	// ob.rootPath can = ("")!
	if len(path) > 1 {
		if path[0:2] == "//" {path = path[1:]}
	}
	idx, path = ob.getDirIndex(nodeId, path)
	ob.dirPath = path
	//log.Println("idx =", idx)
	//log.Println("dirPath =", ob.dirPath)
	listViewItem = ob.dirView.Objects()[idx].(*ui.GoListViewItemObj)
	if listViewItem.IsExpanded() {
		ob.deleteListViewItems(listViewItem, idx)
		
		listViewItem.SetExpanded(false)
	} else {
		//log.Println("ReadDir path =", ob.dirPath)
		
		files, err := os.ReadDir(ob.dirPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			hidden, err := isHidden(path, file.Name())
			if err != nil {
				//log.Print(err)
				continue
			}
			if !hidden || ob.showHiddenFiles {
				if file.IsDir() {
					idx++
					listItem := listViewItem.InsertListItem(icons.FileFolder, file.Name(), idx)
					listItem.SetWidth(listItem.ListView().MaxWidth)
					listViewItem.SetExpanded(true)
					listViewItem.ClearFocus()
				}
			}
		}
	}
	ob.ParentWindow().Refresh()
}

func (ob *GoFileDialogObj) deleteListViewItems(listViewItem *ui.GoListViewItemObj, idx int) {
	if listViewItem.IsExpanded() {
		listViewItemSize := len(listViewItem.Objects())
		for x := 0; x < listViewItemSize; x++ {
			item := listViewItem.Objects()[0]
			if item.(*ui.GoListViewItemObj).IsExpanded() {
				ob.deleteListViewItems(item.(*ui.GoListViewItemObj), idx + 1)
			}
			//log.Println("RemoveListItem(", item.(*ui.GoListViewItemObj).Id(), ")")
			listViewItem.RemoveListItem(item, idx)
		}
	}
}

func (ob *GoFileDialogObj) getDirIndex(nodeId []int, rootPath string) (index int, path string) {
	var listViewItem *ui.GoListViewItemObj
	var idx int = -1
	path = rootPath
	for level := 0; level < len(nodeId); level++ {
		//log.Println("level :", level)
		idx++
		for id := 0; id < nodeId[level]; id++ {
			//log.Println("id :", id)
			//log.Println("idx :", idx)
			listViewItem = ob.dirView.Objects()[idx].(*ui.GoListViewItemObj)
			idx += ob.getItemCount(listViewItem)
		}
		listViewItem = ob.dirView.Objects()[idx].(*ui.GoListViewItemObj)
		path = path + listViewItem.Text() + "/"
	}
	ob.currentDirItem = listViewItem
	return idx, path
}

func (ob *GoFileDialogObj) getFileIndex(nodeId []int) (index int, path string) {
	var listViewItem *ui.GoListViewItemObj
	var idx int = -1
	path = ob.rootPath
	for level := 0; level < len(nodeId); level++ {
		//log.Println("level :", level)
		idx++
		for id := 0; id < nodeId[level]; id++ {
			//log.Println("id :", id)
			//log.Println("idx :", idx)
			listViewItem = ob.fileView.Objects()[idx].(*ui.GoListViewItemObj)
			idx += ob.getItemCount(listViewItem)
		}
		listViewItem = ob.dirView.Objects()[idx].(*ui.GoListViewItemObj)
		path = path + listViewItem.Text() + "/"
	}
	return idx, path
}

func (ob *GoFileDialogObj) getFilePath(listViewItem *ui.GoListViewItemObj) (path string) {
	path = "/" + listViewItem.Text()
	for {
		listViewItem = listViewItem.ParentControl().(*ui.GoListViewItemObj)
		if listViewItem.ObjectType() == "GoListViewItem" {
			path = "/" + listViewItem.Text() + path
		}
	}
	return path
}

func (ob *GoFileDialogObj)getItemCount(listViewItem *ui.GoListViewItemObj) (count int) {
	if listViewItem.IsExpanded() {
		for id := 0; id < listViewItem.ItemCount(); id++ {
			item := listViewItem.Objects()[id].(*ui.GoListViewItemObj)
			count += ob.getItemCount(item)
		}
	}
	count++
	//log.Println("getItemCount :", listViewItem.Text(), count)
	return count
}
