package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
	"log"
	"math/big"
	"os"
)

func main() {
	var (
		PrivateKey        = os.Getenv("PRIVATE_KEY")
		ZkSyncEraProvider = "https://sepolia.era.zksync.dev"
		EthereumProvider  = "https://rpc.ankr.com/eth_sepolia"
	)

	// Connect to zkSync network
	client, err := clients.Dial(ZkSyncEraProvider)
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()

	// Connect to Ethereum network
	ethClient, err := ethclient.Dial(EthereumProvider)
	if err != nil {
		log.Panic(err)
	}
	defer ethClient.Close()

	// Create wallet
	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, ethClient)
	if err != nil {
		log.Panic(err)
	}

	// Show balance before deposit
	balance, err := wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Balance before deposit: ", balance)

	// Perform deposit
	tx, err := wallet.Deposit(nil, accounts.DepositTransaction{
		Token:  utils.EthAddress,
		Amount: big.NewInt(1_000_000_000),
		To:     wallet.Address(),
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("L1 transaction: ", tx.Hash())

	// Wait for deposit transaction to be finalized on L1 network
	fmt.Println("Waiting for deposit transaction to be finalized on L1 network")
	_, err = bind.WaitMined(context.Background(), ethClient, tx)
	if err != nil {
		log.Panic(err)
	}

	// Get transaction receipt for deposit transaction on L1 network
	l1Receipt, err := ethClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Panic(err)
	}
	// Get deposit transaction on L2 network
	l2Tx, err := client.L2TransactionFromPriorityOp(context.Background(), l1Receipt)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("L2 transaction", l2Tx.Hash)

	// Wait for deposit transaction to be finalized on L2 network (5-7 minutes)
	fmt.Println("Waiting for deposit transaction to be finalized on L2 network (5-7 minutes)")
	_, err = client.WaitMined(context.Background(), l2Tx.Hash)
	if err != nil {
		log.Panic(err)
	}

	balance, err = wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Balance after deposit: ", balance)

	/*
		// ClaimFailedDeposit is used when transaction on L2 has failed.
		cfdTx, err := wallet.ClaimFailedDeposit(nil, l2Tx.Hash)
		if err != nil {
			fmt.Println(err) // this should be triggered if deposit was successful
		}
		fmt.Println("ClaimFailedDeposit hash: ", cfdTx.Hash())
	*/
}
