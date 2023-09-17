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
import Web3Core
#else
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
#endif

class WithdrawManager: BaseManager {
    func withdraw(callback: @escaping (() -> Void)) {
        Task {
            let contract = self.zkSync.web3.contract(Web3.Utils.IEthToken)!

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

            var estimate = CodableTransaction.createFunctionCallTransaction(
                from: EthereumAddress(self.signer.address)!,
                to: EthereumAddress.L2EthTokenAddress,
                gasPrice: BigUInt.zero,
                gasLimit: BigUInt.zero,
                value: value,
                data: calldata
            )

            estimate.chainID = signer.domain.chainId

            let fee = try! await zkSync.estimateFee(estimate)

            estimate.eip712Meta?.gasPerPubdata = BigUInt(160000)

            var transaction = await CodableTransaction(
                type: .eip712,
                to: estimate.to,
                nonce: self.getNonce(),
                chainID: self.signer.domain.chainId,
                data: estimate.data
            )
            transaction.from = EthereumAddress(signer.address)!
            transaction.maxPriorityFeePerGas = fee.maxPriorityFeePerGas
            transaction.maxFeePerGas = fee.maxFeePerGas
            transaction.value = value
            transaction.eip712Meta = estimate.eip712Meta

            signTransaction(&transaction)

            let result = try! await contract.web3.eth.send(transaction)

            let txHash = result.hash
            let index = 0

            guard let receipt = await self.transactionReceiptProcessor.waitForTransactionReceipt(hash: txHash) else {
                fatalError("Transaction failed.")
            }

            assert(receipt.status == .ok)

            let l1ERC20Bridge = self.zkSync.web3.contract(
                Web3.Utils.IL1Bridge,
                at: EthereumAddress(self.signer.address)
            )!

            zkSync.mainContract { result in
                DispatchQueue.global().async {
                    switch result {
                    case .success(let address):
                        let zkSyncContract = self.web.contract(
                            Web3.Utils.IZkSync,
                            at: EthereumAddress(address)
                        )!

                        //444let defaultEthereumProvider = DefaultEthereumProvider(self.ethereum, l1ERC20Bridge: l1ERC20Bridge, zkSyncContract: zkSyncContract, gasProvider: DefaultGasProvider())

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

                        self.zkSync.getL2ToL1LogProof(txHash, logIndex: Int(l2tol1log.logIndex)) { result in
                            DispatchQueue.global().async {
                                switch result {
                                case .success(let proof):
                                    let contract = self.zkSync.web3.contract(Web3.Utils.IL1Messenger)!

                                    let eventData = contract.contract.parseEvent(log).eventData
                                    let message = eventData?["_message"] as? Data ?? Data()

//444                                    Task {
//                                        _ = try! await defaultEthereumProvider.finalizeEthWithdrawal(
//                                            receipt.l1BatchNumber,
//                                            l2MessageIndex: BigUInt(proof.id),
//                                            l2TxNumberInBlock: receipt.l1BatchTxIndex,
//                                            message: message,
//                                            proof: proof.proof.compactMap({ Data(hex: $0) }),
//                                            nonce: self.getNonce()
//                                        )
//
//                                        callback()
//                                    }
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
    }

    func withdrawViaWallet(callback: @escaping (() -> Void)) {
        Task {
            let amount = BigUInt(1_000_000_000_000)
            
            _ = await walletL2.withdraw(
                signer.address,
                amount: amount,
                token: Token.ETH
            )
            
            callback()
        }
    }
}
