package cfg

import (
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"vkiss-lib/pkg/util/errutil"
	"vkiss-lib/pkg/util/fileutil"
)

const DefaultConfPath = "/etc/vkiss/config.toml"

//go:embed config.toml
var DefaultConfig string

var (
	LogLevel = NewFlag[string]("log-level", "log.level",
		"log level").SetDefault("info").SetPersistent(true)
	LogPath = NewFlag[string]("log-path", "log.path",
		"log path (default ./app.log)").SetPersistent(true)

	DdnsListen = NewFlag[string]("listen", "ddns.listen",
		"listen address").SetDefault(":5801")
	DdnsServer = NewFlag[string]("endpoint", "ddns.endpoint",
		"endpoint address")

	DdnsTcSecretId = NewFlag[string]("secret-id", "ddns.tencent_cloud.secret_id",
		"tencent_cloud ddns secret id")
	DdnsTcSecretKey = NewFlag[string]("secret-key", "ddns.tencent_cloud.secret_key",
		"tencent_cloud ddns secret key")
	DdnsTcDomain = NewFlag[string]("domain", "ddns.tencent_cloud.domain",
		"tencent_cloud ddns domain")
	DdnsTcSubDomain = NewFlag[string]("sub-domain", "ddns.tencent_cloud.sub_domain",
		"tencent_cloud ddns sub domain")
	DdnsTcRecordId = NewFlag[uint64]("record-id", "ddns.tencent_cloud.record_id",
		"tencent_cloud ddns record id")
	DdnsTcRecordLine = NewFlag[string]("record-line", "ddns.tencent_cloud.record_line",
		"tencent_cloud ddns record line")
	DdnsTcValue = NewFlag[string]("value", "ddns.tencent_cloud.value",
		"tencent_cloud ddns value")
)

func Save() error {
	err := viper.WriteConfig()
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}

func Init(path string, defaultConfig string) {
	if !fileutil.Exist(path) {
		err := fileutil.MkDir(filepath.Dir(path))
		cobra.CheckErr(err)
		err = fileutil.Write(path, []byte(defaultConfig))
		cobra.CheckErr(err)
	}

	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
