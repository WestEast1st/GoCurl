package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"sort"
	"strconv"
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
	var httpclient http.Client
	httpinfo := c.Info
	//
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest(httpinfo.Method, httpinfo.URL, strings.NewReader(httpinfo.Data.Encode()))
	if httpinfo.Method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	for k, v := range httpinfo.Header.HeaderInfo {
		req.Header.Add(k, strings.Join(v, ","))
	}
	if httpinfo.Header.ReadFlag {
		httpclient = http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	} else {
		httpclient = http.Client{Jar: jar}
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	// cookieの保存
	/*
		defer func() {
			set_cookie_url, err := url.Parse(httpinfo.URL)
			if err != nil {
				log.Fatal(err)
			}
			cookiejars := jar.Cookies(set_cookie_url)
		}()
	*/
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
	if httpinfo.Header.ReadFlag {
		var h string
		keys := []string{}
		for k, _ := range resp.Header {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return strings.Compare(keys[i], keys[j]) < 0
		})
		h += fmt.Sprintln(resp.Proto, resp.Status)
		for _, k := range keys {
			h += fmt.Sprintln(k, strings.Join(resp.Header[k], " "))
		}
		return string(h), nil
	}
	return string(byteArray), nil
}

func New(h infomation.HttpInfomation) Client {
	return &client{
		Info: h,
	}
}
