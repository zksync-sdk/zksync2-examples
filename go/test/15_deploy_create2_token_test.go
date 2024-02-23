package test

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"math/big"
	"testing"
	"zksync2-examples/contracts/token"
)

func TestDeployToken(t *testing.T) {
	client, err := clients.Dial(ZkSyncEraProvider)
	defer client.Close()
	assert.NoError(t, err, "clients.Dial should not return an error")

	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, nil)
	assert.NoError(t, err, "NewWallet should not return an error")

	tokenAbi, err := token.TokenMetaData.GetAbi()
	assert.NoError(t, err, "GetAbi should not return an error")

	tokenConstructor, err := tokenAbi.Pack("", "Crown", "Crown", uint8(18))
	assert.NoError(t, err, "Pack should not return an error")

	tokenDeployHash, err := wallet.Deploy(nil, accounts.Create2Transaction{
		Bytecode: common.FromHex(token.TokenMetaData.Bin),
		Calldata: tokenConstructor,
	})
	assert.NoError(t, err, "Deploy should not return an error")

	tokenDeployReceipt, err := client.WaitMined(context.Background(), tokenDeployHash)
	assert.NoError(t, err, "client.WaitMined should not return an error")

	tokenAddress := tokenDeployReceipt.ContractAddress
	assert.NotNil(t, tokenAddress, "Contract should be deployed")

	tokenContract, err := token.NewToken(tokenAddress, client)
	assert.NoError(t, err, "NewToken should not return an error")

	opts, err := bind.NewKeyedTransactorWithChainID(wallet.Signer().PrivateKey(), wallet.Signer().Domain().ChainId)
	assert.NoError(t, err, "NewKeyedTransactorWithChainID should not return an error")

	mint, err := tokenContract.Mint(opts, wallet.Address(), big.NewInt(5))
	assert.NoError(t, err, "Mint should not return an error")

	_, err = client.WaitMined(context.Background(), mint.Hash())
	assert.NoError(t, err, "client.WaitMined should not return an error")
}
