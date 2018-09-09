package rubberducky

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"github.com/fatih/color"
)

// PackageDBJSON temporary string for hack
const PackageDBJSON = `{"contractName":"PackageDB","abi":[{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"setOwner","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x13af4035"},{"constant":false,"inputs":[{"name":"newAuthority","type":"address"}],"name":"setAuthority","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x7a9e5e4b"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x8da5cb5b"},{"constant":true,"inputs":[],"name":"authority","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xbf7e214f"},{"anonymous":false,"inputs":[{"indexed":true,"name":"nameHash","type":"bytes32"},{"indexed":true,"name":"releaseHash","type":"bytes32"}],"name":"PackageReleaseAdd","type":"event","signature":"0xe7fefc8743abe9f1741edd85e0dbcd31cb5a67f7fd10af5857dcc0a45e0c7906"},{"anonymous":false,"inputs":[{"indexed":true,"name":"nameHash","type":"bytes32"},{"indexed":true,"name":"releaseHash","type":"bytes32"}],"name":"PackageReleaseRemove","type":"event","signature":"0x1d5ef20970faaf2bd4543c4b7acd1ab2747031b93fe095340e20aacb6c9b7a0e"},{"anonymous":false,"inputs":[{"indexed":true,"name":"nameHash","type":"bytes32"}],"name":"PackageCreate","type":"event","signature":"0x94d68ac0a5dee0e8dd504e7e82e1fb1eb122682ceb9fc6aa6647f203fee26f1e"},{"anonymous":false,"inputs":[{"indexed":true,"name":"nameHash","type":"bytes32"},{"indexed":false,"name":"reason","type":"string"}],"name":"PackageDelete","type":"event","signature":"0x188d63b2c009063a155fbcf0c8121b521638675d3d54561c1955bbec5b9ea6bb"},{"anonymous":false,"inputs":[{"indexed":true,"name":"nameHash","type":"bytes32"},{"indexed":true,"name":"oldOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"PackageOwnerUpdate","type":"event","signature":"0xfe2ec6b3a2236fea1f48069f386e0daac1b7b56b918998a3c3a2821594618817"},{"anonymous":false,"inputs":[{"indexed":true,"name":"oldOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"OwnerUpdate","type":"event","signature":"0x343765429aea5a34b3ff6a3785a98a5abb2597aca87bfbb58632c173d585373a"},{"anonymous":false,"inputs":[{"indexed":true,"name":"oldAuthority","type":"address"},{"indexed":true,"name":"newAuthority","type":"address"}],"name":"AuthorityUpdate","type":"event","signature":"0xa1d9e0b26ffdd95159e4605308c755be7b756e3e5dd5c5756b4c77f644a52364"},{"constant":false,"inputs":[{"name":"name","type":"string"}],"name":"setPackage","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x083ae1fe"},{"constant":false,"inputs":[{"name":"nameHash","type":"bytes32"},{"name":"reason","type":"string"}],"name":"removePackage","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x001f8d11"},{"constant":false,"inputs":[{"name":"nameHash","type":"bytes32"},{"name":"newPackageOwner","type":"address"}],"name":"setPackageOwner","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function","signature":"0x2406cedb"},{"constant":true,"inputs":[{"name":"nameHash","type":"bytes32"}],"name":"packageExists","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xa9b35240"},{"constant":true,"inputs":[],"name":"getNumPackages","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x7370a38d"},{"constant":true,"inputs":[{"name":"idx","type":"uint256"}],"name":"getPackageNameHash","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x95f0684b"},{"constant":true,"inputs":[{"name":"nameHash","type":"bytes32"}],"name":"getPackageData","outputs":[{"name":"packageOwner","type":"address"},{"name":"createdAt","type":"uint256"},{"name":"updatedAt","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function","signature":"0xb4d6d4c7"},{"constant":true,"inputs":[{"name":"nameHash","type":"bytes32"}],"name":"getPackageName","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function","signature":"0x06fe1fd7"},{"constant":true,"inputs":[{"name":"name","type":"string"}],"name":"hashName","outputs":[{"name":"","type":"bytes32"}],"payable":false,"stateMutability":"pure","type":"function","signature":"0xaf9a3f9b"}],"address":"0xEBC69aC98ef276e64e163a876b26a7cfCdD32222"}`

// SendPackage TODO some comments
func SendPackage(packagename string, t *big.Int, ec *ethclient.Client) error {
	var artifactObject *ArtifactInfo
	var ethabi abi.ABI
	var stx *types.Transaction
	var err error
	dir := node.DefaultDataDir()

	artifact := []byte(PackageDBJSON)
	err = json.Unmarshal(artifact, &artifactObject)
	if err != nil {
		err = fmt.Errorf("Internal error: '%v'", err)
		return err
	}
	jsonBytes, _ := json.Marshal(artifactObject.ABI)
	br := bytes.NewReader(jsonBytes)
	ethabi, _ = abi.JSON(br)
	nb, _ := ethabi.Pack("setPackage", packagename)
	ks := keystore.NewKeyStore(dir+"/rinkeby/keystore", 262144, 1)
	wa := ks.Accounts()
	color.Magenta("\n\nYou have the following addresses available : \n")
	for _, w := range wa {
		fmt.Println(w.Address.String())
	}
	color.Magenta("\nPlease enter the index (starting with '1') of the one you'd like to use. Keep in mind it needs to be funded with Rinkeby ether: \n")
	reader := bufio.NewReader(os.Stdin)

	i, _ := reader.ReadString('\n')
	num := strings.TrimSpace(i)
	var index int
	index, err = strconv.Atoi(num)
	if err != nil {
		err = fmt.Errorf("Internal error: '%v'", err)
		return err
	}

	password := GetPass()

	var non uint64
	non, err = ec.NonceAt(context.Background(), wa[index-1].Address, nil)
	if err != nil {
		err = fmt.Errorf("Internal error: '%v'", err)
		return err
	}
	ca := common.HexToAddress(artifactObject.Address)
	tx := types.NewTransaction(non, ca, big.NewInt(0), 250000, big.NewInt(21000000000), nb)

	stx, err = ks.SignTxWithPassphrase(wa[index-1], password, tx, t)
	if err != nil {
		err = fmt.Errorf("Signing failed: '%v'", err)
		return err
	}
	err = ec.SendTransaction(context.Background(), stx)
	if err != nil {
		err = fmt.Errorf("Internal error: '%v'", err)
	}
	return err
}
