package systemctl

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vksir/vkiss-lib/pkg/template"
	"testing"
)

func TestTestService(t *testing.T) {
	s := &Service{
		Name:             "s1",
		Description:      "d1",
		WorkingDirectory: "/bin/",
		ExecStart:        "/bin/sleep 60",
		RestartOnFailure: true,
		User:             "root",
		Group:            "root",
	}
	res, err := template.ExecuteString(s)
	assert.Nil(t, err)
	assert.Equal(t, `[Unit]
Description = d1
After = network.target
Wants = network.target

[Service]
Type = simple
WorkingDirectory=/bin/
ExecStart=/bin/sleep 60
Restart=on-failure
User=root
Group=root

[Install]
WantedBy = multi-user.target
`, res)
}

func TestTestService2(t *testing.T) {
	s := &Service{
		Name:        "s1",
		Description: "d1",
		ExecStart:   "/bin/sleep 60",
	}
	res, err := template.ExecuteString(s)
	assert.Nil(t, err)
	fmt.Println(res)
	assert.Equal(t, `[Unit]
Description = d1
After = network.target
Wants = network.target

[Service]
Type = simple
ExecStart=/bin/sleep 60
Restart=on-failure

[Install]
WantedBy = multi-user.target
`, res)
}
