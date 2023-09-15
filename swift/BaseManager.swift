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
    
    let zkSync: ZkSyncClient = ZkSyncClientImpl(Config.zkSyncProviderUrl)
    let ethClient: EthereumClient = EthereumClientImpl(Config.ethereumProviderUrl)
    let web: web3 = try! Web3.new(Config.ethereumProviderUrl)
    
    init() {
        let keystoreManager = KeystoreManager([credentials])
        self.web.addKeystoreManager(keystoreManager)
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
    
    var signer: ETHSigner {
        EthSignerImpl(credentials, chainId: chainId)
    }
    
    var deployer: DeployerImpl {
        DeployerImpl(zkSync, web3: web, ethSigner: signer)
    }
    
    var walletL1: WalletL1 {
        WalletL1(zkSync, ethClient: ethClient, web3: web, ethSigner: signer)
    }
    
    var walletL2: WalletL2 {
        WalletL2(zkSync, ethClient: ethClient, web3: web, ethSigner: signer)
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
