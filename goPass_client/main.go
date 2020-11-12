package main

import (
	"os"

	"github.com/ayushsherpa111/goPass/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
