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
        let manager = KeystoreManager.init([credentials])
        zkSync.web3.eth.web3.addKeystoreManager(manager)
        self.eth.addKeystoreManager(manager)
        
        guard let path = Bundle.main.path(forResource: "IEthToken", ofType: "json") else { return }
        
        let data = try! Data(contentsOf: URL(fileURLWithPath: path), options: .mappedIfSafe)
        let jsonResult = try! JSONSerialization.jsonObject(with: data, options: .mutableLeaves)
        guard let json = jsonResult as? [String: Any], let abi = json["abi"] as? [[String: Any]] else { return }
        
        guard let abiData = try? JSONSerialization.data(withJSONObject: abi, options: []) else { return }
        let abiString = String(data: abiData, encoding: .utf8)!
        
        let contract = zkSync.web3.contract(abiString)!
        
        let value = BigUInt(1000000000000000000)
        
        let inputs = [
            ABI.Element.InOut(name: "_l1Receiver", type: .address)
        ]
        
        let function = ABI.Element.Function(name: "withdraw",
                                            inputs: inputs,
                                            outputs: [],
                                            constant: false,
                                            payable: true)
        
        let withdrawFunction: ABI.Element = .function(function)
        
        let parameters: [AnyObject] = [
            EthereumAddress("0x000000000000000000000000000000000000800a")! as AnyObject,
        ]
        
        let calldata = withdrawFunction.encodeParameters(parameters)!
        
        var estimate = EthereumTransaction.createFunctionCallTransaction(from: EthereumAddress(signer.address)!, to: EthereumAddress.L2EthTokenAddress, gasPrice: BigUInt.zero, gasLimit: BigUInt.zero, value: value, data: calldata)
        
        estimate.envelope.parameters.chainID = signer.domain.chainId
        
        let nonce = try! zkSync.web3.eth.getTransactionCountPromise(address: EthereumAddress(signer.address)!, onBlock: ZkBlockParameterName.committed.rawValue).wait()
        
        let fee = try! (zkSync as! JsonRpc2_0ZkSync).zksEstimateFee(estimate).wait()
        
        estimate.parameters.EIP712Meta?.gasPerPubdata = BigUInt(160000)
        
        var transactionOptions = TransactionOptions.defaultOptions
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
        
        var transaction = EthereumTransaction(type: .eip712,
                                              to: estimate.to,
                                              nonce: nonce,
                                              chainID: chainId,
                                              data: estimate.data,
                                              parameters: ethereumParameters)
        
        let signature = signer.signTypedData(signer.domain, typedData: transaction).addHexPrefix()
        
        let unmarshalledSignature = SECP256K1.unmarshalSignature(signatureData: Data(fromHex: signature)!)!
        transaction.envelope.r = BigUInt(fromHex: unmarshalledSignature.r.toHexString().addHexPrefix())!
        transaction.envelope.s = BigUInt(fromHex: unmarshalledSignature.s.toHexString().addHexPrefix())!
        transaction.envelope.v = BigUInt(unmarshalledSignature.v)
        
        let result = try! contract.web3.eth.sendRawTransactionPromise(transaction).wait()
        
        self.finalizeWithdraw(txHash: result.hash, index: 0) {
            callback()
        }
    }
    
    func finalizeWithdraw(txHash: String, index: Int, callback: @escaping (() -> Void)) {
        guard let receipt = transactionReceiptProcessor.waitForTransactionReceipt(hash: txHash) else {
            fatalError("Transaction failed.")
        }
        
        assert(receipt.status == .ok)
        
        let l1ERC20Bridge = zkSync.web3.contract(
            Web3.Utils.IL1Bridge,
            at: EthereumAddress("0x36615cf349d7f6344891b1e7ca7c72883f5dc049")
        )!
        
        zkSync.zksMainContract { result in
            DispatchQueue.global().async {
                switch result {
                case .success(let address):
                    let zkSyncContract = self.eth.contract(
                        Web3.Utils.IZkSync,
                        at: EthereumAddress(address)
                    )!
                    
                    let defaultEthereumProvider = DefaultEthereumProvider(self.eth, l1ERC20Bridge: l1ERC20Bridge, zkSyncContract: zkSyncContract, gasProvider: DefaultGasProvider())
                    
                    let topic = "L1MessageSent(address,bytes32,bytes)"
                    let log = receipt.logs.filter({
                        if $0.address.address == ZkSyncAddresses.MessengerAddress && $0.topics.first == EIP712.keccak256(topic) {
                            return true
                        }
                        return false
                    })[index]
                    
                    guard let l2tol1log = receipt.l2ToL1Logs?.filter({
                        if $0.sender.address == ZkSyncAddresses.MessengerAddress {
                            return true
                        }
                        return false
                    })[index] else {
                        fatalError("No l2 to l1 log found.")
                    }
                    
                    (self.zkSync as! JsonRpc2_0ZkSync).zksGetL2ToL1LogProof(txHash, logIndex: Int(l2tol1log.logIndex)) { result in
                        DispatchQueue.global().async {
                            switch result {
                            case .success(let proof):
                                guard let path = Bundle.main.path(forResource: "IL1Messenger", ofType: "json") else { return }
                                
                                let data = try! Data(contentsOf: URL(fileURLWithPath: path), options: .mappedIfSafe)
                                let jsonResult = try! JSONSerialization.jsonObject(with: data, options: .mutableLeaves)
                                guard let json = jsonResult as? [String: Any], let abi = json["abi"] as? [[String: Any]] else { return }
                                
                                guard let abiData = try? JSONSerialization.data(withJSONObject: abi, options: []) else { return }
                                let abiString = String(data: abiData, encoding: .utf8)!
                                
                                let contract = self.zkSync.web3.contract(abiString)!
                                
                                let eventData = contract.parseEvent(log).eventData
                                let message = eventData?["_message"] as? Data ?? Data()
                                
                                let nonce = try! self.zkSync.web3.eth.getTransactionCountPromise(address: EthereumAddress(self.signer.address)!, onBlock: ZkBlockParameterName.committed.rawValue).wait()
                                
                                let result = try! defaultEthereumProvider.finalizeEthWithdrawal(
                                    receipt.l1BatchNumber,
                                    l2MessageIndex: BigUInt(proof.id),
                                    l2TxNumberInBlock: receipt.l1BatchTxIndex,
                                    message: message,
                                    proof: proof.proof.compactMap({ Data(fromHex: $0) }),
                                    nonce: nonce
                                ).wait()
                                
                                print(result.hash)
                                
                                callback()
                            case .failure(let error):
                                fatalError(error.localizedDescription)
                            }
                        }
                    }
                default: return
                }
            }
        }
    }
    
    func withdrawViaWallet(callback: @escaping (() -> Void)) {
        let amount = BigUInt(1000000000000000000)
        
        let token: Token = Token(l1Address: Token.DefaultAddress, l2Address: Token.DefaultAddress, symbol: "ETH", decimals: 18)
        
        _ = try! wallet.withdraw("0x000000000000000000000000000000000000800a", amount: amount, token: token).wait()
        
        callback()
    }
}
