package main

import (
	"log"

	. "github.com/Binbiubiubiu/launch-editor"
)

func main() {
	err := LaunchEditorWithName("guess.go:59:20", "code")
	if err != nil {
		log.Fatalln(err)
	}
}
