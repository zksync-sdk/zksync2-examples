import {Provider, types, Wallet, ContractFactory} from "zksync2-js";
import {Contract, ethers} from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Goerli);
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider);

async function main() {
    const conf = require("../../solidity/demo/build/combined.json");
    const abi = conf.contracts["Demo.sol:Demo"].abi;
    const bytecode: string = conf.contracts["Demo.sol:Demo"].bin;

    const factory = new ContractFactory(abi, bytecode, wallet, "create2");
    const contract = await factory.deploy({
        customData: {
            salt: ethers.hexlify(ethers.randomBytes(32)),
            factoryDeps: [conf.contracts["Foo.sol:Foo"].bin]
        }
    });
    const contractAddress = await contract.getAddress();
    console.log(`Contract address: ${contractAddress}`);

    const demo = new Contract(contractAddress, abi, wallet);
    console.log(`Value: ${await demo.getFooName()}`);
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})