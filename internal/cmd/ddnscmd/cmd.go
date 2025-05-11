package ddnscmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"time"
	"vkiss-lib/internal/ddns"
	"vkiss-lib/pkg/cfg"
	"vkiss-lib/pkg/log"
	"vkiss-lib/pkg/util"
	"vkiss-lib/pkg/util/errutil"
	"vkiss-lib/thirdpkg/systemctl"
	"vkiss-lib/thirdpkg/tencentcloud"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ddns",
	}

	serverCmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serve(cfg.DdnsListen.Get())
		},
	}
	cfg.DdnsListen.Bind(serverCmd)

	monitorCmd := &cobra.Command{
		Use: "monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return monitor(cfg.DdnsEndpoint.Get())
		},
	}
	cfg.DdnsEndpoint.Bind(monitorCmd)
	cfg.DdnsInterval.Bind(monitorCmd)
	addTencentCloudDdnsFlags(monitorCmd)

	refreshCmd := &cobra.Command{
		Use: "refresh",
		RunE: func(cmd *cobra.Command, args []string) error {
			return refresh(cfg.DdnsTcValue.Get())
		},
	}
	cfg.DdnsTcValue.Bind(refreshCmd)
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
	cfg.DdnsTcSecretId.Bind(cmd)
	cfg.DdnsTcSecretKey.Bind(cmd)
	cfg.DdnsTcDomain.Bind(cmd)
	cfg.DdnsTcSubDomain.Bind(cmd)
	cfg.DdnsTcRecordId.Bind(cmd)
	cfg.DdnsTcRecordLine.Bind(cmd)
}

func serve(listen string) error {
	e := gin.Default()
	g := e.Group("/")
	ddns.LoadRouter(g)
	log.Info("starting serv", "listen", listen)
	return e.Run(listen)
}

func monitor(endpoint string) error {
	interval := cfg.DdnsInterval.Get()
	log.Info("starting monitor", "endpoint", endpoint, "interval", interval)

	curMyIp := ""
	isFirst := true
	for {
		if isFirst {
			isFirst = false
		} else {
			time.Sleep(time.Duration(interval) * time.Minute)
		}

		myIp, err := ddns.GetMyIp(endpoint)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		if myIp == curMyIp {
			log.Debug("myIp has not changed, do nothing", "myIp", myIp)
			continue
		}

		err = refresh(myIp)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		curMyIp = myIp
	}
}

func refresh(myIp string) error {
	log.Warn("begin refresh myIp", "myIp", myIp)
	req := &tencentcloud.ModifyDynamicDNSRequest{
		Domain:     cfg.DdnsTcDomain.Get(),
		SubDomain:  cfg.DdnsTcSubDomain.Get(),
		RecordId:   cfg.DdnsTcRecordId.Get(),
		RecordLine: cfg.DdnsTcRecordLine.Get(),
		Value:      myIp,
	}
	secret := &tencentcloud.Secret{
		Id:  cfg.DdnsTcSecretId.Get(),
		Key: cfg.DdnsTcSecretKey.Get(),
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
		ExecStart:        fmt.Sprintf("%s ddns server -c %s", util.ExePath, cfg.DefaultConfPath),
		RestartOnFailure: true,
	}
	monitorService = &systemctl.Service{
		Name:             "ddns-monitor",
		Description:      "ddns monitor",
		ExecStart:        fmt.Sprintf("%s ddns monitor -c %s", util.ExePath, cfg.DefaultConfPath),
		RestartOnFailure: true,
	}
)

func installServer() error {
	err := util.InstallSelf()
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
	err := util.InstallSelf()
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
