// Package main is the CLI entry point for Nous.
package main

import (
	"os"

	"github.com/cajohnson0125/Nous/internal"
	"github.com/cajohnson0125/Nous/cmd/nous/root"
)

func main() {
	app := internal.NewApp()

	if err := root.Execute(app); err != nil {
		os.Exit(1)
	}
}
