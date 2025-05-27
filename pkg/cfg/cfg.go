package cfg

import (
	_ "embed"
	"github.com/spf13/viper"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"github.com/vksir/vkiss-lib/pkg/util/fileutil"
	"path/filepath"
)

var (
	LogLevel = NewFlag[string]("log-level", "log.level",
		"log level").SetDefault("info").SetPersistent(true)
	LogPath = NewFlag[string]("log-path", "log.path",
		"log path").SetPersistent(true)
)

func Init(path string, defaultConfig string) {
	if !fileutil.Exist(path) {
		err := fileutil.MkDir(filepath.Dir(path))
		errutil.Check(err)
		err = fileutil.Write(path, []byte(defaultConfig))
		errutil.Check(err)
	}

	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	errutil.Check(err)
}
