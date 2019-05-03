package main

import (
	"os"

	"./client"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Go cURL"
	app.Usage = "I tried to make Go URL"
	app.Version = "0.0.1"

	app.Action = func(c *cli.Context) error {
		client := client.New(c.Args().Get(0))
		client.Get()
		return nil
	}

	app.HideHelp = true

	app.Run(os.Args)
}
