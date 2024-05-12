package com.zk.android.app.presentation.transfer.model

import com.zk.android.app.presentation.transfer.widget.FromWalletViewModel
import com.zk.android.app.presentation.transfer.widget.ToWalletViewModel


data class BalanceModel(
    val fromWalletModel: FromWalletViewModel,
    val toWalletModel: ToWalletViewModel,
)