package test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"os"
	"testing"
)

func TestDeployCreateWithDeps(t *testing.T) {
	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, nil)
	assert.NoError(t, err, "NewWallet should not return an error")

	demoBytecode, err := os.ReadFile("../../solidity/demo/build/Demo.zbin")
	assert.NoError(t, err, "ReadFile should not return an error")

	fooBytecode, err := os.ReadFile("../../solidity/demo/build/Foo.zbin")
	assert.NoError(t, err, "ReadFile should not return an error")

	hash, err := wallet.DeployWithCreate(nil, accounts.CreateTransaction{
		Bytecode:     demoBytecode,
		Dependencies: [][]byte{fooBytecode},
	})
	assert.NoError(t, err, "DeployWithCreate should not return an error")

	_, err = client.WaitMined(context.Background(), hash)
	assert.NoError(t, err, "client.WaitMined should not return an error")

	receipt, err := client.TransactionReceipt(context.Background(), hash)
	assert.NoError(t, err, "client.TransactionReceipt should not return an error")

	contractAddress := receipt.ContractAddress
	assert.NotNil(t, contractAddress, "Contract should be deployed")
}
