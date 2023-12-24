package ui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (u *UI) setupKeyboard() {
	//	focusCycle := []tview.Primitive{
	//		u.history,
	//		u.inputMethod,
	//		u.inputUrl,
	//		u.responseText,
	//	}
	//
	//	u.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	//		curFocus := slices.Index(focusCycle, u.app.GetFocus())
	//		nextFocus := (curFocus + 1) % len(focusCycle)
	//		prevFocus := (curFocus - 1 + len(focusCycle)) % len(focusCycle)
	//		switch event.Key() {
	//		case tcell.KeyTab:
	//			u.app.SetFocus(focusCycle[nextFocus])
	//			return nil
	//		case tcell.KeyBacktab:
	//			u.app.SetFocus(focusCycle[prevFocus])
	//			return nil
	//		}
	//		return event
	//	})
	u.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlT:
			u.app.SetRoot(u.responseSwitchModal, true).SetFocus(u.responseSwitchModal)
			return nil
		case tcell.KeyCtrlH:
			if _, ok := u.app.GetFocus().(*tview.InputField); !ok {
				u.app.SetRoot(u.rootView, true).SetFocus(u.historyViewModel.historyField)
				return nil
			}
		case tcell.KeyCtrlF:
			if _, ok := u.app.GetFocus().(*tview.InputField); !ok {
				u.app.SetRoot(u.rootView, true).SetFocus(u.favoritesViewModel.favoriteField)
				return nil
			}
		case tcell.KeyCtrlR:
			u.app.SetRoot(u.rootView, true).SetFocus(u.requestViewModel.requestForm)
			return nil
		case tcell.KeyCtrlS:
			u.send(u.requestViewModel.Request)
			return nil
		}

		switch event.Rune() {
		case 'q':
			if _, ok := u.app.GetFocus().(*tview.InputField); !ok {
				u.app.Stop()
			}
		}
		return event
	})

	u.historyViewModel.historyField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 's':
			curIndex := u.historyViewModel.historyField.GetCurrentItem()
			err := u.favoritesViewModel.Add(u.historyViewModel.Histories[curIndex].Request)
			if err != nil {
				fmt.Printf("cannot save favorite %v\n", err)
			}
			return nil
		}

		return event
	})
	u.responseViewModel.responseField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'y':
			clipboard.WriteAll(u.responseViewModel.responseField.GetText(true))
			return nil
		}

		return event
	})

	u.requestViewModel.headerList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		headerList := u.requestViewModel.headerList
		switch event.Rune() {
		case 'd':
			curIdx := headerList.GetCurrentItem()
			headerList.RemoveItem(curIdx)
			u.requestViewModel.Request.Header = append(u.requestViewModel.Request.Header[:curIdx], u.requestViewModel.Request.Header[curIdx+1:]...)
			return nil
		}

		switch event.Key() {
		case tcell.KeyEnter:
			curIdx := headerList.GetCurrentItem()
			curText, _ := headerList.GetItemText(curIdx)

			u.showInputDialog(headerList,
				func(form *tview.Form) {
					form.AddInputField("Name", curText, 20, nil, func(text string) {
					})
					form.AddInputField("Value", "", 20, nil, func(text string) {
					})
				},
				func(form *tview.Form) {
					name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
					value := form.GetFormItemByLabel("Value").(*tview.InputField).GetText()

					headerItem := fmt.Sprintf("%s:%s", name, value)

					headerList.RemoveItem(curIdx)
					headerList.InsertItem(curIdx, headerItem, "", 20, nil)
					u.requestViewModel.Request.Header[curIdx] = headerItem
				},
			)
		case tcell.KeyEsc, tcell.KeyTab:
			u.app.SetFocus(u.requestViewModel.requestForm)
		}

		return event
	})

}
