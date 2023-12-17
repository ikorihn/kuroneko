package ui

import (
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
		switch event.Rune() {
		case 'q':
			if _, ok := u.app.GetFocus().(*tview.InputField); !ok {
				u.app.Stop()
			}
		}
		return event
	})

	u.responseText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'y':
			clipboard.WriteAll(u.responseText.GetText(false))
			return nil
		}

		return event
	})
	u.headerList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'd':
			curIdx := u.headerList.GetCurrentItem()
			u.headerList.RemoveItem(curIdx)
			u.headers = append(u.headers[:curIdx], u.headers[curIdx+1:]...)
			return nil
		}

		switch event.Key() {
		case tcell.KeyEnter:
			curIdx := u.headerList.GetCurrentItem()
			curText, _ := u.headerList.GetItemText(curIdx)
			u.showInputDialog(curText, "edit header", 20, u.headerList, func(text string) {
				u.headerList.RemoveItem(curIdx)
				u.headerList.InsertItem(curIdx, text, "", 20, nil)
				u.headers[curIdx] = text
			})
		case tcell.KeyEsc, tcell.KeyTab:
			u.app.SetFocus(u.inputForm)
		}

		return event
	})

	u.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlT:
			u.app.SetRoot(u.responseSwitchModal, true).SetFocus(u.responseSwitchModal)
			return nil
		}

		return event
	})

}
