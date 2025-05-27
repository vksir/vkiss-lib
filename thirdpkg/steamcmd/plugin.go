package steamcmd

import (
	"github.com/vksir/vkiss-lib/pkg/cfg"
)

var (
	CfgSteamcmdPath = cfg.NewFlag[string]("steamcmd", "steamcmd",
		"steamcmd executable path").SetDefault("steamcmd")
)
