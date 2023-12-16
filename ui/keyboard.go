package ui

import (
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (u *UI) setupKeyboard() {
	focusCycle := []tview.Primitive{
		u.history,
		u.inputMethod,
		u.inputUrl,
		u.responseText,
	}

	// Setup app level keyboard shortcuts.
	u.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		curFocus := slices.Index(focusCycle, u.app.GetFocus())
		nextFocus := (curFocus + 1) % len(focusCycle)
		prevFocus := (curFocus - 1 + len(focusCycle)) % len(focusCycle)
		switch event.Key() {
		case tcell.KeyTab:
			u.app.SetFocus(focusCycle[nextFocus])
			return nil
		case tcell.KeyBacktab:
			u.app.SetFocus(focusCycle[prevFocus])
			return nil
		}
		return event
	})

}
