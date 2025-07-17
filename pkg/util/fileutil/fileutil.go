package fileutil

import (
	"archive/zip"
	"errors"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var Home = func() string {
	p, err := os.UserHomeDir()
	errutil.Check(err)
	return p
}()

var Executable = func() string {
	p, err := os.Executable()
	errutil.Check(err)
	return p
}()

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
			return errutil.WrapPath("MkdirAll", p, err)
		}
	}
	return nil
}

func MkdirAndReturn(path string) (string, error) {
	if !Exist(path) {
		err := MkDir(path)
		if err != nil {
			return "", errutil.Wrap(err)
		}
	}
	return path, nil
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
			return err
		}

		// 获取相对路径
		relativePath := filepath.ToSlash(strings.TrimPrefix(path, src))
		if info.IsDir() {
			relativePath += "/"
		}

		dstWriter, err := zipWriter.Create(relativePath)
		if err != nil {
			return errutil.WrapPath("zipWriter.Create", path, err)
		}

		srcReader, err := os.Open(path)
		if err != nil {
			return errutil.WrapPath("os.Open", path, err)
		}
		defer func() { _ = srcReader.Close() }()

		_, err = io.Copy(dstWriter, srcReader)
		if err != nil {
			return errutil.WrapPath("io.Copy", path, err)
		}
		return err
	})
	return err
}

func Unzip(src string, dst string) error {
	err := MkDir(dst)
	if err != nil {
		return errutil.WrapPath("MkDir", dst, err)
	}

	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return errutil.WrapPath("zip.OpenReader", src, err)
	}
	defer func() { _ = zipReader.Close() }()

	// 遍历 ZIP 文件中的所有条目
	for _, srcFile := range zipReader.File {
		dstPath := filepath.Join(dst, srcFile.Name)

		// 如果是目录，仅创建目录
		if srcFile.FileInfo().IsDir() {
			err = MkDir(dstPath)
			if err != nil {
				return errutil.WrapPath("MkDir", dstPath, err)
			}
			continue
		}

		err = func() error {
			// 创建目标文件
			dstWriter, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcFile.Mode())
			if err != nil {
				return errutil.WrapPath("os.OpenFile", dstPath, err)
			}
			defer func() { _ = dstWriter.Close() }()

			// 打开 ZIP 文件中的条目
			srcReader, err := srcFile.Open()
			if err != nil {
				return errutil.WrapPath("zip.File.Open", dstPath, err)
			}
			defer func() { _ = srcReader.Close() }()

			// 复制内容到目标文件中
			_, err = io.Copy(dstWriter, srcReader)
			if err != nil {
				return errutil.WrapPath("io.Copy", dstPath, err)
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
	return nil
}

func init() {
	var err error
	Home, err = os.UserHomeDir()
	errutil.Check(err)
}
