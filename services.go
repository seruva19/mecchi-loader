package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ServicesUi(input *widget.Entry, loading func(bool)) *fyne.Container {
	reset := widget.NewButton("reset components", func() {
		mecchiFolder := MecchiFolderExists()

		if !mecchiFolder.success {
			OutputResult(input, "Mecchi does not exist!")
		} else {
			OutputResult(input, "Launching Mecchi...")
			resetResult := ResetExtensions(mecchiFolder.result, input)

			if resetResult.success {
				OutputResult(input, "Now components will be reinitialized on next Mecchi launch.")
			}
		}
	})

	content := container.NewVBox(
		container.NewAdaptiveGrid(1, layout.NewSpacer()),
		container.NewAdaptiveGrid(1, reset),
	)

	return content
}
