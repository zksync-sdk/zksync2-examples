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
import Web3Core
#else
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
#endif

class TokenManager: BaseManager {
    func deployToken(callback: @escaping (() -> Void)) {
        Task {
            guard let path = Bundle.main.path(forResource: "Token", ofType: "json") else { return }
            
            let data = try! Data(contentsOf: URL(fileURLWithPath: path), options: .mappedIfSafe)
            let jsonResult = try! JSONSerialization.jsonObject(with: data, options: .mutableLeaves)
            guard let json = jsonResult as? [String: Any], let bytecode = json["bytecode"] as? String else { return }
            
            let bytecodeData = Data(hex: bytecode)
            
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
            
            let contractTransaction = CodableTransaction.create2ContractTransaction(from: EthereumAddress(signer.address)!, gasPrice: BigUInt.zero, gasLimit: BigUInt.zero, bytecode: bytecodeData, deps: [bytecodeData], calldata: encodedCallData, salt: Data(), chainId: self.signer.domain.chainId)
            
            let chainID = self.signer.domain.chainId
            
            var estimate = CodableTransaction.createFunctionCallTransaction(from: EthereumAddress(self.signer.address)!, to: contractTransaction.to, gasPrice: BigUInt.zero, gasLimit: BigUInt.zero, data: contractTransaction.data)
            
            estimate.eip712Meta?.factoryDeps = [bytecodeData]
            
            let fee = try! await self.zkSync.estimateFee(estimate)
            
            var transaction = await CodableTransaction(
                type: .eip712,
                to: estimate.to,
                nonce: self.getNonce(),
                chainID: self.signer.domain.chainId,
                data: estimate.data
            )
            transaction.value = contractTransaction.value
            transaction.gasLimit = fee.gasLimit
            transaction.maxPriorityFeePerGas = fee.maxPriorityFeePerGas
            transaction.maxFeePerGas = fee.maxFeePerGas
            transaction.from = contractTransaction.from
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

    func mintToken(tokenAddress: String, callback: @escaping (() -> Void)) {
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
            writeTransaction.transaction.from = EthereumAddress(self.signer.address)!
            
            var estimate = CodableTransaction.createFunctionCallTransaction(
                from: EthereumAddress(self.signer.address)!,
                to: writeTransaction.transaction.to,
                gasPrice: BigUInt.zero,
                gasLimit: BigUInt.zero,
                data: writeTransaction.transaction.data
            )
            
            let fee = try! await zkSync.estimateFee(estimate)
            
            estimate.eip712Meta?.gasPerPubdata = BigUInt(160000)
            
            let estimateGas = try! await self.zkSync.web3.eth.estimateGas(for: estimate)
            
            var transaction = await CodableTransaction(
                type: .eip712,
                to: estimate.to,
                nonce: self.getNonce(),
                chainID: self.signer.domain.chainId,
                data: estimate.data
            )
            transaction.gasLimit = estimateGas
            transaction.from = EthereumAddress(self.signer.address)!
            transaction.to = estimate.to
            transaction.maxPriorityFeePerGas = fee.maxPriorityFeePerGas
            transaction.maxFeePerGas = fee.maxFeePerGas
            transaction.eip712Meta = estimate.eip712Meta
            
            signTransaction(&transaction)
            
            _ = try! await zkSync.web3.eth.send(transaction)
            
            callback()
        }
    }

    func getAllTokens(callback: @escaping (() -> Void)) {
        zkSync.confirmedTokens(0, limit: 255) { result in
            print("tokens:", result)

            callback()
        }
    }

    func tokenBalance(tokenAddress: String, callback: @escaping (() -> Void)) {
        Task {
            let balance = try! await walletL2.getBalance(Token(l1Address: "", l2Address: tokenAddress, symbol: "", decimals: 18))
            
            print("balance:", balance)
            
            callback()
        }
    }
}
