import {Provider, types, Wallet, ContractFactory, Contract} from "zksync2-js";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider);

async function main() {
    const conf = require("../../solidity/incrementer/build/combined.json");
    const abi = conf.contracts["Incrementer.sol:Incrementer"].abi;
    const bytecode: string = conf.contracts["Incrementer.sol:Incrementer"].bin;

    const factory = new ContractFactory(abi, bytecode, wallet);
    const incrementer = await factory.deploy(2) as Contract;
    console.log(`Contract address: ${await incrementer.getAddress()}`);

    console.log(`Value before Increment method execution: ${await incrementer.get()}`);

    const tx = await incrementer.increment();
    await tx.wait();

    console.log(`Value after Increment method execution: ${await incrementer.get()}`);
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})