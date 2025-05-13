package main

import (
	"github.com/spf13/cobra"
	"github.com/vksir/vkiss-lib/internal/cmd"
)

func main() {
	err := cmd.Execute()
	cobra.CheckErr(err)
}
