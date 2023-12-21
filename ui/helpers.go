package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (u *UI) send(r *request) error {
	u.responseViewModel.Clear()

	url := r.Url
	method := r.Method
	contentType := r.ContentType
	u.footerText.SetText(fmt.Sprintf("Execute %s %s", method, url))

	headerMap := make(map[string]string, 0)
	for _, v := range r.Header {
		sp := strings.Split(v, ":")
		headerMap[sp[0]] = sp[1]
	}

	ctx := context.Background()
	res, err := u.controller.Send(ctx, method, url, contentType, headerMap, r.Body)
	if err != nil {
		u.showErr(err)
		return err
	}

	u.responseViewModel.Update(res)
	u.historyViewModel.Add(*res)
	return nil
}

func (u *UI) showErr(err error) {
	u.footerText.SetText(err.Error()).SetTextColor(tcell.ColorRed)
}

func (u *UI) showInfo(msg string) {
	u.footerText.SetText(msg).SetTextColor(tcell.ColorGreen)
}

func (u *UI) showInputDialog(backTo tview.Primitive, formTransformer func(*tview.Form), okHandler func(*tview.Form)) {
	input := tview.NewForm()
	input.AddButton("OK", func() {
		okHandler(input)

		u.rootView.RemovePage("input")
		u.app.SetRoot(u.rootView, true).SetFocus(backTo)
	})
	formTransformer(input)

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
