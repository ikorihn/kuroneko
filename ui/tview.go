package ui

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/ikorihn/kuroneko/controller"
	"github.com/rivo/tview"
)

var httpMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

var contentTypes = []string{
	"",
	"application/json",
	"application/xml",
	"text/plain",
}

type request struct {
	Method      string
	Url         string
	ContentType string
	Header      []string
	Body        []byte
}

type UI struct {
	controller *controller.Controller

	app      *tview.Application
	rootView *tview.Pages
	rootGrid *tview.Grid

	historyViewModel   *historyViewModel
	inputForm          *tview.Form
	inputMethod        *tview.DropDown
	inputUrl           *tview.InputField
	requestContentType *tview.DropDown
	requestViewModel   *requestViewModel

	headerList          *tview.List
	responseText        *tview.TextView
	responseSwitchModal *tview.Modal
	footerText          *tview.TextView

	headers     []string
	requestBody []byte
	response    *controller.History
}

func NewUi() *UI {
	ui := &UI{}

	ui.controller = controller.NewController()
	ui.app = tview.NewApplication()

	ui.historyViewModel = &historyViewModel{
		Parent:       ui,
		History:      []controller.History{},
		HistoryField: tview.NewList().ShowSecondaryText(true).SetSecondaryTextColor(tcell.ColorDimGray),
	}
	ui.historyViewModel.HistoryField.SetTitle("History (Ctrl+H)").SetBorder(true)

	ui.headers = make([]string, 0)

	ui.inputMethod = tview.NewDropDown().
		SetLabel("Method: ").
		SetOptions(httpMethods, nil).
		SetCurrentOption(0)
	ui.inputUrl = tview.NewInputField()
	ui.inputUrl.SetLabel("URL: ")
	ui.requestContentType = tview.NewDropDown().
		SetLabel("Content-Type: ").
		SetOptions(contentTypes, nil).
		SetCurrentOption(0)

	ui.responseSwitchModal = tview.NewModal().
		SetText("Select response field you want to see").
		AddButtons([]string{"Body", "Header", "Status"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Body":
				ui.responseText.SetText(string(ui.response.Body))
			case "Header":
				txt := make([]string, 0)
				for k := range ui.response.Header {
					txt = append(txt, fmt.Sprintf("%s: %s", k, ui.response.Header.Get(k)))
				}
				ui.responseText.SetText(strings.Join(txt, "\n"))
			case "Status":
				ui.responseText.SetText(strconv.Itoa(ui.response.StatusCode))
			}

			ui.app.SetRoot(ui.rootView, true).SetFocus(ui.responseText)
		})

	ui.inputForm = tview.NewForm().
		AddFormItem(ui.inputMethod).
		AddFormItem(ui.inputUrl).
		AddFormItem(ui.requestContentType).
		AddButton("Add header", func() {
			ui.showInputDialog("", "header to add", 20, ui.inputForm, func(text string) {
				ui.headerList.AddItem(text, "", 20, nil)
				ui.headers = append(ui.headers, text)
			})
		}).
		AddButton("Edit header", func() {
			ui.app.SetFocus(ui.headerList)
		}).
		AddButton("Body", func() {
			ui.app.Suspend(func() {
				body, err := ui.controller.EditBody()
				if err != nil {
					ui.showErr(err)
					return
				}

				ui.requestBody = body
			})
		}).
		AddButton("Send", func() {
			ui.send()
			ui.app.SetFocus(ui.responseText)
		})
	ui.inputForm.SetTitle("Request form (Ctrl+R)").SetBorder(true)

	ui.headerList = tview.NewList().ShowSecondaryText(false)
	ui.headerList.SetTitle("Header (d->delete, Enter->edit)").SetBorder(true)

	ui.responseText = tview.NewTextView()
	ui.responseText.SetTitle("Response (Ctrl+T)").SetBorder(true)

	ui.footerText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("footer").SetTextColor(tcell.ColorGray)

	navigation := tview.NewGrid().SetRows(0).
		AddItem(ui.historyViewModel.HistoryField, 0, 0, 1, 1, 0, 0, true)
	reqAndRes := tview.NewGrid().
		SetRows(20, 0).
		SetColumns(20, 0).
		AddItem(ui.inputForm, 0, 0, 1, 15, 10, 0, false).
		AddItem(ui.headerList, 0, 15, 1, 5, 10, 0, false).
		AddItem(ui.responseText, 1, 0, 1, 20, 0, 0, false)
	ui.rootGrid = tview.NewGrid().
		SetRows(0, 2).
		SetColumns(40, 0).
		SetBorders(false).
		AddItem(navigation, 0, 0, 1, 1, 0, 0, true).
		AddItem(reqAndRes, 0, 1, 1, 1, 0, 0, false).
		AddItem(ui.footerText, 1, 0, 1, 2, 0, 0, false)

	ui.rootView = tview.NewPages().AddPage("main", ui.rootGrid, true, true)
	ui.setupKeyboard()

	return ui
}

func (u *UI) Run() error {
	return u.app.SetRoot(u.rootView, true).SetFocus(u.inputForm).Run()
}
