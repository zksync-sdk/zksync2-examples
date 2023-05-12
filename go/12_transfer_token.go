package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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
		PrivateKey1    = os.Getenv("PRIVATE_KEY")
		PublicKey2     = "81E9D85b65E9CC8618D85A1110e4b1DF63fA30d9"
		ZkSyncProvider = "https://testnet.era.zksync.dev"
		TokenL2Address = common.HexToAddress("Cd9BDa1d0FC539043D4C80103bdF4f9cb108931B")
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

	// Get token contract on zkSync network
	token, err := erc20.NewIERC20(TokenL2Address, zp)
	if err != nil {
		log.Panic(err)
	}

	account1TokenBalance, err := token.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}
	account2TokenBalance, err := token.BalanceOf(nil, common.HexToAddress(PublicKey2))
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Account 1 token balance before transfer: ", account1TokenBalance)
	fmt.Println("Account 2 token balance before transfer: ", account2TokenBalance)

	// Perform transfer
	hash, err := w.Transfer(
		common.HexToAddress(PublicKey2),
		big.NewInt(3),
		&types.Token{L2Address: TokenL2Address},
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

	account1TokenBalance, err = token.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}
	account2TokenBalance, err = token.BalanceOf(nil, common.HexToAddress(PublicKey2))
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Account 1 token balance after transfer: ", account1TokenBalance)
	fmt.Println("Account 2 token balance after transfer: ", account2TokenBalance)

}
