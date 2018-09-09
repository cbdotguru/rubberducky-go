package main

import (
	"bufio"
	"bytes"
	"log"
	"math/big"
	"os"

	"github.com/Hackdom/rubberducky-go/internal/app/rubberducky"
	"github.com/dimiro1/banner"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func main() {
	var ec *ethclient.Client
	var t *big.Int
	var packagename string
	var err error

	app := cli.NewApp()
	app.Name = "rubberducky"
	app.Usage = "finally share code on Ethereum!"
	app.Action = func(c *cli.Context) error {
		banner.Init(os.Stdout, true, true, bytes.NewBufferString("{{ .AnsiColor.Yellow }}{{ .AnsiBackground.Black }}      ,~~.                                                                          \n{{ .AnsiColor.Yellow }}{{ .AnsiBackground.Black }}     (  9 )-_,\n{{ .AnsiColor.Yellow }}{{ .AnsiBackground.Black }}(\\___ )=='-'\n{{ .AnsiColor.Yellow }}{{ .AnsiBackground.Black }} \\ .   ) )\n{{ .AnsiColor.Yellow }}{{ .AnsiBackground.Black }}  \\ `-' /\n{{ .AnsiColor.Yellow }}{{ .AnsiBackground.Black }}   `~j-'   hjw\n{{ .AnsiColor.Yellow }}{{ .AnsiBackground.Black }}    \"=:{{ .AnsiBackground.White }}\n"))
		dir, _ := os.Getwd()
		reader := bufio.NewReader(os.Stdin)
		color.Blue("\nPlease enter a relative path to a truffle directory (no trailing slash), currently at '" + dir + "' : ")
		directoryname, _ := reader.ReadString('\n')
		directoryname = directoryname[:len(directoryname)-1]

		reader = bufio.NewReader(os.Stdin)
		color.Blue("\nPlease enter a package name (it has to be unique, shouldn't be an issue right now, but just sayin'): ")
		packagename, _ = reader.ReadString('\n')
		packagename = packagename[:len(packagename)-1]

		packagename, err = rubberducky.BuildAndWritePackage(directoryname, packagename)
		if err != nil {
			log.Fatalf("Packaging failure : '%v'", err)
		}

		ec, t, err = rubberducky.ConnectGeth()
		if err != nil {
			log.Fatalf("Geth connection failed : '%v'", err)
		}
		err = rubberducky.SendPackage(packagename, t, ec)
		if err != nil {
			log.Fatalf("Transaction failed : '%v'", err)
		}
		color.Cyan("\n\nYou have successfully deployed your package to the rubberducky.test domain on Rinkeby testnet. You should take a moment to look in the truffle project you chose. In there, you will see a shiny new file named ethpm.json tightly packed for ipfs and beyond. Prettify that file and look into the future!")
		return nil
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
