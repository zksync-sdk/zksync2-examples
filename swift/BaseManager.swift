//
//  BaseManager.swift
//  zkSync-Demo
//
//  Created by Bojan on 14.5.23..
//

import Foundation
import CryptoKit
import BigInt
#if canImport(web3swift)
import web3swift
#else
import web3swift_zksync
#endif
#if canImport(ZkSync2)
import ZkSync2
#endif

class BaseManager {
    static let privateKey = "0x7726827caac94a7f9e1b160f7ea819f172f7b6f9d2a97f992c38edeab82d4110"
    
    let credentials = Credentials(BaseManager.privateKey)
    
    let zkSync: ZkSync = ZkSyncImpl(URL(string: "http://127.0.0.1:3050")!)
    let ethereum: web3 = try! Web3.new(URL(string: "http://127.0.0.1:8545")!)
    
    init() {
        let keystoreManager = KeystoreManager([credentials])
        self.ethereum.addKeystoreManager(keystoreManager)
    }
    
    var chainId: BigUInt {
        try! zkSync.web3.eth.getChainIdPromise().wait()
    }
    
    var nonce: BigUInt {
        try! zkSync.web3.eth.getTransactionCountPromise(
            address: EthereumAddress(signer.address)!,
            onBlock: ZkBlockParameterName.committed.rawValue
        ).wait()
    }
    
    var signer: EthSigner {
        PrivateKeyEthSigner(credentials, chainId: chainId)
    }
    
    var wallet: ZkSyncWallet {
        ZkSyncWallet(zkSync, ethSigner: signer, feeToken: Token.ETH)
    }
    
    var transactionReceiptProcessor: ZkSyncTransactionReceiptProcessor {
        ZkSyncTransactionReceiptProcessor(zkSync: zkSync)
    }
    
    func signTransaction(_ transaction: inout EthereumTransaction) {
        let signature = signer.signTypedData(signer.domain, typedData: transaction).addHexPrefix()
        
        let unmarshalledSignature = SECP256K1.unmarshalSignature(signatureData: Data(fromHex: signature)!)!
        transaction.envelope.r = BigUInt(fromHex: unmarshalledSignature.r.toHexString().addHexPrefix())!
        transaction.envelope.s = BigUInt(fromHex: unmarshalledSignature.s.toHexString().addHexPrefix())!
        transaction.envelope.v = BigUInt(unmarshalledSignature.v)
    }
}
