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
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
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
        
        let contractTransaction = EthereumTransaction.create2ContractTransaction(from: EthereumAddress(signer.address)!, gasPrice: BigUInt.zero, gasLimit: BigUInt.zero, bytecode: bytecodeData, deps: [bytecodeData], calldata: encodedCallData, salt: Data(), chainId: signer.domain.chainId)
        
        let chainID = signer.domain.chainId
        
        var estimate = EthereumTransaction.createFunctionCallTransaction(from: EthereumAddress(signer.address)!, to: contractTransaction.to, gasPrice: BigUInt.zero, gasLimit: BigUInt.zero, data: contractTransaction.data)
        
        estimate.parameters.EIP712Meta?.factoryDeps = [bytecodeData]
        
        let fee = try! zkSync.estimateFee(estimate).wait()
        
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
        
        signTransaction(&transaction)
        
        guard let message = transaction.encode(for: .transaction) else {
            fatalError("Failed to encode transaction.")
        }
        
        let result = try! zkSync.web3.eth.sendRawTransactionPromise(message).wait()
        
        let receipt = transactionReceiptProcessor.waitForTransactionReceipt(hash: result.hash)
        
        assert(receipt?.status == .ok)
        
        callback()
    }
    
    func mintToken(tokenAddress: String, callback: (() -> Void)) {
        let contract = zkSync.web3.contract(Web3.Utils.IToken, at: EthereumAddress(tokenAddress)!)!
        
        let value = BigUInt(1_000)
        
        let parameters = [
            EthereumAddress(signer.address)! as AnyObject,
            value as AnyObject
        ] as [AnyObject]
        
        var transactionOptions = TransactionOptions.defaultOptions
        transactionOptions.from = EthereumAddress(signer.address)!
        
        guard let writeTransaction = contract.write(
            "mint",
            parameters: parameters,
            transactionOptions: transactionOptions
        ) else {
            return
        }
        
        var estimate = EthereumTransaction.createFunctionCallTransaction(
            from: EthereumAddress(signer.address)!,
            to: writeTransaction.transaction.to,
            gasPrice: BigUInt.zero,
            gasLimit: BigUInt.zero,
            data: writeTransaction.transaction.data
        )
        
        let fee = try! zkSync.estimateFee(estimate).wait()
        
        estimate.parameters.EIP712Meta?.gasPerPubdata = BigUInt(160000)
        
        transactionOptions = TransactionOptions.defaultOptions
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
        
        signTransaction(&transaction)
        
        _ = try! zkSync.web3.eth.sendRawTransactionPromise(transaction).wait()
        
        callback()
    }
    
    func getAllTokens(callback: @escaping (() -> Void)) {
        zkSync.confirmedTokens(0, limit: 255) { result in
            print("tokens:", result)
            
            callback()
        }
    }
    
    func tokenBalance(tokenAddress: String, callback: (() -> Void)) {
        let balance = try! walletL2.getBalance(Token(l1Address: "", l2Address: tokenAddress, symbol: "", decimals: 18)).wait()
        
        print("balance:", balance)
        
        callback()
    }
}
