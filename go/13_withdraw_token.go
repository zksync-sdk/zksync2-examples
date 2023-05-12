package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/contracts/erc20"
	"github.com/zksync-sdk/zksync2-go/types"
	"log"
	"math/big"
	"os"
)

func main() {
	var (
		PrivateKey       = os.Getenv("PRIVATE_KEY")
		ZkSyncProvider   = "https://testnet.era.zksync.dev"
		EthereumProvider = "https://rpc.ankr.com/eth_goerli"
		TokenL2Address   = common.HexToAddress("1958F3a8246B526796DdE3F37fB2b9E04660Bf33")
		TokenL1Address   = common.HexToAddress("c8F8cE6491227a6a2Ab92e67a64011a4Eba1C6CF")
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

	// Get token contract on Ethereum network
	tokenL1, err := erc20.NewIERC20(TokenL1Address, ep.GetClient())
	if err != nil {
		log.Panic(err)
	}
	// Get token contract on zkSync network
	tokenL2, err := erc20.NewIERC20(TokenL2Address, zp)
	if err != nil {
		log.Panic(err)
	}

	tokenL1Balance, err := tokenL1.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}
	tokenL2Balance, err := tokenL2.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Token balance before withdrawal on L1 network: ", tokenL1Balance)
	fmt.Println("Token balance before withdrawal on L2 network: ", tokenL2Balance)

	// Perform withdraw
	wHash, err := w.Withdraw(
		w.GetAddress(),
		big.NewInt(1),
		&types.Token{L1Address: TokenL1Address, L2Address: TokenL2Address},
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Withdraw transaction: ", wHash)

	// Wait until transaction is finalized
	fmt.Println("Waiting for finalize transaction to be finalized on L2 network")
	_, err = zp.WaitFinalized(context.Background(), wHash)
	if err != nil {
		panic(err)
	}

	// Perform finalize withdraw
	fwHash, err := w.FinalizeWithdraw(wHash, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finalize withdraw transaction: ", fwHash)

	// Wait for finalize withdraw transaction to be finalized on L1 network
	fmt.Println("Waiting for finalize withdraw transaction to be finalized on L1 network")
	_, err = ep.WaitMined(context.Background(), fwHash)
	if err != nil {
		log.Panic(err)
	}

	tokenL1Balance, err = tokenL1.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}
	tokenL2Balance, err = tokenL2.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Token balance after withdrawal on L1 network: ", tokenL1Balance)
	fmt.Println("Token balance after withdrawal on L2 network: ", tokenL2Balance)

}
