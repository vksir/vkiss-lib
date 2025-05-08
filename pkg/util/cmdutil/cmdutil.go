package cmdutil

import (
	"fmt"
	"os/exec"
)

func RunCmd(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("RunCmd failed: "+
			"cmd=%s, out=%s, err=%w", cmd.String(), out, err)
	}
	return string(out), nil
}
