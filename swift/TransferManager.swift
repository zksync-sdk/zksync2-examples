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
import web3swift_zksync
#endif
#if canImport(ZkSync2)
import ZkSync2
#endif

class TransferManager: BaseManager {
    func transfer(callback: (() -> Void)) {
        let value = BigUInt(1_000_000_000_000_000_000)
        
        var estimate = EthereumTransaction.createFunctionCallTransaction(
            from: EthereumAddress(signer.address)!,
            to: EthereumAddress("0xa61464658AfeAf65CccaaFD3a512b69A83B77618")!,
            gasPrice: BigUInt.zero,
            gasLimit: BigUInt.zero,
            data: Data()
        )
        
        let fee = try! zkSync.zksEstimateFee(estimate).wait()
        
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
    
    func transferViaWallet(callback: (() -> Void)) {
        let amount = BigUInt(1_000_000_000_000_000_000)
        
        _ = try! wallet.transfer(
            "0xa61464658AfeAf65CccaaFD3a512b69A83B77618",
            amount: amount
        ).wait()
        
        callback()
    }
}
