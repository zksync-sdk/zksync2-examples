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

func TestDepositToken(t *testing.T) {
	amount := big.NewInt(5)

	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	ethClient, err := ethclient.Dial(EthereumProvider)
	assert.NoError(t, err, "ethclient.Dial should not return an error")
	defer ethClient.Close()

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, ethClient)
	assert.NoError(t, err, "NewWallet should not return an error")

	l2BalanceBeforeDeposit, err := wallet.Balance(context.Background(), L2Dai, nil)
	assert.NoError(t, err, "Balance should not return an error")

	l1BalanceBeforeDeposit, err := wallet.BalanceL1(nil, L1Dai)
	assert.NoError(t, err, "BalanceL1 should not return an error")

	tx, err := wallet.Deposit(nil, accounts.DepositTransaction{
		To:              wallet.Address(),
		Token:           L1Dai,
		Amount:          amount,
		ApproveERC20:    true,
		RefundRecipient: wallet.Address(),
	})
	assert.NoError(t, err, "Deposit should not return an error")

	l1Receipt, err := bind.WaitMined(context.Background(), ethClient, tx)
	assert.NoError(t, err, "bind.WaitMined should not return an error")

	l2Tx, err := client.L2TransactionFromPriorityOp(context.Background(), l1Receipt)
	assert.NoError(t, err, "L2TransactionFromPriorityOp should not return an error")

	l2Receipt, err := client.WaitMined(context.Background(), l2Tx.Hash)
	assert.NoError(t, err, "bind.WaitMined should not return an error")
	assert.NotNil(t, l2Receipt.BlockHash, "Transaction should be mined")

	l2BalanceAfterDeposit, err := wallet.Balance(context.Background(), L2Dai, nil)
	assert.NoError(t, err, "Balance should not return an error")

	l1BalanceAfterDeposit, err := wallet.BalanceL1(nil, L1Dai)
	assert.NoError(t, err, "BalanceL1 should not return an error")

	assert.True(t, new(big.Int).Sub(l2BalanceAfterDeposit, l2BalanceBeforeDeposit).Cmp(amount) >= 0, "Balance on L2 should be increased")
	assert.True(t, new(big.Int).Sub(l1BalanceBeforeDeposit, l1BalanceAfterDeposit).Cmp(amount) >= 0, "Balance on L1 should be decreased")
}
