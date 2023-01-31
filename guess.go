package launch_editor

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/Binbiubiubiu/launch-editor/editor_info"
	"github.com/google/shlex"
	"github.com/samber/lo"
)

func execCmd(cmd string) (output string, err error) {
	shellArgs, err := shlex.Split(cmd)
	if err != nil {
		return
	}

	cp := spawn(shellArgs[0], shellArgs[1:]...)
	buf, err := cp.CombinedOutput()
	if err != nil {
		return
	}
	output = string(buf)
	return
}

func guessEditor(specifiedEditor string) (editor string, args []string) {
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
	if isOsx {
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
	} else if isWindows {
		output, _ = execCmd(`wmic process where "executablepath is not null" get executablepath`)
		runningProcesses := strings.Split(output, "\r\n")
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
	} else if isLinux {
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
