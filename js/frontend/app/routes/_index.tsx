import type { MetaFunction } from "@remix-run/node";
import {BrowserProvider, Provider, types, utils } from "zksync-ethers";
import {ethers} from "ethers";

export const meta: MetaFunction = () => {
  return [
    { title: "zkSync Era ethers SDK" },
    { name: "description", content: "zkSync Era ethers SDK" },
  ];
};

const receiver = "0x81E9D85b65E9CC8618D85A1110e4b1DF63fA30d9";
const token = "0x6a4Fb925583F7D4dF82de62d98107468aE846FD1";
const approvalToken = "0x927488F48ffbc32112F1fF721759649A89721F8F"; // Crown token which can be minted for free
const paymaster = "0x13D0D8550769f59aa241a41897D4859c87f7Dd46"; // Paymaster for Crown token

async function transferETH() {
  // @ts-ignore
  const browserProvider = new BrowserProvider(window.ethereum);
  const signer = await browserProvider.getSigner();
  const provider = new Provider('https://sepolia.era.zksync.dev');

  console.log(`Account1 balance before transfer: ${await signer.getBalance()}`);
  console.log(`Account2 balance before transfer: ${await provider.getBalance(receiver)}`);

  const tx = await signer.transfer({
    to: receiver,
    amount: ethers.parseEther("0.01"),
  });
  const receipt = await tx.wait();
  console.log(`Tx: ${receipt.hash}`);

  console.log(`Account1 balance after transfer: ${await signer.getBalance()}`);
  console.log(`Account2 balance after transfer: ${await provider.getBalance(receiver)}`);
}

async function transferToken() {
  // @ts-ignore
  const browserProvider = new BrowserProvider(window.ethereum);
  const signer = await browserProvider.getSigner();
  const provider = new Provider('https://sepolia.era.zksync.dev');

  console.log(`Account1 balance before transfer: ${await signer.getBalance(token)}`);
  console.log(
      `Account2 balance before transfer: ${await provider.getBalance(receiver, "latest", token)}`,
  );

  const tx = await signer.transfer({
    token: token,
    to: receiver,
    amount: 5,
  });
  const receipt = await tx.wait();
  console.log(`Tx: ${receipt.hash}`);

  console.log(`Account1 balance after transfer: ${await signer.getBalance(token)}`);
  console.log(
      `Account2 balance after transfer: ${await provider.getBalance(receiver, "latest", token)}`,
  );
}

async function transferETHUsingApprovalPaymaster() {
  // @ts-ignore
  const browserProvider = new BrowserProvider(window.ethereum);
  const signer = await browserProvider.getSigner();
  const provider = new Provider('https://sepolia.era.zksync.dev');

  console.log(`Account1 balance before transfer: ${await signer.getBalance()}`);
  console.log(`Account2 balance before transfer: ${await provider.getBalance(receiver)}`);

  const tx = await signer.transfer({
    to: receiver,
    amount: ethers.parseEther("0.01"),
    paymasterParams: utils.getPaymasterParams(paymaster, {
      type: "ApprovalBased",
      token: approvalToken,
      minimalAllowance: 1,
      innerInput: new Uint8Array(),
    }),
  });
  const receipt = await tx.wait();
  console.log(`Tx: ${receipt.hash}`);

  console.log(`Account1 balance after transfer: ${await signer.getBalance()}`);
  console.log(`Account2 balance after transfer: ${await provider.getBalance(receiver)}`);
}

async function transferTokenUsingApprovalPaymaster() {
  // @ts-ignore
  const browserProvider = new BrowserProvider(window.ethereum);
  const signer = await browserProvider.getSigner();
  const provider = new Provider('https://sepolia.era.zksync.dev');

  console.log(`Account1 balance before transfer: ${await signer.getBalance(token)}`);
  console.log(
      `Account2 balance before transfer: ${await provider.getBalance(receiver, "latest", token)}`,
  );

  const tx = await signer.transfer({
    token: token,
    to: receiver,
    amount: 5,
    paymasterParams: utils.getPaymasterParams(paymaster, {
      type: "ApprovalBased",
      token: approvalToken,
      minimalAllowance: 1,
      innerInput: new Uint8Array(),
    }),
  });
  const receipt = await tx.wait();
  console.log(`Tx: ${receipt.hash}`);

  console.log(`Account1 balance after transfer: ${await signer.getBalance(token)}`);
  console.log(
      `Account2 balance after transfer: ${await provider.getBalance(receiver, "latest", token)}`,
  );
}

async function withdraw() {

}

export default function Index() {
  return (
      <div style={{fontFamily: "system-ui, sans-serif", lineHeight: "1.8"}}>
        <h1>Welcome to zksync-ethers SDK</h1>
        <button onClick={transferETH}>Transfer ETH</button>
        <button onClick={transferToken}>Transfer token</button>
        <button onClick={transferETHUsingApprovalPaymaster}>Transfer ETH using paymaster</button>
        <button onClick={transferTokenUsingApprovalPaymaster}>Transfer token using paymaster</button>
      </div>
  );
}
