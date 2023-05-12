package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"log"
	"math/big"
	"os"
	"zksync2-examples/contracts/storage"
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

	// Read smart contract bytecode
	bytecode, err := os.ReadFile("../solidity/storage/build/Storage.zbin")
	if err != nil {
		log.Panic(err)
	}

	//Deploy smart contract
	hash, err := w.DeployWithCreate(bytecode, nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Transaction: ", hash)

	// Wait unit transaction is finalized
	_, err = zp.WaitMined(context.Background(), hash)
	if err != nil {
		log.Panic(err)
	}

	receipt, err := zp.GetTransactionReceipt(hash)
	if err != nil {
		log.Panic(err)
	}
	contractAddress := receipt.ContractAddress

	if err != nil {
		panic(err)
	}
	fmt.Println("Smart contract address", contractAddress.String())

	// INTERACT WITH SMART CONTRACT

	// Create instance of Storage smart contract
	storageContract, err := storage.NewStorage(contractAddress, zp)
	if err != nil {
		log.Panic(err)
	}

	// Execute Get method from storage smart contract
	value, err := storageContract.Get(nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Value:", value)

	// Read the private key
	privateKey, err := crypto.ToECDSA(common.Hex2Bytes(PrivateKey))
	if err != nil {
		log.Panic(err)
	}

	// Start configuring transaction parameters
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, w.GetEthSigner().GetDomain().ChainId)
	if err != nil {
		log.Panic(err)
	}

	// Execute Set method from storage smart contract with configured transaction parameters
	tx, err := storageContract.Set(opts, big.NewInt(200))
	if err != nil {
		log.Panic(err)
	}
	// Wait for transaction to be finalized
	_, err = zp.WaitMined(context.Background(), tx.Hash())
	if err != nil {
		log.Panic(err)
	}

	// Execute Get method again to check if state is changed
	value, err = storageContract.Get(nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Value after Set method execution: ", value)

	// INTERACT WITH SMART CONTRACT USING EIP712 TRANSACTIONS
	abi, err := storage.StorageMetaData.GetAbi()
	if err != nil {
		log.Panic(err)
	}
	// Encode set function arguments
	setArguments, err := abi.Pack("set", big.NewInt(500))
	if err != nil {
		log.Panic(err)
	}
	// Execute set function
	execute, err := w.Execute(contractAddress, setArguments, nil, nil)
	if err != nil {
		log.Panic(err)
	}

	_, err = zp.WaitMined(context.Background(), execute)
	if err != nil {
		log.Panic(err)
	}
}
