package com.zk.android.app.domain.deposit

import com.zk.android.app.domain.balance.BalanceModelDomain


data class DepositModelDomain (
    val newBalance: BalanceModelDomain
)