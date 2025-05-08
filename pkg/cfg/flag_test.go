package cfg

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFlag(t *testing.T) {
	os.Args = []string{os.Args[0], "--k1", "v1", "--k2", "2", "--k3"}
	cmd := &cobra.Command{}
	f1 := NewFlag[string]("k1", "k1", "")
	f2 := NewFlag[int64]("k2", "k2", "")
	f3 := NewFlag[bool]("k3", "k3", "")
	f4 := NewFlag[bool]("k4", "k4", "")
	f1.Bind(cmd)
	f2.Bind(cmd)
	f3.Bind(cmd)
	f4.Bind(cmd)
	err := cmd.Execute()
	assert.Nil(t, err)
	assert.Equal(t, f1.Get(), "v1")
	assert.Equal(t, f2.Get(), int64(2))
	assert.Equal(t, f3.Get(), true)
	assert.Equal(t, f4.Get(), false)
}
