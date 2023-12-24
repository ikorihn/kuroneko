package ui

import (
	"fmt"
	"log"
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
	Request *request
	Grid    *tview.Grid

	requestForm      *tview.Form
	headerList       *tview.List
	bodyText         *tview.TextView
	inputMethod      *tview.DropDown
	inputUrl         *tview.InputField
	inputContentType *tview.DropDown
}

func NewRequestViewModel(ui *UI) *requestViewModel {
	request := &request{
		Header: make([]string, 0),
	}

	headerList := tview.NewList().ShowSecondaryText(false)
	headerList.SetTitle("Header (d->delete, Enter->edit)").SetBorder(true)

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
		AddButton("Add header", func() {
			ui.showInputDialog(ui.requestViewModel.requestForm,
				func(form *tview.Form) {
					form.AddInputField("Name", "", 20, nil, func(text string) {
					})
					form.AddInputField("Value", "", 20, nil, func(text string) {
					})
				},
				func(form *tview.Form) {
					name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
					value := form.GetFormItemByLabel("Value").(*tview.InputField).GetText()

					headerItem := fmt.Sprintf("%s:%s", name, value)

					headerList.AddItem(headerItem, "", 20, nil)

					request.Header = append(request.Header, headerItem)
				},
			)
		}).
		AddButton("Edit header", func() {
			ui.app.SetFocus(headerList)
		}).
		AddButton("Body", func() {
			ui.app.Suspend(func() {
				body, err := ui.controller.EditBody()
				if err != nil {
					ui.showErr(err)
					return
				}

				request.Body = body
				bodyText.SetText(string(body))
			})
		}).
		AddButton("Send", func() {
			ui.send(ui.requestViewModel.Request)
			ui.app.SetFocus(ui.responseViewModel.responseField)
		})

	grid := tview.NewGrid().
		SetRows(10, 20).
		AddItem(form, 0, 0, 1, 15, 0, 0, false).
		AddItem(headerList, 1, 0, 9, 15, 10, 0, false).
		AddItem(bodyText, 0, 15, 10, 5, 0, 0, false)
	grid.SetBorder(true).SetTitle("Request form (Ctrl+R)")

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

func (r *requestViewModel) Update(req *request) {
	r.Request = req
	r.inputMethod.SetCurrentOption(slices.Index(httpMethods, req.Method))
	r.inputUrl.SetText(req.Url)
	r.inputContentType.SetCurrentOption(slices.Index(contentTypes, req.ContentType))
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
		SetRows(20).
		AddItem(responseText, 0, 0, 1, 15, 0, 0, false).
		AddItem(statsText, 0, 15, 1, 10, 0, 0, false)
	grid.SetTitle("Response (Ctrl+T)").SetBorder(true)

	return &responseViewModel{
		Grid:          grid,
		responseField: responseText,
		statsText:     statsText,
	}
}

func (r *responseViewModel) Update(response *controller.History) {
	r.Response = response
	r.responseField.SetText(string(response.Body))
	r.statsText.SetText(fmt.Sprintf("%+v", response.HttpStat))
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
	historyField := tview.NewList().ShowSecondaryText(true).SetSecondaryTextColor(tcell.ColorGray)
	historyField.SetTitle("History (Ctrl+H)").SetBorder(true)
	return &historyViewModel{
		Parent:       ui,
		Histories:    []controller.History{},
		historyField: historyField,
	}
}

func (h *historyViewModel) Add(history controller.History) {
	h.Histories = append(h.Histories, history)
	h.historyField.AddItem(fmt.Sprintf("%s %s", history.Request.Method, history.Request.Url), history.ExecutionTime.Format(time.RFC3339), 0, func() {
		h.Parent.requestViewModel.Update(&request{
			Method:      history.Request.Method,
			Url:         history.Request.Url,
			ContentType: history.Request.ContentType,
		})
		h.Parent.responseViewModel.Update(&history)
		h.Parent.app.SetFocus(h.Parent.requestViewModel.requestForm)
	})
}

type favoriteViewModel struct {
	Parent        *UI
	favorite      controller.Favorite
	favoriteField *tview.List
}

func NewFavoriteViewModel(ui *UI) *favoriteViewModel {
	favoriteField := tview.NewList().ShowSecondaryText(false).SetSelectedFocusOnly(true)
	favoriteField.SetTitle("Favorites (Ctrl+F)").SetBorder(true)
	favorite, err := ui.controller.LoadFavorite()
	if err != nil {
		log.Fatalf("failed to load favorite %v\n", err)
	}

	for _, req := range favorite.Request {
		req := req
		favoriteField.AddItem(fmt.Sprintf("%s %s", req.Method, req.Url), "", 0, func() {
			ui.requestViewModel.Update(&request{
				Method:      req.Method,
				Url:         req.Url,
				ContentType: req.ContentType,
			})
			ui.app.SetFocus(ui.requestViewModel.requestForm)
		})
	}

	return &favoriteViewModel{
		Parent:        ui,
		favorite:      favorite,
		favoriteField: favoriteField,
	}
}

func (f *favoriteViewModel) Add(req controller.Request) error {
	err := f.Parent.controller.SaveFavorite(req)
	if err != nil {
		return fmt.Errorf("cannot save favorite: %w", err)
	}

	f.favoriteField.AddItem(fmt.Sprintf("%s %s", req.Method, req.Url), "", 0, func() {
		f.Parent.requestViewModel.Update(&request{
			Method:      req.Method,
			Url:         req.Url,
			ContentType: req.ContentType,
		})
		f.Parent.app.SetFocus(f.Parent.requestViewModel.requestForm)
	})
	return nil
}
