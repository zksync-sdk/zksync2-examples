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

    var zkSync: ZkSyncClient
    let ethClient: EthereumClient
    var web: Web3

    init() {
        zkSync = ZkSyncClientImpl(Config.zkSyncProviderUrl)
        ethClient = EthereumClientImpl(Config.ethereumProviderUrl)
        web = Web3(provider: Web3HttpProvider(url: Config.ethereumProviderUrl, network: .Mainnet))
        
        let keystoreManager = KeystoreManager([credentials])
        self.web.addKeystoreManager(keystoreManager)
    }

    var chainId: BigUInt {
        return zkSync.web3.eth.provider.network!.chainID
    }

    func getNonce() async -> BigUInt {
        try! await zkSync.web3.eth.getTransactionCount(
            for: EthereumAddress(signer.address)!,
            onBlock: .pending
        )
    }

    var signer: ETHSigner {
        BaseSigner(credentials, chainId: chainId)
    }
    
    var deployer: BaseDeployer {
        BaseDeployer(zkSync, web3: web, ethSigner: signer)
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

    func signTransaction(_ transaction: inout CodableTransaction) {
        try! transaction.sign(privateKey: credentials.privateKey)
    }
}
