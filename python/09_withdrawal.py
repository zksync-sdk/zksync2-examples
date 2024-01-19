import os

from eth_account import Account
from eth_account.signers.local import LocalAccount
from eth_typing import HexStr
from web3 import Web3

from zksync2.account.wallet import Wallet
from zksync2.core.types import Token, WithdrawTransaction
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

    l2_balance_before = wallet.get_balance()
    amount = 0.005

    withdraw_tx_hash = wallet.withdraw(
        WithdrawTransaction(token=Token.create_eth().l1_address, amount=Web3.to_wei(amount, "ether")))

    tx_receipt = zk_web3.zksync.wait_for_transaction_receipt(
        withdraw_tx_hash, timeout=240, poll_latency=0.5
    )

    l2_balance_after = wallet.get_balance()

    print(f"Withdraw transaction receipt: {tx_receipt}")
    print("Wait for withdraw transaction to be finalized on L2 network (11-24 hours)")
    print("Read more about withdrawal delay: https://era.zksync.io/docs/dev/troubleshooting/withdrawal-delay.html")
    print("When withdraw transaction is finalized, execute 10_finalize_withdrawal.py script  "
          "with WITHDRAW_TX_HASH environment variable set")
