package steam

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vksir/vkiss-lib/pkg/util/convutil"
	"testing"
)

func TestGetPublishedFileDetails(t *testing.T) {
	data, err := GetPublishedFileDetails("3485622209", "1699194522")
	assert.Nil(t, err)
	fmt.Println(convutil.MustJsonString(data))
}

func TestGetNewsForApp(t *testing.T) {
	data, err := GetNewsForApp("343050")
	assert.Nil(t, err)
	fmt.Println(convutil.MustJsonString(data))
}
