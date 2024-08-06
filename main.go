package main

import (
	"os"

	"github.com/raystack/raccoon/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
