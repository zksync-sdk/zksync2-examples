import {Provider, types, Wallet, ContractFactory} from "zksync2-js";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider);

const tokenAddress = "0xA70dF8446A6AeA0017D60e97e816e141aa28759b";

async function main() {
    const conf = require("../../solidity/custom_paymaster/paymaster/build/Paymaster.json");
    const abi = conf.abi;
    const bytecode: string = conf.bytecode;

    const factory = new ContractFactory(abi, bytecode, wallet, 'createAccount');
    const account = await factory.deploy(tokenAddress);
    console.log(`Account address: ${await account.getAddress()}`);
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})