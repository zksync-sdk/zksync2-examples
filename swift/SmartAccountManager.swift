//
//  SmartAccountManager.swift
//  zkSync-Demo
//
//  Created by Bojan on 30.6.23..
//

import Foundation
import BigInt
#if canImport(web3swift)
import web3swift
import Web3Core
#else
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
#endif

//222class SmartAccountManager: BaseManager {
//    func deploySmartAccount(tokenAddress: String, callback: (() -> Void)) {
//        guard let path = Bundle.main.path(forResource: "Paymaster", ofType: "json") else { return }
//
//        let data = try! Data(contentsOf: URL(fileURLWithPath: path), options: .mappedIfSafe)
//        let jsonResult = try! JSONSerialization.jsonObject(with: data, options: .mutableLeaves)
//        guard let json = jsonResult as? [String: Any], let bytecode = json["bytecode"] as? String else { return }
//
//        let bytecodeData = Data(fromHex: bytecode)!
//
//        let inputs = [
//            ABI.Element.InOut(name: "erc20", type: .address),
//        ]
//
//        let function = ABI.Element.Function(
//            name: "",
//            inputs: inputs,
//            outputs: [],
//            constant: false,
//            payable: false)
//
//        let elementFunction: ABI.Element = .function(function)
//
//        let parameters: [AnyObject] = [
//            EthereumAddress(tokenAddress)! as AnyObject
//        ]
//
//        guard var encodedCallData = elementFunction.encodeParameters(parameters) else {
//            fatalError("Failed to encode function.")
//        }
//
//        // Removing signature prefix, which is first 4 bytes
//        for _ in 0..<4 {
//            encodedCallData = encodedCallData.dropFirst()
//        }
//
//        let estimate = EthereumTransaction.create2AccountTransaction(
//            from: EthereumAddress(signer.address)!,
//            gasPrice: BigUInt.zero,
//            gasLimit: BigUInt.zero,
//            bytecode: bytecodeData,
//            deps: [bytecodeData],
//            calldata: encodedCallData,
//            salt: Data(),
//            chainId: signer.domain.chainId
//        )
//
//        let chainID = signer.domain.chainId
//        let gasPrice = try! zkSync.web3.eth.getGasPrice()
//
//        var transactionOptions = TransactionOptions.defaultOptions
//        transactionOptions.gasPrice = .manual(BigUInt.zero)
//        transactionOptions.type = .eip712
//        transactionOptions.chainID = chainID
//        transactionOptions.nonce = .manual(nonce)
//        transactionOptions.to = estimate.to
//        transactionOptions.value = BigUInt.zero
//        transactionOptions.maxPriorityFeePerGas = .manual(BigUInt(100000000))
//        transactionOptions.maxFeePerGas = .manual(gasPrice)
//        transactionOptions.from = estimate.parameters.from
//
//        let estimateGas = try! zkSync.web3.eth.estimateGas(estimate, transactionOptions: transactionOptions)
//        transactionOptions.gasLimit = .manual(estimateGas)
//
//        var ethereumParameters = EthereumParameters(from: transactionOptions)
//        ethereumParameters.EIP712Meta = estimate.parameters.EIP712Meta
//        ethereumParameters.EIP712Meta?.factoryDeps = [bytecodeData]
//
//        var transaction = EthereumTransaction(
//            type: .eip712,
//            to: estimate.to,
//            nonce: nonce,
//            chainID: chainId,
//            data: estimate.data,
//            parameters: ethereumParameters
//        )
//
//        signTransaction(&transaction)
//
//        guard let message = transaction.encode(for: .transaction) else {
//            fatalError("Failed to encode transaction.")
//        }
//
//        let result = try! zkSync.web3.eth.sendRawTransactionPromise(message).wait()
//
//        let receipt = transactionReceiptProcessor.waitForTransactionReceipt(hash: result.hash)
//
//        assert(receipt?.status == .ok)
//
//        callback()
//    }
//}
