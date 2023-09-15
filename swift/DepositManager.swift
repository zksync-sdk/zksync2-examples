//
//  DepositManager.swift
//  zkSync-Demo
//
//  Created by Bojan on 14.5.23..
//

import Foundation
import BigInt
#if canImport(web3swift)
import web3swift
#else
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
#endif

class DepositManager: BaseManager {
    func depositViaWallet(callback: @escaping (() -> Void)) {
        let amount = BigUInt(1_000_000_000_000)
        
        _ = try! walletL1.deposit(
            signer.address,
            amount: amount,
            token: Token.ETH
        ).wait()
        
        callback()
    }
}
