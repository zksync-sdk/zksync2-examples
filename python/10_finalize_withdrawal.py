import os

from eth_account import Account
from eth_account.signers.local import LocalAccount
from eth_typing import HexStr
from hexbytes import HexBytes
from web3 import Web3
from web3.middleware import geth_poa_middleware
from web3.types import TxReceipt

from zksync2.account.wallet import Wallet
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

    # Finalize withdraw of previous successful withdraw transaction
    tx_receipt = wallet.finalize_withdrawal("<TX_HASH>")

    print(f"Finalize withdraw transaction: {tx_receipt['transactionHash'].hex()}")
