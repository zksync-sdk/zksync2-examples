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
		PrivateKey        = os.Getenv("PRIVATE_KEY")
		ZkSyncEraProvider = "https://sepolia.era.zksync.dev"
		TokenAddress      = "0x927488F48ffbc32112F1fF721759649A89721F8F" // Crown token which can be minted for free
	)

	// Connect to zkSync network
	client, err := clients.Dial(ZkSyncEraProvider)
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()

	// Create wallet
	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, nil)
	if err != nil {
		log.Panic(err)
	}

	_, paymasterAbi, bytecode, err := utils.ReadStandardJson("../solidity/custom_paymaster/paymaster/build/Paymaster.json")
	if err != nil {
		log.Panic(err)
	}

	constructor, err := paymasterAbi.Pack("", common.HexToAddress(TokenAddress))
	if err != nil {
		log.Panic(err)
	}

	// Deploy paymaster contract
	hash, err := wallet.DeployAccountWithCreate(nil, accounts.CreateTransaction{
		Bytecode: bytecode,
		Calldata: constructor,
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Transaction: ", hash)

	// Wait unit transaction is finalized
	_, err = client.WaitMined(context.Background(), hash)
	if err != nil {
		log.Panic(err)
	}

	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		log.Panic(err)
	}
	contractAddress := receipt.ContractAddress
	fmt.Println("Paymaster address", contractAddress.String())
}
