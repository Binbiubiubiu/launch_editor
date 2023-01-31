package main

import (
	"log"

	launch_editor "github.com/Binbiubiubiu/launch-editor"
)

func main() {
	err := launch_editor.LaunchEditor("guess.go:10:20")
	if err != nil {
		log.Fatalln(err)
	}
}
