package ui

import (
	"net/http"

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

var responseSwitchButtons = []string{
	"body", "header", "status", "curl",
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
	favoritesViewModel  *favoriteViewModel
	requestViewModel    *requestViewModel
	responseViewModel   *responseViewModel
	responseSwitchModal *tview.Modal
	footerText          *tview.TextView
}

func NewUi() *UI {
	ui := &UI{}

	ui.controller = controller.NewController()
	ui.app = tview.NewApplication()

	ui.historyViewModel = NewHistoryViewModel(ui)
	ui.favoritesViewModel = NewFavoriteViewModel(ui)
	ui.requestViewModel = NewRequestViewModel(ui)
	ui.responseViewModel = NewResponseViewModel(ui)

	ui.responseSwitchModal = tview.NewModal().
		SetText("Select response field you want to see").
		AddButtons(responseSwitchButtons).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.responseViewModel.Show(buttonIndex)

			ui.app.SetRoot(ui.rootView, true).SetFocus(ui.responseViewModel.responseField)
		})

	ui.footerText = tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("footer").SetTextColor(tcell.ColorGray)

	navigation := tview.NewGrid().
		SetRows(20, 0).
		AddItem(ui.historyViewModel.historyField, 0, 0, 1, 20, 0, 0, true).
		AddItem(ui.favoritesViewModel.favoriteField, 1, 0, 1, 20, 0, 0, true)
	reqAndRes := tview.NewGrid().
		SetRows(20, 0).
		AddItem(ui.requestViewModel.Grid, 0, 0, 1, 20, 0, 0, false).
		AddItem(ui.responseViewModel.Grid, 1, 0, 1, 20, 0, 0, false)
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
