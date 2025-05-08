package systemctl

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTestService(t *testing.T) {
	s := Service{
		Name:             "s1",
		Description:      "d1",
		WorkingDirectory: "/bin/",
		ExecStart:        "/bin/sleep 60",
		RestartOnFailure: true,
		User:             "root",
		Group:            "root",
	}
	buf, err := s.genConfig()
	assert.Nil(t, err)
	fmt.Println(buf.String())
}

func TestTestService2(t *testing.T) {
	s := Service{
		Name:        "s1",
		Description: "d1",
		ExecStart:   "/bin/sleep 60",
	}
	buf, err := s.genConfig()
	assert.Nil(t, err)
	fmt.Println(buf.String())
}
