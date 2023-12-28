package ui

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/ikorihn/kuroneko/controller"
	"github.com/rivo/tview"
	"moul.io/http2curl"
)

type requestViewModel struct {
	Parent  *UI
	Request *controller.Request
	Grid    *tview.Grid

	requestForm      *tview.Form
	headerList       *tview.List
	bodyText         *tview.TextView
	inputMethod      *tview.DropDown
	inputUrl         *tview.InputField
	inputContentType *tview.DropDown
}

func NewRequestViewModel(ui *UI) *requestViewModel {
	request := controller.NewRequest()
	request.Method = httpMethods[0]

	headerList := tview.NewList().ShowSecondaryText(false).SetSelectedFocusOnly(true)
	headerList.SetTitle("Header (a->add, d->delete, e->edit)").SetBorder(true)

	inputMethod := tview.NewDropDown().
		SetLabel("Method: ").
		SetOptions(httpMethods, nil).
		SetCurrentOption(0).
		SetSelectedFunc(func(text string, index int) {
			ui.requestViewModel.Request.Method = text
		})
	inputUrl := tview.NewInputField().SetLabel("URL: ").
		SetChangedFunc(func(text string) {
			ui.requestViewModel.Request.Url = text
		})
	inputContentType := tview.NewDropDown().
		SetLabel("Content-Type: ").
		SetOptions(contentTypes, nil).
		SetCurrentOption(0).
		SetSelectedFunc(func(text string, index int) {
			ui.requestViewModel.Request.ContentType = text
		})
	bodyText := tview.NewTextView()
	bodyText.SetTitle("Body").SetBorder(true)

	form := tview.NewForm().
		AddFormItem(inputMethod).
		AddFormItem(inputUrl).
		AddFormItem(inputContentType).
		AddButton("Header", func() {
			ui.app.SetFocus(headerList)
		}).
		AddButton("Body", func() {
			ui.app.Suspend(func() {
				body, err := ui.controller.EditBody(ui.requestViewModel.Request.Body)
				if err != nil {
					ui.showErr(err)
					return
				}

				ui.requestViewModel.Request.Body = body
				bodyText.SetText(string(body))
			})
		}).
		AddButton("Send", func() {
			ui.send(ui.requestViewModel.Request)
			ui.app.SetFocus(ui.responseViewModel.responseField)
		})

	grid := tview.NewGrid().
		SetRows(10, 0).
		SetColumns(0, 10).
		AddItem(form, 0, 0, 1, 8, 0, 0, false).
		AddItem(headerList, 1, 0, 1, 8, 0, 0, false).
		AddItem(bodyText, 0, 8, 2, 2, 0, 0, false)
	grid.SetBorder(true).SetTitle("Request[Ctrl+R]")

	return &requestViewModel{
		Parent:           ui,
		Request:          request,
		Grid:             grid,
		requestForm:      form,
		headerList:       headerList,
		bodyText:         bodyText,
		inputMethod:      inputMethod,
		inputUrl:         inputUrl,
		inputContentType: inputContentType,
	}
}

func (r *requestViewModel) Update(req controller.Request) {
	r.Request = &req
	r.inputMethod.SetCurrentOption(slices.Index(httpMethods, req.Method))
	r.inputUrl.SetText(req.Url)
	r.inputContentType.SetCurrentOption(slices.Index(contentTypes, req.ContentType))
	r.bodyText.SetText(string(req.Body))

	r.headerList.Clear()
	for k, v := range req.Headers {
		r.headerList.AddItem(fmt.Sprintf("%s:%s", k, v), "", 0, nil)
	}
}

type responseViewModel struct {
	Response *controller.History
	Grid     *tview.Grid

	responseField *tview.TextView
	statsText     *tview.TextView
}

func NewResponseViewModel(ui *UI) *responseViewModel {
	responseText := tview.NewTextView()

	statsText := tview.NewTextView()
	statsText.SetTitle("Stats").SetBorder(true)

	grid := tview.NewGrid().
		SetColumns(0, 30).
		AddItem(responseText, 0, 0, 1, 1, 0, 0, false).
		AddItem(statsText, 0, 1, 1, 1, 0, 0, false)
	grid.SetTitle("Response[Ctrl+T] (y->Copy to clipboard)").SetBorder(true)

	return &responseViewModel{
		Grid:          grid,
		responseField: responseText,
		statsText:     statsText,
	}
}

func (r *responseViewModel) Update(response *controller.History) {
	r.Response = response
	r.responseField.SetText(string(response.Body))

	stats := fmt.Sprintf("Status: %d\n", response.StatusCode)
	stats += fmt.Sprintf("---------\n")
	stats += fmt.Sprintf("%+v", response.HttpStat)
	r.statsText.SetText(stats)
}

func (r *responseViewModel) Clear() {
	r.Response = nil
	r.responseField.Clear()
	r.statsText.Clear()
}

func (r *responseViewModel) Show(buttonIndex int) {
	i := 0
	if i == buttonIndex {
		r.responseField.SetText(string(r.Response.Body))
		return
	}
	i++
	if i == buttonIndex {
		txt := make([]string, 0)
		for k := range r.Response.Header {
			txt = append(txt, fmt.Sprintf("%s: %s", k, r.Response.Header.Get(k)))
		}
		r.responseField.SetText(strings.Join(txt, "\n"))
		return
	}
	i++
	if i == buttonIndex {
		r.responseField.SetText(strconv.Itoa(r.Response.StatusCode))
		return
	}
	i++
	if i == buttonIndex {
		req := r.Response.Request.ToHttpReq()
		command, _ := http2curl.GetCurlCommand(req)
		r.responseField.SetText(command.String())
		return
	}
}

type historyViewModel struct {
	Parent       *UI
	Histories    []controller.History
	historyField *tview.List
}

func NewHistoryViewModel(ui *UI) *historyViewModel {
	historyField := tview.NewList().ShowSecondaryText(true).SetSecondaryTextColor(tcell.ColorGray).SetSelectedFocusOnly(true)
	historyField.SetTitle("History[Ctrl+H] (s->Add favorite)").SetBorder(true)
	return &historyViewModel{
		Parent:       ui,
		Histories:    []controller.History{},
		historyField: historyField,
	}
}

func (h *historyViewModel) Add(history controller.History) {
	h.Histories = slices.Insert(h.Histories, 0, history)
	h.historyField.InsertItem(0, fmt.Sprintf("%s %s", history.Request.Method, history.Request.Url), history.ExecutionTime.Format(time.RFC3339), 0, func() {
		h.Parent.requestViewModel.Update(history.Request)
		h.Parent.responseViewModel.Update(&history)
		h.Parent.app.SetFocus(h.Parent.requestViewModel.requestForm)
	})
}

type favoriteViewModel struct {
	Parent        *UI
	favoriteField *tview.List
}

func NewFavoriteViewModel(ui *UI) *favoriteViewModel {
	favoriteField := tview.NewList().ShowSecondaryText(false).SetSelectedFocusOnly(true)
	favoriteField.SetTitle("Favorites[Ctrl+F]").SetBorder(true)
	favorite := ui.controller.Favorites

	for _, req := range favorite.Request {
		req := req
		favoriteField.AddItem(fmt.Sprintf("%s %s", req.Method, req.Url), "", 0, func() {
			ui.requestViewModel.Update(req)
			ui.app.SetFocus(ui.requestViewModel.requestForm)
		})
	}

	return &favoriteViewModel{
		Parent:        ui,
		favoriteField: favoriteField,
	}
}

func (f *favoriteViewModel) Add(req controller.Request) error {
	favReqs := append(f.Parent.controller.Favorites.Request, req)
	err := f.Parent.controller.SaveFavorite(favReqs)
	if err != nil {
		return fmt.Errorf("cannot save favorite: %w", err)
	}

	f.favoriteField.AddItem(fmt.Sprintf("%s %s", req.Method, req.Url), "", 0, func() {
		f.Parent.requestViewModel.Update(req)
		f.Parent.app.SetFocus(f.Parent.requestViewModel.requestForm)
	})
	return nil
}
