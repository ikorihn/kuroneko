package ui

import (
	"fmt"
	"slices"
	"time"

	"github.com/ikorihn/kuroneko/controller"
	"github.com/rivo/tview"
)

type requestViewModel struct {
	Parent       *UI
	Request      *request
	RequestField *tview.Form
}

func NewRequestViewModel(ui *UI) *requestViewModel {
	request := &request{
		Header: make([]string, 0),
	}

	form := tview.NewForm().
		AddDropDown("Method", httpMethods, 0, func(selected string, optionIndex int) {
			request.Method = selected
		}).
		AddInputField("URL", "", 0, nil, func(text string) {
			request.Url = text
		}).
		AddDropDown("Content-Type", contentTypes, 0, func(selected string, optionIndex int) {
			request.ContentType = selected
		}).
		AddButton("Add header", func() {
			ui.showInputDialog("", "header to add", 20, ui.inputForm, func(text string) {
				ui.headerList.AddItem(text, "", 20, nil)

				request.Header = append(request.Header, text)
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

				request.Body = body
			})
		}).
		AddButton("Send", func() {
			ui.send()
			ui.app.SetFocus(ui.responseText)
		})

	form.SetTitle("Request form (Ctrl+R)").SetBorder(true)

	return &requestViewModel{
		Parent:       ui,
		Request:      request,
		RequestField: form,
	}
}

type responseViewModel struct {
	Response      *controller.History
	ResponseField *tview.TextView
}

func (r *responseViewModel) Update(response *controller.History) {
	r.Response = response
	r.ResponseField.SetText(string(response.Body))
}

type historyViewModel struct {
	Parent       *UI
	History      []controller.History
	HistoryField *tview.List
}

func (h *historyViewModel) Add(history controller.History) {
	h.History = append(h.History, history)
	h.HistoryField.AddItem(fmt.Sprintf("%s %s", history.Request.Method, history.Request.URL.String()), history.ExecutionTime.Format(time.RFC3339), 0, func() {
		h.Parent.inputMethod.SetCurrentOption(slices.Index(httpMethods, history.Request.Method))
		h.Parent.inputUrl.SetText(history.Request.URL.String())
		h.Parent.requestContentType.SetCurrentOption(slices.Index(contentTypes, history.Request.Header.Get("Content-Type")))
		h.Parent.responseText.SetText(string(history.Body))
		h.Parent.app.SetFocus(h.Parent.inputForm)
	})
}
