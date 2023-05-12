package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
	"log"
	"os"
	"zksync2-examples/contracts/demo"
)

func main() {
	var (
		PrivateKey     = os.Getenv("PRIVATE_KEY")
		ZkSyncProvider = "https://testnet.era.zksync.dev"
	)

	// Connect to zkSync network
	zp, err := clients.NewDefaultProvider(ZkSyncProvider)
	if err != nil {
		log.Panic(err)
	}
	defer zp.Close()

	// Create singer object from private key for appropriate chain
	chainID, err := zp.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	es, err := accounts.NewEthSignerFromRawPrivateKey(common.Hex2Bytes(PrivateKey), chainID.Int64())
	if err != nil {
		log.Fatal(err)
	}

	// Create wallet
	w, err := accounts.NewWallet(es, zp)
	if err != nil {
		log.Panic(err)
	}

	// Read bytecode of Demo contract
	demoBytecode, err := os.ReadFile("../solidity/demo/build/Demo.zbin")
	if err != nil {
		log.Panic(err)
	}

	// Read bytecode of Foo contract
	fooBytecode, err := os.ReadFile("../solidity/demo/build/Foo.zbin")
	if err != nil {
		log.Panic(err)
	}

	// Deploy smart contract
	hash, err := w.Deploy(demoBytecode, nil, nil, [][]byte{fooBytecode}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Transaction: ", hash)

	// Wait unit transaction is finalized
	_, err = zp.WaitMined(context.Background(), hash)
	if err != nil {
		log.Panic(err)
	}

	// Get address of deployed smart contract
	contractAddress, err := utils.ComputeL2Create2Address(
		w.GetAddress(),
		demoBytecode,
		nil,
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Smart contract address: ", contractAddress.String())

	// INTERACT WITH SMART CONTRACT

	// Create instance of Demo contract
	demoContract, err := demo.NewDemo(contractAddress, zp)
	if err != nil {
		log.Panic(err)
	}

	// Execute GetFooName method
	value, err := demoContract.GetFooName(nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Value:", value)
}
