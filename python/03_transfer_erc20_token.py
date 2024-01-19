import os

from eth_account import Account
from eth_account.signers.local import LocalAccount
from eth_typing import HexAddress, HexStr
from web3 import Web3

from zksync2.account.wallet import Wallet
from zksync2.core.types import EthBlockParams, TransferTransaction
from zksync2.module.module_builder import ZkSyncBuilder


def transfer_erc20(
        token_contract,
        account: LocalAccount,
        address: HexAddress,
        amount: float) -> HexStr:
    """
       Transfer ERC20 token to a desired address on zkSync network

       :param token_contract:
           Instance of ERC20 contract

       :param account:
           From which account the transfer will be made

       :param address:
         Desired ETH address that you want to transfer to.

       :param amount:
         Desired ETH amount that you want to transfer.

       :return:
         The transaction hash of the transfer transaction.

       """
    tx = token_contract.functions.transfer(address, amount).build_transaction({
        "nonce": zk_web3.zksync.get_transaction_count(account.address, EthBlockParams.LATEST.value),
        "from": account.address,
        "maxPriorityFeePerGas": 1_000_000,
        "maxFeePerGas": zk_web3.zksync.gas_price,
    })

    signed = account.sign_transaction(tx)

    # Send transaction to zkSync network
    tx_hash = zk_web3.zksync.send_raw_transaction(signed.rawTransaction)
    print(f"Tx: {tx_hash.hex()}")

    tx_receipt = zk_web3.zksync.wait_for_transaction_receipt(
        tx_hash, timeout=240, poll_latency=0.5
    )
    print(f"Tx status: {tx_receipt['status']}")

    return tx_hash


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

    # Get l2 token address from l1
    l1_token_address = "0x70a0F165d6f8054d0d0CF8dFd4DD2005f0AF6B55"
    l2_token_address = zk_web3.zksync.l2_token_address(l1_token_address)

    print(wallet.get_balance(token_address=l2_token_address))

    tx_hash = wallet.transfer(TransferTransaction(to=Web3.to_checksum_address(wallet.address),
                                                  token_address=Web3.to_checksum_address(l2_token_address),
                                                  amount=1))

    tx_receipt = zk_web3.zksync.wait_for_transaction_receipt(
        tx_hash, timeout=240, poll_latency=0.5
    )
    print(wallet.get_balance(token_address=l2_token_address))

    print(f"Tx receipt : {tx_receipt}")
