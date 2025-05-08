package ddns

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMyIp(t *testing.T) {
	ip, err := GetMyIp("http://127.0.0.1:5801")
	assert.Nil(t, err)
	fmt.Println(ip)
}
