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
        let amount: BigUInt = 1
        
        let manager = KeystoreManager.init([credentials])
        self.eth.addKeystoreManager(manager)
        
        let l1ERC20Bridge = zkSync.web3.contract(
            Web3.Utils.IL1Bridge,
            at: EthereumAddress("0x4ee775658259028d399f4cf9d637b14773472988")
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
                    
                    let token = Token(l1Address: "0x36615Cf349d7F6344891B1e7CA7C72883F5dc049", l2Address: Token.DefaultAddress, symbol: "ETH", decimals: 18)
                    
                    let result = try! defaultEthereumProvider.deposit(with: token, amount: amount, operatorTips: BigUInt(0), to: "0x36615Cf349d7F6344891B1e7CA7C72883F5dc049").wait()
                    
                    let receipt = self.transactionReceiptProcessor.waitForTransactionReceipt(hash: result.hash)
                    
                    assert(receipt?.status == .ok)
                    
                    callback()
                default: return
                }
            }
        }
    }
}
