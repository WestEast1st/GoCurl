package main

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"./client"
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
		cli.BoolFlag{
			Name:  "head, I",
			Usage: "HTTP/FTP/FILEなどのヘッダーファイル情報を表示します",
		},
	}

	app.Action = func(c *cli.Context) {
		var urls infomation.HttpInfomation
		var outputconf infomation.Output
		var headerconf infomation.Header
		var method string
		var filename string
		method = "GET"
		wg := &sync.WaitGroup{}
		// output関連のオプション設定
		go func() {
			if c.Bool("O") || c.String("output") != "" {
				wg.Add(1)
				defer wg.Done()
				outputconf.Flag = true
				if c.String("output") != "" {
					path := strings.Split(c.String("output"), "/")
					outputconf.Filepath = strings.Join(path[:len(path)-1], "/")
					outputconf.Filename = path[len(path)-1]
				}
			}
		}()
		// header情報を取得
		go func() {
			if c.Bool("head") {
				wg.Add(1)
				defer wg.Done()
				headerconf.ReadFlag = c.Bool("head")
			}
		}()

		go func() {
			if len(c.StringSlice("header")) > 0 {
				wg.Add(1)
				defer wg.Done()
				m := map[string][]string{
					"Accept-Encoding": {"chunked"},
				}
				for _, v := range c.StringSlice("header") {
					h := strings.Split(v, ":")
					m[h[0]] = []string{h[1]}
				}
				headerconf.HeaderInfo = m
			}
		}()
		//利用可能なスキーム
		r := regexp.MustCompile(`^(http|https|ftp|ftps|dns|file)$`)
		wg.Wait()
		for _, arg := range c.Args() {
			u, _ := url.Parse(arg)
			if len(r.FindStringSubmatch(u.Scheme)) > 0 {
				// -O remote-name用の箇所
				if outputconf.Flag && outputconf.Filename == "" {
					path := strings.Split(u.Path, "/")
					outputconf.Filepath = "./"
					if path[len(path)-1] == "" {
						filename = "index.html"
					} else {
						filename = path[len(path)-1]
					}
					outputconf.Filename = filename
				}
				// -d --dataの送信データの構造変更
				values := url.Values{}
				if len(c.StringSlice("data")) > 0 {
					method = "POST"
					for _, data := range c.StringSlice("data") {
						slice := strings.Split(data, "=")
						values.Set(slice[0], slice[1])
					}
				}
				// URL関連の構造体
				urls = infomation.HttpInfomation{
					URL:      u.String(),
					URI:      u.RequestURI(),
					Method:   method,
					Port:     u.Port(),
					Hostname: u.Hostname(),
					Path:     u.Path,
					Query:    u.Query(),
					Data:     values,
					Fragment: u.Fragment,
					Output:   outputconf,
					Header:   headerconf,
				}
			}
		}
		if len(urls.URL) > 0 {
			client := client.New(urls)
			res, _ := client.Request()
			if res != "" {
				fmt.Println(res)
			}
		}
	}

	app.HideHelp = true

	app.Run(os.Args)
}
