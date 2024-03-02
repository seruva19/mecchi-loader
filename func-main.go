package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"fyne.io/fyne/v2/widget"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneMecchi() LoaderResult[string] {
	currentDir, _err := os.Getwd()
	if _err != nil {
		return LoaderResult[string]{success: false, error: _err.Error()}
	}

	fullPath := filepath.Join(currentDir, MecchiFolder)

	_, err := git.PlainClone(fullPath, false, &git.CloneOptions{
		URL: MecchiRepo,
	})

	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	return LoaderResult[string]{success: true, result: fullPath}
}

func InstallMecchi(mecchiFolder string, firstTimeInstall bool, input *widget.Entry) LoaderResult[string] {
	packageJSON := GetScripts(mecchiFolder)

	if !packageJSON.success {
		return LoaderResult[string]{success: false, error: packageJSON.error}
	}

	var shell string
	var shellArgs []string

	if os.Getenv("SHELL") != "" {
		shell = os.Getenv("SHELL")
		shellArgs = []string{"-c", packageJSON.result.InstallLinux}
	} else {
		shell = "cmd"
		shellArgs = []string{"/c", packageJSON.result.Install}
	}

	cmd := exec.Command(shell, shellArgs...)
	cmd.Dir = mecchiFolder

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	err = cmd.Start()
	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	scanner := bufio.NewScanner(stdout)

	go func() {
		for scanner.Scan() {
			OutputResult(input, scanner.Text())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	return LoaderResult[string]{success: true}
}

func LaunchMecchi(mecchiFolder string, input *widget.Entry) LoaderResult[int] {
	packageJSON := GetScripts(mecchiFolder)

	if !packageJSON.success {
		return LoaderResult[int]{success: false, error: packageJSON.error}
	}

	var shell string
	var shellArgs []string

	if os.Getenv("SHELL") != "" {
		shell = os.Getenv("SHELL")
		shellArgs = []string{"-c", packageJSON.result.LaunchLinux}
	} else {
		shell = "cmd"
		shellArgs = []string{"/c", packageJSON.result.Launch}
	}

	cmd := exec.Command(shell, shellArgs...)
	cmd.Dir = mecchiFolder

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_err := cmd.Start()
	if _err != nil {
		return LoaderResult[int]{success: false, error: _err.Error()}
	}

	OpenInBrowser(MecchiUrl)
	return LoaderResult[int]{success: true, result: cmd.Process.Pid}
}

func KillMecchi(pid int, input *widget.Entry) error {
	if runtime.GOOS == "linux" {
		OutputResult(input, "Not implemented yet for Linux.")
	} else if runtime.GOOS == "windows" {
		process, _ := os.FindProcess(pid)

		killCmd := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(process.Pid))
		killCmd.Stderr = os.Stderr
		killCmd.Stdout = os.Stdout
		return killCmd.Run()
	} else {
		OutputResult(input, "Unsupported operating system:")
	}

	return nil
}

func UpdateMecchi(input *widget.Entry) LoaderResult[string] {
	mecchiFolder := MecchiFolderExists()
	OutputResult(input, "Trying to locate existing Mecchi folder...")

	if !mecchiFolder.success {
		return LoaderResult[string]{success: false, error: "Mecchi folder does not exist"}
	} else {
		OutputResult(input, "Found Mecchi folder at "+mecchiFolder.result)
	}

	repo, err := git.PlainOpen(mecchiFolder.result)
	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	err = worktree.Pull(&git.PullOptions{
		ReferenceName: plumbing.Main,
	})

	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	packageJSON := GetScripts(mecchiFolder.result)

	if !packageJSON.success {
		return LoaderResult[string]{success: false, error: packageJSON.error}
	}

	var shell string
	var shellArgs []string

	if os.Getenv("SHELL") != "" {
		shell = os.Getenv("SHELL")
		shellArgs = []string{"-c", packageJSON.result.UpdateLinux}
	} else {
		shell = "cmd"
		shellArgs = []string{"/c", packageJSON.result.Update}
	}

	cmd := exec.Command(shell, shellArgs...)
	cmd.Dir = mecchiFolder.result

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	err = cmd.Start()
	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	scanner := bufio.NewScanner(stdout)

	go func() {
		for scanner.Scan() {
			OutputResult(input, scanner.Text())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return LoaderResult[string]{success: false, error: err.Error()}
	}

	return LoaderResult[string]{success: true}
}
