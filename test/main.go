package main

import (
	"log"

	. "github.com/Binbiubiubiu/launch-editor"
)

func main() {
	err := LaunchEditor("guess.go:10:20")
	if err != nil {
		log.Fatalln(err)
	}
}
