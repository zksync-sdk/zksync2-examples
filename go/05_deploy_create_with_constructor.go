package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"log"
	"math/big"
	"os"
	"zksync2-examples/contracts/incrementer"
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
	bytecode, err := os.ReadFile("../solidity/incrementer/build/Incrementer.zbin")
	if err != nil {
		log.Panic(err)
	}

	// Get ABI
	abi, err := incrementer.IncrementerMetaData.GetAbi()
	if err != nil {
		log.Panic(err)
	}

	// Encode constructor arguments
	constructor, err := abi.Pack("", big.NewInt(2))
	if err != nil {
		log.Panicf("error while encoding constructor arguments: %s", err)
	}

	// Use salt if there is need to deploy same contract twice, otherwise remove salt
	salt := make([]byte, 32)
	_, err = rand.Read(salt)
	if err != nil {
		log.Panicf("error while generating salt: %s", err)
	}

	// Deploy smart contract
	hash, err := w.DeployWithCreate(bytecode, constructor, nil, nil)
	if err != nil {
		panic(err)
		// When contract is deployed twice without salt the following error occurs:
		// panic: failed to EstimateGas712: failed to query eth_estimateGas: execution reverted: Code hash is non-zero
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
	fmt.Println("Smart contract address: ", contractAddress.String())

	// INTERACT WITH SMART CONTRACT

	// Create instance of Incrementer contract
	incrementerContract, err := incrementer.NewIncrementer(contractAddress, zp)
	if err != nil {
		log.Panic(err)
	}

	// Execute Get method
	value, err := incrementerContract.Get(nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Value before Increment method execution: ", value)

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
	tx, err := incrementerContract.Increment(opts)
	if err != nil {
		log.Panic(err)
	}
	// Wait for transaction to be finalized
	_, err = zp.WaitMined(context.Background(), tx.Hash())
	if err != nil {
		log.Panic(err)
	}

	// Execute Get method again to check if state is changed
	value, err = incrementerContract.Get(nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Value after Increment method execution: ", value)
}
