package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/types"
	"github.com/zksync-sdk/zksync2-go/utils"
	"log"
	"math/big"
	"os"
	"zksync2-examples/contracts/crown"
)

/*
This example demonstrates how to use a paymaster to facilitate fee payment with an ERC20 token.
The user initiates a mint transaction that is configured to be paid with an ERC20 token through the paymaster.
During transaction execution, the paymaster receives the ERC20 token from the user and covers the transaction fee using ETH.
*/
func main() {
	var (
		PrivateKey       = os.Getenv("PRIVATE_KEY")
		ZkSyncProvider   = "https://testnet.era.zksync.dev"
		TokenAddress     = "Cd9BDa1d0FC539043D4C80103bdF4f9cb108931B" // Crown token which can be minted for free
		PaymasterAddress = "d660c2F92d3d0634e5A20f26821C43F1b09298fe" // Paymaster for Crown token
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

	// Get token contract
	token, err := crown.NewCrown(common.HexToAddress(TokenAddress), zp)
	if err != nil {
		log.Panic(token)
	}

	// Transfer some ETH to paymaster, so it can pay fee with ETH
	hash, err := w.Transfer(common.HexToAddress(PaymasterAddress),
		big.NewInt(2_00_000_000_000_000_000), // 0,2 ETH
		utils.CreateETH(),
		nil)
	if err != nil {
		return
	}

	_, err = zp.WaitMined(context.Background(), hash)
	if err != nil {
		log.Panic(err)
	}

	// Also mint some tokens to user account, so user can pay fee with tokens
	privateKey, err := crypto.ToECDSA(common.Hex2Bytes(PrivateKey))
	if err != nil {
		log.Panic(err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, w.GetEthSigner().GetDomain().ChainId)
	if err != nil {
		log.Panic(err)
	}
	tx, err := token.Mint(opts, w.GetAddress(), big.NewInt(10))
	if err != nil {
		log.Panic(err)
	}
	_, err = zp.WaitMined(context.Background(), tx.Hash())
	if err != nil {
		log.Panic(err)
	}

	// Read token and ETH balances from user and paymaster accounts
	balance, err := w.GetBalance()
	if err != nil {
		return
	}
	fmt.Println("Account balance before mint: ", balance)

	tokenBalance, err := token.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Account token balance before mint: ", tokenBalance)

	balance, err = zp.GetBalance(common.HexToAddress(PaymasterAddress), types.BlockNumberLatest)
	if err != nil {
		return
	}
	fmt.Println("Paymaster balance before mint: ", balance)

	tokenBalance, err = token.BalanceOf(nil, common.HexToAddress(TokenAddress))
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Paymaster token balance before mint: ", tokenBalance)

	abi, err := crown.CrownMetaData.GetAbi()
	if err != nil {
		log.Panic(err)
	}

	// Encode mint function from token contract
	calldata, err := abi.Pack("mint", w.GetAddress(), big.NewInt(7))
	if err != nil {
		log.Panic(err)
	}

	// Create paymaster parameters with encoded paymaster input
	paymasterParams, err := utils.GetPaymasterParams(
		common.HexToAddress(PaymasterAddress),
		&types.ApprovalBasedPaymasterInput{
			Token:            common.HexToAddress(TokenAddress),
			MinimalAllowance: big.NewInt(1),
			InnerInput:       []byte{},
		})
	if err != nil {
		log.Panic(err)
	}

	// In order to use paymaster, EIP712 transaction
	// need to be crated with configured paymaster parameters
	mintTx := utils.CreateFunctionCallTransaction(
		w.GetAddress(),
		common.HexToAddress(TokenAddress),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		calldata,
		nil,
		paymasterParams,
	)

	nonce, err := w.GetNonce()
	if err != nil {
		log.Panic(nonce)
	}

	// Execute transaction
	hash, err = w.EstimateAndSend(mintTx, nonce)
	if err != nil {
		log.Panic(err)
	}

	_, err = zp.WaitMined(context.Background(), hash)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Tx: ", hash)

	balance, err = w.GetBalance()
	if err != nil {
		return
	}
	fmt.Println("Account balance after mint: ", balance)

	tokenBalance, err = token.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Account token balance after mint: ", tokenBalance)

	balance, err = zp.GetBalance(common.HexToAddress(PaymasterAddress), types.BlockNumberLatest)
	if err != nil {
		return
	}
	fmt.Println("Paymaster balance after mint: ", balance)

	tokenBalance, err = token.BalanceOf(nil, common.HexToAddress(TokenAddress))
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Paymaster token balance after mint: ", tokenBalance)

}
