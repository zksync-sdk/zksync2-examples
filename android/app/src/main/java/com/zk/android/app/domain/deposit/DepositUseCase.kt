package com.zk.android.app.domain.deposit

import com.zk.android.app.domain.balance.GetBalanceUseCase
import com.zk.android.app.network.WalletProvider
import com.zk.android.app.utils.toWei
import com.zk.android.app.utils.waitForTransactionReceiptSafe
import io.zksync.transaction.type.DepositTransaction
import io.zksync.utils.ZkSyncAddresses
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import java.math.BigInteger
import javax.inject.Inject


class DepositUseCase @Inject constructor(
    private val walletProvider: WalletProvider,
    private val balanceUseCase: GetBalanceUseCase
) {

    suspend fun deposit(amount: Float): Result<DepositModelDomain> = withContext(Dispatchers.IO) {
        val wallet = walletProvider.create().getOrElse {
            return@withContext Result.failure(it)
        }

        val sendAmount: BigInteger = amount.toWei()

        val transaction = DepositTransaction(ZkSyncAddresses.ETH_ADDRESS, sendAmount)
        val hash = wallet.deposit(transaction).sendAsync().join().result
        wallet.transactionReceiptProcessorL1.waitForTransactionReceiptSafe(hash)

        val newBalance = balanceUseCase.getBalance().getOrElse {
            return@withContext Result.failure(it)
        }

        return@withContext Result.success(DepositModelDomain(newBalance))
    }
}