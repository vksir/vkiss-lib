package template

import (
	"fmt"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"io"
	"os"
	"strings"
	"text/template"
)

type Template interface {
	Template() string
}

func Execute(t Template, w io.Writer) error {
	tmpl, err := template.New("template").Parse(t.Template())
	if err != nil {
		return errutil.Wrap(err)
	}
	err = tmpl.Execute(w, t)
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}

func ExecuteString(t Template) (string, error) {
	buf := &strings.Builder{}
	err := Execute(t, buf)
	return buf.String(), err
}

func ExecuteFile(t Template, path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o640)
	if err != nil {
		return errutil.Wrap(err)
	}
	defer func() { _ = f.Close() }()
	return Execute(t, f)
}

func ExecuteFiles(m map[string]Template) error {
	for p, t := range m {
		err := ExecuteFile(t, p)
		if err != nil {
			return fmt.Errorf("execute template %s failed: %w", p, err)
		}
	}
	return nil
}
