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
		PrivateKey        = os.Getenv("PRIVATE_KEY")
		ZkSyncEraProvider = "https://sepolia.era.zksync.dev"
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

	_, tokenAbi, bytecode, err := utils.ReadStandardJson("../solidity/custom_paymaster/token/build/Token.json")
	if err != nil {
		log.Panic(err)
	}

	constructor, err := tokenAbi.Pack("", "Crown", "Crown", uint8(18))
	if err != nil {
		log.Panic(err)
	}

	//Deploy smart contract
	hash, err := wallet.DeployWithCreate(nil, accounts.CreateTransaction{Bytecode: bytecode, Calldata: constructor})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Transaction: ", hash)

	// Wait unit transaction is finalized
	tx, err := client.WaitMined(context.Background(), hash)
	if err != nil {
		log.Panic(err)
	}

	// Get address of deployed smart contract
	tokenAddress := tx.ContractAddress
	fmt.Println("Token address", tokenAddress.String())

	// Create instance of token contract
	token, err := crown.NewCrown(tokenAddress, client)
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
