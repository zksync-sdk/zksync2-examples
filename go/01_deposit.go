package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
	"log"
	"math/big"
	"os"
)

func main() {
	var (
		PrivateKey       = os.Getenv("PRIVATE_KEY")
		ZkSyncProvider   = "https://testnet.era.zksync.dev"
		EthereumProvider = "https://rpc.ankr.com/eth_goerli"
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

	// Connect to Ethereum network
	ethRpc, err := rpc.Dial(EthereumProvider)
	if err != nil {
		log.Panic(err)
	}
	ep, err := w.CreateEthereumProvider(ethRpc)
	if err != nil {
		log.Panic(err)
	}

	// Show balance before deposit
	balance, err := w.GetBalance()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Balance before deposit: ", balance)

	// Perform deposit
	tx, err := ep.Deposit(
		utils.CreateETH(),
		big.NewInt(1_000_000_000),
		w.GetAddress(),
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("L1 transaction: ", tx.Hash())

	// Wait for deposit transaction to be finalized on L1 network
	fmt.Println("Waiting for deposit transaction to be finalized on L1 network")
	_, err = ep.WaitMined(context.Background(), tx.Hash())
	if err != nil {
		log.Panic(err)
	}

	// Get transaction receipt for deposit transaction on L1 network
	l1Receipt, err := ep.GetClient().TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Panic(err)
	}

	// Get deposit transaction hash on L2 network
	l2Hash, err := ep.GetL2HashFromPriorityOp(l1Receipt)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("L2 transaction", l2Hash)

	// Wait for deposit transaction to be finalized on L2 network (5-7 minutes)
	fmt.Println("Waiting for deposit transaction to be finalized on L2 network (5-7 minutes)")
	_, err = zp.WaitMined(context.Background(), l2Hash)
	if err != nil {
		log.Panic(err)
	}

	balance, err = w.GetBalance()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Balance after deposit: ", balance)

	/*
		// ClaimFailedDeposit is used when transaction on L2 has failed
		cfdHash, err := w.ClaimFailedDeposit(l2Hash, ep)
		if err != nil {
			fmt.Println(err) // this should trigger if deposit succeed
		}
		fmt.Println("ClaimFailedDeposit hash", cfdHash)
	*/
}
