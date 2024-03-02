package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2/widget"
)

type LoaderResult[T any] struct {
	error   string
	success bool
	result  T
}

type PackageJSONScripts struct {
	Install         string `json:"install"`
	InstallLinux    string `json:"install:linux"`
	Launch          string `json:"launch"`
	LaunchLinux     string `json:"launch:linux"`
	Update          string `json:"update"`
	UpdateLinux     string `json:"update:linux"`
	ResetExtensions string `json:"reset"`
}

type PackageJSON struct {
	LoaderScripts PackageJSONScripts `json:"loader"`
}

func GetScripts(mecchiFolder string) LoaderResult[PackageJSONScripts] {
	packageJsonPath := filepath.Join(mecchiFolder, "package.json")
	packageJsonData, err := os.ReadFile(packageJsonPath)

	if err != nil {
		return LoaderResult[PackageJSONScripts]{success: false, error: err.Error()}
	}

	var packageJSON PackageJSON

	err = json.Unmarshal(packageJsonData, &packageJSON)
	if err != nil {
		return LoaderResult[PackageJSONScripts]{success: false, error: err.Error()}
	}

	return LoaderResult[PackageJSONScripts]{success: true, result: packageJSON.LoaderScripts}
}

func MecchiFolderExists() LoaderResult[string] {
	currentDir, err := os.Getwd()
	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	subdirectory, err := FindFile(currentDir, "src", "mecchi.py")

	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	if subdirectory != "" {
		return LoaderResult[string]{success: true, result: subdirectory}
	} else {
		return LoaderResult[string]{success: false, error: "mecchi folder not found"}
	}
}

func FindFile(rootPath string, foldername string, filename string) (string, error) {
	var resultPath string

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == foldername {
			targetFilePath := filepath.Join(path, filename)
			if _, err := os.Stat(targetFilePath); err == nil {
				resultPath = filepath.Dir(path)
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return resultPath, nil
}

func OpenInBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return os.ErrInvalid
	}

	return cmd.Start()
}

func OutputResult(input *widget.Entry, text string) {
	if input.Text == "" {
		input.Text += text
	} else {
		input.Text += fmt.Sprintf("\n%s", text)
	}

	input.CursorRow = len(input.Text)
	input.Refresh()
}
