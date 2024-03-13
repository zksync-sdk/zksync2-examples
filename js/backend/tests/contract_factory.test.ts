import * as chai from "chai";
import "./custom-matchers";
import { Provider, types, ContractFactory, Wallet, Contract } from "zksync-ethers";
import { ethers, Typed } from "ethers";

const { expect } = chai;

describe("ContractFactory", () => {
    const PRIVATE_KEY = "0x7726827caac94a7f9e1b160f7ea819f172f7b6f9d2a97f992c38edeab82d4110";

    const provider = Provider.getDefaultProvider(types.Network.Localhost);
    const wallet = new Wallet(PRIVATE_KEY, provider);

    const l1Tokens = require("./tokens.json");
    const l1DAI = l1Tokens[0].address;

    it("should deploy a contract without constructor using CREATE opcode", async () => {
        const conf = require("../../../solidity/storage/build/combined.json");
        const abi = conf.contracts["Storage.sol:Storage"].abi;
        const bytecode: string = conf.contracts["Storage.sol:Storage"].bin;

        const factory = new ContractFactory(abi, bytecode, wallet);
        const storage = (await factory.deploy()) as Contract;

        const code = await provider.getCode(await storage.getAddress());
        expect(code).not.to.be.null;

        let value: bigint = await storage.get();
        expect(value == BigInt(0)).to.be.true;

        const tx = await storage.set(Typed.uint256(200));
        await tx.wait();

        value = await storage.get();
        expect(value == BigInt(200)).to.be.true;
    }).timeout(10_000);

    it("should deploy a contract with a constructor using CREATE opcode", async () => {
        const conf = require("../../../solidity/incrementer/build/combined.json");
        const abi = conf.contracts["Incrementer.sol:Incrementer"].abi;
        const bytecode: string = conf.contracts["Incrementer.sol:Incrementer"].bin;

        const factory = new ContractFactory(abi, bytecode, wallet);
        const incrementer = (await factory.deploy(2)) as Contract;

        const code = await provider.getCode(await incrementer.getAddress());
        expect(code).not.to.be.null;

        let value: bigint = await incrementer.get();
        expect(value == BigInt(0)).to.be.true;

        const tx = await incrementer.increment();
        await tx.wait();

        value = await incrementer.get();
        expect(value == BigInt(2)).to.be.true;
    }).timeout(10_000);

    it("should deploy a contract with dependencies using CREATE opcode", async () => {
        const conf = require("../../../solidity/demo/build/combined.json");
        const abi = conf.contracts["Demo.sol:Demo"].abi;
        const bytecode: string = conf.contracts["Demo.sol:Demo"].bin;

        const factory = new ContractFactory(abi, bytecode, wallet);
        const demo = (await factory.deploy({
            customData: { factoryDeps: [conf.contracts["Foo.sol:Foo"].bin] },
        })) as Contract;

        const code = await provider.getCode(await demo.getAddress());
        expect(code).not.to.be.null;

        const value: string = await demo.getFooName();
        expect(value).to.be.equal("Foo");
    }).timeout(10_000);

    it("should deploy a token using CREATE opcode", async () => {
        const conf = require("../../../solidity/custom_paymaster/token/build/Token.json");
        const abi = conf.abi;
        const bytecode: string = conf.bytecode;

        const factory = new ContractFactory(abi, bytecode, wallet);
        const token = (await factory.deploy("Crown", "Crown", 18)) as Contract;
        const tokenAddress = await token.getAddress();

        const code = await provider.getCode(await tokenAddress);
        expect(code).not.to.be.null;
    }).timeout(10_000);

    it("should deploy an account using CREATE opcode", async () => {
        const conf = require("../../../solidity/custom_paymaster/paymaster/build/Paymaster.json");
        const abi = conf.abi;
        const bytecode: string = conf.bytecode;

        const accountFactory = new ContractFactory(abi, bytecode, wallet, "createAccount");
        const paymasterContract = await accountFactory.deploy(await provider.l2TokenAddress(l1DAI));

        const code = await provider.getCode(await paymasterContract.getAddress());
        expect(code).not.to.be.null;
    }).timeout(10_000);

    it("should deploy a contract without constructor using CREATE2 opcode", async () => {
        const conf = require("../../../solidity/storage/build/combined.json");
        const abi = conf.contracts["Storage.sol:Storage"].abi;
        const bytecode: string = conf.contracts["Storage.sol:Storage"].bin;

        const factory = new ContractFactory(abi, bytecode, wallet, "create2");
        const storage = (await factory.deploy({
            customData: { salt: ethers.hexlify(ethers.randomBytes(32)) },
        })) as Contract;

        const code = await provider.getCode(await storage.getAddress());
        expect(code).not.to.be.null;

        let value: bigint = await storage.get();
        expect(value == BigInt(0)).to.be.true;

        const tx = await storage.set(Typed.uint256(200));
        await tx.wait();

        value = await storage.get();
        expect(value == BigInt(200)).to.be.true;
    }).timeout(10_000);

    it("should deploy a contract with a constructor using CREATE2 opcode", async () => {
        const conf = require("../../../solidity/incrementer/build/combined.json");
        const abi = conf.contracts["Incrementer.sol:Incrementer"].abi;
        const bytecode: string = conf.contracts["Incrementer.sol:Incrementer"].bin;

        const factory = new ContractFactory(abi, bytecode, wallet, "create2");
        const incrementer = (await factory.deploy(2, {
            customData: { salt: ethers.hexlify(ethers.randomBytes(32)) },
        })) as Contract;

        const code = await provider.getCode(await incrementer.getAddress());
        expect(code).not.to.be.null;

        let value: bigint = await incrementer.get();
        expect(value == BigInt(0)).to.be.true;

        const tx = await incrementer.increment();
        await tx.wait();

        value = await incrementer.get();
        expect(value == BigInt(2)).to.be.true;
    }).timeout(10_000);

    it("should deploy a contract with dependencies using CREATE2 opcode", async () => {
        const conf = require("../../../solidity/demo/build/combined.json");
        const abi = conf.contracts["Demo.sol:Demo"].abi;
        const bytecode: string = conf.contracts["Demo.sol:Demo"].bin;

        const factory = new ContractFactory(abi, bytecode, wallet, "create2");
        const demo = (await factory.deploy({
            customData: {
                salt: ethers.hexlify(ethers.randomBytes(32)),
                factoryDeps: [conf.contracts["Foo.sol:Foo"].bin],
            },
        })) as Contract;

        const code = await provider.getCode(await demo.getAddress());
        expect(code).not.to.be.null;

        const value: string = await demo.getFooName();
        expect(value).to.be.equal("Foo");
    }).timeout(10_000);

    it("should deploy a token using CREATE2 opcode", async () => {
        const conf = require("../../../solidity/custom_paymaster/token/build/Token.json");
        const abi = conf.abi;
        const bytecode: string = conf.bytecode;

        const factory = new ContractFactory(abi, bytecode, wallet, "create2");
        const token = (await factory.deploy("Crown", "Crown", 18, {
            customData: { salt: ethers.hexlify(ethers.randomBytes(32)) },
        })) as Contract;
        const tokenAddress = await token.getAddress();

        const code = await provider.getCode(await tokenAddress);
        expect(code).not.to.be.null;
    }).timeout(10_000);

    it("should deploy an account using CREATE2 opcode", async () => {
        const conf = require("../../../solidity/custom_paymaster/paymaster/build/Paymaster.json");
        const abi = conf.abi;
        const bytecode: string = conf.bytecode;

        const factory = new ContractFactory(abi, bytecode, wallet, "create2Account");
        const account = await factory.deploy(await provider.l2TokenAddress(l1DAI), {
            customData: { salt: ethers.hexlify(ethers.randomBytes(32)) },
        });

        const code = await provider.getCode(await account.getAddress());
        expect(code).not.to.be.null;
    }).timeout(10_000);
});
