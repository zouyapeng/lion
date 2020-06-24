package lion

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func doRequest(lionURL string, body map[string]string) (respData []byte, err error) {
	httpClient := &http.Client{}
	httpClient.Timeout = 10 * time.Second
	data := url.Values{}
	for key, value := range body {
		data.Set(key, value)
	}

	req, err := http.NewRequest("POST", lionURL, strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	defer func() { _ = resp.Body.Close() }()

	respData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
