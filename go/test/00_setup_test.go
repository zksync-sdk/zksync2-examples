package test

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"log"
	"math/big"
	"os"
	"testing"
	"time"
)

const EthereumProvider = "http://localhost:8545"
const ZkSyncEraProvider = "http://localhost:3050"
const PrivateKey = "7726827caac94a7f9e1b160f7ea819f172f7b6f9d2a97f992c38edeab82d4110"

var Receiver = common.HexToAddress("0xa61464658AfeAf65CccaaFD3a512b69A83B77618")

var L1Tokens []TokenData
var L2Dai common.Address
var L1Dai common.Address

type TokenData struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Decimals int            `json:"decimals"`
	Address  common.Address `json:"address"`
}

const TokenPath = "./testdata/tokens.json"

func readTokens() []TokenData {
	file, err := os.Open(TokenPath)
	if err != nil {
		log.Printf("Could not find tokens.json")
		return nil
	}

	var tokens []TokenData
	decoder := json.NewDecoder(file)
	if errDecode := decoder.Decode(&tokens); err != nil {
		log.Fatalf("Error decoding JSON: %s", errDecode)
	}
	return tokens
}

func createTokenL2(wallet *accounts.Wallet, client clients.Client, ethClient *ethclient.Client, l1Token common.Address) (common.Address, common.Hash, common.Hash) {
	tx, err := wallet.Deposit(nil, accounts.DepositTransaction{
		Token:           l1Token,
		Amount:          big.NewInt(30),
		To:              wallet.Address(),
		ApproveERC20:    true,
		RefundRecipient: wallet.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = bind.WaitMined(context.Background(), ethClient, tx)
	if err != nil {
		log.Fatal(err)
	}

	l1Receipt, err := ethClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}

	l2Tx, err := client.L2TransactionFromPriorityOp(context.Background(), l1Receipt)
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.WaitMined(context.Background(), l2Tx.Hash)
	if err != nil {
		log.Fatal(err)
	}

	tokenL2Address, err := client.L2TokenAddress(context.Background(), l1Token)
	if err != nil {
		log.Fatal(err)
	}

	return tokenL2Address, tx.Hash(), l2Tx.Hash
}

func wait() {
	const maxAttempts = 30

	nodeURL := "http://localhost:3050"
	client, err := clients.Dial(nodeURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	for i := 0; i < maxAttempts; i++ {
		_, err := client.NetworkID(context.Background())
		if err == nil {
			log.Println("Node is ready to receive traffic.")
			return
		}

		log.Println("Node not ready yet. Retrying...")
		time.Sleep(20 * time.Second)
	}

	log.Fatal("Maximum retries exceeded.")
}

func TestMain(m *testing.M) {
	wait()

	client, err := clients.Dial(ZkSyncEraProvider)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ethClient, err := ethclient.Dial(EthereumProvider)
	if err != nil {
		log.Fatal(err)
	}

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, ethClient)
	if err != nil {
		log.Fatal(err)
	}

	L1Tokens = readTokens()
	L1Dai = L1Tokens[0].Address

	L2Dai, _, _ = createTokenL2(wallet, client, ethClient, L1Dai)

	os.Exit(m.Run())
}
