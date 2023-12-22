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

func TestWithdraw(t *testing.T) {
	amount := big.NewInt(7_000_000_000)

	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	ethClient, err := ethclient.Dial(EthereumProvider)
	assert.NoError(t, err, "ethclient.Dial should not return an error")
	defer ethClient.Close()

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, ethClient)
	assert.NoError(t, err, "NewWallet should not return an error")

	l2BalanceBeforeWithdrawal, err := wallet.Balance(context.Background(), utils.EthAddress, nil)
	assert.NoError(t, err, "Balance should not return an error")

	//l1BalanceBeforeWithdrawal, err := wallet.BalanceL1(nil, l1TokenAddress)
	//assert.NoError(t, err, "BalanceL1 should not return an error")

	tx, err := wallet.Withdraw(nil, accounts.WithdrawalTransaction{
		To:     wallet.Address(),
		Amount: amount,
		Token:  utils.EthAddress,
	})
	assert.NoError(t, err, "Withdraw should not return an error")

	receipt, err := client.WaitMined(context.Background(), tx.Hash())
	assert.NoError(t, err, "client.WaitMined should not return an error")
	assert.NotNil(t, receipt.BlockHash, "Transaction should be mined")

	l2BalanceAfterWithdrawal, err := wallet.Balance(context.Background(), utils.EthAddress, nil)
	assert.NoError(t, err, "Balance should not return an error")

	//l1BalanceAfterWithdrawal, err := wallet.BalanceL1(nil, l1TokenAddress)
	//assert.NoError(t, err, "BalanceL1 should not return an error")

	assert.True(t, new(big.Int).Sub(l2BalanceBeforeWithdrawal, l2BalanceAfterWithdrawal).Cmp(amount) >= 0, "Balance on L2 should be decreased")
	//assert.True(t, new(big.Int).Sub(l1BalanceBeforeDeposit, l1BalanceAfterDeposit).Cmp(amount) >= 0, "Balance on L1 should be decreased")
}
