package main

import (
	"os"

	"github.com/pangu-studio/awesome-skills/cmd/weather/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
