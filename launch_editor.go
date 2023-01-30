package launch_editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/matishsiao/goInfo"
	"github.com/samber/lo"
)

var OS goInfo.GoInfoObject

const (
	IS_LINUX   = runtime.GOOS == "linux"
	IS_OSX     = runtime.GOOS == "darwin"
	IS_WINDOWS = runtime.GOOS == "windows"
)

func init() {
	OS, _ = goInfo.GetInfo()
}

type ErrorCallback func(filename string, errorMessage string)

func wrapErrorCallback(cb ErrorCallback) ErrorCallback {
	return func(fileName string, errorMessage string) {
		fmt.Println()
		color.Red("Could not open %s in the editor.", filepath.Base(fileName))
		if errorMessage != "" {
			color.Red("The editor process exited with an error: %s", errorMessage)
		}
		if cb != nil {
			cb(fileName, errorMessage)
		}

	}
}

func isTerminalEditor(editor string) bool {
	switch editor {
	case "vim", "emacs", "nano":
		return true
	}
	return false
}

var positionRE = regexp.MustCompile(`:(\d+)(:(\d+))?$`)

func parseFile(file string) (fileName string, lineNumber int, columnNumber int) {
	fileName = positionRE.ReplaceAllLiteralString(file, "")
	match := positionRE.FindAllStringSubmatch(file, 1)
	matchSlice := match[0]
	var err error
	lineNumber, err = strconv.Atoi(matchSlice[1])
	if err != nil {
		lineNumber = 0
	}
	columnNumber, err = strconv.Atoi(matchSlice[3])
	if err != nil {
		columnNumber = 0
	}
	return
}

var childProcess *exec.Cmd

func LaunchEditor(file string, specifiedEditor string, onErrorCallback ErrorCallback) (err error) {
	fileName, lineNumber, columnNumber := parseFile(file)

	if _, err = os.Stat(fileName); err != nil && os.IsNotExist(err) {
		return
	}

	onErrorCallback = wrapErrorCallback(onErrorCallback)

	editor, args := GuessEditor(specifiedEditor)
	if lo.IsEmpty(editor) {
		onErrorCallback(fileName, "")
		return
	}

	if IS_LINUX && strings.HasPrefix(fileName, "mnt") && regexp.MustCompile(`(?!)Microsoft`).MatchString(OS.Core) {
		// Assume WSL / "Bash on Ubuntu on Windows" is being used, and
		// that the file exists on the Windows file system.
		// `os.release()` is "4.4.0-43-Microsoft" in the current release
		// build of WSL, see: https://github.com/Microsoft/BashOnWindows/issues/423#issuecomment-221627364
		// When a Windows editor is specified, interop functionality can
		// handle the path translation, but only if a relative path is used.
		fileName, _ = filepath.Rel("", fileName)
	}

	if lineNumber > 0 {
		extraArgs := GetArgumentsForPosition(editor, fileName, lineNumber, columnNumber)
		args = append(args, extraArgs...)
	} else {
		args = append(args, fileName)
	}

	if childProcess != nil && isTerminalEditor(editor) {
		syscall.Kill(childProcess.Process.Pid, syscall.SIGKILL)
	}

	if IS_WINDOWS {
		args = append([]string{"/C", editor}, args...)
		childProcess = exec.Command("cmd.exe", args...)

	} else {
		childProcess = exec.Command(editor, args...)
	}
	childProcess.Stdin = os.Stdin
	childProcess.Stdout = os.Stdout
	childProcess.Stderr = os.Stderr

	if err := childProcess.Start(); err != nil {
		onErrorCallback(fileName, err.Error())
	}

	if re := childProcess.Wait(); re != nil {
		childProcess = nil
		err = re
		re := re.(*exec.ExitError)
		code := re.ProcessState.ExitCode()
		if code > 0 {
			onErrorCallback(fileName, fmt.Sprintf(`(code %v)`, code))
		}
	}
	return
}
