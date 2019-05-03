package main

import (
	"fmt"
	"os"

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
			Usage: "保存名を指定してください",
		},
		cli.BoolTFlag{
			Name:  "remote-name, O",
			Usage: "リモートファイルのファイル部分を名前に利用します",
		},
		cli.StringSliceFlag{
			Name:  "data, d",
			Usage: "POST送信するデータを入力してください",
		},
	}

	app.Action = func(c *cli.Context) error {
		client := client.New(c)
		response, err := client.Get()
		if err != nil {
			return err
		}
		response = ""
		fmt.Println(response)
		return nil
	}

	app.HideHelp = true

	app.Run(os.Args)
}
