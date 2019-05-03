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
		cli.BoolFlag{
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
		if c.Bool("O") != true && c.String("output") == "" {
			response, err := client.Request()
			if err != nil {
				return err
			}
			fmt.Println(response)
		} else {
			client.GetImage()
		}
		return nil
	}

	app.HideHelp = true

	app.Run(os.Args)
}
