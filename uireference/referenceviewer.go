// SPDX-License-Identifier: Unlicense OR MIT

/* github.com/utopiagio/utopia-x/uireference/referenceviewer.go */

package uireference

import (
    "log"
    "strings"
    
    ui "github.com/utopiagio/utopia"
    "github.com/utopiagio/utopia/history"
    "github.com/utopiagio/docs"
    "github.com/pkg/browser"
)

type PageObj struct {
    Layout *ui.GoLayoutObj
    hdrLogo *ui.GoLabelObj
    hdrSection *ui.GoLabelObj
    navigator *ui.GoListViewObj
    mainpanel *ui.GoLayoutObj
    mainview *ui.GoLayoutObj
    pageview *ui.GoLayoutObj
    currentDoc *ui.GoRichTextObj
    popup *GoPagePopupObj

    docHistory *history.GoHistoryObj
}

func Page(parent ui.GoObject, logo string, section string, content string) (hObj *PageObj) {
    var startDoc string = "ove.Overview"
    
    //popup := ui.GoWindow("UtopiaGio: Reference Documentation")
    layout := ui.GoVFlexBoxLayout(parent)
    hUIRef := &PageObj{
        Layout: layout,        
    }
    hdrLayout := ui.GoHFlexBoxLayout(layout)
    hdrLayout.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)

    /*hUIRef.popup.SetSize(700, 200)
    hUIRef.popup.SetLayoutStyle(ui.VFlexBoxLayout)
    hUIRef.popup.Layout().SetPadding(10,10,10,10)*/

    hUIRef.hdrLogo = ui.H2Label(hdrLayout, logo)
    hUIRef.hdrLogo.SetFontBold(true)
    hUIRef.hdrLogo.SetFontItalic(true)
    hUIRef.hdrLogo.SetSizePolicy(ui.FixedWidth, ui.PreferredHeight)
    hUIRef.hdrLogo.SetWidth(200)

    hUIRef.hdrSection = ui.H4Label(hdrLayout, section)
    hUIRef.hdrSection.SetFontBold(true)
    hUIRef.hdrSection.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
    hUIRef.hdrSection.SetMargin(0,5,0,0)

    hUIRef.mainpanel = ui.GoHFlexBoxLayout(layout)
    hUIRef.mainpanel.SetMargin(0,5,0,0)

    navpanel := ui.GoVBoxLayout(hUIRef.mainpanel)
    navpanel.SetSizePolicy(ui.FixedWidth, ui.ExpandingHeight)
    navpanel.SetWidth(200)
    navpanel.SetPadding(5,10,5,0)
    //navpanel.SetBorder(ui.BorderSingleLine, 1, 5, ui.Color_Blue)

    hUIRef.navigator = ui.GoListView(navpanel)
    hUIRef.navigator.SetLayoutMode(ui.Vertical)
    hUIRef.navigator.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
    //hUIRef.navigator.SetBorder(ui.BorderSingleLine, 1, 0, ui.Color_Black)
    //hUIRef.navigator.SetOnItemClicked(hUIRef.NavLink_Clicked)
    hUIRef.setupNavigator()
    hUIRef.navigator.SetOnItemClicked(hUIRef.NavLink_Clicked)
    hUIRef.navigator.SetOnItemDoubleClicked(hUIRef.NavLink_DoubleClicked)

    hUIRef.mainview = ui.GoVBoxLayout(hUIRef.mainpanel)
    hUIRef.mainview.SetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)
    //mainview.SetBorder(ui.BorderSingleLine, 1, 5, ui.Color_Blue)
    hUIRef.pageview = ui.GoHFlexBoxLayout(hUIRef.mainview)
    hUIRef.pageview.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
    hUIRef.pageview.SetPadding(10,10,10,10)
    //hUIRef.richText[0] = ui.GoRichText(pageview)
    //hUIRef.richText[0].SetOnLinkClick(hUIRef.Link_Clicked)
    //mainview.SetMargin(10,10,10,10)
    //mainview.SetBorder(ui.BorderSingleLine, 2, 6, ui.Color_LightBlue)
    
    title, name, docContent := docs.Page(startDoc)
    hUIRef.hdrSection.SetText(title)
    for x := 0; x < len(docContent); x++ {
        hUIRef.currentDoc = ui.GoRichText(hUIRef.pageview, name)
        hUIRef.currentDoc.LoadMarkDown(docContent[x])
        hUIRef.currentDoc.SetOnLinkClick(hUIRef.Link_Clicked)
    }
    hUIRef.docHistory = history.GoFileHistory(startDoc + "#")
    return hUIRef
}

func (ob *PageObj) ActionBackHistory() {
    link := ob.docHistory.Back()
    doc, anchor, found := strings.Cut(link, "#")
    if found {
        //log.Println("ActionBackHistory doc=", doc, "anchor=", anchor)
        // load document into viewer
        if doc != "" && doc != ob.currentDoc.Name() {
            //log.Println("Load doc......................")
            ob.Load(doc)
        }
        if anchor == "" {
            ob.mainview.ScrollToOffset(0) 
        } else {
            offset := ob.currentDoc.AnchorTable(anchor)
            //log.Println("mainview.ScrollToOffset......................", offset)
            ob.mainview.ScrollToOffset(offset)
        }
        ob.switchNavigatorFocus(doc)
    }
}

func (ob *PageObj) ActionForwardHistory() {
   link := ob.docHistory.Forward()
    doc, anchor, found := strings.Cut(link, "#")
    if found {
        //log.Println("ActionForwardHistory doc=", doc, "anchor=", anchor)
        // load document into viewer
        if doc != "" && doc != ob.currentDoc.Name() {
            //log.Println("Load doc......................")
            ob.Load(doc)
        }
        if anchor == "" {
            ob.mainview.ScrollToOffset(0) 
        } else {
            offset := ob.currentDoc.AnchorTable(anchor)
            //log.Println("mainview.ScrollToOffset......................", offset)
            ob.mainview.ScrollToOffset(offset)
        }
        ob.switchNavigatorFocus(doc)
    }
}

func (ob *PageObj) JumpTo(anchor string) {

}

func (ob *PageObj) Link_Clicked(link string) {
    // link can be a web address or a golang repository.
    // a link to a golang repository is defined as 'package.page#anchor'
    // it can be just the package page eg. api.GoWindow#
    // or an anchor link in the page eg. api.GoWindow#Close
    // if the package page is already loaded then if an anchor link is
    // provided the page scrolls to the anchor link,
    // otherwise the package page is loaded first and then if an anchor link
    // is provided then the page scrolls to the anchor link.
    if strings.HasPrefix(link, "https://") {
        // Open link in default Web Browser
        err := browser.OpenURL(link)
        if err != nil {
            log.Println("Failed to open link in browser.")
        }
    } else {
        doc, anchor, found := strings.Cut(link, "#")
        if found {
            if doc == "" {
                doc = ob.currentDoc.Name()
                link = doc + anchor
            }
            log.Println("Link_Clicked:", link)
            log.Println("Link doc =", doc, "anchor =", anchor)
            log.Println("CurrentPath:", ob.docHistory.CurrentPath())
            log.Println("CurrentDoc:", ob.currentDoc.Name())
            if link == ob.docHistory.CurrentPath() {
                return
            }
            // load document into viewer
            if doc != ob.currentDoc.Name() {
                log.Println("Load doc......................")
                ob.Load(doc)
            }
            if anchor != "" {
                offset := ob.currentDoc.AnchorTable(anchor)
                //log.Println("mainview.ScrollToOffset......................", offset)
                ob.mainview.ScrollToOffset(offset)
            }
            //log.Println("Add FileHistory: (", link, " )")
            ob.docHistory.Add(link)
        }
    }
}

func (ob *PageObj) PopupClosed() {
    ob.popup = nil
}

func (ob *PageObj) Load(doc string) {
    var title string
    ob.pageview.Clear()
    packg, docName, docContent := docs.Page(doc)
    if docName != "" {
        title = packg + " - " + docName
    } else {
        title = packg
        docName = packg
    }
    ob.hdrSection.SetText(title)
    for x := 0; x < len(docContent); x++ {
        ob.currentDoc = ui.GoRichText(ob.pageview, doc)
        ob.currentDoc.LoadMarkDown(docContent[x])
        ob.currentDoc.SetOnLinkClick(ob.Link_Clicked)
    }
    ob.mainview.ScrollToOffset(0)
    log.Println("Load doc.................EXIT.....")
}

func (ob *PageObj) NavLink_Clicked(nodeId []int) {
    var link string
    listItem := ob.navigator.Item(nodeId)
    doc := listItem.Text()
    switch doc {
        case "Overview":
            link = "ove.Overview#"
        case "Reference":
            link = "ref.Index#"
        case "API Reference":
            link = "api.Index#"
        default: 
            if listItem.ParentControl().ObjectType() == "GoListViewItemObj" {
                //log.Println("listItem has ParentControl")
                //log.Println("listItem.Parent.Text =", listItem.ParentControl().(*ui.GoListViewItemObj).Text())
                switch listItem.ParentControl().(*ui.GoListViewItemObj).Text() {
                    case "Overview":
                        link = "ove." + doc + "#"
                    case "Reference":
                        link = "ref." + doc + "#"
                    case "API Reference":
                        link = "api." + doc + "#"
                    default:
                        // Show error page not available
                }
            }
    }
    //log.Println("Link_Clicked(", link, ")")
    ob.Link_Clicked(link)
}

func (ob *PageObj) NavLink_DoubleClicked(nodeId []int) {
    listItem := ob.navigator.Item(nodeId)
    if listItem.IsExpanded() {
        listItem.SetExpanded(false)
    } else {
        listItem.SetExpanded(true)
    }
    ob.navigator.SwitchFocus(listItem)
}

func (ob *PageObj) switchNavigatorFocus(docRequest string) {
    var nodeLabel []string = nil
    packg, doc, found := strings.Cut(docRequest, ".")
    if found {
        switch packg {
            case "api":
                if doc == "Index" {
                    nodeLabel = []string{"API Reference"}
                } else {
                    nodeLabel = []string{"API Reference", doc}
                }
            case "ref":
                if doc == "Index" {
                    nodeLabel = []string{"Reference"}
                } else {
                    nodeLabel = []string{"Reference", doc}
                }
            case "ove":
                if doc == "Overview" {
                    nodeLabel = []string{"Overview"}
                } else {
                    nodeLabel = []string{"Overview", doc}
                }    
            default:  
        }
    }
    if nodeLabel != nil {
        item := ob.navigator.ItemByLabel(nodeLabel)
        ob.navigator.SwitchFocus(item)
    }
}

func (ob *PageObj) setupNavigator() {
    /*overview := */ob.navigator.AddListItem(nil, "Overview")
    api_reference := ob.navigator.AddListItem(nil, "API Reference")
    api_reference.AddListItem(nil, "GoApplication")
    api_reference.AddListItem(nil, "GoWindow")
    api_reference.AddListItem(nil, "GoButton")
    api_reference.AddListItem(nil, "GoButtonGroup")
    api_reference.AddListItem(nil, "GoCanvas")
    api_reference.AddListItem(nil, "GoCheckBox")
    api_reference.AddListItem(nil, "GoImage")
    api_reference.AddListItem(nil, "GoLabel")
    api_reference.AddListItem(nil, "GoLayout")
    api_reference.AddListItem(nil, "GoList")
    api_reference.AddListItem(nil, "GoMenu")
    api_reference.AddListItem(nil, "GoMenuBar")
    api_reference.AddListItem(nil, "GoMenuItem")
    api_reference.AddListItem(nil, "GoRadioButton")
    api_reference.AddListItem(nil, "GoRichText")
    api_reference.AddListItem(nil, "GoSlider")
    api_reference.AddListItem(nil, "GoSpacer")
    api_reference.AddListItem(nil, "GioObject")
    api_reference.AddListItem(nil, "GioWidget")
    /*reference := */ob.navigator.AddListItem(nil, "Reference")
}


type GoPagePopupObj struct {
    Window *ui.GoWindowObj
    page *ui.GoRichTextObj
    content string
}

func GoPagePopup(title string, content string) (hPopup *GoPagePopupObj) {
    window := ui.GoWindow(title)
    mainview := ui.GoVBoxLayout(window.Layout())
    mainview.SetSizePolicy(ui.ExpandingWidth, ui.ExpandingHeight)
    pageview := ui.GoHFlexBoxLayout(mainview)
    pageview.SetSizePolicy(ui.ExpandingWidth, ui.PreferredHeight)
    pageview.SetPadding(10,10,10,10)
    
    hPopup = &GoPagePopupObj{
        Window: window,
    }
    hPopup.page = ui.GoRichText(pageview, "")
    hPopup.AddContent(content)
    //mainview.SetBorder(ui.BorderSingleLine, 1, 5, ui.Color_Blue)
    hPopup.Window.Show()
    return hPopup
}

func (ob *GoPagePopupObj) AddContent(content string) {
    ob.page.Clear()
    ob.content = ob.content + content
    ob.page.LoadMarkDown(ob.content)
}