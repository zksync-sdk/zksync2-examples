package test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/contracts/erc20"
	"math/big"
	"testing"
)

func TestTransferToken(t *testing.T) {
	tokenData := readToken()
	amount := big.NewInt(5)
	tokenAddress := common.HexToAddress(tokenData.L2Address)
	receiver := common.HexToAddress(Receiver)

	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	ethClient, err := ethclient.Dial(EthereumProvider)
	assert.NoError(t, err, "ethclient.Dial should not return an error")
	defer ethClient.Close()

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, ethClient)
	assert.NoError(t, err, "NewWallet should not return an error")

	tokenContract, err := erc20.NewIERC20(tokenAddress, client)
	assert.NoError(t, err, "NewIERC20 should not return an error")

	balanceBeforeTransferSender, err := wallet.Balance(context.Background(), tokenAddress, nil)
	assert.NoError(t, err, "Balance should not return an error")

	balanceBeforeTransferReceiver, err := tokenContract.BalanceOf(nil, receiver)
	assert.NoError(t, err, "BalanceOf should not return an error")

	tx, err := wallet.Transfer(nil, accounts.TransferTransaction{
		To:     receiver,
		Amount: amount,
		Token:  tokenAddress,
	})
	assert.NoError(t, err, "Transfer should not return an error")

	receipt, err := client.WaitMined(context.Background(), tx.Hash())
	assert.NoError(t, err, "client.WaitMined should not return an error")
	assert.NotNil(t, receipt.BlockHash, "Transaction should be mined")

	balanceAfterTransferSender, err := wallet.Balance(context.Background(), tokenAddress, nil)
	assert.NoError(t, err, "Balance should not return an error")

	balanceAfterTransferReceiver, err := tokenContract.BalanceOf(nil, receiver)
	assert.NoError(t, err, "BalanceOf should not return an error")

	assert.True(t, new(big.Int).Sub(balanceBeforeTransferSender, balanceAfterTransferSender).Cmp(amount) >= 0, "Sender balance should be decreased")
	assert.True(t, new(big.Int).Sub(balanceAfterTransferReceiver, balanceBeforeTransferReceiver).Cmp(amount) >= 0, "Receiver balance should be increased")
}
