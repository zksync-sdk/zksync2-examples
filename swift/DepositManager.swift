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
import web3swift_zksync
#endif
#if canImport(ZkSync2)
import ZkSync2
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
                        at: EthereumAddress("0x4ee775658259028d399f4cf9d637b14773472988")
                    )!
                    
                    let defaultEthereumProvider = DefaultEthereumProvider(
                        self.ethereum,
                        l1ERC20Bridge: l1ERC20Bridge,
                        zkSyncContract: zkSyncContract,
                        gasProvider: DefaultGasProvider()
                    )
                    
                    let amount = BigUInt(1_000_000_000_000)
                    
                    _ = try! defaultEthereumProvider.deposit(
                        with: Token.ETH,
                        amount: amount,
                        operatorTips: BigUInt(0),
                        to: self.signer.address
                    ).wait()
                    
                    callback()
                case .failure(let error):
                    print("Error:", error)
                }
            }
        }
    }
    
    func depositViaWallet(callback: @escaping (() -> Void)) {
        let amount = BigUInt(1_000_000_000_000)
        
        _ = try! wallet.deposit(
            signer.address,
            amount: amount
        ).wait()
        
        callback()
    }
}
