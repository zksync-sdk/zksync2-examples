package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"log"
	"math/big"
	"os"
)

func main() {
	var (
		PrivateKey1    = os.Getenv("PRIVATE_KEY")
		PublicKey2     = "81E9D85b65E9CC8618D85A1110e4b1DF63fA30d9"
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
	es, err := accounts.NewEthSignerFromRawPrivateKey(common.Hex2Bytes(PrivateKey1), chainID.Int64())
	if err != nil {
		log.Fatal(err)
	}

	// Create wallet
	w, err := accounts.NewWallet(es, zp)
	if err != nil {
		log.Panic(err)
	}

	// Show balances before transfer for both accounts
	account1Balance, err := w.GetBalance()
	if err != nil {
		log.Panic(err)
	}
	account2Balance, err := zp.BalanceAt(context.Background(), common.HexToAddress(PublicKey2), nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Account1 balance before transfer: ", account1Balance)
	fmt.Println("Account2 balance before transfer: ", account2Balance)

	// Perform transfer
	hash, err := w.Transfer(
		common.HexToAddress(PublicKey2),
		big.NewInt(1_000_000_000),
		nil,
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Transaction: ", hash)

	// Wait for transaction to be finalized on L2 network
	_, err = zp.WaitMined(context.Background(), hash)
	if err != nil {
		log.Panic(err)
	}

	// Show balances after transfer for both accounts
	account1Balance, err = w.GetBalance()
	if err != nil {
		log.Panic(err)
	}
	account2Balance, err = zp.BalanceAt(context.Background(), common.HexToAddress(PublicKey2), nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Account1 balance after transfer: ", account1Balance)
	fmt.Println("Account2 balance after transfer: ", account2Balance)
}
