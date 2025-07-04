package main

import (
	"os"

	"github.com/hbelmiro/kgrep/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
