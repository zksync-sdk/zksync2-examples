package test

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"math/big"
	"testing"
)

func TestWithdrawToken(t *testing.T) {
	amount := big.NewInt(5)

	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	ethClient, err := ethclient.Dial(EthereumProvider)
	assert.NoError(t, err, "ethclient.Dial should not return an error")
	defer ethClient.Close()

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, ethClient)
	assert.NoError(t, err, "NewWallet should not return an error")

	l2BalanceBeforeWithdrawal, err := wallet.Balance(context.Background(), L2Dai, nil)
	assert.NoError(t, err, "Balance should not return an error")

	withdrawTx, err := wallet.Withdraw(nil, accounts.WithdrawalTransaction{
		To:     wallet.Address(),
		Amount: amount,
		Token:  L2Dai,
	})
	assert.NoError(t, err, "Withdraw should not return an error")

	withdrawReceipt, err := client.WaitFinalized(context.Background(), withdrawTx.Hash())
	assert.NoError(t, err, "client.WaitMined should not return an error")
	assert.NotNil(t, withdrawReceipt.BlockHash, "Withdraw transaction should be mined")

	isWithdrawFinalized, err := wallet.IsWithdrawFinalized(nil, withdrawTx.Hash(), 0)
	assert.NoError(t, err, "IsWithdrawFinalized should not return an error")
	assert.False(t, isWithdrawFinalized, "Withdraw transaction should not be finalized")

	finalizeWithdrawTx, err := wallet.FinalizeWithdraw(nil, withdrawTx.Hash(), 0)
	assert.NoError(t, err, "FinalizeWithdraw should not return an error")

	finalizeWithdrawReceipt, err := bind.WaitMined(context.Background(), ethClient, finalizeWithdrawTx)
	assert.NoError(t, err, "bind.WaitMined should not return an error")
	assert.NotNil(t, finalizeWithdrawReceipt.BlockHash, "Finalize withdraw transaction should be mined")

	l2BalanceAfterWithdrawal, err := wallet.Balance(context.Background(), L2Dai, nil)
	assert.NoError(t, err, "Balance should not return an error")

	assert.True(t, new(big.Int).Sub(l2BalanceBeforeWithdrawal, l2BalanceAfterWithdrawal).Cmp(amount) >= 0, "Balance on L2 should be decreased")
}
