package com.zk.android.app.network

import androidx.annotation.WorkerThread
import io.zksync.protocol.ZkSync
import io.zksync.protocol.account.Wallet
import org.web3j.crypto.Credentials
import org.web3j.protocol.Web3j
import javax.inject.Inject


class WalletProvider @Inject constructor(
    private val web3: Web3j,
    private val zkSync: ZkSync,
    private val credentials: Credentials
) {

    @WorkerThread
    //TODO can we use single Wallet instance?
    fun create(): Result<Wallet> {
        return runCatching { Wallet(web3, zkSync, credentials) }
    }
}