package steam

import (
	"context"
	"github.com/vksir/vkiss-lib/pkg/cfg"
	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/subprocess"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"os/exec"
	"path/filepath"
)

var (
	CfgSteamcmdPath = cfg.NewFlag[string]("steamcmd", "steamcmd",
		"steamcmd executable path").SetDefault("steamcmd")
)

type Steamcmd struct {
	executePath     string
	logger          *log.Logger
	forceInstallDir string
}

func NewSteamcmd(path string, log *log.Logger) *Steamcmd {
	return &Steamcmd{executePath: path, logger: log.With("tag", "steamcmd")}
}

func (s *Steamcmd) SetForceInstallDir(dir string) *Steamcmd {
	s.forceInstallDir = dir
	return s
}

func (s *Steamcmd) DoWorkShopDownloadItem(ctx context.Context, appId, publishedFileId string) error {
	s.logger.InfoC(ctx, "begin DoWorkShopDownloadItem", "appId", appId, "publishedFileId", publishedFileId)
	extArgs := []string{"+workshop_download_item", appId, publishedFileId}
	err := s.execute(ctx, extArgs)
	if err != nil {
		s.logger.ErrorC(ctx, "DoWorkShopDownloadItem failed", "appId", appId, "publishedFileId", publishedFileId, "err", err)
		return errutil.Wrap(err)
	}
	s.logger.InfoC(ctx, "DoWorkShopDownloadItem success", "appId", appId, "publishedFileId", publishedFileId)
	return nil
}

func (s *Steamcmd) DoAppUpdate(ctx context.Context, appId string) error {
	s.logger.InfoC(ctx, "begin DoAppUpdate", "appId", appId)
	extArgs := []string{"+app_update", appId, "validate"}
	err := s.execute(ctx, extArgs)
	if err != nil {
		s.logger.ErrorC(ctx, "DoAppUpdate failed", "appId", appId, "err", err)
		return errutil.Wrap(err)
	}
	s.logger.InfoC(ctx, "DoAppUpdate success", "appId", appId)
	return nil
}

func (s *Steamcmd) execute(ctx context.Context, extArgs []string) error {
	args := []string{"+login", "anonymous"}
	if s.forceInstallDir != "" {
		args = append(args, "+force_install_dir", s.forceInstallDir)
	}
	args = append(args, extArgs...)
	args = append(args, "+quit")

	exe, err := exec.LookPath(s.executePath)
	if err != nil {
		return errutil.Wrap(err)
	}
	process := subprocess.New("steamcmd", exe, args).
		SetLogger(s.logger).
		SetDir(filepath.Dir(exe))

	logger := s.logger.With("output", "steamcmd")
	process.RegisterOutFunc("log", func(line *string) {
		logger.Debug(*line)
	})

	err = process.Start(ctx)
	if err != nil {
		return errutil.Wrap(err)
	}

	err = process.Wait()
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}
