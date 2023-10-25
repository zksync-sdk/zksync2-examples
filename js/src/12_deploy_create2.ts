import {Provider, types, Wallet, ContractFactory} from "zksync2-js";
import {ethers, Contract, Typed} from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Goerli);
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider);

async function main() {
    const conf = require("../../solidity/storage/build/combined.json");
    const abi = conf.contracts["Storage.sol:Storage"].abi;
    const bytecode: string = conf.contracts["Storage.sol:Storage"].bin;

    const factory = new ContractFactory(abi, bytecode, wallet, "create2");
    const contract = await factory.deploy({
        customData: {salt: ethers.hexlify(ethers.randomBytes(32))}
    });
    const contractAddress = await contract.getAddress();
    console.log(`Contract address: ${contractAddress}`);

    const storage = new Contract(contractAddress, abi, wallet);
    console.log(`Value: ${await storage.get()}`);

    const tx = await storage.set(Typed.uint256(200));
    await tx.wait();

    console.log(`Value: ${await storage.get()}`);
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})