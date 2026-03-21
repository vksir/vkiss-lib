package steam

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vksir/vkiss-lib/pkg/util/convutil"
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

func TestGetSteamCmdAppInfo(t *testing.T) {
	data, err := GetSteamCmdAppInfo("343050")
	assert.Nil(t, err)
	fmt.Println(convutil.MustJsonString(data))
	fmt.Println(data.Data.Field1.Depots.Branches.Public.Buildid)
}
