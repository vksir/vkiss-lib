package ddns

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"net/http"
	"net/url"
)

func GetMyIp(endpoint string) (string, error) {
	u, err := url.JoinPath(endpoint, "/my_ip")
	if err != nil {
		return "", errutil.Wrap(err)
	}

	resp, err := resty.New().R().Get(u)
	if err != nil {
		return "", errutil.Wrap(err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("request failed: code=%d, body=%s", resp.StatusCode(), string(resp.Body()))
	}

	var data MyIpResponse
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return "", errutil.Wrap(err)
	}
	return data.Ip, nil
}
