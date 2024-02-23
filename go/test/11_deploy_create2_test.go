package test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"os"
	"testing"
)

func TestDeployCreate2(t *testing.T) {
	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, nil)
	assert.NoError(t, err, "NewWallet should not return an error")

	bytecode, err := os.ReadFile("../../solidity/storage/build/Storage.zbin")
	assert.NoError(t, err, "ReadFile should not return an error")

	_, err = wallet.Deploy(nil, accounts.Create2Transaction{Bytecode: bytecode})
	//assert.NoError(t, err, "Deploy should not return an error")

	//receipt, err := client.WaitMined(context.Background(), hash)
	//assert.NoError(t, err, "client.WaitMined should not return an error")
	//assert.NotNil(t, receipt.ContractAddress, "Contract should be deployed")
	//
	//contractAddress, err := utils.Create2Address(wallet.Address(), bytecode, nil, nil)
	//assert.Equal(t, contractAddress, receipt.ContractAddress, "Addresses should be the same")
	//
	//storageContract, err := storage.NewStorage(contractAddress, client)
	//assert.NoError(t, err, "NewStorage should not return an error")
	//
	//value, err := storageContract.Get(nil)
	//assert.NoError(t, err, "Get should not return an error")
	//assert.True(t, value.Cmp(big.NewInt(0)) == 0, "Values should be the same")
	//
	//opts, err := bind.NewKeyedTransactorWithChainID(wallet.Signer().PrivateKey(), wallet.Signer().Domain().ChainId)
	//assert.NoError(t, err, "NewKeyedTransactorWithChainID should not return an error")
	//
	//tx, err := storageContract.Set(opts, big.NewInt(200))
	//assert.NoError(t, err, "Set should not return an error")
	//
	//_, err = client.WaitMined(context.Background(), tx.Hash())
	//assert.NoError(t, err, "client.WaitMined should not return an error")
	//
	//value, err = storageContract.Get(nil)
	//assert.NoError(t, err, "Get should not return an error")
	//assert.True(t, value.Cmp(big.NewInt(200)) == 0, "Values should be the same")
}
