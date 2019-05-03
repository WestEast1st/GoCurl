package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//Client ...
//clientはhttpリクエストを送信するためのモジュールです。
type Client interface {
	Get()
}

type client struct {
	URL string
}

func (c *client) Get() {
	req, _ := http.NewRequest("GET", c.URL, nil)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}

func New(url string) Client {
	return &client{URL: url}
}
