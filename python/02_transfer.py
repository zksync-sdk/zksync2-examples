import os

from eth_account import Account
from eth_account.signers.local import LocalAccount
from eth_typing import HexStr
from web3 import Web3

from zksync2.account.wallet import Wallet
from zksync2.core.types import TransferTransaction
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

    # Show balance before ETH transfer
    print(f"Balance before transfer : {wallet.get_l1_balance()} ETH")

    # Perform the ETH transfer
    tx_hash = wallet.transfer(TransferTransaction(
        to=HexStr("0x81E9D85b65E9CC8618D85A1110e4b1DF63fA30d9"),
        amount=Web3.to_wei(0.001, "ether"),

    ))
    zk_web3.zksync.wait_for_transaction_receipt(
        tx_hash, timeout=240, poll_latency=0.5
    )

    # Show balance after ETH transfer
    print(f"Balance after transfer : {wallet.get_balance()} ETH")
