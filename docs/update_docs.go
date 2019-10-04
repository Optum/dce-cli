package main

import (
	"fmt"

	"github.com/Optum/dce-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./docs")
	if err != nil {
		fmt.Println(err)
	}
}
