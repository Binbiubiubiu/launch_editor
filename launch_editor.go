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

	"github.com/matishsiao/goInfo"
	"github.com/samber/lo"
)

const (
	isLinux   = runtime.GOOS == "linux"
	isOsx     = runtime.GOOS == "darwin"
	isWindows = runtime.GOOS == "windows"
)

var osInfo goInfo.GoInfoObject
var childProcess *spawnProcess = nil
var positionRE = regexp.MustCompile(`:(\d+)(:(\d+))?$`)

func init() {
	osInfo, _ = goInfo.GetInfo()
}

func LaunchEditor(file string) error {
	return LaunchEditorWithName(file, "")
}

func LaunchEditorWithName(file string, specifiedEditor string) (err error) {
	fileName, lineNumber, columnNumber := parseFile(file)

	if _, err = os.Stat(fileName); err != nil && os.IsNotExist(err) {
		return
	}

	editor, args := guessEditor(specifiedEditor)
	if lo.IsEmpty(editor) {
		return &EditorProcessError{fileName: fileName}
	}

	if isLinux && strings.HasPrefix(fileName, "mnt") && regexp.MustCompile(`(?i)Microsoft`).MatchString(osInfo.Core) {
		// Assume WSL / "Bash on Ubuntu on Windows" is being used, and
		// that the file exists on the Windows file system.
		// `os.release()` is "4.4.0-43-Microsoft" in the current release
		// build of WSL, see: https://github.com/Microsoft/BashOnWindows/issues/423#issuecomment-221627364
		// When a Windows editor is specified, interop functionality can
		// handle the path translation, but only if a relative path is used.
		fileName, _ = filepath.Rel("", fileName)
	}

	if lineNumber > 0 {
		extraArgs := getArgumentsForPosition(editor, fileName, lineNumber, columnNumber)
		args = append(args, extraArgs...)
	} else {
		args = append(args, fileName)
	}

	if childProcess != nil && isTerminalEditor(editor) {
		childProcess.cancel()
	}

	if isWindows {
		args = append([]string{"/C", editor}, args...)
		childProcess = spawn("cmd.exe", args...)

	} else {
		childProcess = spawn(editor, args...)
	}
	childProcess.Stdin = os.Stdin
	childProcess.Stdout = os.Stdout
	childProcess.Stderr = os.Stderr

	if err = childProcess.Start(); err != nil {
		return &EditorProcessError{fileName: fileName, errorMessage: err.Error()}
	}

	if re := childProcess.Wait(); re != nil {
		childProcess = nil
		code := re.(*exec.ExitError).ExitCode()
		if code > 0 {
			return &EditorProcessError{fileName: fileName, errorMessage: fmt.Sprintf(`(code %v)`, code)}
		}
	}
	return
}

type EditorProcessError struct {
	fileName     string
	errorMessage string
}

func (e *EditorProcessError) Error() string {
	var msg string
	if lo.IsNotEmpty(e.errorMessage) {
		msg = fmt.Sprintf("The editor process exited with an error: %s", e.errorMessage)
	}
	return fmt.Sprintf("Could not open %s in the editor.%s", filepath.Base(e.fileName), msg)
}

func isTerminalEditor(editor string) bool {
	switch editor {
	case "vim", "emacs", "nano":
		return true
	}
	return false
}

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
