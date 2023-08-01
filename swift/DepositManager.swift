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
import Web3Core
#else
import web3swift_zksync2
#endif
#if canImport(zkSync2_swift)
import zkSync2_swift
#endif

class DepositManager: BaseManager {
    func deposit(callback: @escaping (() -> Void)) {
        zkSync.zksMainContract { result in
            DispatchQueue.global().async {
                switch result {
                case .success(let address):
                    let zkSyncContract = self.ethereum.contract(
                        Web3.Utils.IZkSync,
                        at: EthereumAddress(address)
                    )!

                    let l1ERC20Bridge = self.zkSync.web3.contract(
                        Web3.Utils.IL1Bridge,
                        at: EthereumAddress(self.signer.address)
                    )!

                    let defaultEthereumProvider = DefaultEthereumProvider(
                        self.ethereum,
                        l1ERC20Bridge: l1ERC20Bridge,
                        zkSyncContract: zkSyncContract,
                        gasProvider: DefaultGasProvider()
                    )

                    let amount = BigUInt(1_000_000_000_000)

                    Task {
                        _ = try! await defaultEthereumProvider.deposit(
                            with: Token.ETH,
                            amount: amount,
                            operatorTips: BigUInt(0),
                            to: self.signer.address
                        )
                        
                        callback()
                    }
                case .failure(let error):
                    print("Error:", error)
                }
            }
        }
    }

    func depositViaWallet(callback: @escaping (() -> Void)) {
        Task {
            let amount = BigUInt(1_000_000_000_000)
            
            _ = await wallet.deposit(
                signer.address,
                amount: amount
            )
            
            callback()
        }
    }
}
