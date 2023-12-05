import {Provider, types, Wallet, ContractFactory, Contract} from "zksync2-js";
import {ethers, Typed} from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider);

async function main() {
    const conf = require("../../solidity/storage/build/combined.json");
    const abi = conf.contracts["Storage.sol:Storage"].abi;
    const bytecode: string = conf.contracts["Storage.sol:Storage"].bin;

    const factory = new ContractFactory(abi, bytecode, wallet);
    const storage = await factory.deploy() as Contract;
    console.log(`Contract address: ${await storage.getAddress()}`);

    console.log(`Value: ${await storage.get()}`);

    const tx = await storage.set(Typed.uint256(200));
    await tx.wait();

    console.log(`Value: ${await storage.get()}`);
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})