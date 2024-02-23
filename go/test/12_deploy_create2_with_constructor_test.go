package test

//
//import (
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/stretchr/testify/assert"
//	"github.com/zksync-sdk/zksync2-go/accounts"
//	"github.com/zksync-sdk/zksync2-go/clients"
//	"math/big"
//	"os"
//	"testing"
//	"zksync2-examples/contracts/incrementer"
//)
//
//func TestDeployCreate2WithConstructor(t *testing.T) {
//	client, err := clients.Dial(ZkSyncEraProvider)
//	defer client.Close()
//	assert.NoError(t, err, "clients.Dial should not return an error")
//
//	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, nil)
//	assert.NoError(t, err, "NewWallet should not return an error")
//
//	bytecode, err := os.ReadFile("../../solidity/incrementer/build/Incrementer.zbin")
//	assert.NoError(t, err, "ReadFile should not return an error")
//
//	abi, err := incrementer.IncrementerMetaData.GetAbi()
//	assert.NoError(t, err, "GetAbi should not return an error")
//
//	constructor, err := abi.Pack("", big.NewInt(2))
//	assert.NoError(t, err, "Pack should not return an error")
//
//	hash, err := wallet.Deploy(nil, accounts.Create2Transaction{
//		Bytecode: bytecode,
//		Calldata: constructor,
//	})
//	assert.NoError(t, err, "Deploy should not return an error")
//	assert.NotNil(t, hash, "Deployment transaction hash should not be nil")
//
//	//receipt, err := client.WaitMined(context.Background(), hash)
//	//assert.NoError(t, err, "client.WaitMined should not return an error")
//	//
//	//contractAddress := receipt.ContractAddress
//	//assert.NotNil(t, contractAddress, "Contract should be deployed")
//}
