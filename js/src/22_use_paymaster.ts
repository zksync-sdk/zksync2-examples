import { Provider, types, utils, Wallet } from "zksync-ethers";
import { ethers } from "ethers";

const provider = Provider.getDefaultProvider(types.Network.Sepolia);
const ethProvider = ethers.getDefaultProvider("sepolia");
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);

const tokenAddress = "0x927488F48ffbc32112F1fF721759649A89721F8F"; // Crown token which can be minted for free
const paymasterAddress = "0x13D0D8550769f59aa241a41897D4859c87f7Dd46"; // Paymaster for Crown token
const token = new ethers.Interface(
    require("../../solidity/custom_paymaster/token/build/Token.json").abi,
);

/*
This example demonstrates how to use a paymaster to facilitate fee payment with an ERC20 token.
The user initiates a mint transaction that is configured to be paid with an ERC20 token through the paymaster.
During transaction execution, the paymaster receives the ERC20 token from the user and covers the transaction fee using ETH.
 */
async function main() {
    console.log(`Paymaster balance before mint: ${await provider.getBalance(paymasterAddress)}`);
    console.log(`Wallet address before mint: ${await wallet.getBalance()}`);
    console.log(
        `Paymaster Crown tokens before mint: ${await provider.getBalance(
            paymasterAddress,
            "latest",
            tokenAddress,
        )}`,
    );
    console.log(`Wallet Crown tokens before mint: ${await wallet.getBalance(tokenAddress)}`);

    // Also mint some tokens to user account, so it can offer to pay fee with it
    const mintTx = await wallet.sendTransaction({
        to: tokenAddress,
        data: token.encodeFunctionData("mint", [await wallet.getAddress(), 10]),
    });
    await mintTx.wait();

    // Transfer some ETH to paymaster, so it can pay fee with ETH
    const faucetTx = await wallet.transfer({
        token: utils.ETH_ADDRESS,
        to: paymasterAddress,
        amount: ethers.parseEther("0.1"),
    });
    await faucetTx.wait();

    // In order to use paymaster, EIP712 transaction
    // need to be created with configured paymaster parameters
    const tx = await wallet.sendTransaction({
        type: utils.EIP712_TX_TYPE,
        to: tokenAddress,
        data: token.encodeFunctionData("mint", [await wallet.getAddress(), 3]),
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
    await tx.wait();

    console.log(`Paymaster balance after mint: ${await provider.getBalance(paymasterAddress)}`);
    console.log(`Wallet address after mint: ${await wallet.getBalance()}`);
    console.log(
        `Paymaster Crown tokens after mint: ${await provider.getBalance(
            paymasterAddress,
            "latest",
            tokenAddress,
        )}`,
    );
    console.log(`Wallet Crown tokens after mint: ${await wallet.getBalance(tokenAddress)}`);
}

main()
    .then()
    .catch((error) => {
        console.log(`Error: ${error}`);
    });
