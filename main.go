package main

import (
	"github.com/spf13/cobra"
	"vkiss-lib/internal/cmd"
)

func main() {
	err := cmd.Execute()
	cobra.CheckErr(err)
}
