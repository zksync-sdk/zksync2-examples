import {Provider, types, Wallet} from "zksync2-js";
import {ethers} from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const ethProvider = ethers.getDefaultProvider("sepolia");
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);

const WITHDRAW_TX = process.env.WITHDRAW_TX;

async function main() {

    // Perform the finalize withdrawal if it has not been already finalized.
    // On the testnet, withdrawals are automatically finalized. For additional information, please refer
    // to the documentation: https://era.zksync.io/docs/reference/concepts/bridging-asset.html#withdrawals-to-l1.
    // There is no need to execute FinalizeWithdraw, otherwise, an error with code `jj` would occur.
    if (! (await wallet.isWithdrawalFinalized(WITHDRAW_TX))) {
        const finalizeWithdrawTx = await wallet.finalizeWithdrawal(WITHDRAW_TX);
        const receipt = await finalizeWithdrawTx.wait();
        console.log(`Tx: ${receipt.hash}`);

        console.log(`L2 balance after withdraw: ${await wallet.getBalance()}`);
        console.log(`L1 balance after withdraw: ${await wallet.getBalanceL1()}`);
    }
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})