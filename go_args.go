package launch_editor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func GetArgumentsForPosition(editor string, fileName string, lineNumber int, columnNumber int) []string {
	eidtorBasename := regexp.MustCompile(`\.(exe|cmd|bat)$`).ReplaceAllString(filepath.Base(editor), "")
	switch eidtorBasename {
	case "atom", "Atom", "Atom Beta", "subl", "sublime", "sublime_text", "wstorm", "charm":
		return []string{fmt.Sprintf("%s:%v:%v", fileName, lineNumber, columnNumber)}
	case "notepad++":
		return []string{fmt.Sprintf("-n%v", lineNumber), fmt.Sprintf("-c%v", columnNumber), fileName}
	case "vim", "mvim":
		return []string{fmt.Sprintf("+call cursor(%v,%v)", lineNumber, columnNumber), fileName}
	case "joe", "gvim":
		return []string{fmt.Sprintf("+%v", lineNumber), fileName}
	case "emacs", "emacsclient":
		return []string{fmt.Sprintf("+%v:%v", lineNumber, columnNumber), fileName}
	case "rmate", "mate", "mine":
		return []string{"--line", strconv.Itoa(lineNumber), fileName}
	case "code", "Code", "code-insiders", "Code - Insiders", "codium", "vscodium", "VSCodium":
		return []string{"-r", "-g", fmt.Sprintf("%s:%v:%v", fileName, lineNumber, columnNumber)}
	case "appcode", "clion", "clion64", "idea", "idea64", "phpstorm", "phpstorm64", "pycharm", "pycharm64", "rubymine", "rubymine64", "webstorm", "webstorm64", "goland", "goland64", "rider", "rider64":
		return []string{"--line", strconv.Itoa(lineNumber), "--column", strconv.Itoa(columnNumber), fileName}
	}

	if _, ok := os.LookupEnv("LAUNCH_EDITOR"); ok {
		return []string{fileName, strconv.Itoa(lineNumber), strconv.Itoa(columnNumber)}
	}

	// For all others, drop the lineNumber until we have
	// a mapping above, since providing the lineNumber incorrectly
	// can result in errors or confusing behavior.
	return []string{fileName}
}
