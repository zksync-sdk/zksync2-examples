package com.zk.android.app.domain.balance

import com.zk.android.app.network.WalletProvider
import com.zk.android.app.utils.sendSafe
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.withContext
import java.math.BigInteger
import javax.inject.Inject


class GetBalanceUseCase @Inject constructor(
    private val walletProvider: WalletProvider
) {

    suspend fun getBalance(): Result<BalanceModelDomain> = withContext(Dispatchers.IO) {
        val wallet = walletProvider.create().getOrElse {
            return@withContext Result.failure(it)
        }

        val l1Balance: BigInteger = wallet.balanceL1.sendSafe().getOrElse {
            return@withContext Result.failure(it)
        }
        val l2Balance: BigInteger = wallet.balance.sendSafe().getOrElse {
            return@withContext Result.failure(it)
        }

        return@withContext Result.success(
            BalanceModelDomain(
                l1Balance,
                l2Balance
            )
        )
    }
}