package systemctl

import (
	"bytes"
	_ "embed"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
	"vkiss-lib/pkg/log"
	"vkiss-lib/pkg/util/cmdutil"
	"vkiss-lib/pkg/util/errutil"
	"vkiss-lib/pkg/util/fileutil"
)

const (
	systemConfDir = "/etc/systemd/system/"
)

//go:embed service.tmpl
var service string

type Service struct {
	Name             string
	Description      string
	WorkingDirectory string
	ExecStart        string
	RestartOnFailure bool
	User             string
	Group            string
}

func (s *Service) ServiceName() string {
	return fmt.Sprintf("%s.service", s.Name)
}

func (s *Service) Deploy() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("no support %s", runtime.GOOS)
	}

	buf, err := s.genConfig()
	if err != nil {
		return errutil.Wrap(err)
	}

	log.Info("config service", "service", s.Name, "config", buf.String())
	path := filepath.Join(systemConfDir, s.ServiceName())
	err = fileutil.Write(path, buf.Bytes())
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}

func (s *Service) Restart() error {
	return s.action("restart")
}

func (s *Service) Enable() error {
	return s.action("enable")
}

func (s *Service) action(a string) error {
	cmd := exec.Command("systemctl", a, s.ServiceName())
	_, err := cmdutil.RunCmd(cmd)
	return err
}

func (s *Service) genConfig() (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	t := template.Must(template.New("service").Parse(service))
	err := t.Execute(buf, s)
	if err != nil {
		return nil, errutil.Wrap(err)
	}
	return buf, nil
}
