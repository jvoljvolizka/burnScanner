package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"sync"

	"burnScan/factory"
	"burnScan/pair"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	var rpc string
	var fac string
	var silent bool
	flag.StringVar(&rpc, "rpc", "", "")
	flag.StringVar(&fac, "factory", "", "")
	flag.BoolVar(&silent, "silent", false, "")
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
	var wg sync.WaitGroup

	for i := 0; pairLen.Cmp(big.NewInt(int64(i))) > 0; i++ {
		wg.Add(1)
		go doStuff(client, factoryInstance, int64(i), &wg, silent)
	}
	wg.Wait()
}

func doStuff(client *ethclient.Client, factoryInstance *factory.Factory, i int64, wg *sync.WaitGroup, silent bool) {
	defer wg.Done()
	pairAddress, err := factoryInstance.AllPairs(&bind.CallOpts{}, big.NewInt(int64(i)))
	if err != nil {
		if !silent {
			fmt.Printf("bad :/ %s\n", big.NewInt(int64(i)).String())
		}

		return
	}
	pairInstance, err := pair.NewPair(pairAddress, client)
	balance, err := pairInstance.BalanceOf(&bind.CallOpts{}, pairAddress)
	if err != nil {
		if !silent {
			fmt.Printf("bad : %s\n", pairAddress.String())
		}
		return
	}
	if balance.Cmp(big.NewInt(int64(0))) > 0 {
		fmt.Printf("found it! : ")
		fmt.Println(pairAddress.String())
	} else {
		if !silent {
			fmt.Println(pairAddress.String() + " :/")
		}

	}
}
