//
//  TransferManager.swift
//  zkSync-Demo
//
//  Created by Bojan on 14.5.23..
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

class TransferManager: BaseManager {
    func transfer(toAddress: String, value: BigUInt, callback: (() -> Void)) {
        var estimate = EthereumTransaction.createFunctionCallTransaction(
            from: EthereumAddress(signer.address)!,
            to: EthereumAddress(toAddress)!,
            gasPrice: BigUInt.zero,
            gasLimit: BigUInt.zero,
            data: Data()
        )
        
        let fee = try! zkSync.estimateFee(estimate).wait()
        
        estimate.parameters.EIP712Meta?.gasPerPubdata = fee.gasPerPubdataLimit
        
        var transactionOptions = TransactionOptions.defaultOptions
        transactionOptions.type = .eip712
        transactionOptions.from = EthereumAddress(signer.address)!
        transactionOptions.to = estimate.to
        transactionOptions.gasLimit = .manual(fee.gasLimit)
        transactionOptions.maxPriorityFeePerGas = .manual(fee.maxPriorityFeePerGas)
        transactionOptions.maxFeePerGas = .manual(fee.maxFeePerGas)
        transactionOptions.value = value
        transactionOptions.nonce = .manual(nonce)
        transactionOptions.chainID = chainId
        
        var ethereumParameters = EthereumParameters(from: transactionOptions)
        ethereumParameters.EIP712Meta = estimate.parameters.EIP712Meta
        
        var transaction = EthereumTransaction(
            type: .eip712,
            to: estimate.to,
            nonce: nonce,
            chainID: chainId,
            value: value,
            data: estimate.data,
            parameters: ethereumParameters
        )
        
        signTransaction(&transaction)
        
        let result = try! zkSync.web3.eth.sendRawTransactionPromise(transaction).wait()
        
        let receipt = transactionReceiptProcessor.waitForTransactionReceipt(hash: result.hash)
        
        assert(receipt?.status == .ok)
        
        callback()
    }
    
    func transferViaWallet(toAddress: String, value: BigUInt, callback: (() -> Void)) {
        _ = try! walletL2.transfer(
            toAddress,
            amount: value
        ).wait()
        
        callback()
    }
}
