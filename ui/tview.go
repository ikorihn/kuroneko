package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ui struct {
	app     *tview.Application
	curTime *tview.TableCell

	inputField *tview.InputField
}

func NewUi() *ui {
	app := tview.NewApplication()
	u := &ui{
		app: app,
	}
	u.buildApplication()

	return u
}

func (u *ui) Run() error {
	return u.app.Run()
}

func (u *ui) buildApplication() {
	infoPanel := u.createInput()
	commandList := u.createNavigation()
	layout := u.createLayout(commandList, infoPanel)
	pages := tview.NewPages()
	pages.AddPage("main", layout, true, true)

	u.app.SetRoot(pages, true)

}

func (u *ui) createInput() *tview.Flex {
	textView := tview.NewTextView()
	textView.SetTitle("textView")
	textView.SetBorder(true)

	inputField := tview.NewInputField()
	inputField.SetLabel("input: ")
	inputField.SetTitle("inputField").
		SetBorder(true)

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			textView.SetText(textView.GetText(true) + inputField.GetText() + "\n")
			inputField.SetText("")
			return nil
		}
		return event
	})

	u.inputField = inputField

	infoPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(inputField, 3, 0, true).
		AddItem(textView, 0, 1, false)
	return infoPanel
}

func (u *ui) createNavigation() (commandList *tview.List) {
	commandList = tview.NewList()
	commandList.SetBorder(true).SetTitle("Command")

	commandList.AddItem("Test", "", 'p', func() {
		u.app.SetFocus(u.inputField)
	})
	commandList.AddItem("Quit", "", 'q', func() {
		u.app.Stop()
	})
	return commandList
}

func (u *ui) createInfoPanel() (infoPanel *tview.Flex) {

	infoTable := tview.NewTable()
	infoTable.SetBorder(true).SetTitle("Information")

	cnt := 0
	infoTable.SetCellSimple(cnt, 0, "Data1:")
	infoTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	info1 := tview.NewTableCell("aaa")
	infoTable.SetCell(cnt, 1, info1)
	cnt++

	infoTable.SetCellSimple(cnt, 0, "Data2:")
	infoTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	info2 := tview.NewTableCell("bbb")
	infoTable.SetCell(cnt, 1, info2)
	cnt++

	infoTable.SetCellSimple(cnt, 0, "Time:")
	infoTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	u.curTime = tview.NewTableCell("0")
	infoTable.SetCell(cnt, 1, u.curTime)
	cnt++

	outputTable := tview.NewTable()
	outputTable.SetBorder(true).SetTitle("Information")

	cnt = 0
	outputTable.SetCellSimple(cnt, 0, "Output:")
	outputTable.GetCell(cnt, 0).SetAlign(tview.AlignRight)
	output := tview.NewTableCell("123")
	outputTable.SetCell(cnt, 1, output)

	infoPanel = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(infoTable, 0, 1, false).
		AddItem(outputTable, 0, 1, false)
	return infoPanel
}

func (u *ui) createLayout(cList tview.Primitive, recvPanel tview.Primitive) (layout *tview.Flex) {
	bodyLayout := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(cList, 20, 1, true).
		AddItem(recvPanel, 0, 1, false)

	header := tview.NewTextView()
	header.SetBorder(true)
	header.SetText("tview study")
	header.SetTextAlign(tview.AlignCenter)

	layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 3, 1, false).
		AddItem(bodyLayout, 0, 1, true)

	return layout
}
