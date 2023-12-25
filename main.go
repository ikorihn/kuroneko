package main

import (
	"github.com/ikorihn/kuroneko/ui"
	"github.com/mattn/go-runewidth"
)

func main() {
	runewidth.DefaultCondition = &runewidth.Condition{
		EastAsianWidth: false,
	}

	app, err := ui.NewUi()
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
