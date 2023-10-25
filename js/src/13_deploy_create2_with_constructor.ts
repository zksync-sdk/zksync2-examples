import {Provider, types, Wallet, ContractFactory} from "zksync2-js";
import {Contract, ethers} from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Goerli);
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider);

async function main() {
    const conf = require("../../solidity/incrementer/build/combined.json");
    const abi = conf.contracts["Incrementer.sol:Incrementer"].abi;
    const bytecode: string = conf.contracts["Incrementer.sol:Incrementer"].bin;

    const factory = new ContractFactory(abi, bytecode, wallet, "create2");
    const contract = await factory.deploy(2, {
        customData: {salt: ethers.hexlify(ethers.randomBytes(32))}
    });
    const contractAddress = await contract.getAddress();
    console.log(`Contract address: ${contractAddress}`);

    const incrementer = new Contract(contractAddress, abi, wallet);
    console.log(`Value before Increment method execution: ${await incrementer.get()}`);

    const tx = await incrementer.increment();
    await tx.wait();

    console.log(`Value after Increment method execution: ${await incrementer.get()}`);
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})