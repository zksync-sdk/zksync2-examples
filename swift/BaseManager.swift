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
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
#endif

class BaseManager {
    let credentials = Credentials(Config.privateKey)
    
    let zkSync: ZkSync = ZkSyncImpl(Config.zkSyncProviderUrl)
    let ethereum: web3 = try! Web3.new(Config.ethereumProviderUrl)
    
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
        ZkSyncWallet(zkSync, ethereum: ethereum, ethSigner: signer, feeToken: Token.ETH)
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
