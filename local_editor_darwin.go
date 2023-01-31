package launch_editor

func getLocalEditor() (editor string) {
	output, _ := execCmd("ps x -o comm=")
	processNames := lo.Keys(editorInfo)
	processList := strings.Split(output, "\n")
	for _, processName := range processNames {
		if strings.Contains(output, processName) {
			editor = editorInfo[processName]
			return
		}
		processNameWithoutApplications := strings.Replace(processName, "/Applications", "", 1)
		// Find editor installation not in /Applications.
		if strings.Contains(output, processNameWithoutApplications) {
			// Use the CLI command if one is specified
			if processName != editorInfo[processName] {
				editor = editorInfo[processName]
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
	return
}
