package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/types"
	"log"
	"math/big"
	"os"
	"zksync2-examples/contracts/storage"
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

	// Read smart contract bytecode
	//bytecode, err := os.ReadFile("../solidity/storage/build/Storage.zbin")
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	////Deploy smart contract
	//hash, err := wallet.DeployWithCreate(nil, accounts.CreateTransaction{Bytecode: bytecode})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Transaction: ", hash)
	//
	//// Wait unit transaction is finalized
	//_, err = client.WaitMined(context.Background(), hash)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//receipt, err := client.TransactionReceipt(context.Background(), hash)
	//if err != nil {
	//	log.Panic(err)
	//}
	//contractAddress := receipt.ContractAddress
	//
	//if err != nil {
	//	log.Panic(err)
	//}
	//fmt.Println("Smart contract address", contractAddress.String())

	// INTERACT WITH SMART CONTRACT

	// Create instance of Storage smart contract
	//storageContract, err := storage.NewStorage(contractAddress, client)
	//if err != nil {
	//	log.Panic(err)
	//}

	contractAddress := common.HexToAddress("0x99B4da9890d62eC523851A345e2b216f1216EDB4")
	abi, err := storage.StorageMetaData.GetAbi()
	if err != nil {
		log.Panic(err)
	}
	// Encode set function arguments
	setArguments, err := abi.Pack("set", big.NewInt(700))
	if err != nil {
		log.Panic(err)
	}
	gas, err := client.EstimateGasL2(context.Background(), types.CallMsg{
		CallMsg: ethereum.CallMsg{
			To:   &contractAddress,
			From: wallet.Address(),
			Data: setArguments,
		},
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Gas: ", gas)

	result, err := wallet.CallContract(context.Background(), accounts.CallMsg{
		To:   &contractAddress,
		Data: setArguments,
	}, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Result: ", result)

	//// Execute Get method from storage smart contract
	//value, err := storageContract.Get(nil)
	//if err != nil {
	//	log.Panic(err)
	//}
	//fmt.Println("Value:", value)
	//
	//// Start configuring transaction parameters
	//opts, err := bind.NewKeyedTransactorWithChainID(wallet.Signer().PrivateKey(), wallet.Signer().Domain().ChainId)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//// Execute Set method from storage smart contract with configured transaction parameters
	//tx, err := storageContract.Set(opts, big.NewInt(200))
	//if err != nil {
	//	log.Panic(err)
	//}
	//// Wait for transaction to be finalized
	//_, err = client.WaitMined(context.Background(), tx.Hash())
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//// Execute Get method again to check if state is changed
	//value, err = storageContract.Get(nil)
	//if err != nil {
	//	log.Panic(err)
	//}
	//fmt.Println("Value after first Set method execution: ", value)
	//
	//// INTERACT WITH SMART CONTRACT USING EIP-712 TRANSACTIONS
	//abi, err := storage.StorageMetaData.GetAbi()
	//if err != nil {
	//	log.Panic(err)
	//}
	//// Encode set function arguments
	//setArguments, err := abi.Pack("set", big.NewInt(500))
	//if err != nil {
	//	log.Panic(err)
	//}
	//// Execute set function
	//execute, err := wallet.SendTransaction(context.Background(), &accounts.Transaction{
	//	To:   &contractAddress,
	//	Data: setArguments,
	//})
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//_, err = client.WaitMined(context.Background(), execute)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//// Execute Get method again to check if state is changed
	//value, err = storageContract.Get(nil)
	//if err != nil {
	//	log.Panic(err)
	//}
	//fmt.Println("Value after second Set method execution: ", value)
}
