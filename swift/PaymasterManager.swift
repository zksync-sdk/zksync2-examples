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
import Web3Core
#else
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
#endif

class PaymasterManager: BaseManager {
    func mintTokenUsingPaymaster(tokenAddress: String, paymasterAddress: String, callback: @escaping (() -> Void)) {
        Task {
            let contract = self.zkSync.web3.contract(Web3.Utils.IToken, at: EthereumAddress(tokenAddress)!)!
            
            let value = BigUInt(1_000)
            
            let parameters = [
                EthereumAddress(self.signer.address)! as AnyObject,
                value as AnyObject
            ] as [AnyObject]
            
            guard let writeTransaction = contract.createWriteOperation(
                "mint",
                parameters: parameters
            ) else {
                return
            }
            
            writeTransaction.transaction.from = EthereumAddress(signer.address)!
            
            var estimate = CodableTransaction.createFunctionCallTransaction(
                from: EthereumAddress(self.signer.address)!,
                to: writeTransaction.transaction.to,
                gasPrice: BigUInt.zero,
                gasLimit: BigUInt.zero,
                data: writeTransaction.transaction.data
            )
            
            //444let fee = try! self.zkSync.estimateFee(estimate).wait()
            
            //444estimate.eip712Meta?.gasPerPubdata = BigUInt(160000)
            
            let paymasterAddress = EthereumAddress(paymasterAddress)!
            let paymasterInput = Paymaster.encodeApprovalBased(
                EthereumAddress(tokenAddress)!,
                minimalAllowance: BigUInt(1),
                paymasterInput: Data()
            )
            
            //444estimate.eip712Meta?.paymasterParams = PaymasterParams(paymaster: paymasterAddress, paymasterInput: paymasterInput)
            
            let estimateGas = try! await self.zkSync.web3.eth.estimateGas(for: estimate)
            
            var transaction = await CodableTransaction(
                //444type: .eip712,
                to: estimate.to,
                nonce: self.getNonce(),
                chainID: self.signer.domain.chainId,
                data: estimate.data
            )
            transaction.from = EthereumAddress(signer.address)!
            transaction.gasLimit = estimateGas
//444            transaction.maxPriorityFeePerGas = fee.maxPriorityFeePerGas
//            transaction.maxFeePerGas = fee.maxFeePerGas
            //444transaction.eip712Meta = estimate.eip712Meta
            
            signTransaction(&transaction)
            
            _ = try! await zkSync.web3.eth.send(transaction)
            
            callback()
        }
    }
}
