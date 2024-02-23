package test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
	"math/big"
	"testing"
)

func TestTransfer(t *testing.T) {
	amount := big.NewInt(7_000_000_000)

	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	ethClient, err := ethclient.Dial(EthereumProvider)
	assert.NoError(t, err, "ethclient.Dial should not return an error")
	defer ethClient.Close()

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, ethClient)
	assert.NoError(t, err, "NewWallet should not return an error")

	balanceBeforeTransferSender, err := wallet.Balance(context.Background(), utils.EthAddress, nil)
	assert.NoError(t, err, "Balance should not return an error")

	balanceBeforeTransferReceiver, err := client.BalanceAt(context.Background(), Receiver, nil)
	assert.NoError(t, err, "BalanceAt should not return an error")

	tx, err := wallet.Transfer(nil, accounts.TransferTransaction{
		To:     Receiver,
		Amount: amount,
		Token:  utils.EthAddress,
	})
	assert.NoError(t, err, "Transfer should not return an error")

	receipt, err := client.WaitMined(context.Background(), tx.Hash())
	assert.NoError(t, err, "client.WaitMined should not return an error")
	assert.NotNil(t, receipt.BlockHash, "Transaction should be mined")

	balanceAfterTransferSender, err := wallet.Balance(context.Background(), utils.EthAddress, nil)
	assert.NoError(t, err, "Balance should not return an error")

	balanceAfterTransferReceiver, err := wallet.BalanceL1(nil, utils.EthAddress)
	assert.NoError(t, err, "BalanceL1 should not return an error")

	assert.True(t, new(big.Int).Sub(balanceBeforeTransferSender, balanceAfterTransferSender).Cmp(amount) >= 0, "Sender balance should be decreased")
	assert.True(t, new(big.Int).Sub(balanceAfterTransferReceiver, balanceBeforeTransferReceiver).Cmp(amount) >= 0, "Receiver balance should be increased")
}
