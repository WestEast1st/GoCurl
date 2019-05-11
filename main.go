package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"./client"
	"./cookie"
	"./infomation"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Go cURL"
	app.Usage = "I tried to make Go URL"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "output, o",
			Usage: "保存名を指定してください",
		},
		cli.BoolFlag{
			Name:  "remote-name, O",
			Usage: "リモートファイルのファイル部分を名前に利用します",
		},
		cli.StringSliceFlag{
			Name:  "data, d",
			Usage: "POST送信するデータを入力してください",
		},
		cli.StringSliceFlag{
			Name:  "header, H",
			Usage: "POST送信するデータを入力してください",
		},
		cli.StringSliceFlag{
			Name:  "cookie, b",
			Usage: "ヘッダーに乗せて送信するcookie",
		},
		cli.BoolFlag{
			Name:  "head, I",
			Usage: "HTTP/FTP/FILEなどのヘッダーファイル情報を表示します",
		},
	}

	app.Action = func(c *cli.Context) {
		var (
			urls       = infomation.HttpInfomation{}
			outputconf infomation.Output
			headerconf infomation.Header
			method     = "GET"
			filename   string
			values     = url.Values{}
			domain     string
		)
		urls.Cookie = cookie.New()
		wg := &sync.WaitGroup{}
		wg.Add(4)
		// output関連のオプション設定
		go func() {
			defer wg.Done()
			if c.Bool("O") || c.String("output") != "" {
				outputconf.Flag = true
				if c.String("output") != "" {
					path := strings.Split(c.String("output"), "/")
					outputconf.Filepath = strings.Join(path[:len(path)-1], "/")
					outputconf.Filename = path[len(path)-1]
				}
			}
		}()
		go func() {
			defer wg.Done()
			r := regexp.MustCompile(`^(http|https|ftp|ftps|dns|file)$`)
			for _, arg := range c.Args() {
				u, _ := url.Parse(arg)
				if len(r.FindStringSubmatch(u.Scheme)) > 0 {
					urls.URL = u.String()
					urls.URI = u.RequestURI()
					urls.Port = u.Port()
					urls.Hostname = u.Hostname()
					urls.Path = u.Path
					urls.Query = u.Query()
					urls.Fragment = u.Fragment
					domain = u.Host
				}
			}
		}()
		go func() {
			defer wg.Done()
			m := map[string][]string{
				"Accept-Encoding": {"chunked"},
			}
			if len(c.StringSlice("header")) > 0 {
				for _, v := range c.StringSlice("header") {
					h := strings.Split(v, ":")
					m[h[0]] = append(m[h[0]], h[1])
				}
			}
			headerconf.HeaderInfo = m
		}()
		go func() {
			// -d --dataの送信データの構造変更
			defer wg.Done()
			if len(c.StringSlice("data")) > 0 {
				method = "POST"
				for _, data := range c.StringSlice("data") {
					slice := strings.Split(data, "=")
					values.Set(slice[0], slice[1])
				}
			}
		}()
		headerconf.ReadFlag = c.Bool("head")
		wg.Wait()
		if len(c.StringSlice("cookie")) > 0 {
			for _, v := range c.StringSlice("cookie") {
				slice := strings.Split(v, "; ")
				for _, ck := range slice {
					c := strings.Split(ck, "=")
					urls.Cookie.Add(http.Cookie{
						Domain: domain,
						Name:   c[0],
						Value:  c[1],
					})
				}
			}
		}
		// 実行
		if len(urls.URL) > 0 {
			urls.Method = method
			urls.Header = headerconf
			urls.Data = values
			if outputconf.Flag {
				path := strings.Split(urls.Path, "/")
				outputconf.Filepath = "./"
				if path[len(path)-1] == "" {
					filename = "index.html"
				} else {
					filename = path[len(path)-1]
				}
				outputconf.Filename = filename
				urls.Output = outputconf
			}
			client := client.New(urls)
			client.Requests()
			if outputconf.Flag {
				client.WriteFile()
			} else if c.Bool("head") {
				s, _ := client.Header()
				fmt.Println(s)
			} else {
				s, _ := client.Body()
				fmt.Println(s)
			}
		}
	}

	app.HideHelp = true

	app.Run(os.Args)
}
