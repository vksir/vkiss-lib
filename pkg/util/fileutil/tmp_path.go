package fileutil

import (
	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"path/filepath"
)

type TmpPath string

func NewTmpPath(elem ...string) (TmpPath, error) {
	path := filepath.Join(elem...)
	if Exist(path) {
		log.Error("tmp path already exists, clear it", "path", path)
		err := Rm(path)
		if err != nil {
			return "", errutil.Wrap(err)
		}
	}
	return TmpPath(path), nil
}

func (p TmpPath) String() string {
	return string(p)
}

func (p TmpPath) Clear() {
	err := Rm(string(p))
	if err != nil {
		log.Error("clear tmp path failed", "path", p, "err", err)
	} else {
		log.Info("clear tmp path success", "path", p)
	}
}
