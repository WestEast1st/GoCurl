package client

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

//Client ...
//clientはhttpリクエストを送信するためのモジュールです。
type Client interface {
	Get() (string, error)
}

type client struct {
	URL      string
	Method   string
	FileName string
	Data     string
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

func getUrl(c *cli.Context) string {
	var urls []string
	r := regexp.MustCompile(`^(http|https)://([\w-]+\.)+[\w-]+(/[\w-./?%&=]*)?$`)
	for _, arg := range c.Args() {
		if len(r.FindStringSubmatch(arg)) > 0 {
			urls = append(urls, r.FindStringSubmatch(arg)[0])
		}
	}
	return urls[len(urls)-1]
}

func New(c *cli.Context) Client {
	var method string
	if c.String("data") != "" {
		method = "POST"
	} else {
		method = "GET"
	}
	return &client{
		URL:    getUrl(c),
		Method: method,
		Data:   strings.Join(c.StringSlice("data"), "&"),
	}
}
