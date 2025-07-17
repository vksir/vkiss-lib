package main

import (
	"github.com/spf13/cobra"
	"github.com/vksir/vkiss-lib/assets"
	"github.com/vksir/vkiss-lib/internal/cmd"
	"github.com/vksir/vkiss-lib/internal/constant"
	"github.com/vksir/vkiss-lib/pkg/cfg"
	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
)

func main() {
	root := cmd.NewRootCmd()

	var cfgFile string
	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", constant.ConfPath,
		"config file")

	cobra.OnInitialize(func() {
		cfg.Init(cfgFile, assets.DefaultConfig)
		log.Init(cfg.LogPath.Get(), cfg.LogLevel.Get())
	})
	err := root.Execute()
	errutil.Check(err)
}
