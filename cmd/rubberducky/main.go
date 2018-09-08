package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rubberducky"
	app.Usage = "finally share code on Ethereum!"
	app.Action = func(c *cli.Context) error {
		fmt.Printf("Hello %v", c.Args().Get(0))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
