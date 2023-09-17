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
import Web3Core
#else
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
#endif

class TransferManager: BaseManager {
    func transfer(toAddress: String, value: BigUInt, callback: @escaping (() -> Void)) {
        Task {
            var estimate = CodableTransaction.createFunctionCallTransaction(
                from: EthereumAddress(self.signer.address)!,
                to: EthereumAddress(toAddress)!,
                gasPrice: BigUInt.zero,
                gasLimit: BigUInt.zero,
                data: Data()
            )

            let fee = try! await zkSync.estimateFee(estimate)

            estimate.eip712Meta?.gasPerPubdata = fee.gasPerPubdataLimit

            var transaction = await CodableTransaction(
                type: .eip712,
                to: estimate.to,
                nonce: self.getNonce(),
                chainID: self.signer.domain.chainId,
                value: value,
                data: estimate.data
            )
            transaction.from = EthereumAddress(signer.address)!
            transaction.gasLimit = fee.gasLimit
            transaction.maxPriorityFeePerGas = fee.maxPriorityFeePerGas
            transaction.maxFeePerGas = fee.maxFeePerGas
            transaction.eip712Meta = estimate.eip712Meta

            signTransaction(&transaction)

            let result = try! await zkSync.web3.eth.send(transaction)
            
            let receipt = await transactionReceiptProcessor.waitForTransactionReceipt(hash: result.hash)
            
            assert(receipt?.status == .ok)
            
            callback()
        }
    }

    func transferViaWallet(toAddress: String, value: BigUInt, callback: @escaping (() -> Void)) {
        Task {
            _ = await walletL2.transfer(
                toAddress,
                amount: value
            )
            
            callback()
        }
    }
}
