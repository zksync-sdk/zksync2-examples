package com.zk.android.app.presentation.transfer

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import com.zk.android.app.R
import com.zk.android.app.domain.balance.GetBalanceUseCase
import com.zk.android.app.domain.deposit.DepositUseCase
import com.zk.android.app.presentation.transfer.model.ModelMapper
import com.zk.android.app.presentation.transfer.widget.FromWalletViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import org.orbitmvi.orbit.ContainerHost
import org.orbitmvi.orbit.container
import org.orbitmvi.orbit.syntax.simple.intent
import org.orbitmvi.orbit.syntax.simple.postSideEffect
import org.orbitmvi.orbit.syntax.simple.reduce
import org.orbitmvi.orbit.syntax.simple.runOn
import javax.inject.Inject

@HiltViewModel
internal class TransferViewModel @Inject constructor(
    private val application: Application,
    private val mapper: ModelMapper,
    private val balanceUseCase: GetBalanceUseCase,
    private val depositUseCase: DepositUseCase
) : AndroidViewModel(application), ContainerHost<ViewState, Effect> {

    override val container = viewModelScope.container<ViewState, Effect>(ViewState.Connection.Loading)

    init {
        connect()
    }

    private fun connect() = intent {
        reduce { ViewState.Connection.Loading }
        val result = balanceUseCase.getBalance().getOrElse {
            reduce { ViewState.Connection.Error }
            return@intent
        }
        val state = ViewState.Data(
            mapper.map(result)
        )
        reduce { state }
    }

    fun obtainAction(action: Action) {
        when (action) {
            is Action.ClickConnect -> processClickConnect()
            is Action.ClickSend -> processClickSend(action)
        }
    }

    private fun processClickConnect() = intent {
        connect()
    }

    private fun processClickSend(action: Action.ClickSend) = intent {
        runOn(ViewState.Data::class) {
            val amount = action.amount
            reduce {
                state.copy(
                    balance = state.balance.copy(fromWalletModel = FromWalletViewModel(amount.toString())),
                    inProgress = true
                )
            }
            val result = depositUseCase.deposit(action.amount).getOrElse {
                reduce { ViewState.Connection.Error }
                return@runOn
            }
            val state = ViewState.Data(
                mapper.map(result.newBalance)
            )
            reduce { state }
            postSideEffect(ViewEffect.ShowToast(R.string.deposit_completed))
        }
    }
}