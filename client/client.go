package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"../infomation"
	"../request"
)

//Client ...
//clientはhttpリクエストを送信するためのモジュールです。
type Client interface {
	Requests() error
	Header() (string, error)
	Body() (string, error)
	WriteFile() error
}

type client struct {
	Info     infomation.HttpInfomation
	Response *http.Response
}

func (c *client) Header() (string, error) {
	var h string
	keys := []string{}
	for k := range c.Response.Header {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})
	h += fmt.Sprintln(c.Response.Proto, c.Response.Status)
	for _, k := range keys {
		h += fmt.Sprintln(k, strings.Join(c.Response.Header[k], " "))
	}
	fmt.Println(c.Response.ContentLength)
	return string(h), nil
}

func (c *client) Body() (string, error) {
	b, _ := ioutil.ReadAll(c.Response.Body)
	return string(b), nil
}

func (c *client) WriteFile() error {
	if _, err := os.Stat(path.Dir(c.Info.Output.Filename)); os.IsNotExist(err) {
		os.Mkdir(path.Dir(c.Info.Output.Filename), 0777)
	}
	file, err := os.Create(c.Info.Output.Filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(file, c.Response.Body)
	return nil
}

func (c *client) Requests() error {
	req := request.New(c.Info.Method, c.Info.URL, c.Info.Data.Encode())
	header := []string{}
	if c.Info.Method == "POST" {
		header = append(header, "Content-Type: application/x-www-form-urlencoded")
	}
	for k, v := range c.Info.Header.HeaderInfo {
		header = append(header, k+": "+strings.Join(v, ","))
	}
	req.UpdataIsRedirect(c.Info.Header.ReadFlag)
	req.SetHeader(header)
	c.Response, _ = req.Do()
	return nil
}

//New is Client return
func New(h infomation.HttpInfomation) Client {
	return &client{
		Info: h,
	}
}
