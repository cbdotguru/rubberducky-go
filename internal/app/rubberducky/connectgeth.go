package rubberducky

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
)

// ConnectGeth TODO some comments
func ConnectGeth() (*ethclient.Client, *big.Int, error) {
	var ec *ethclient.Client
	var t *big.Int
	var err error
	dir := node.DefaultDataDir()
	ec, err = ethclient.Dial(dir + "/rinkeby/geth.ipc")
	if err != nil {
		err = fmt.Errorf("Could not find geth connection to '%v': '%v'", dir+"/rinkeby/geth.ipc", err)
		return nil, nil, err
	}
	t, err = ec.NetworkID(context.Background())
	if err != nil {
		err = fmt.Errorf("Geth ain't acting right: '%v'", err)
		return nil, nil, err
	}
	return ec, t, nil
}
