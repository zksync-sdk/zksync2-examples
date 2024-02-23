import os

from eth_account import Account
from eth_account.signers.local import LocalAccount
from eth_typing import HexStr
from web3 import Web3

from zksync2.account.wallet import Wallet
from zksync2.module.module_builder import ZkSyncBuilder
from zksync2.core.types import EthBlockParams


def check_balance():
    # Get the private key from OS environment variables
    PRIVATE_KEY = HexStr("0xc273a8616a4c58de9e58750fd2672d07b10497d64cd91b5942cce0909aaa391a")

    # Set a provider
    ZKSYNC_PROVIDER = "https://sepolia.era.zksync.dev"
    ETH_PROVIDER = "https://rpc.ankr.com/eth_sepolia"

    account: LocalAccount = Account.from_key(PRIVATE_KEY)
    zksync_web3 = ZkSyncBuilder.build(ZKSYNC_PROVIDER)
    eth_web3 = Web3(Web3.HTTPProvider(ETH_PROVIDER))

    wallet = Wallet(zksync_web3, eth_web3, account)

    # Balance using provider
    zk_balance = zksync_web3.zksync.get_balance(account.address, EthBlockParams.LATEST.value)
    print(f"Balance: {zk_balance}")

    # Blance using wallet
    balance = wallet.get_balance()
    l1_balance = wallet.get_l1_balance()
    print(f"Balance: {balance}")
    print(f"L1 balance: {l1_balance}")

    # Token l2 balance
    l1_address = "0x70a0F165d6f8054d0d0CF8dFd4DD2005f0AF6B55"
    l2_address = zksync_web3.zksync.l2_token_address(l1_address)
    token_balance = wallet.get_balance(token_address=l2_address)
    print(f"Token balance: {token_balance}")


if __name__ == "__main__":
    check_balance()
