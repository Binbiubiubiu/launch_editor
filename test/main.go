package main

import (
	"log"

	. "github.com/Binbiubiubiu/launch-editor"
)

func main() {
	err := LaunchEditor("guess.go:59:20", "code")
	if err != nil {
		log.Fatalln(err)
	}
}
