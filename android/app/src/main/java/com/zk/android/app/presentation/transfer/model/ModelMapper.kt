package com.zk.android.app.presentation.transfer.model

import com.zk.android.app.domain.balance.BalanceModelDomain
import com.zk.android.app.presentation.transfer.widget.FromWalletViewModel
import com.zk.android.app.presentation.transfer.widget.ToWalletViewModel
import com.zk.android.app.utils.weiToEther
import javax.inject.Inject


class ModelMapper @Inject constructor() {
    fun map(balanceModelDomain: BalanceModelDomain): BalanceModel {
        return with(balanceModelDomain) {
            BalanceModel(
                FromWalletViewModel(l1Balance.weiToEther().toString()),
                ToWalletViewModel(l2Balance.weiToEther().toString())
            )
        }
    }
}