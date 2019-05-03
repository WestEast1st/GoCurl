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

	app.Action = func(c *cli.Context) error {
		url := c.Args().Get(0)
		client := client.New(url)
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
