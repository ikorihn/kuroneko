package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (u *UI) send() error {
	u.responseText.Clear()

	url := u.inputUrl.GetText()
	_, method := u.inputMethod.GetCurrentOption()
	_, contentType := u.requestContentType.GetCurrentOption()
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

	u.historyViewModel.Add(*res)

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

	// inputForm := tview.NewForm().
	// 	AddInputField(label, "", width, nil, callback).
	// 	AddButton("OK", func() {
	// 		callback(input.GetText())
	// 	})

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).
			SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	u.rootView.AddAndSwitchToPage(
		"input",
		modal(input, 40, 10),
		true,
	)

}
