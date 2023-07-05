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
import web3swift_zksync
#endif
#if canImport(ZkSync2)
import ZkSync2
#endif

class PaymasterManager: BaseManager {
    func mintTokenUsingPaymaster(callback: (() -> Void)) {
        let manager = KeystoreManager.init([credentials])
        zkSync.web3.eth.web3.addKeystoreManager(manager)
        self.eth.addKeystoreManager(manager)
        
        guard let path = Bundle.main.path(forResource: "Token", ofType: "json") else { return }
        
        let data = try! Data(contentsOf: URL(fileURLWithPath: path), options: .mappedIfSafe)
        let jsonResult = try! JSONSerialization.jsonObject(with: data, options: .mutableLeaves)
        guard let json = jsonResult as? [String: Any], let abi = json["abi"] as? [[String: Any]] else { return }
        
        guard let abiData = try? JSONSerialization.data(withJSONObject: abi, options: []) else { return }
        let abiString = String(data: abiData, encoding: .utf8)!
        
        let contractAddress = EthereumAddress("0xbc6b677377598a79fa1885e02df1894b05bc8b33")!
        
        let contract = zkSync.web3.contract(abiString, at: contractAddress)!
        
        let value = BigUInt(100)
        
        let parameters = [
            EthereumAddress("0x36615Cf349d7F6344891B1e7CA7C72883F5dc049")! as AnyObject,
            value as AnyObject
        ] as [AnyObject]
        
        var transactionOptions1 = TransactionOptions.defaultOptions
        transactionOptions1.from = EthereumAddress(signer.address)!
        
        guard let writeTransaction = contract.write(
            "mint",
            parameters: parameters,
            transactionOptions: transactionOptions1
        ) else {
            return
        }
        
        var estimate = EthereumTransaction.createFunctionCallTransaction(from: EthereumAddress(signer.address)!, to: writeTransaction.transaction.to, gasPrice: BigUInt.zero, gasLimit: BigUInt.zero, data: writeTransaction.transaction.data)
        
        let nonce = try! zkSync.web3.eth.getTransactionCountPromise(address: EthereumAddress(signer.address)!, onBlock: ZkBlockParameterName.committed.rawValue).wait()
        
        let fee = try! zkSync.zksEstimateFee(estimate).wait()
        
        estimate.parameters.EIP712Meta?.gasPerPubdata = BigUInt(160000)
        
        let to = EthereumAddress("0xbc6b677377598a79fa1885e02df1894b05bc8b33")!
        
        let paymasterAddress = EthereumAddress("0x49720d21525025522040f73da5b3992112bbec00")!
        let paymasterInput = Paymaster.encodeApprovalBased(
            to,
            minimalAllowance: BigUInt(1),
            input: Data()
        )
        
        estimate.parameters.EIP712Meta?.paymasterParams = PaymasterParams(paymaster: paymasterAddress, paymasterInput: paymasterInput)
        
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
        
        var transaction = EthereumTransaction(
            type: .eip712,
            to: estimate.to,
            nonce: nonce,
            chainID: chainId,
            data: estimate.data,
            parameters: ethereumParameters
        )
        
        let signature = signer.signTypedData(signer.domain, typedData: transaction).addHexPrefix()
        
        let unmarshalledSignature = SECP256K1.unmarshalSignature(signatureData: Data(fromHex: signature)!)!
        transaction.envelope.r = BigUInt(fromHex: unmarshalledSignature.r.toHexString().addHexPrefix())!
        transaction.envelope.s = BigUInt(fromHex: unmarshalledSignature.s.toHexString().addHexPrefix())!
        transaction.envelope.v = BigUInt(unmarshalledSignature.v)
        
        _ = try! zkSync.web3.eth.sendRawTransactionPromise(transaction).wait()
        
        callback()
    }
}
