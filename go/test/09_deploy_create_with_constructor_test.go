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
	"zksync2-examples/contracts/incrementer"
)

func TestDeployCreateWithConstructor(t *testing.T) {
	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, nil)
	assert.NoError(t, err, "NewWallet should not return an error")

	bytecode, err := os.ReadFile("../../solidity/incrementer/build/Incrementer.zbin")
	assert.NoError(t, err, "ReadFile should not return an error")

	abi, err := incrementer.IncrementerMetaData.GetAbi()
	assert.NoError(t, err, "GetAbi should not return an error")

	constructor, err := abi.Pack("", big.NewInt(2))
	assert.NoError(t, err, "Pack should not return an error")

	hash, err := wallet.DeployWithCreate(nil, accounts.CreateTransaction{
		Bytecode: bytecode,
		Calldata: constructor,
	})
	assert.NoError(t, err, "DeployWithCreate should not return an error")

	receipt, err := client.WaitMined(context.Background(), hash)
	assert.NoError(t, err, "client.WaitMined should not return an error")

	contractAddress := receipt.ContractAddress
	assert.NotNil(t, contractAddress, "Contract should be deployed")

	incrementerContract, err := incrementer.NewIncrementer(contractAddress, client)
	assert.NoError(t, err, "NewIncrementer should not return an error")

	value, err := incrementerContract.Get(nil)
	assert.NoError(t, err, "Get should not return an error")
	assert.True(t, value.Cmp(big.NewInt(0)) == 0, "Values should be the same")

	opts, err := bind.NewKeyedTransactorWithChainID(wallet.Signer().PrivateKey(), wallet.Signer().Domain().ChainId)
	assert.NoError(t, err, "NewKeyedTransactorWithChainID should not return an error")

	tx, err := incrementerContract.Increment(opts)
	assert.NoError(t, err, "Increment should not return an error")

	_, err = client.WaitMined(context.Background(), tx.Hash())
	assert.NoError(t, err, "client.WaitMined should not return an error")

	value, err = incrementerContract.Get(nil)
	assert.NoError(t, err, "Get should not return an error")
	assert.True(t, value.Cmp(big.NewInt(2)) == 0, "Values should be the same")
}
