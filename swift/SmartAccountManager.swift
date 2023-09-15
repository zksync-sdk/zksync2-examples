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

class SmartAccountManager: BaseManager {
    func deploySmartAccount(tokenAddress: String, callback: @escaping (() -> Void)) {
        Task {
            guard let path = Bundle.main.path(forResource: "Paymaster", ofType: "json") else { return }
            
            let data = try! Data(contentsOf: URL(fileURLWithPath: path), options: .mappedIfSafe)
            let jsonResult = try! JSONSerialization.jsonObject(with: data, options: .mutableLeaves)
            guard let json = jsonResult as? [String: Any], let bytecode = json["bytecode"] as? String else { return }
            
            let bytecodeData = Data(hex: bytecode)
            
            let inputs = [
                ABI.Element.InOut(name: "erc20", type: .address),
            ]
            
            let function = ABI.Element.Function(
                name: "",
                inputs: inputs,
                outputs: [],
                constant: false,
                payable: false)
            
            let elementFunction: ABI.Element = .function(function)
            
            let parameters: [AnyObject] = [
                EthereumAddress(tokenAddress)! as AnyObject
            ]
            
            guard var encodedCallData = elementFunction.encodeParameters(parameters) else {
                fatalError("Failed to encode function.")
            }
            
            // Removing signature prefix, which is first 4 bytes
            for _ in 0..<4 {
                encodedCallData = encodedCallData.dropFirst()
            }
            
            let estimate = CodableTransaction.create2AccountTransaction(
                from: EthereumAddress(self.signer.address)!,
                gasPrice: BigUInt.zero,
                gasLimit: BigUInt.zero,
                bytecode: bytecodeData,
                deps: [bytecodeData],
                calldata: encodedCallData,
                salt: Data(),
                chainId: self.signer.domain.chainId
            )
            
            let gasPrice = try! await zkSync.web3.eth.gasPrice()
            
            let estimateGas = try! await zkSync.web3.eth.estimateGas(for: estimate)
            
            var transaction = await CodableTransaction(
                type: .eip712,
                to: estimate.to,
                nonce: self.getNonce(),
                chainID: self.signer.domain.chainId,
                data: estimate.data
            )
            transaction.value = BigUInt.zero
            transaction.gasPrice = BigUInt.zero
            transaction.maxPriorityFeePerGas = BigUInt(100000000)
            transaction.maxFeePerGas = gasPrice
            transaction.gasLimit = estimateGas
            transaction.from = estimate.from
            transaction.eip712Meta = estimate.eip712Meta
            transaction.eip712Meta?.factoryDeps = [bytecodeData]
            
            signTransaction(&transaction)
            
            guard let message = transaction.encode(for: .transaction) else {
                fatalError("Failed to encode transaction.")
            }
            
            let result = try! await zkSync.web3.eth.send(raw: message)
            
            let receipt = await transactionReceiptProcessor.waitForTransactionReceipt(hash: result.hash)
            
            assert(receipt?.status == .ok)
            
            callback()
        }
    }
}
