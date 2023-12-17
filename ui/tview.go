package ui

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	"none",
	"application/json",
	"application/xml",
	"text/plain",
}

type UI struct {
	controller *controller.Controller

	app      *tview.Application
	rootView *tview.Pages
	rootGrid *tview.Grid

	history             *tview.List
	inputForm           *tview.Form
	inputMethod         *tview.DropDown
	inputUrl            *tview.InputField
	requestContentType  *tview.DropDown
	headerList          *tview.List
	responseText        *tview.TextView
	responseSwitchModal *tview.Modal
	footerText          *tview.TextView

	headers     []string
	requestBody []byte
	response    *controller.Response
}

func NewUi() *UI {
	app := tview.NewApplication()
	ui := &UI{
		app: app,
	}

	ui.controller = controller.NewController()

	ui.app = tview.NewApplication()

	ui.history = tview.NewList().ShowSecondaryText(true).SetSecondaryTextColor(tcell.ColorDimGray)
	ui.history.SetTitle("History").SetBorder(true)

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
		})
	ui.inputForm.SetTitle("Request form").SetBorder(true)

	ui.headerList = tview.NewList().ShowSecondaryText(false)
	ui.headerList.SetTitle("Request Header").SetBorder(true)

	ui.responseText = tview.NewTextView()
	ui.responseText.SetTitle("Response").SetBorder(true)

	ui.footerText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("footer").SetTextColor(tcell.ColorGray)

	navigation := tview.NewGrid().SetRows(0).
		AddItem(ui.history, 0, 0, 1, 1, 0, 0, true)
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

func (u *UI) send() error {
	u.responseText.Clear()

	url := u.inputUrl.GetText()
	methodIdx, method := u.inputMethod.GetCurrentOption()
	contentTypeIdx, contentType := u.requestContentType.GetCurrentOption()
	u.footerText.SetText(fmt.Sprintf("Execute %s %s", method, url))

	headerMap := make(map[string]string, 0)
	for _, v := range u.headers {
		sp := strings.Split(v, ":")
		headerMap[sp[0]] = sp[1]
	}

	res, err := u.controller.Send(method, url, contentType, headerMap, u.requestBody)
	if err != nil {
		u.showErr(err)
		return err
	}

	u.response = res

	u.responseText.SetText(string(res.Body))

	executedTime := time.Now()

	u.history.AddItem(fmt.Sprintf("%s %s", method, url), executedTime.Format(time.RFC3339), 0, func() {
		u.inputMethod.SetCurrentOption(methodIdx)
		u.inputUrl.SetText(url)
		u.requestContentType.SetCurrentOption(contentTypeIdx)
		u.responseText.SetText(string(res.Body))
	})

	u.app.SetFocus(u.responseText)

	u.requestBody = nil

	return nil
}

func (u *UI) showErr(err error) {
	u.footerText.SetText(err.Error()).SetTextColor(tcell.ColorRed)
}

func (u *UI) showInfo(msg string) {
	u.footerText.SetText(msg).SetTextColor(tcell.ColorGreen)
}

func (u *UI) showInputDialog(text string, label string, width int, backTo tview.Primitive, callback func(text string)) {
	input := tview.NewInputField().SetText(text)
	input.SetLabel(label).SetLabelWidth(width).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			callback(input.GetText())
		}
		switch key {
		case tcell.KeyEnter, tcell.KeyEsc:
			u.rootView.RemovePage("input")
			u.app.SetRoot(u.rootView, true).SetFocus(backTo)
		}
	})

	u.rootView.AddAndSwitchToPage(
		"input",
		tview.NewGrid().
			SetColumns(0, 0, 0).
			SetRows(0, 3, 0).
			AddItem(input, 1, 1, 1, 1, 0, 0, true),
		true,
	)
}
