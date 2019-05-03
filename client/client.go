package client

import (
	"io/ioutil"
	"net/http"
)

//Client ...
//clientはhttpリクエストを送信するためのモジュールです。
type Client interface {
	Get() (string, error)
}

type client struct {
	URL string
}

func (c *client) Get() (string, error) {
	req, _ := http.NewRequest("GET", c.URL, nil)

	httpclient := new(http.Client)
	resp, err := httpclient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray), nil
}

func New(url string) Client {
	return &client{URL: url}
}
