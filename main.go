package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"

	"burnScan/factory"
	"burnScan/pair"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	var rpc string
	var fac string

	flag.StringVar(&rpc, "rpc", "", "")
	flag.StringVar(&fac, "factory", "", "")

	flag.Parse()

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	factoryAddress := common.HexToAddress(fac)

	factoryInstance, err := factory.NewFactory(factoryAddress, client)
	if err != nil {
		//	fmt.Printf("bad : %s\n", address)
		log.Fatal(err)
	}
	pairLen, err := factoryInstance.AllPairsLength(&bind.CallOpts{})
	if err != nil {
		//	fmt.Printf("bad : %s\n", address)
		log.Fatal(err)
	}

	for i := 0; pairLen.Cmp(big.NewInt(int64(i))) > 0; i++ {
		pairAddress, err := factoryInstance.AllPairs(&bind.CallOpts{}, big.NewInt(int64(i)))
		if err != nil {
			fmt.Printf("bad :/ %s\n", big.NewInt(int64(i)).String())
			continue
		}
		pairInstance, err := pair.NewPair(pairAddress, client)
		balance, err := pairInstance.BalanceOf(&bind.CallOpts{}, pairAddress)
		if err != nil {
			//	fmt.Printf("bad : %s\n", address)
			fmt.Printf("bad : %s\n", pairAddress.String())
			continue
		}
		if balance.Cmp(big.NewInt(int64(0))) > 0 {
			fmt.Printf("found it! : ")
			fmt.Println(pairAddress.String())
		} else {
			fmt.Println(pairAddress.String() + " :/")
		}
	}
}
