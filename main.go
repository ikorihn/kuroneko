package main

import (
	"github.com/ikorihn/kuroneko/ui"
	"github.com/mattn/go-runewidth"
)

func main() {
	runewidth.DefaultCondition = &runewidth.Condition{
		EastAsianWidth: false,
	}

	app := ui.NewUi()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
