package main

import (
	"log"
	"math/big"
	"os"

	"github.com/Hackdom/rubberducky-go/internal/app/rubberducky"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
)

func main() {
	var ec *ethclient.Client
	var t *big.Int
	var err error

	app := cli.NewApp()
	app.Name = "rubberducky"
	app.Usage = "finally share code on Ethereum!"
	app.Action = func(c *cli.Context) error {
		password := rubberducky.GetPass()
		ec, t, err = rubberducky.ConnectGeth()
		if err != nil {
			log.Fatalf("Geth connection failed : '%v'", err)
		}
		err = rubberducky.SendPackage("newtest-new", password, t, ec)
		if err != nil {
			log.Fatalf("Transaction failed : '%v'", err)
		}
		return nil
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
