package steam

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/vksir/vkiss-lib/pkg/util/errutil"
	"net/http"
	"strconv"
)

type GetPublishedFileDetailsResponse struct {
	Response struct {
		Result               int `json:"result"`
		Resultcount          int `json:"resultcount"`
		Publishedfiledetails []struct {
			Publishedfileid       string `json:"publishedfileid"`
			Result                int    `json:"result"`
			Creator               string `json:"creator"`
			CreatorAppId          int    `json:"creator_app_id"`
			ConsumerAppId         int    `json:"consumer_app_id"`
			Filename              string `json:"filename"`
			FileSize              string `json:"file_size"`
			FileUrl               string `json:"file_url"`
			HcontentFile          string `json:"hcontent_file"`
			PreviewUrl            string `json:"preview_url"`
			HcontentPreview       string `json:"hcontent_preview"`
			Title                 string `json:"title"`
			Description           string `json:"description"`
			TimeCreated           int    `json:"time_created"`
			TimeUpdated           int    `json:"time_updated"`
			Visibility            int    `json:"visibility"`
			Banned                int    `json:"banned"`
			BanReason             string `json:"ban_reason"`
			Subscriptions         int    `json:"subscriptions"`
			Favorited             int    `json:"favorited"`
			LifetimeSubscriptions int    `json:"lifetime_subscriptions"`
			LifetimeFavorited     int    `json:"lifetime_favorited"`
			Views                 int    `json:"views"`
			Tags                  []struct {
				Tag string `json:"tag"`
			} `json:"tags"`
		} `json:"publishedfiledetails"`
	} `json:"response"`
}

func GetPublishedFileDetails(workshopId ...string) (GetPublishedFileDetailsResponse, error) {
	url := "https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/"
	data := map[string]string{
		"itemcount": strconv.Itoa(len(workshopId)),
	}
	for i, id := range workshopId {
		key := fmt.Sprintf("publishedfileids[%d]", i)
		data[key] = id
	}
	var res GetPublishedFileDetailsResponse
	err := request(http.MethodPost, url, data, &res)
	return res, err
}

type GetNewsForAppResponse struct {
	Appnews struct {
		Appid     int `json:"appid"`
		Newsitems []struct {
			Gid           string   `json:"gid"`
			Title         string   `json:"title"`
			Url           string   `json:"url"`
			IsExternalUrl bool     `json:"is_external_url"`
			Author        string   `json:"author"`
			Contents      string   `json:"contents"`
			Feedlabel     string   `json:"feedlabel"`
			Date          int      `json:"date"`
			Feedname      string   `json:"feedname"`
			FeedType      int      `json:"feed_type"`
			Appid         int      `json:"appid"`
			Tags          []string `json:"tags"`
		} `json:"newsitems"`
		Count int `json:"count"`
	} `json:"appnews"`
}

func GetNewsForApp(appId string) (GetNewsForAppResponse, error) {
	url := "https://api.steampowered.com/ISteamNews/GetNewsForApp/v2/"
	data := map[string]string{
		"appid": appId,
	}
	var res GetNewsForAppResponse
	err := request(http.MethodGet, url, data, &res)
	return res, err
}

func request(method, url string, data map[string]string, a any) error {
	client := resty.New().R()
	if method == http.MethodGet {
		client.SetQueryParams(data)
	} else {
		client.SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetFormData(data)
	}
	resp, err := client.Execute(method, url)
	if err != nil {
		return errutil.Wrap(err)
	}
	if resp.StatusCode() != http.StatusOK {
		return errutil.WrapF("request failed: code=%d, body=%s", resp.StatusCode(), resp.Body())
	}
	err = json.Unmarshal(resp.Body(), a)
	if err != nil {
		return errutil.Wrap(err)
	}
	return nil
}
