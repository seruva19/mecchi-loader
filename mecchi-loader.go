package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	mecchiLoader := app.New()

	window := mecchiLoader.NewWindow("mecchi loader")
	window.Resize(fyne.NewSize(WindowWidth, WindowHeight))

	input := widget.NewMultiLineEntry()
	input.MultiLine = true
	input.SetMinRowsVisible(10)

	input.Wrapping = fyne.TextWrapWord

	infinite := widget.NewProgressBarInfinite()
	infinite.Hide()

	loading := func(show bool) {
		if show {
			infinite.Show()
		} else {
			infinite.Hide()
		}
	}

	tabs := container.NewAppTabs(
		container.NewTabItem("Main", MainUi(input, loading)),
		container.NewTabItem("Services", ServicesUi(input, loading)),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	window.SetContent(container.NewVBox(tabs, layout.NewSpacer(), infinite, input))
	window.ShowAndRun()
}
