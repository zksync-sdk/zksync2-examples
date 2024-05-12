package com.zk.android.app.presentation.transfer

import androidx.annotation.RawRes
import androidx.annotation.StringRes
import com.zk.android.app.R
import com.zk.android.app.presentation.transfer.model.BalanceModel

internal sealed class ViewState {
    sealed class Connection(
        @RawRes val animationRes: Int,
        @StringRes val descRes: Int
    ) : ViewState() {
        data object Loading : Connection(
            R.raw.lottie_loading,
            R.string.connection
        )

        data object Error : Connection(
            R.raw.lottie_connection_error,
            R.string.connection_error
        )
    }

    data class Data(val balance: BalanceModel, val inProgress: Boolean = false) : ViewState()
}

internal sealed interface Effect

internal sealed class Navigation : Effect {

}

internal sealed class ViewEffect : Effect {
    class ShowToast(@StringRes val textResId: Int) : ViewEffect()
}

internal sealed class Action {
    data object ClickConnect : Action()
    class ClickSend(val amount: Float) : Action()
}