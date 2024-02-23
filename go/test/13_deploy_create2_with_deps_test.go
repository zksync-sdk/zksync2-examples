package test

//func TestDeployCreate2WithDeps(t *testing.T) {
//	client, err := clients.Dial(ZkSyncEraProvider)
//	defer client.Close()
//	assert.NoError(t, err, "clients.Dial should not return an error")
//
//	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), &client, nil)
//	assert.NoError(t, err, "NewWallet should not return an error")
//
//	demoBytecode, err := os.ReadFile("../../solidity/demo/build/Demo.zbin")
//	assert.NoError(t, err, "ReadFile should not return an error")
//
//	fooBytecode, err := os.ReadFile("../../solidity/demo/build/Foo.zbin")
//	assert.NoError(t, err, "ReadFile should not return an error")
//
//	hash, err := wallet.Deploy(nil, accounts.Create2Transaction{
//		Bytecode:     demoBytecode,
//		Dependencies: [][]byte{fooBytecode},
//	})
//	assert.NoError(t, err, "DeployWithCreate should not return an error")
//
//	_, err = client.WaitMined(context.Background(), hash)
//	assert.NoError(t, err, "client.WaitMined should not return an error")
//
//	contractAddress, err := utils.Create2Address(
//		wallet.Address(),
//		demoBytecode,
//		nil,
//		nil,
//	)
//	assert.NoError(t, err, "utils.Create2Address should not return an error")
//	assert.NotNil(t, contractAddress, "Contract should be deployed")
//
//	demoContract, err := demo.NewDemo(contractAddress, client)
//	assert.NoError(t, err, "NewDemo should not return an error")
//
//	value, err := demoContract.GetFooName(nil)
//	assert.NoError(t, err, "Get should not return an error")
//	assert.Equal(t, "Foo", value, "Values should be the same")
//}
