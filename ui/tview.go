package ui

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/gdamore/tcell/v2"
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
	inputView    *tview.Grid
	inputMethod  *tview.DropDown
	inputUrl     *tview.InputField
	responseText *tview.TextView
	footerText   *tview.TextView
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
		SetLabel("Select method").
		SetOptions(httpMethods, nil)
	ui.inputUrl = tview.NewInputField()
	ui.inputUrl.
		SetLabel("URL: ").
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				url := ui.inputUrl.GetText()
				_, method := ui.inputMethod.GetCurrentOption()
				ui.footerText.SetText(fmt.Sprintf("Execute %s %s", method, url))

				req, err := http.NewRequest(method, url, nil)
				if err != nil {
					ui.footerText.SetText(err.Error()).SetTextColor(tcell.ColorRed)
				}

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					ui.footerText.SetText(err.Error()).SetTextColor(tcell.ColorRed)
				}

				defer res.Body.Close()
				b, err := io.ReadAll(res.Body)
				if err != nil {
					ui.footerText.SetText(err.Error()).SetTextColor(tcell.ColorRed)
				}

				ui.responseText.SetText(string(b))

				executedTime := time.Now()

				ui.history.AddItem(fmt.Sprintf("%s %s", method, url), executedTime.Format(time.RFC3339), 0, func() {
					ui.inputMethod.SetCurrentOption(slices.Index(httpMethods, method))
					ui.inputUrl.SetText(url)
					ui.responseText.SetText(string(b))
				})
			}
		})

	ui.inputView = tview.NewGrid().SetRows(1).
		AddItem(ui.inputMethod, 0, 0, 1, 1, 0, 0, true).
		AddItem(ui.inputUrl, 0, 1, 1, 1, 0, 0, true)

	ui.responseText = tview.NewTextView()
	ui.responseText.SetTitle("Response").SetBorder(true)

	ui.footerText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("footer").SetTextColor(tcell.ColorGray)

	navigation := tview.NewGrid().SetRows(0).
		AddItem(ui.history, 0, 0, 1, 1, 0, 0, true)
	reqAndRes := tview.NewGrid().SetRows(0, 0).
		AddItem(ui.inputView, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.responseText, 1, 0, 9, 1, 0, 0, false)
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
	return u.app.SetRoot(u.rootView, true).Run()
}
