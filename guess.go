package launch_editor

import (
	"os"

	"github.com/google/shlex"
)

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

	editor = getLocalEditor()

	if VISUAL, ok := os.LookupEnv("VISUAL"); ok {
		editor = VISUAL
		return
	} else if EDITOR, ok := os.LookupEnv("EDITOR"); ok {
		editor = EDITOR
		return
	}
	return
}
