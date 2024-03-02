package main

import (
	"fmt"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MainUi(input *widget.Entry, loading func(bool)) *fyne.Container {
	var pid int

	install := widget.NewButton("install", func() {
		loading(true)
		OutputResult(input, "Trying to locate existing Mecchi folder...")

		mecchiFolder := MecchiFolderExists()

		firstTimeInstall := false
		mecchiFolderPath := ""

		if mecchiFolder.success {
			OutputResult(input, fmt.Sprintf("Mecchi already exists: %s", mecchiFolder.result))
			mecchiFolderPath = mecchiFolder.result
		} else {
			OutputResult(input, "Mecchi does not exist, downloading from Github...")

			newMecchiFolder := CloneMecchi()
			if newMecchiFolder.success {
				OutputResult(input, fmt.Sprintf("Mecchi downloaded successfully to: %s", newMecchiFolder.result))
				mecchiFolderPath = newMecchiFolder.result
			} else {
				OutputResult(input, fmt.Sprintf("Mecchi download failed: %s", newMecchiFolder.error))
			}
		}

		OutputResult(input, fmt.Sprintf("Now installing dependencies to %s...", mecchiFolderPath))
		mecchiInstallation := InstallMecchi(mecchiFolderPath, firstTimeInstall, input)

		if mecchiInstallation.success {
			OutputResult(input, "Successfully installed dependencies...")
		} else {
			OutputResult(input, mecchiInstallation.error)
		}

		loading(false)
	})

	var launch *widget.Button
	var stop *widget.Button

	launch = widget.NewButton("launch", func() {
		go func() {
			mecchiFolder := MecchiFolderExists()

			if !mecchiFolder.success {
				OutputResult(input, "Mecchi does not exist, install it first! "+mecchiFolder.error)
			} else {
				OutputResult(input, "Launching Mecchi...")
				mecchiProcess := LaunchMecchi(mecchiFolder.result, input)

				if mecchiProcess.success {
					pid = mecchiProcess.result
				}

				launch.Hidden = true
				stop.Hidden = false
			}
		}()
	})

	stop = widget.NewButton("stop", func() {
		OutputResult(input, "Stopping Mecchi...")
		KillMecchi(pid, input)
		OutputResult(input, "Mecchi stopped.")

		launch.Hidden = false
		stop.Hidden = true
	})

	update := widget.NewButton("update", func() {
		loading(true)

		updateResult := UpdateMecchi(input)
		if updateResult.success {
			OutputResult(input, "Mecchi successfully updated")
		} else {
			OutputResult(input, updateResult.error)
		}

		loading(false)
	})

	stop.Hidden = true

	content := container.NewVBox(
		container.NewAdaptiveGrid(1, layout.NewSpacer()),
		container.NewAdaptiveGrid(1, launch, stop),
		container.NewGridWithColumns(2, install, update),
	)

	return content
}
