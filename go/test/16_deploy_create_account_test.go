package test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"github.com/zksync-sdk/zksync2-go/utils"
	"testing"
)

func TestDeployAccountWithCreate(t *testing.T) {
	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, nil)
	assert.NoError(t, err, "NewWallet should not return an error")

	_, paymasterAbi, bytecode, err := utils.ReadStandardJson("../../solidity/custom_paymaster/paymaster/build/Paymaster.json")
	assert.NoError(t, err, "ReadStandardJson should not return an error")

	constructor, err := paymasterAbi.Pack("", L2Dai)
	assert.NoError(t, err, "Pack should not return an error")

	hash, err := wallet.DeployAccountWithCreate(nil, accounts.CreateTransaction{
		Bytecode: bytecode,
		Calldata: constructor,
	})
	assert.NoError(t, err, "DeployWithCreate should not return an error")

	receipt, err := client.WaitMined(context.Background(), hash)
	assert.NoError(t, err, "client.WaitMined should not return an error")

	contractAddress := receipt.ContractAddress
	assert.NotNil(t, contractAddress, "Contract should be deployed")
}
