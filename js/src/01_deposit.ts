import {Provider, types, utils, Wallet} from "zksync2-js";
import {ethers} from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const ethProvider = ethers.getDefaultProvider("sepolia");
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);


async function main() {
    console.log(`L2 balance before deposit: ${await wallet.getBalance()}`);
    console.log(`L1 balance before deposit: ${await wallet.getBalanceL1()}`);

    const tx = await wallet.deposit({
        token: utils.ETH_ADDRESS,
        to: await wallet.getAddress(),
        amount: ethers.parseEther("0.00020"),
        refundRecipient: await wallet.getAddress()
    });
    const receipt =  await tx.wait();
    console.log(`Tx: ${receipt.hash}`);

    console.log(`L2 balance after deposit: ${await wallet.getBalance()}`);
    console.log(`L1 balance after deposit: ${await wallet.getBalanceL1()}`);
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})