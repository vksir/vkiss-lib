package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vksir/vkiss-lib/internal/cmd/ddnscmd"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "vkiss",
		Short: "Vkiss Tool",
		Long:  `Vkiss Tool`,
	}

	root.AddCommand(ddnscmd.NewCmd())
	return root
}
