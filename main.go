package main

import (
	"os"

	"github.com/goark/gocli/rwi"
	"github.com/goark/gpt-cli/facade"
)

func main() {
	facade.Execute(
		rwi.New(
			rwi.WithReader(os.Stdin),
			rwi.WithWriter(os.Stdout),
			rwi.WithErrorWriter(os.Stderr),
		),
		os.Args[1:],
	).ExitIfNotNormal()
}
