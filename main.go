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
		var (
			urls       infomation.HttpInfomation
			outputconf infomation.Output
			headerconf infomation.Header
			method     = "GET"
			filename   string
			values     = url.Values{}
		)
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
					urls = infomation.HttpInfomation{
						URL:      u.String(),
						URI:      u.RequestURI(),
						Port:     u.Port(),
						Hostname: u.Hostname(),
						Path:     u.Path,
						Query:    u.Query(),
						Fragment: u.Fragment,
					}
				}
			}
		}()
		go func() {
			defer wg.Done()
			m := map[string][]string{
				"Accept-Encoding": {"chunked", "gzip"},
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
