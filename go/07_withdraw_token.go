package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"log"
	"math/big"
	"os"
)

func main() {
	var (
		PrivateKey        = os.Getenv("PRIVATE_KEY")
		ZkSyncEraProvider = "https://sepolia.era.zksync.dev"
		EthereumProvider  = "https://rpc.ankr.com/eth_sepolia"
		TokenL2Address    = common.HexToAddress("0x1958F3a8246B526796DdE3F37fB2b9E04660Bf33")
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

	// Perform withdraw
	tx, err := wallet.Withdraw(nil, accounts.WithdrawalTransaction{
		To:     wallet.Address(),
		Amount: big.NewInt(1),
		Token:  TokenL2Address,
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Withdraw transaction: ", tx.Hash())

	// The duration for submitting a withdrawal transaction to L1 can last up to 24 hours. For additional information,
	// please refer to the documentation: https://era.zksync.io/docs/reference/troubleshooting/withdrawal-delay.html.
	// Once the withdrawal transaction is submitted on L1, it needs to be finalized.
	// To learn more about how to achieve this, please take a look at the 04_finalize_withdraw.go script.
}
