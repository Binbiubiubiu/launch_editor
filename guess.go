package launch_editor

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	. "github.com/Binbiubiubiu/launch-editor/editor_info"
	"github.com/google/shlex"
	"github.com/samber/lo"
)

func execCmd(cmd string) (output string, err error) {
	shellArgs, err := shlex.Split(cmd)
	if err != nil {
		return
	}
	binary, err := exec.LookPath(shellArgs[0])
	if err != nil {
		return
	}

	err = syscall.Exec(binary, shellArgs, os.Environ())
	if err != nil {
		return
	}
	buf, err := io.ReadAll(os.Stdout)
	if err != nil {
		return
	}
	output = string(buf)
	return
}

func GuessEditor(specifiedEditor string) (editor string, args []string) {
	args, err := shlex.Split(specifiedEditor)
	if err != nil {
		editor = specifiedEditor
		return
	}
	if len(args) > 0 {
		editor = args[0]
		args = args[1:]
		return
	}

	if LAUNCH_EDITOR, ok := os.LookupEnv("LAUNCH_EDITOR"); ok {
		editor = LAUNCH_EDITOR
		return
	}

	if _, isWebContainer := os.LookupEnv("webcontainer"); isWebContainer {
		if EDITOR, ok := os.LookupEnv("EDITOR"); ok {
			editor = EDITOR
		} else {
			editor = "code"
		}
		return
	}

	var output string
	if IS_OSX {
		output, _ = execCmd("ps x -o comm=")
		processNames := lo.Keys(COMMON_EDITORS_OSX)
		processList := strings.Split(output, "\n")
		for _, processName := range processNames {
			if strings.Contains(output, processName) {
				editor = COMMON_EDITORS_OSX[processName]
				return
			}
			processNameWithoutApplications := strings.Replace(processName, "/Applications", "", 1)
			// Find editor installation not in /Applications.
			if strings.Contains(output, processNameWithoutApplications) {
				// Use the CLI command if one is specified
				if processName != COMMON_EDITORS_OSX[processName] {
					editor = COMMON_EDITORS_OSX[processName]
					return
				}
				// Use a partial match to find the running process path.  If one is found, use the
				// existing path since it can be running from anywhere.
				runningProcess, ok := lo.Find(processList, func(procName string) bool {
					return strings.HasSuffix(procName, processNameWithoutApplications)
				})
				if ok {
					editor = runningProcess
					return
				}
			}
		}
	} else if IS_WINDOWS {
		output, _ = execCmd(`powershell -NoProfile -Command "Get-CimInstance -Query \\"select executablepath from win32_process where executablepath is not null\\" | % { $_.ExecutablePath }"`)
		runningProcesses := strings.Split(output, `\r\n`)
		for i := 0; i < len(runningProcesses); i++ {
			fullProcessPath := strings.TrimSpace(runningProcesses[i])
			shortProcessName := filepath.Base(fullProcessPath)

			for _, v := range COMMON_EDITORS_WIN {
				if v == shortProcessName {
					editor = fullProcessPath
					return
				}
			}
		}
	} else if IS_LINUX {
		output, _ = execCmd("ps x --no-heading -o comm --sort=comm")
		processNames := lo.Keys(COMMON_EDITORS_LINUX)
		for _, processName := range processNames {
			if strings.Contains(output, processName) {
				editor = COMMON_EDITORS_LINUX[processName]
				return
			}
		}
	}

	if VISUAL, ok := os.LookupEnv("VISUAL"); ok {
		editor = VISUAL
		return
	} else if EDITOR, ok := os.LookupEnv("EDITOR"); ok {
		editor = EDITOR
		return
	}

	return
}
