package installutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"github.com/vksir/vkiss-lib/pkg/util/fileutil"
	"github.com/vksir/vkiss-lib/thirdpkg/systemctl"
)

func InstallExec(dst string) error {
	src, err := os.Executable()
	if err != nil {
		return errutil.Wrap(err)
	}

	err = fileutil.MkDir(filepath.Dir(dst))
	if err != nil {
		return errutil.Wrap(err)
	}
	err = fileutil.Remove(dst)
	if err != nil {
		return errutil.Wrap(err)
	}
	err = fileutil.Cp(src, dst)
	if err != nil {
		return errutil.Wrap(err)
	}
	log.Info("copy exec", "src", src, "dst", dst)
	err = os.Chmod(dst, 0755)
	if err != nil {
		return errutil.Wrap(err)
	}

	err = fileutil.Remove(src)
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}

func InstallService(svc *systemctl.Service, exec string) error {
	err := svc.Stop()
	if err != nil {
		log.Warn("stop service failed", "err", err)
	}

	err = InstallExec(exec)
	if err != nil {
		return errutil.Wrap(err)
	}
	err = svc.Deploy()
	if err != nil {
		return errutil.Wrap(err)
	}
	err = svc.DaemonReload()
	if err != nil {
		return errutil.Wrap(err)
	}
	err = svc.Enable()
	if err != nil {
		return errutil.Wrap(err)
	}
	err = svc.Start()
	if err != nil {
		cmdStr := fmt.Sprintf("systemctl restart %s", svc.Name)
		log.Error("start server failed", "cmd", cmdStr, "err", err)
	} else {
		cmdStr := fmt.Sprintf("systemctl status %s", svc.Name)
		log.Info("install and start server success", "cmd", cmdStr)
	}
	return nil
}
