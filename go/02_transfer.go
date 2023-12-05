package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
	"log"
	"math/big"
	"os"
)

func main() {
	var (
		PrivateKey1       = os.Getenv("PRIVATE_KEY")
		PublicKey2        = "0x81E9D85b65E9CC8618D85A1110e4b1DF63fA30d9"
		ZkSyncEraProvider = "https://sepolia.era.zksync.dev"
	)

	// Connect to zkSync network
	client, err := clients.Dial(ZkSyncEraProvider)
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()

	// Create wallet
	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey1), &client, nil)
	if err != nil {
		log.Panic(err)
	}

	// Show balances before transfer for both accounts
	account1Balance, err := wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Panic(err)
	}
	account2Balance, err := client.BalanceAt(context.Background(), common.HexToAddress(PublicKey2), nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Account1 balance before transfer: ", account1Balance)
	fmt.Println("Account2 balance before transfer: ", account2Balance)

	// Perform transfer
	tx, err := wallet.Transfer(nil, accounts.TransferTransaction{
		To:     common.HexToAddress(PublicKey2),
		Amount: big.NewInt(1_000_000_000),
		Token:  utils.EthAddress,
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Transaction: ", tx.Hash())

	// Wait for transaction to be finalized on L2 network
	_, err = client.WaitMined(context.Background(), tx.Hash())
	if err != nil {
		log.Panic(err)
	}

	// Show balances after transfer for both accounts
	account1Balance, err = wallet.Balance(context.Background(), utils.EthAddress, nil)
	if err != nil {
		log.Panic(err)
	}
	account2Balance, err = client.BalanceAt(context.Background(), common.HexToAddress(PublicKey2), nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Account1 balance after transfer: ", account1Balance)
	fmt.Println("Account2 balance after transfer: ", account2Balance)
}
