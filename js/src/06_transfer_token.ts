import {Provider, types, Wallet} from "zksync-ethers";
import {ethers} from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const ethProvider = ethers.getDefaultProvider("sepolia");
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);


async function main() {
    const token = "0xCd9BDa1d0FC539043D4C80103bdF4f9cb108931B";
    const receiver = "0x81E9D85b65E9CC8618D85A1110e4b1DF63fA30d9";

    console.log(`Account1 balance before transfer: ${await wallet.getBalance()}`);
    console.log(`Account2 balance before transfer: ${await provider.getBalance(receiver)}`);

    const tx = await wallet.transfer({
        token: token,
        to: receiver,
        amount: 5,
    });
    const receipt =  await tx.wait();
    console.log(`Tx: ${receipt.hash}`);

    console.log(`Account1 balance after transfer: ${await wallet.getBalance()}`);
    console.log(`Account2 balance after transfer: ${await provider.getBalance(receiver)}`);
}

main().then().catch(error => {
    console.log(`Error: ${error}`);
})