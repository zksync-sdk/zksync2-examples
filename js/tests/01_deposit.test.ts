import * as chai from "chai";
import "./custom-matchers";
import { Provider, types, utils, Wallet } from "zksync-ethers";
import { ethers } from "ethers";
import * as fs from "fs";

const { expect } = chai;

describe("Deposit", () => {
    const PRIVATE_KEY = "0x7726827caac94a7f9e1b160f7ea819f172f7b6f9d2a97f992c38edeab82d4110";

    const provider = Provider.getDefaultProvider(types.Network.Localhost);
    const ethProvider = ethers.getDefaultProvider("http://localhost:8545");
    const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);

    it("should deposit ETH to L2 network", async () => {
        const amount = 7_000_000_000;
        const l2BalanceBeforeDeposit = await wallet.getBalance();
        const l1BalanceBeforeDeposit = await wallet.getBalanceL1();
        const tx = await wallet.deposit({
            token: utils.ETH_ADDRESS,
            to: await wallet.getAddress(),
            amount: amount,
            refundRecipient: await wallet.getAddress(),
        });
        const result = await tx.wait();
        const l2BalanceAfterDeposit = await wallet.getBalance();
        const l1BalanceAfterDeposit = await wallet.getBalanceL1();
        expect(result).not.to.be.null;
        expect(l2BalanceAfterDeposit - l2BalanceBeforeDeposit >= BigInt(amount)).to.be.true;
        expect(l1BalanceBeforeDeposit - l1BalanceAfterDeposit >= BigInt(amount)).to.be.true;
    }).timeout(10_000);

});
