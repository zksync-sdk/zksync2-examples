//
//  PaymasterManager.swift
//  zkSync-Demo
//
//  Created by Bojan on 20.6.23..
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

class PaymasterManager: BaseManager {
    func mintTokenUsingPaymaster(callback: (() -> Void)) {
        let contract = zkSync.web3.contract(Web3.Utils.IToken, at: EthereumAddress("0xbc6b677377598a79fa1885e02df1894b05bc8b33")!)!
        
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
        
        let fee = try! zkSync.zksEstimateFee(estimate).wait()
        
        estimate.parameters.EIP712Meta?.gasPerPubdata = BigUInt(160000)
        
        let paymasterAddress = EthereumAddress("0x49720d21525025522040f73da5b3992112bbec00")!
        let paymasterInput = Paymaster.encodeApprovalBased(
            EthereumAddress("0xbc6b677377598a79fa1885e02df1894b05bc8b33")!,
            minimalAllowance: BigUInt(1),
            paymasterInput: Data()
        )
        
        estimate.parameters.EIP712Meta?.paymasterParams = PaymasterParams(paymaster: paymasterAddress, paymasterInput: paymasterInput)
        
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
}
