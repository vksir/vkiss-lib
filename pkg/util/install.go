package util

import (
	"os"
	"vkiss-lib/pkg/log"
	"vkiss-lib/pkg/util/errutil"
	"vkiss-lib/pkg/util/fileutil"
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
