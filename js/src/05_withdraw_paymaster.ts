import { Provider, types, utils, Wallet } from "zksync-ethers";
import { ethers } from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const ethProvider = ethers.getDefaultProvider("sepolia");
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);

const tokenAddress = "0x927488F48ffbc32112F1fF721759649A89721F8F"; // Crown token which can be minted for free
const paymasterAddress = "0x13D0D8550769f59aa241a41897D4859c87f7Dd46"; // Paymaster for Crown token

/*
This example demonstrates how to use a paymaster to facilitate fee payment with an ERC20 token.
The user initiates a withdrawal transaction that is configured to be paid with an ERC20 token through the paymaster.
During transaction execution, the paymaster receives the ERC20 token from the user and covers the transaction fee using ETH.
 */
async function main() {
    console.log(`L2 balance before withdraw: ${await wallet.getBalance()}`);
    console.log(`L1 balance before withdraw: ${await wallet.getBalanceL1()}`);

    const tx = await wallet.withdraw({
        token: utils.ETH_ADDRESS,
        to: await wallet.getAddress(),
        amount: ethers.parseEther("0.00020"),
        paymasterParamas: utils.getPaymasterParams(paymasterAddress, {
            type: "ApprovalBased",
            token: tokenAddress,
            minimalAllowance: 1,
            innerInput: new Uint8Array(),
        }),
    });
    const receipt = await tx.wait();
    console.log(`Tx: ${receipt.hash}`);

    // The duration for submitting a withdrawal transaction to L1 can last up to 24 hours. For additional information,
    // please refer to the documentation: https://era.zksync.io/docs/reference/troubleshooting/withdrawal-delay.html.
    // Once the withdrawal transaction is submitted on L1, it needs to be finalized.
    // To learn more about how to achieve this, please take a look at the 04_finalize_withdraw.ts script.
}

main()
    .then()
    .catch((error) => {
        console.log(`Error: ${error}`);
    });
