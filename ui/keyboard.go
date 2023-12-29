package ui

import (
	"fmt"
	"slices"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/ikorihn/kuroneko/controller"
	"github.com/rivo/tview"
)

func (u *UI) setupKeyboard() {
	u.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			u.send(u.requestViewModel.Request)
			return nil
		}

		switch event.Rune() {
		case '1':
			if _, ok := u.app.GetFocus().(*tview.InputField); !ok {
				u.app.SetRoot(u.rootView, true).SetFocus(u.historyViewModel.historyField)
				return nil
			}
		case '2':
			u.app.SetRoot(u.rootView, true).SetFocus(u.requestViewModel.requestForm)
			return nil
		case '3':
			if _, ok := u.app.GetFocus().(*tview.InputField); !ok {
				u.app.SetRoot(u.rootView, true).SetFocus(u.favoritesViewModel.favoriteField)
				return nil
			}
		case '4':
			u.app.SetRoot(u.responseSwitchModal, true).SetFocus(u.responseSwitchModal)
			return nil
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
		case 'f':
			resp := u.responseViewModel.Response
			formatted := controller.FormatBody(*resp)
			u.responseViewModel.responseField.SetText(formatted)
			return nil
		}

		return event
	})

	u.requestViewModel.headerList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		headerList := u.requestViewModel.headerList
		switch event.Rune() {
		case 'd':
			if headerList.GetItemCount() == 0 {
				return nil
			}
			curIdx := headerList.GetCurrentItem()
			headerItem, _ := headerList.GetItemText(curIdx)
			headerList.RemoveItem(curIdx)
			u.requestViewModel.Request.Headers.RemoveNameValue(headerItem)
			return nil
		case 'e':
			if headerList.GetItemCount() == 0 {
				return nil
			}
			curIdx := headerList.GetCurrentItem()
			u.showInputHeaderDialog(headerList, curIdx)

		case 'a':
			u.showInputHeaderDialog(headerList, -1)
		}

		switch event.Key() {
		case tcell.KeyEsc, tcell.KeyTab:
			u.app.SetFocus(u.requestViewModel.requestForm)
		}

		return event
	})

	u.favoritesViewModel.favoriteField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		list := u.favoritesViewModel.favoriteField
		switch event.Rune() {
		case 'd':
			if list.GetItemCount() == 0 {
				return nil
			}
			curIdx := list.GetCurrentItem()
			list.RemoveItem(curIdx)
			u.controller.SaveFavorite(slices.Delete(u.controller.Favorites.Request, curIdx, curIdx+1))
			return nil
		}

		return event
	})

}
