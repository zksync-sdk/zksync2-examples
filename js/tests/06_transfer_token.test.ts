import * as chai from "chai";
import "./custom-matchers";
import { Provider, types, utils, Wallet } from "zksync-ethers";
import { ethers } from "ethers";
import * as fs from "fs";

const { expect } = chai;

describe("Transfer token", () => {
    const PRIVATE_KEY = "0x7726827caac94a7f9e1b160f7ea819f172f7b6f9d2a97f992c38edeab82d4110";
    const RECEIVER = "0xa61464658AfeAf65CccaaFD3a512b69A83B77618";

    const provider = Provider.getDefaultProvider(types.Network.Localhost);
    const ethProvider = ethers.getDefaultProvider("http://localhost:8545");
    const wallet = new Wallet(PRIVATE_KEY, provider, ethProvider);

    const DAI = require("./token.json");

    it("should transfer DAI", async () => {
        const amount = 5;
        const l2DAI = await provider.l2TokenAddress(DAI.l1Address);
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

});
