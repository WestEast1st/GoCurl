package client

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"../infomation"
)

//Client ...
//clientはhttpリクエストを送信するためのモジュールです。
type Client interface {
	Request() (string, error)
}

type client struct {
	Info infomation.HttpInfomation
}

func (c *client) Request() (string, error) {
	httpinfo := c.Info
	req, _ := http.NewRequest(httpinfo.Method, httpinfo.URL, strings.NewReader(httpinfo.Data.Encode()))
	if httpinfo.Method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	httpclient := http.Client{}
	resp, err := httpclient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if httpinfo.Output.Flag {
		if _, err := os.Stat(path.Dir(httpinfo.Output.Filename)); os.IsNotExist(err) {
			os.Mkdir(path.Dir(httpinfo.Output.Filename), 0777)
		}
		file, err := os.Create(httpinfo.Output.Filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		io.Copy(file, resp.Body)
		return "", nil
	}
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray), nil
}

func New(h infomation.HttpInfomation) Client {
	return &client{
		Info: h,
	}
}
