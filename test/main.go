package main

import (
	"log"

	"github.com/Binbiubiubiu/launch_editor"
)

func main() {
	err := launch_editor.LaunchEditor("guess.go:10:20")
	if err != nil {
		log.Fatalln(err)
	}
}
