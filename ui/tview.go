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

	historyViewModel    *historyViewModel
	requestViewModel    *requestViewModel
	responseText        *tview.TextView
	responseSwitchModal *tview.Modal
	footerText          *tview.TextView

	response *controller.History
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

	ui.requestViewModel = NewRequestViewModel(ui)

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

	ui.responseText = tview.NewTextView()
	ui.responseText.SetTitle("Response (Ctrl+T)").SetBorder(true)

	ui.footerText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("footer").SetTextColor(tcell.ColorGray)

	navigation := tview.NewGrid().SetRows(0).
		AddItem(ui.historyViewModel.HistoryField, 0, 0, 1, 1, 0, 0, true)
	reqAndRes := tview.NewGrid().
		SetRows(20, 0).
		AddItem(ui.requestViewModel.Grid, 0, 0, 1, 20, 10, 0, false).
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
	return u.app.SetRoot(u.rootView, true).SetFocus(u.requestViewModel.requestForm).Run()
}
