package cmd

import (
	"github.com/spf13/cobra"
	"vkiss-lib/internal/cmd/ddnscmd"
	"vkiss-lib/pkg/cfg"
	"vkiss-lib/pkg/log"
)

func Execute() error {
	var cfgFile string

	root := &cobra.Command{
		Use:   "vkiss",
		Short: "Vkiss Tool",
		Long:  `Vkiss Tool`,
	}

	cobra.OnInitialize(func() {
		cfg.Init(cfgFile, cfg.DefaultConfig)
		log.Init(cfg.LogPath.Get(), cfg.LogLevel.Get())
	})

	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfg.DefaultConfPath,
		"config file")
	cfg.LogPath.Bind(root)
	cfg.LogLevel.Bind(root)
	root.AddCommand(ddnscmd.NewCmd())

	return root.Execute()
}
