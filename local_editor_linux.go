package launch_editor

func getLocalEditor() (editor string) {
	output, _ := execCmd("ps x --no-heading -o comm --sort=comm")
	processNames := lo.Keys(editorInfo)
	for _, processName := range processNames {
		if strings.Contains(output, processName) {
			editor = editorInfo[processName]
			return
		}
	}
}
