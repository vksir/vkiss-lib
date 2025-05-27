package fileutil

import (
	"archive/zip"
	"errors"
	"github.com/vksir/vkiss-lib/pkg/log"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Exist(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func Write(path string, content []byte) error {
	return os.WriteFile(path, content, 0640)
}

func Rm(path string) error {
	return os.RemoveAll(path)
}

func MkDir(paths ...string) error {
	for _, p := range paths {
		if err := os.MkdirAll(p, 0o755); err != nil {
			return errutil.Wrap(err)
		}
	}
	return nil
}

func ClearDir(path string) error {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, d := range dirs {
		err = os.RemoveAll(filepath.Join(path, d.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func Cp(src, dst string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return errutil.Wrap(err)
	}

	if srcStat.IsDir() {
		srcFs := os.DirFS(src)
		err = os.CopyFS(dst, srcFs)
		if err != nil {
			return errutil.Wrap(err)
		}
		return nil
	}

	// 如果目标是目录，则将文件复制到目录中
	newDst := dst
	dstStat, err := os.Stat(dst)
	if err == nil && dstStat.IsDir() {
		newDst = filepath.Join(dst, srcStat.Name())
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return errutil.Wrap(err)
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.OpenFile(newDst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, srcStat.Mode())
	if err != nil {
		return errutil.Wrap(err)
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}

func Mv(src, dst string) error {
	return os.Rename(src, dst)
}

func Zip(src string, dst string) error {
	zipFile, err := os.Create(dst)
	if err != nil {
		return errutil.Wrap(err)
	}
	defer func() { _ = zipFile.Close() }()

	zipWriter := zip.NewWriter(zipFile)
	defer func() { _ = zipWriter.Close() }()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errutil.WrapPathErr("WalkFunc", path, err)
		}

		// 获取相对路径
		relativePath := filepath.ToSlash(strings.TrimPrefix(path, src))
		if info.IsDir() {
			relativePath += "/"
		}

		w, err := zipWriter.Create(relativePath)
		if err != nil {
			return errutil.WrapPathErr("zipWriter.Create", path, err)
		}

		r, err := os.Open(path)
		if err != nil {
			return errutil.WrapPathErr("os.Open", path, err)
		}
		defer func() { _ = r.Close() }()

		_, err = io.Copy(w, r)
		if err != nil {
			return errutil.WrapPathErr("io.Copy", path, err)
		}
		return err
	})
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}

func Unzip(src string, dst string) error {
	zipFile, err := zip.OpenReader(src)
	if err != nil {
		return errutil.WrapPathErr("zip.OpenReader", src, err)
	}
	defer func() { _ = zipFile.Close() }()

	// 遍历 ZIP 文件中的所有条目
	for _, f := range zipFile.File {
		path := filepath.Join(dst, f.Name)

		// 如果是目录，仅创建目录
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return errutil.WrapPathErr("os.MkdirAll", path, err)
			}
			continue
		}

		err = func() error {
			// 创建目标文件
			w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return errutil.WrapPathErr("os.OpenFile", path, err)
			}
			defer func() { _ = w.Close() }()

			// 打开 ZIP 文件中的条目并复制内容到目标文件中
			r, err := f.Open()
			if err != nil {
				return errutil.WrapPathErr("Open", path, err)
			}
			defer func() { _ = r.Close() }()

			_, err = io.Copy(w, r)
			if err != nil {
				return errutil.WrapPathErr("io.Copy", path, err)
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

func InstallSelf(path string) error {
	src, err := os.Executable()
	if err != nil {
		return errutil.Wrap(err)
	}

	dst := path

	err = Rm(dst)
	if err != nil {
		return errutil.Wrap(err)
	}
	err = Cp(src, dst)
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
