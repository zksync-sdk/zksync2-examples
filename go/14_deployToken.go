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
	"zksync2-examples/contracts/crown"
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

	// Read token contract from standard json
	_, tokenAbi, bytecode, err := utils.ReadStandardJson("../solidity/custom_paymaster/token/build/Token.json")
	if err != nil {
		log.Panic(err)
	}

	// Encode token constructor
	constructor, err := tokenAbi.Pack("", "Crown", "Crown", uint8(18))
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(common.Bytes2Hex(constructor))

	//Deploy smart contract
	//hash, err := w.Deploy(bytecode, constructor, nil, nil, nil)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Transaction: ", hash)

	// Wait unit transaction is finalized
	//_, err = zp.WaitMined(context.Background(), hash)
	//if err != nil {
	//	log.Panic(err)
	//}

	// Get address of deployed smart contract
	tokenAddress, err := utils.ComputeL2Create2Address(
		w.GetAddress(),
		bytecode,
		constructor,
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Token address", tokenAddress.String())

	// Create instance of token contract
	token, err := crown.NewCrown(tokenAddress, zp)
	if err != nil {
		log.Panic(err)
	}

	symbol, err := token.Symbol(nil)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Symbol: ", symbol)

	decimals, err := token.Decimals(nil)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Decimals: ", decimals)
}
