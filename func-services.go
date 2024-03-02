package main

import (
	"os"
	"os/exec"

	"fyne.io/fyne/v2/widget"
)

func ResetExtensions(mecchiFolder string, input *widget.Entry) LoaderResult[any] {
	packageJSON := GetScripts(mecchiFolder)

	if !packageJSON.success {
		return LoaderResult[any]{success: false, error: packageJSON.error}
	}

	var shell string
	var shellArgs []string

	if os.Getenv("SHELL") != "" {
		shell = os.Getenv("SHELL")
	} else {
		shell = "cmd"
	}

	shellArgs = []string{"-c", packageJSON.result.ResetExtensions}

	cmd := exec.Command(shell, shellArgs...)
	cmd.Dir = mecchiFolder

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_err := cmd.Start()
	if _err != nil {
		return LoaderResult[any]{success: false, error: _err.Error()}
	}

	return LoaderResult[any]{success: true}
}
