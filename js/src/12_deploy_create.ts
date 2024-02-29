import { Provider, types, Wallet, ContractFactory, Contract, utils } from "zksync-ethers";
import { Typed } from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider);

const tokenAddress = "0x927488F48ffbc32112F1fF721759649A89721F8F"; // Crown token which can be minted for free
const paymasterAddress = "0x13D0D8550769f59aa241a41897D4859c87f7Dd46"; // Paymaster for Crown token

async function main() {
    const conf = require("../../solidity/storage/build/combined.json");
    const abi = conf.contracts["Storage.sol:Storage"].abi;
    const bytecode: string = conf.contracts["Storage.sol:Storage"].bin;

    const factory = new ContractFactory(abi, bytecode, wallet);
    const storage = (await factory.deploy()) as Contract;
    console.log(`Contract address: ${await storage.getAddress()}`);

    console.log(`Value: ${await storage.get()}`);

    const tx = await storage.set(Typed.uint256(200));
    await tx.wait();

    console.log(`Value: ${await storage.get()}`);

    const paymasterTx = await storage.set(Typed.uint256(500), {
        customData: {
            gasPerPubdata: utils.DEFAULT_GAS_PER_PUBDATA_LIMIT,
            paymasterParams: utils.getPaymasterParams(paymasterAddress, {
                type: "ApprovalBased",
                token: tokenAddress,
                minimalAllowance: 1,
                innerInput: new Uint8Array(),
            }),
        },
    });
    await paymasterTx.wait();
    console.log(`Value: ${await storage.get()}`);
}

main()
    .then()
    .catch((error) => {
        console.log(`Error: ${error}`);
    });
