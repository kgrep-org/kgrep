package main

import (
	"github.com/hbelmiro/kgrep/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
