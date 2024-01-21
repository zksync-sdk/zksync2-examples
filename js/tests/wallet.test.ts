import * as chai from "chai";
import "./custom-matchers";
import { Provider, types, utils, Wallet } from "zksync-ethers";
import { ethers } from "ethers";

const { expect } = chai;

describe("Wallet", () => {
    const PRIVATE_KEY = "0x7726827caac94a7f9e1b160f7ea819f172f7b6f9d2a97f992c38edeab82d4110";
    const RECEIVER = "0xa61464658AfeAf65CccaaFD3a512b69A83B77618";

    const provider = Provider.getDefaultProvider(types.Network.Localhost);
    const ethProvider = ethers.getDefaultProvider("http://localhost:8545");
    const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);

    const l1Tokens = require("./tokens.json");
    const l1DAI = l1Tokens[0].address;

    it("should deposit ETH to the L2 network", async () => {
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

    it("should transfer ETH", async () => {
        const amount = 7_000_000_000;
        const balanceBeforeTransfer = await provider.getBalance(RECEIVER);
        const tx = await wallet.transfer({
            token: utils.ETH_ADDRESS,
            to: RECEIVER,
            amount: amount,
        });
        const result = await tx.wait();
        const balanceAfterTransfer = await provider.getBalance(RECEIVER);
        expect(result).not.to.be.null;
        expect(balanceAfterTransfer - balanceBeforeTransfer).to.be.equal(BigInt(amount));
    }).timeout(25_000);

    it("should withdraw ETH to the L1 network", async () => {
        const amount = 7_000_000_000;
        const l2BalanceBeforeWithdrawal = await wallet.getBalance();
        const withdrawTx = await wallet.withdraw({
            token: utils.ETH_ADDRESS,
            to: await wallet.getAddress(),
            amount: amount,
        });
        await withdrawTx.waitFinalize();
        expect(await wallet.isWithdrawalFinalized(withdrawTx.hash)).to.be.false;

        const finalizeWithdrawTx = await wallet.finalizeWithdrawal(withdrawTx.hash);
        const result = await finalizeWithdrawTx.wait();
        const l2BalanceAfterWithdrawal = await wallet.getBalance();
        expect(result).not.to.be.null;
        expect(l2BalanceBeforeWithdrawal - l2BalanceAfterWithdrawal >= BigInt(amount)).to.be.true;
    }).timeout(25_000);

    it("should deposit DAI to the L2 network", async () => {
        const amount = 5;
        const l2DAI = await provider.l2TokenAddress(l1DAI);
        const l2BalanceBeforeDeposit = await wallet.getBalance(l2DAI);
        const l1BalanceBeforeDeposit = await wallet.getBalanceL1(l1DAI);
        const tx = await wallet.deposit({
            token: l1DAI,
            to: await wallet.getAddress(),
            amount: amount,
            approveERC20: true,
            refundRecipient: await wallet.getAddress(),
        });
        const result = await tx.wait();
        const l2BalanceAfterDeposit = await wallet.getBalance(l2DAI);
        const l1BalanceAfterDeposit = await wallet.getBalanceL1(l1DAI);
        expect(result).not.to.be.null;
        expect(l2BalanceAfterDeposit - l2BalanceBeforeDeposit === BigInt(amount)).to.be.true;
        expect(l1BalanceBeforeDeposit - l1BalanceAfterDeposit === BigInt(amount)).to.be.true;
    }).timeout(25_000);

    it("should transfer DAI", async () => {
        const amount = 5;
        const l2DAI = await provider.l2TokenAddress(l1DAI);
        const balanceBeforeTransfer = await provider.getBalance(RECEIVER, "latest", l2DAI);
        const tx = await wallet.transfer({
            token: l2DAI,
            to: RECEIVER,
            amount: amount,
        });
        const result = await tx.wait();
        const balanceAfterTransfer = await provider.getBalance(RECEIVER, "latest", l2DAI);
        expect(result).not.to.be.null;
        expect(balanceAfterTransfer - balanceBeforeTransfer).to.be.equal(BigInt(amount));
    }).timeout(25_000);

    it("should withdraw DAI to the L1 network", async () => {
        const amount = 5;
        const l2DAI = await provider.l2TokenAddress(l1DAI);
        const l2BalanceBeforeWithdrawal = await wallet.getBalance(l2DAI);
        const l1BalanceBeforeWithdrawal = await wallet.getBalanceL1(l1DAI);
        const withdrawTx = await wallet.withdraw({
            token: l2DAI,
            to: await wallet.getAddress(),
            amount: amount,
        });
        await withdrawTx.waitFinalize();
        expect(await wallet.isWithdrawalFinalized(withdrawTx.hash)).to.be.false;

        const finalizeWithdrawTx = await wallet.finalizeWithdrawal(withdrawTx.hash);
        const result = await finalizeWithdrawTx.wait();
        const l2BalanceAfterWithdrawal = await wallet.getBalance(l2DAI);
        const l1BalanceAfterWithdrawal = await wallet.getBalanceL1(l1DAI);
        expect(result).not.to.be.null;
        expect(l2BalanceBeforeWithdrawal - l2BalanceAfterWithdrawal == BigInt(amount)).to.be.true;
        expect(l1BalanceAfterWithdrawal - l1BalanceBeforeWithdrawal == BigInt(amount)).to.be.true;
    }).timeout(25_000);
});
