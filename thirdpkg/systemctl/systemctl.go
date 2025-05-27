package systemctl

import (
	_ "embed"
	"fmt"
	"github.com/vksir/vkiss-lib/pkg/template"
	"github.com/vksir/vkiss-lib/pkg/util/cmdutil"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"os/exec"
	"path/filepath"
	"runtime"
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

func (s *Service) Template() string {
	return service
}

func (s *Service) ServiceName() string {
	return fmt.Sprintf("%s.service", s.Name)
}

func (s *Service) Deploy() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("no support %s", runtime.GOOS)
	}

	path := filepath.Join(systemConfDir, s.ServiceName())
	err := template.ExecuteFile(s, path)
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
