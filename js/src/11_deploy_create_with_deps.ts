import {Provider, types, Wallet, ContractFactory, Contract} from "zksync2-js";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider);

async function main() {
    const conf = require("../../solidity/demo/build/combined.json");
    const abi = conf.contracts["Demo.sol:Demo"].abi;
    const bytecode: string = conf.contracts["Demo.sol:Demo"].bin;

    const factory = new ContractFactory(abi, bytecode, wallet);
    const demo = await factory.deploy({
        customData: {factoryDeps: [conf.contracts["Foo.sol:Foo"].bin]}
    }) as Contract;
    console.log(`Contract address: ${await demo.getAddress()}`);

    console.log(`Value: ${await demo.getFooName()}`);
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})