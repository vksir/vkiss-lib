package util

import (
	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"github.com/vksir/vkiss-lib/pkg/util/fileutil"
	"os"
)

const (
	ExePath = "/bin/vkiss"
)

func InstallSelf() error {
	src, err := os.Executable()
	if err != nil {
		return errutil.Wrap(err)
	}

	dst := ExePath

	err = fileutil.Rm(dst)
	if err != nil {
		return errutil.Wrap(err)
	}
	err = fileutil.Cp(src, dst)
	if err != nil {
		return errutil.Wrap(err)
	}
	err = os.Chmod(dst, 0755)
	if err != nil {
		return errutil.Wrap(err)
	}

	log.Info("install success", "bin", dst)
	return nil
}
