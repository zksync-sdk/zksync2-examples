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
)

func main() {
	var (
		PrivateKey     = os.Getenv("PRIVATE_KEY")
		ZkSyncProvider = "https://testnet.era.zksync.dev"
		TokenAddress   = "Cd9BDa1d0FC539043D4C80103bdF4f9cb108931B" // Crown token which can be minted for free
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

	// Read paymaster contract from standard json
	_, paymasterAbi, bytecode, err := utils.ReadStandardJson("../solidity/custom_paymaster/paymaster/build/Paymaster.json")
	if err != nil {
		log.Panic(err)
	}

	// Encode paymaster constructor
	constructor, err := paymasterAbi.Pack("", common.HexToAddress(TokenAddress))
	if err != nil {
		log.Panic(err)
	}

	// Deploy paymaster contract
	hash, err := w.DeployAccount(bytecode, constructor, nil, nil)
	if err != nil {
		log.Panic(err)
	}
	if err != nil {
		log.Panic(err)
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
		bytecode,
		constructor,
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Paymaster address", contractAddress.String())
}
