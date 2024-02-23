package test

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"math/big"
	"os"
	"testing"
	"zksync2-examples/contracts/storage"
)

func TestDeployCreate(t *testing.T) {
	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, nil)
	assert.NoError(t, err, "NewWallet should not return an error")

	bytecode, err := os.ReadFile("../../solidity/storage/build/Storage.zbin")
	assert.NoError(t, err, "ReadFile should not return an error")

	hash, err := wallet.DeployWithCreate(nil, accounts.CreateTransaction{Bytecode: bytecode})
	assert.NoError(t, err, "DeployWithCreate should not return an error")

	receipt, err := client.WaitMined(context.Background(), hash)
	assert.NoError(t, err, "client.WaitMined should not return an error")

	contractAddress := receipt.ContractAddress
	assert.NotNil(t, contractAddress, "Contract should be deployed")

	storageContract, err := storage.NewStorage(contractAddress, client)
	assert.NoError(t, err, "NewStorage should not return an error")

	abi, err := storage.StorageMetaData.GetAbi()
	assert.NoError(t, err, "GetAbi should not return an error")

	setArguments, err := abi.Pack("set", big.NewInt(700))
	assert.NoError(t, err, "Pack should not return an error")

	result, err := wallet.CallContract(context.Background(), accounts.CallMsg{
		To:   &contractAddress,
		Data: setArguments,
	}, nil)
	assert.NotNil(t, result, "Result from contract should be non-nil value")

	value, err := storageContract.Get(nil)
	assert.NoError(t, err, "Get should not return an error")
	assert.True(t, value.Cmp(big.NewInt(0)) == 0, "Values should be the same")

	opts, err := bind.NewKeyedTransactorWithChainID(wallet.Signer().PrivateKey(), wallet.Signer().Domain().ChainId)
	assert.NoError(t, err, "NewKeyedTransactorWithChainID should not return an error")

	tx, err := storageContract.Set(opts, big.NewInt(200))
	assert.NoError(t, err, "Set should not return an error")

	_, err = client.WaitMined(context.Background(), tx.Hash())
	assert.NoError(t, err, "client.WaitMined should not return an error")

	value, err = storageContract.Get(nil)
	assert.NoError(t, err, "Get should not return an error")
	assert.True(t, value.Cmp(big.NewInt(200)) == 0, "Values should be the same")

	setArguments, err = abi.Pack("set", big.NewInt(500))
	assert.NoError(t, err, "Pack should not return an error")

	execute, err := wallet.SendTransaction(context.Background(), &accounts.Transaction{
		To:   &contractAddress,
		Data: setArguments,
	})
	assert.NoError(t, err, "SendTransaction should not return an error")

	_, err = client.WaitMined(context.Background(), execute)
	assert.NoError(t, err, "client.WaitMined should not return an error")

	value, err = storageContract.Get(nil)
	assert.NoError(t, err, "Get should not return an error")
	assert.True(t, value.Cmp(big.NewInt(500)) == 0, "Values should be the same")
}
