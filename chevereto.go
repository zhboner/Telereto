package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type CheveretoResponse struct {
	OriginUrl string
	SiteUrl   string
	DeleteUrl string
}

func call_chevereto_api(tg_url string) (CheveretoResponse, error) {
	client := &http.Client{}

	form_data := url.Values{
		"source": []string{tg_url},
		"format": []string{string("json")},
	}
	req, err := http.NewRequest("POST", CHEVERETO_HOST_API_URL, strings.NewReader(form_data.Encode()))
	if err != nil {
		return CheveretoResponse{}, err
	}

	req.Header.Add("X-API-Key", CHEVERETO_API_KEY)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return CheveretoResponse{}, err
	}
	defer func() {
		err = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CheveretoResponse{}, err
	}

	var (
		data map[string]interface{}
	)
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return CheveretoResponse{}, err
	}

	status_code := data["status_code"].(float64)
	if status_code != 200 {
		error_mgs := data["error"].(map[string]interface{})

		return CheveretoResponse{}, errors.New(error_mgs["message"].(string))
	}

	image_info := data["image"].(map[string]interface{})

	return CheveretoResponse{
		OriginUrl: image_info["url"].(string),
		DeleteUrl: image_info["delete_url"].(string),
		SiteUrl:   image_info["url_short"].(string),
	}, nil
}
