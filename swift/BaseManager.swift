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
import Web3Core
#else
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
#endif

class BaseManager {
    let credentials = Credentials(Config.privateKey)

    let zkSync: ZkSync
    let ethereum: Web3

    init() async {
        zkSync = await ZkSyncImpl(Config.zkSyncProviderUrl)
        ethereum = try! await Web3.new(Config.ethereumProviderUrl)
        
        let keystoreManager = KeystoreManager([credentials])
        self.ethereum.addKeystoreManager(keystoreManager)
    }

//222    var chainId: BigUInt {
//        try! zkSync.web3.eth.getChainIdPromise().wait()
//    }
//
//    func getNonce() async -> BigUInt {
//        try! await zkSync.web3.eth.getTransactionCount(
//            for: EthereumAddress(signer.address)!,
//            onBlock: .pending
//        )
//    }
//
//    var signer: EthSigner {
//        PrivateKeyEthSigner(credentials, chainId: chainId)
//    }
//
//    var wallet: ZkSyncWallet {
//        ZkSyncWallet(zkSync, ethereum: ethereum, ethSigner: signer, feeToken: Token.ETH)
//    }
//
//    var transactionReceiptProcessor: ZkSyncTransactionReceiptProcessor {
//        ZkSyncTransactionReceiptProcessor(zkSync: zkSync)
//    }
//
//    func signTransaction(_ transaction: inout CodableTransaction) {
//        let signature = signer.signTypedData(signer.domain, typedData: transaction).addHexPrefix()
//
//        let unmarshalledSignature = SECP256K1.unmarshalSignature(signatureData: Data(fromHex: signature)!)!
//        transaction.envelope.r = BigUInt(fromHex: unmarshalledSignature.r.toHexString().addHexPrefix())!
//        transaction.envelope.s = BigUInt(fromHex: unmarshalledSignature.s.toHexString().addHexPrefix())!
//        transaction.envelope.v = BigUInt(unmarshalledSignature.v)
//    }
}
