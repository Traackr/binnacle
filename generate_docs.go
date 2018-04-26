package main

import (
	"log"

	"github.com/Traackr/binnacle/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs/commands")
	if err != nil {
		log.Fatal(err)
	}
}
