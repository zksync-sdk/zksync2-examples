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

	// Show balances before deposit
	balance, err := w.GetBalance()
	if err != nil {
		log.Panic(err)
	}
	tokenBalance, err := tokenL1.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Balance before deposit on L1 network: ", balance)
	fmt.Println("Token balance before deposit on L1 network: ", tokenBalance)

	depositAmount := big.NewInt(5)

	// Bridging ERC20 tokens from Ethereum requires approving the tokens to the zkSync Ethereum smart contract
	tx, err := ep.ApproveDeposit(&types.Token{L1Address: TokenL1Address}, depositAmount, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("L1 approve deposit transaction: ", tx.Hash())

	fmt.Println("Waiting for approve deposit transaction to be finalized on L1 network")
	_, err = ep.WaitMined(context.Background(), tx.Hash())
	if err != nil {
		log.Panic(err)
	}

	// Perform deposit
	tx, err = ep.Deposit(
		&types.Token{L1Address: TokenL1Address},
		depositAmount,
		w.GetAddress(),
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("L1 deposit transaction: ", tx.Hash())

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
	tokenBalance, err = tokenL1.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Balance after deposit on L1 network: ", balance)
	fmt.Println("Token balance after deposit on L1 network: ", tokenBalance)

	tokenL2Address, err := zp.L2TokenAddress(TokenL1Address)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Token L2 address: ", tokenL2Address)
	// Get token contract on zkSync network
	tokenL2, err := erc20.NewIERC20(tokenL2Address, zp)
	if err != nil {
		log.Panic(err)
	}

	tokenL2Balance, err := tokenL2.BalanceOf(nil, w.GetAddress())
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Token balance on L2 network: ", tokenL2Balance)
}
