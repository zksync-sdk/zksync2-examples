import {Provider, types, utils, Wallet} from "zksync2-js";
import {ethers} from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const ethProvider = ethers.getDefaultProvider("sepolia");
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);


async function main() {
    console.log(`L2 balance before withdraw: ${await wallet.getBalance()}`);
    console.log(`L1 balance before withdraw: ${await wallet.getBalanceL1()}`);

    const tx = await wallet.withdraw({
        token: utils.ETH_ADDRESS,
        to: await wallet.getAddress(),
        amount: ethers.parseEther("0.00020"),
    });
    const receipt =  await tx.wait();
    console.log(`Tx: ${receipt.hash}`);

    // The duration for submitting a withdrawal transaction to L1 can last up to 24 hours. For additional information,
    // please refer to the documentation: https://era.zksync.io/docs/reference/troubleshooting/withdrawal-delay.html.
    // Once the withdrawal transaction is submitted on L1, it needs to be finalized.
    // To learn more about how to achieve this, please take a look at the 04_finalize_withdraw.ts script.
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})