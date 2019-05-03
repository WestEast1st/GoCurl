package main

import (
	"fmt"
	"os"
	"regexp"

	"./client"
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
			Value: "./",
			Usage: "保存名を指定してください",
		},
		cli.StringFlag{
			Name:  "data, d",
			Usage: "POST送信するデータを入力してください",
		},
	}
	app.Action = func(c *cli.Context) error {
		urls := getUrls(c)
		client := client.New(urls[0])
		response, err := client.Get()
		if err != nil {
			return err
		}
		fmt.Println(response)
		return nil
	}

	app.HideHelp = true

	app.Run(os.Args)
}

func getUrls(c *cli.Context) []string {
	var urls []string
	r := regexp.MustCompile(`^(http|https)://([\w-]+\.)+[\w-]+(/[\w-./?%&=]*)?$`)
	for _, arg := range c.Args() {
		if len(r.FindStringSubmatch(arg)) > 0 {
			urls = append(urls, r.FindStringSubmatch(arg)[0])
		}
	}
	return urls
}
