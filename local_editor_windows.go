package launch_editor

import (
	"path/filepath"
	"strings"
)

func getLocalEditor() (editor string) {
	output, _ := execCmd(`wmic process where "executablepath is not null" get executablepath`)
	runningProcesses := strings.Split(output, "\r\n")
	for i := 0; i < len(runningProcesses); i++ {
		fullProcessPath := strings.TrimSpace(runningProcesses[i])
		shortProcessName := filepath.Base(fullProcessPath)
		for _, v := range editorInfo {
			if v == shortProcessName {
				editor = fullProcessPath
				return
			}
		}
	}
	return
}
