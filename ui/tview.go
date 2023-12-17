package ui

import (
	"fmt"
	"net/http"
	"slices"
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

type UI struct {
	app      *tview.Application
	rootView *tview.Grid

	history      *tview.List
	inputForm    *tview.Form
	inputMethod  *tview.DropDown
	inputUrl     *tview.InputField
	responseText *tview.TextView
	footerText   *tview.TextView

	requestBody []byte
}

func NewUi() *UI {
	app := tview.NewApplication()
	ui := &UI{
		app: app,
	}

	ui.app = tview.NewApplication()

	ui.history = tview.NewList().ShowSecondaryText(true).SetSecondaryTextColor(tcell.ColorDimGray)
	ui.history.SetTitle("History").SetBorder(true)

	ui.inputMethod = tview.NewDropDown().
		SetLabel("Method: ").
		SetOptions(httpMethods, nil).
		SetCurrentOption(0)
	ui.inputUrl = tview.NewInputField()
	ui.inputUrl.SetLabel("URL: ")

	ui.inputForm = tview.NewForm().
		AddFormItem(ui.inputMethod).
		AddFormItem(ui.inputUrl).
		AddButton("Body", func() {
			ui.app.Suspend(func() {
				body, err := controller.EditBody()
				if err != nil {
					ui.showErr(err)
					return
				}

				ui.requestBody = body
			})
		}).
		AddButton("Send", func() {
			ui.Send()
		})
	ui.inputForm.SetTitle("Request").SetBorder(true)

	ui.responseText = tview.NewTextView()
	ui.responseText.SetTitle("Response").SetBorder(true)

	ui.footerText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("footer").SetTextColor(tcell.ColorGray)

	navigation := tview.NewGrid().SetRows(0).
		AddItem(ui.history, 0, 0, 1, 1, 0, 0, true)
	reqAndRes := tview.NewGrid().SetRows(0, 0).
		AddItem(ui.inputForm, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.responseText, 1, 0, 1, 1, 0, 0, false)
	ui.rootView = tview.NewGrid().
		SetRows(0, 2).
		SetColumns(40, 0).
		SetBorders(false).
		AddItem(navigation, 0, 0, 1, 1, 0, 0, true).
		AddItem(reqAndRes, 0, 1, 1, 1, 0, 0, false).
		AddItem(ui.footerText, 1, 0, 1, 2, 0, 0, false)

	ui.setupKeyboard()

	return ui
}

func (u *UI) Run() error {
	return u.app.SetRoot(u.rootView, true).SetFocus(u.inputForm).Run()
}

func (u *UI) Send() error {
	u.responseText.Clear()

	url := u.inputUrl.GetText()
	_, method := u.inputMethod.GetCurrentOption()
	u.footerText.SetText(fmt.Sprintf("Execute %s %s", method, url))

	res, err := controller.Send(method, url)
	if err != nil {
		u.showErr(err)
		return err
	}

	u.responseText.SetText(string(res.Body))

	executedTime := time.Now()

	u.history.AddItem(fmt.Sprintf("%s %s", method, url), executedTime.Format(time.RFC3339), 0, func() {
		u.inputMethod.SetCurrentOption(slices.Index(httpMethods, method))
		u.inputUrl.SetText(url)
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
