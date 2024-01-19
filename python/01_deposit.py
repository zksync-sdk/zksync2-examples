import os

from eth_account import Account
from eth_account.signers.local import LocalAccount
from eth_typing import HexStr
from web3 import Web3

from zksync2.account.wallet import Wallet
from zksync2.core.types import Token, DepositTransaction
from zksync2.manage_contracts.utils import zksync_abi_default
from zksync2.module.module_builder import ZkSyncBuilder

if __name__ == "__main__":
    # Get the private key from OS environment variables
    PRIVATE_KEY = bytes.fromhex(os.environ.get("PRIVATE_KEY"))

    # Set a provider
    ZKSYNC_PROVIDER = "https://sepolia.era.zksync.dev"
    ETH_PROVIDER = "https://rpc.ankr.com/eth_sepolia"

    # Connect to zkSync network
    zk_web3 = ZkSyncBuilder.build(ZKSYNC_PROVIDER)

    # Connect to Ethereum network
    eth_web3 = Web3(Web3.HTTPProvider(ETH_PROVIDER))

    # Get account object by providing from private key
    account: LocalAccount = Account.from_key(PRIVATE_KEY)

    # Create Ethereum provider
    wallet = Wallet(zk_web3, eth_web3, account)

    #ZkSync contract
    zksync_contract = eth_web3.eth.contract(Web3.to_checksum_address(zk_web3.zksync.main_contract_address),
                                                 abi=zksync_abi_default())

    amount = 0.1
    l2_balance_before = wallet.get_balance()

    tx_hash = wallet.deposit(DepositTransaction(token=Token.create_eth().l1_address,
                                                     amount=Web3.to_wei(amount, "ether"),
                                                     to=wallet.address))

    l1_tx_receipt = eth_web3.eth.wait_for_transaction_receipt(tx_hash)
    l2_hash = zk_web3.zksync.get_l2_hash_from_priority_op(l1_tx_receipt, zksync_contract)
    l2_tx_receipt = zk_web3.zksync.wait_for_transaction_receipt(transaction_hash=l2_hash,
                                                    timeout=360,
                                                    poll_latency=10)

    print(f"L1 transaction: {l1_tx_receipt}")
    print(f"L2 transaction: {l2_tx_receipt}")
