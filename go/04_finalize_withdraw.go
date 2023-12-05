package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"log"
	"os"
)

func main() {
	var (
		PrivateKey        = os.Getenv("PRIVATE_KEY")
		ZkSyncEraProvider = "https://sepolia.era.zksync.dev"
		EthereumProvider  = "https://rpc.ankr.com/eth_sepolia"
		WithdrawTx        = common.HexToHash("<Withdraw tx hash>")
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

	// Perform the finalize withdrawal if it has not been already finalized.
	// On the testnet, withdrawals are automatically finalized. For additional information, please refer
	// to the documentation: https://era.zksync.io/docs/reference/concepts/bridging-asset.html#withdrawals-to-l1.
	// There is no need to execute FinalizeWithdraw, otherwise, an error with code `jj` would occur.
	if isFinalized, errFinalized := wallet.IsWithdrawFinalized(nil, WithdrawTx, 0); errFinalized != nil {
		log.Panic(errFinalized)
	} else if !isFinalized {
		finalizeWithdrawTx, errWithdraw := wallet.FinalizeWithdraw(nil, WithdrawTx, 0)
		if errWithdraw != nil {
			log.Panic(errWithdraw)
		}
		fmt.Println("Finalize withdraw transaction: ", finalizeWithdrawTx.Hash())

		// Wait for finalize withdraw transaction to be finalized on L1 network
		fmt.Println("Waiting for finalize withdraw transaction to be finalized on L1 network")
		_, err = bind.WaitMined(context.Background(), ethClient, finalizeWithdrawTx)
		if err != nil {
			log.Panic(err)
		}
	}

}
