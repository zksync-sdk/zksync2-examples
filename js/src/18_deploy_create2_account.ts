import { Provider, types, Wallet, ContractFactory } from "zksync-ethers";
import { ethers } from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider);

const tokenAddress = "0x765F5AF819D818a8e8ee6ff63D8d0e8056DBE150";

async function main() {
    const conf = require("../../solidity/custom_paymaster/paymaster/build/Paymaster.json");
    const abi = conf.abi;
    const bytecode: string = conf.bytecode;

    const factory = new ContractFactory(abi, bytecode, wallet, "create2Account");
    const account = await factory.deploy(tokenAddress, {
        customData: { salt: ethers.hexlify(ethers.randomBytes(32)) },
    });
    console.log(`Account address: ${await account.getAddress()}`);
}

main()
    .then()
    .catch((error) => {
        console.log(`Error: ${error}`);
    });
