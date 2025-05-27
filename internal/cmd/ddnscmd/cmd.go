package ddnscmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/vksir/vkiss-lib/internal/constant"
	"github.com/vksir/vkiss-lib/internal/ddns"
	"github.com/vksir/vkiss-lib/pkg/cfg"
	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"github.com/vksir/vkiss-lib/pkg/util/fileutil"
	"github.com/vksir/vkiss-lib/thirdpkg/systemctl"
	"github.com/vksir/vkiss-lib/thirdpkg/tencentcloud"
	"time"
)

var (
	DdnsListen = cfg.NewFlag[string]("listen", "ddns.listen",
		"listen address").SetDefault(":5801")
	DdnsEndpoint = cfg.NewFlag[string]("endpoint", "ddns.endpoint",
		"endpoint address")
	DdnsInterval = cfg.NewFlag[int]("interval", "ddns.interval",
		"monitor loop interval (minute)").SetDefault(20)

	DdnsTcSecretId = cfg.NewFlag[string]("secret-id", "ddns.tencent_cloud.secret_id",
		"tencent_cloud ddns secret id")
	DdnsTcSecretKey = cfg.NewFlag[string]("secret-key", "ddns.tencent_cloud.secret_key",
		"tencent_cloud ddns secret key")
	DdnsTcDomain = cfg.NewFlag[string]("domain", "ddns.tencent_cloud.domain",
		"tencent_cloud ddns domain")
	DdnsTcSubDomain = cfg.NewFlag[string]("sub-domain", "ddns.tencent_cloud.sub_domain",
		"tencent_cloud ddns sub domain")
	DdnsTcRecordId = cfg.NewFlag[uint64]("record-id", "ddns.tencent_cloud.record_id",
		"tencent_cloud ddns record id")
	DdnsTcRecordLine = cfg.NewFlag[string]("record-line", "ddns.tencent_cloud.record_line",
		"tencent_cloud ddns record line")
	DdnsTcValue = cfg.NewFlag[string]("value", "ddns.tencent_cloud.value",
		"tencent_cloud ddns value")
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ddns",
	}

	serverCmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serve(DdnsListen.Get())
		},
	}
	DdnsListen.Bind(serverCmd)

	monitorCmd := &cobra.Command{
		Use: "monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return monitor(DdnsEndpoint.Get())
		},
	}
	DdnsEndpoint.Bind(monitorCmd)
	DdnsInterval.Bind(monitorCmd)
	addTencentCloudDdnsFlags(monitorCmd)

	refreshCmd := &cobra.Command{
		Use: "refresh",
		RunE: func(cmd *cobra.Command, args []string) error {
			return refresh(DdnsTcValue.Get())
		},
	}
	DdnsTcValue.Bind(refreshCmd)
	addTencentCloudDdnsFlags(refreshCmd)

	installCmd := newInstallCmd()

	cmd.AddCommand(serverCmd)
	cmd.AddCommand(monitorCmd)
	cmd.AddCommand(refreshCmd)
	cmd.AddCommand(installCmd)
	return cmd
}

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "install",
	}

	installServerCmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return installServer()
		},
	}

	installMonitorCmd := &cobra.Command{
		Use: "monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return installMonitor()
		},
	}

	cmd.AddCommand(installServerCmd)
	cmd.AddCommand(installMonitorCmd)
	return cmd
}

func addTencentCloudDdnsFlags(cmd *cobra.Command) {
	DdnsTcSecretId.Bind(cmd)
	DdnsTcSecretKey.Bind(cmd)
	DdnsTcDomain.Bind(cmd)
	DdnsTcSubDomain.Bind(cmd)
	DdnsTcRecordId.Bind(cmd)
	DdnsTcRecordLine.Bind(cmd)
}

func serve(listen string) error {
	e := gin.Default()
	ddns.LoadRouter(&e.RouterGroup)
	log.Info("starting serv", "listen", listen)
	return e.Run(listen)
}

func monitor(endpoint string) error {
	interval := DdnsInterval.Get()
	log.Info("starting monitor", "endpoint", endpoint, "interval", interval)

	// 失败时快循环，成功时慢循环
	curMyIp := ""
	for {
		myIp, err := ddns.GetMyIp(endpoint)
		if err != nil {
			log.Error(err.Error())
			time.Sleep(time.Minute)
			continue
		}

		if myIp == curMyIp {
			log.Debug("myIp has not changed, do nothing", "myIp", myIp)
			time.Sleep(time.Duration(interval) * time.Minute)
			continue
		}

		err = refresh(myIp)
		if err != nil {
			log.Error(err.Error())
			time.Sleep(time.Minute)
			continue
		}
		curMyIp = myIp
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func refresh(myIp string) error {
	log.Warn("begin refresh myIp", "myIp", myIp)
	req := &tencentcloud.ModifyDynamicDNSRequest{
		Domain:     DdnsTcDomain.Get(),
		SubDomain:  DdnsTcSubDomain.Get(),
		RecordId:   DdnsTcRecordId.Get(),
		RecordLine: DdnsTcRecordLine.Get(),
		Value:      myIp,
	}
	secret := &tencentcloud.Secret{
		Id:  DdnsTcSecretId.Get(),
		Key: DdnsTcSecretKey.Get(),
	}
	info, err := tencentcloud.ModifyDynamicDns(req, secret)
	if err != nil {
		return errutil.Wrap(fmt.Errorf("tencentcloud.ModifyDynamicDns failed: info=%s, err=%w", info, err))
	}
	log.Warn("refresh myIp success", "myIp", myIp)
	return nil
}

var (
	serverService = &systemctl.Service{
		Name:             "ddns-server",
		Description:      "ddns server",
		ExecStart:        fmt.Sprintf("%s ddns server -c %s", constant.ExePath, constant.ConfPath),
		RestartOnFailure: true,
	}
	monitorService = &systemctl.Service{
		Name:             "ddns-monitor",
		Description:      "ddns monitor",
		ExecStart:        fmt.Sprintf("%s ddns monitor -c %s", constant.ExePath, constant.ConfPath),
		RestartOnFailure: true,
	}
)

func installServer() error {
	err := fileutil.InstallSelf(constant.ExePath)
	if err != nil {
		return errutil.Wrap(err)
	}

	err = serverService.Deploy()
	if err != nil {
		return errutil.Wrap(err)
	}

	err = serverService.Enable()
	if err != nil {
		return errutil.Wrap(err)
	}

	err = serverService.Restart()
	if err != nil {
		cmd := fmt.Sprintf("systemctl restart %s", serverService.Name)
		log.Error("start server failed", "cmd", cmd, "err", err)
	} else {
		cmd := fmt.Sprintf("systemctl status %s", serverService.Name)
		log.Info("install and start server success", "cmd", cmd)
	}
	return nil
}
func installMonitor() error {
	err := fileutil.InstallSelf(constant.ExePath)
	if err != nil {
		return errutil.Wrap(err)
	}

	err = monitorService.Deploy()
	if err != nil {
		return errutil.Wrap(err)
	}

	err = monitorService.Enable()
	if err != nil {
		return errutil.Wrap(err)
	}

	err = monitorService.Restart()
	if err != nil {
		cmd := fmt.Sprintf("systemctl restart %s", monitorService.Name)
		log.Error("start monitor failed", "cmd", cmd, "err", err)
	} else {
		cmd := fmt.Sprintf("systemctl status %s", monitorService.Name)
		log.Info("install and start monitor success", "cmd", cmd)
	}
	return nil
}
