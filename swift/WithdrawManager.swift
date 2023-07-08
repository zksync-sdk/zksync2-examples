//
//  WithdrawManager.swift
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

class WithdrawManager: BaseManager {
    func withdraw(callback: @escaping (() -> Void)) {
        let contract = zkSync.web3.contract(Web3.Utils.IEthToken)!
        
        let value = BigUInt(1_000_000_000_000)
        
        let inputs = [
            ABI.Element.InOut(name: "_l1Receiver", type: .address)
        ]
        
        let function = ABI.Element.Function(
            name: "withdraw",
            inputs: inputs,
            outputs: [],
            constant: false,
            payable: true
        )
        
        let withdrawFunction: ABI.Element = .function(function)
        
        let parameters: [AnyObject] = [
            EthereumAddress(signer.address)! as AnyObject,
        ]
        
        let calldata = withdrawFunction.encodeParameters(parameters)!
        
        var estimate = EthereumTransaction.createFunctionCallTransaction(
            from: EthereumAddress(signer.address)!,
            to: EthereumAddress.L2EthTokenAddress,
            gasPrice: BigUInt.zero,
            gasLimit: BigUInt.zero,
            value: value,
            data: calldata
        )
        
        estimate.envelope.parameters.chainID = signer.domain.chainId
        
        let fee = try! zkSync.zksEstimateFee(estimate).wait()
        
        estimate.parameters.EIP712Meta?.gasPerPubdata = BigUInt(160000)
        
        var transactionOptions = TransactionOptions.defaultOptions
        transactionOptions.type = .eip712
        transactionOptions.from = EthereumAddress(signer.address)!
        transactionOptions.to = estimate.to
        transactionOptions.maxPriorityFeePerGas = .manual(fee.maxPriorityFeePerGas)
        transactionOptions.maxFeePerGas = .manual(fee.maxFeePerGas)
        transactionOptions.value = value
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
        
        let result = try! contract.web3.eth.sendRawTransactionPromise(transaction).wait()
        
        let txHash = result.hash
        let index = 0
        
        guard let receipt = transactionReceiptProcessor.waitForTransactionReceipt(hash: txHash) else {
            fatalError("Transaction failed.")
        }
        
        assert(receipt.status == .ok)
        
        let l1ERC20Bridge = zkSync.web3.contract(
            Web3.Utils.IL1Bridge,
            at: EthereumAddress(signer.address)
        )!
        
        zkSync.zksMainContract { result in
            DispatchQueue.global().async {
                switch result {
                case .success(let address):
                    let zkSyncContract = self.ethereum.contract(
                        Web3.Utils.IZkSync,
                        at: EthereumAddress(address)
                    )!
                    
                    let defaultEthereumProvider = DefaultEthereumProvider(self.ethereum, l1ERC20Bridge: l1ERC20Bridge, zkSyncContract: zkSyncContract, gasProvider: DefaultGasProvider())
                    
                    let topic = "L1MessageSent(address,bytes32,bytes)"
                    let log = receipt.logs.filter({
                        if $0.address.address == ZkSyncAddresses.MessengerAddress && $0.topics.first == EIP712.keccak256(topic) {
                            return true
                        }
                        return false
                    })[index]
                    
                    let l2tol1log = receipt.l2ToL1Logs!.filter({
                        if $0.sender.address == ZkSyncAddresses.MessengerAddress {
                            return true
                        }
                        return false
                    })[index]
                    
                    self.zkSync.zksGetL2ToL1LogProof(txHash, logIndex: Int(l2tol1log.logIndex)) { result in
                        DispatchQueue.global().async {
                            switch result {
                            case .success(let proof):
                                let contract = self.zkSync.web3.contract(Web3.Utils.IL1Messenger)!
                                
                                let eventData = contract.parseEvent(log).eventData
                                let message = eventData?["_message"] as? Data ?? Data()
                                
                                _ = try! defaultEthereumProvider.finalizeEthWithdrawal(
                                    receipt.l1BatchNumber,
                                    l2MessageIndex: BigUInt(proof.id),
                                    l2TxNumberInBlock: receipt.l1BatchTxIndex,
                                    message: message,
                                    proof: proof.proof.compactMap({ Data(fromHex: $0) }),
                                    nonce: self.nonce
                                ).wait()
                                
                                callback()
                            case .failure(let error):
                                fatalError(error.localizedDescription)
                            }
                        }
                    }
                case .failure(let error):
                    print("Error:", error)
                }
            }
        }
    }
    
    func withdrawViaWallet(callback: @escaping (() -> Void)) {
        let amount = BigUInt(1_000_000_000_000)
        
        _ = try! wallet.withdraw(
            signer.address,
            amount: amount,
            token: Token.ETH
        ).wait()
        
        callback()
    }
}
