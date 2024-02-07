import { Provider, types, Wallet, utils } from "zksync-ethers";
import { ethers } from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const ethProvider = ethers.getDefaultProvider("sepolia");
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);

const tokenAddress = "0x927488F48ffbc32112F1fF721759649A89721F8F"; // Crown token which can be minted for free
const paymasterAddress = "0x13D0D8550769f59aa241a41897D4859c87f7Dd46"; // Paymaster for Crown token

/*
This example demonstrates how to use a paymaster to facilitate fee payment with an ERC20 token.
The user initiates a transfer transaction that is configured to be paid with an ERC20 token through the paymaster.
During transaction execution, the paymaster receives the ERC20 token from the user and covers the transaction fee using ETH.
 */
async function main() {
    const receiver = "0x81E9D85b65E9CC8618D85A1110e4b1DF63fA30d9";

    console.log(`Account1 balance before transfer: ${await wallet.getBalance()}`);
    console.log(`Account2 balance before transfer: ${await provider.getBalance(receiver)}`);

    const tx = await wallet.transfer({
        to: receiver,
        amount: ethers.parseEther("0.01"),
        paymasterParamas: utils.getPaymasterParams(paymasterAddress, {
            type: "ApprovalBased",
            token: tokenAddress,
            minimalAllowance: 1,
            innerInput: new Uint8Array(),
        }),
    });
    const receipt = await tx.wait();
    console.log(`Tx: ${receipt.hash}`);

    console.log(`Account1 balance after transfer: ${await wallet.getBalance()}`);
    console.log(`Account2 balance after transfer: ${await provider.getBalance(receiver)}`);
}

main()
    .then()
    .catch((error) => {
        console.log(`Error: ${error}`);
    });
