//
//  TokenManager.swift
//  zkSync-Demo
//
//  Created by Bojan on 23.6.23..
//

import Foundation
import BigInt
#if canImport(web3swift)
import web3swift
#else
import web3swift_zksync
#endif
#if canImport(ZkSync2)
import ZkSync2
#endif

class TokenManager: BaseManager {
    func deployToken(callback: (() -> Void)) {
        guard let path = Bundle.main.path(forResource: "Token", ofType: "json") else { return }
        
        let data = try! Data(contentsOf: URL(fileURLWithPath: path), options: .mappedIfSafe)
        let jsonResult = try! JSONSerialization.jsonObject(with: data, options: .mutableLeaves)
        guard let json = jsonResult as? [String: Any], let bytecode = json["bytecode"] as? String else { return }
        
        let bytecodeData = Data(fromHex: bytecode)!
        
        let inputs = [
            ABI.Element.InOut(name: "name_", type: .string),
            ABI.Element.InOut(name: "symbol_", type: .string),
            ABI.Element.InOut(name: "decimals_", type: .uint(bits: 8))
        ]
        
        let function = ABI.Element.Function(
            name: "",
            inputs: inputs,
            outputs: [],
            constant: false,
            payable: false)
        
        let elementFunction: ABI.Element = .function(function)
        
        let name = "USD Coin"
        let symbol = "USDC"
        let decimals: UInt8 = 18
        let parameters: [AnyObject] = [
            name as AnyObject,
            symbol as AnyObject,
            decimals as AnyObject
        ]
        
        guard var encodedCallData = elementFunction.encodeParameters(parameters) else {
            fatalError("Failed to encode function.")
        }
        
        // Removing signature prefix, which is first 4 bytes
        for _ in 0..<4 {
            encodedCallData = encodedCallData.dropFirst()
        }
        
        let nonce = try! zkSync.web3.eth.getTransactionCountPromise(address: EthereumAddress(signer.address)!, onBlock: ZkBlockParameterName.committed.rawValue).wait()
        
        let contractTransaction = EthereumTransaction.create2ContractTransaction(from: EthereumAddress(signer.address)!, gasPrice: BigUInt.zero, gasLimit: BigUInt.zero, bytecode: bytecodeData, deps: [bytecodeData], calldata: encodedCallData, salt: Data(), chainId: signer.domain.chainId)
        
        let chainID = signer.domain.chainId
        
        var estimate = EthereumTransaction.createFunctionCallTransaction(from: EthereumAddress(signer.address)!, to: contractTransaction.to, gasPrice: BigUInt.zero, gasLimit: BigUInt.zero, data: contractTransaction.data)
        
        estimate.parameters.EIP712Meta?.factoryDeps = [bytecodeData]
        
        let fee = try! (zkSync as! JsonRpc2_0ZkSync).zksEstimateFee(estimate).wait()
        
        var transactionOptions = TransactionOptions.defaultOptions
        transactionOptions.type = .eip712
        transactionOptions.chainID = chainID
        transactionOptions.nonce = .manual(nonce)
        transactionOptions.to = contractTransaction.to
        transactionOptions.value = contractTransaction.value
        transactionOptions.gasLimit = .manual(fee.gasLimit)
        transactionOptions.maxPriorityFeePerGas = .manual(fee.maxPriorityFeePerGas)
        transactionOptions.maxFeePerGas = .manual(fee.maxFeePerGas)
        transactionOptions.from = contractTransaction.parameters.from
        
        var ethereumParameters = EthereumParameters(from: transactionOptions)
        ethereumParameters.EIP712Meta = estimate.parameters.EIP712Meta
        ethereumParameters.EIP712Meta?.factoryDeps = [bytecodeData]
        
        var transaction = EthereumTransaction(
            type: .eip712,
            to: estimate.to,
            nonce: nonce,
            chainID: chainId,
            data: estimate.data,
            parameters: ethereumParameters
        )
        
        let signature = signer.signTypedData(signer.domain, typedData: transaction).addHexPrefix()
        
        let unmarshalledSignature = SECP256K1.unmarshalSignature(signatureData: Data(fromHex: signature)!)!
        transaction.envelope.r = BigUInt(fromHex: unmarshalledSignature.r.toHexString().addHexPrefix())!
        transaction.envelope.s = BigUInt(fromHex: unmarshalledSignature.s.toHexString().addHexPrefix())!
        transaction.envelope.v = BigUInt(unmarshalledSignature.v)
        
        guard let message = transaction.encode(for: .transaction) else {
            fatalError("Failed to encode transaction.")
        }
        
        let result = try! zkSync.web3.eth.sendRawTransactionPromise(message).wait()
        
        let receipt = transactionReceiptProcessor.waitForTransactionReceipt(hash: result.hash)
        
        assert(receipt?.status == .ok)
        
        callback()
    }
    
    func mintToken(callback: (() -> Void)) {
        let manager = KeystoreManager.init([credentials])
        zkSync.web3.eth.web3.addKeystoreManager(manager)
        self.eth.addKeystoreManager(manager)
        
        guard let path = Bundle.main.path(forResource: "Token", ofType: "json") else { return }
        
        let data = try! Data(contentsOf: URL(fileURLWithPath: path), options: .mappedIfSafe)
        let jsonResult = try! JSONSerialization.jsonObject(with: data, options: .mutableLeaves)
        guard let json = jsonResult as? [String: Any], let abi = json["abi"] as? [[String: Any]] else { return }
        
        guard let abiData = try? JSONSerialization.data(withJSONObject: abi, options: []) else { return }
        let abiString = String(data: abiData, encoding: .utf8)!
        
        let contractAddress = EthereumAddress("0xbc6b677377598a79fa1885e02df1894b05bc8b33")!
        
        let contract = zkSync.web3.contract(abiString, at: contractAddress)!
        
        let value = BigUInt(100)
        
        let parameters1 = [
            EthereumAddress("0x36615Cf349d7F6344891B1e7CA7C72883F5dc049")! as AnyObject,
            value as AnyObject
        ] as [AnyObject]
        
        var transactionOptions1 = TransactionOptions.defaultOptions
        transactionOptions1.from = EthereumAddress(signer.address)!
        
        guard let writeTransaction = contract.write(
            "mint",
            parameters: parameters1,
            transactionOptions: transactionOptions1
        ) else {
            return
        }
        
        var estimate = EthereumTransaction.createFunctionCallTransaction(from: EthereumAddress(signer.address)!, to: writeTransaction.transaction.to, gasPrice: BigUInt.zero, gasLimit: BigUInt.zero, data: writeTransaction.transaction.data)
        
        let nonce = try! zkSync.web3.eth.getTransactionCountPromise(address: EthereumAddress(signer.address)!, onBlock: ZkBlockParameterName.committed.rawValue).wait()
        
        let fee = try! (zkSync as! JsonRpc2_0ZkSync).zksEstimateFee(estimate).wait()
        
        estimate.parameters.EIP712Meta?.gasPerPubdata = BigUInt(160000)
        
        var transactionOptions = TransactionOptions.defaultOptions
        transactionOptions.type = .eip712
        transactionOptions.from = EthereumAddress(signer.address)!
        transactionOptions.to = estimate.to
        transactionOptions.maxPriorityFeePerGas = .manual(fee.maxPriorityFeePerGas)
        transactionOptions.maxFeePerGas = .manual(fee.maxFeePerGas)
        transactionOptions.nonce = .manual(nonce)
        transactionOptions.chainID = chainId
        
        let estimateGas = try! self.zkSync.web3.eth.estimateGas(estimate, transactionOptions: transactionOptions)
        transactionOptions.gasLimit = .manual(estimateGas)
        
        var ethereumParameters = EthereumParameters(from: transactionOptions)
        ethereumParameters.EIP712Meta = estimate.parameters.EIP712Meta
        
        var transaction = EthereumTransaction(
            type: .eip712,
            to: estimate.to,
            nonce: nonce,
            chainID: chainId,
            data: estimate.data,
            parameters: ethereumParameters
        )
        
        let signature = signer.signTypedData(signer.domain, typedData: transaction).addHexPrefix()
        
        let unmarshalledSignature = SECP256K1.unmarshalSignature(signatureData: Data(fromHex: signature)!)!
        transaction.envelope.r = BigUInt(fromHex: unmarshalledSignature.r.toHexString().addHexPrefix())!
        transaction.envelope.s = BigUInt(fromHex: unmarshalledSignature.s.toHexString().addHexPrefix())!
        transaction.envelope.v = BigUInt(unmarshalledSignature.v)
        
        _ = try! zkSync.web3.eth.sendRawTransactionPromise(transaction).wait()
        
        callback()
    }
    
    func tokenBalance(callback: (() -> Void)) {
        let balance = try! wallet.getBalance(Token(l1Address: "", l2Address: "0xbc6b677377598a79fa1885e02df1894b05bc8b33", symbol: "", decimals: 18)).wait()
        
        print("balance:", balance)
        
        callback()
    }
}
