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

	subdirectory, err := FindMecchiDirectory(currentDir)

	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	if subdirectory != "" {
		return LoaderResult[string]{success: true, result: subdirectory}
	} else {
		return LoaderResult[string]{success: false, error: currentDir}
	}
}

func containsFile(folderPath string, targetFile string) bool {
	filePath := filepath.Join(folderPath, targetFile)
	_, err := os.Stat(filePath)
	return err == nil
}

func FindMecchiDirectory(rootPath string) (string, error) {
	targetFolder := "src"
	targetFile := "mecchi.py"

	folders, err := os.ReadDir(rootPath)

	for _, folder := range folders {
		if folder.IsDir() {
			fmt.Println(folder.Name())
		}
	}

	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}

	for _, folder := range folders {
		if folder.IsDir() {
			subfolderPath := filepath.Join(rootPath, folder.Name())
			srcFilePath := filepath.Join(subfolderPath, targetFolder, targetFile)
			fmt.Println(srcFilePath)

			_, err := os.Stat(srcFilePath)
			if err == nil {
				fmt.Println("found " + subfolderPath)
				return subfolderPath, nil
			}

		}
	}

	return "", nil
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
