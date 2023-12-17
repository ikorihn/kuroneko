package ui

import (
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
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
	//

	u.responseText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'y':
			clipboard.WriteAll(u.responseText.GetText(false))
			return nil
		}

		return event
	})

}
