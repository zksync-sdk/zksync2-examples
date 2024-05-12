package com.zk.android.app.presentation.transfer


import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.core.view.isVisible
import androidx.fragment.app.Fragment
import androidx.fragment.app.viewModels
import com.zk.android.app.databinding.FragmentTransferBinding
import com.zk.android.app.utils.viewLifecycle
import dagger.hilt.android.AndroidEntryPoint
import org.orbitmvi.orbit.viewmodel.observe

@AndroidEntryPoint
internal class TransferFragment : Fragment() {

    private val viewModel: TransferViewModel by viewModels()
    private var viewBinding: FragmentTransferBinding by viewLifecycle()


    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View {
        viewBinding = FragmentTransferBinding.inflate(inflater, container, false)
        return viewBinding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        setupListeners()
        viewModel.observe(
            lifecycleOwner = viewLifecycleOwner,
            state = ::handleState,
            sideEffect = ::handleSideEffect
        )
    }

    private fun setupListeners() {
        viewBinding.btnTransfer.setOnClickListener {
            val amount = viewBinding.balanceFrom.getAmount()
            viewModel.obtainAction(Action.ClickSend(amount))
        }
        viewBinding.btnConnection.setOnClickListener { viewModel.obtainAction(Action.ClickConnect) }
    }

    private fun handleState(state: ViewState) {
        applyStateVisibility(state)
        when (state) {
            is ViewState.Connection -> showConnection(state)
            is ViewState.Data -> showData(state)
        }
    }

    private fun applyStateVisibility(state: ViewState) {
        viewBinding.groupConnection.isVisible = state !is ViewState.Data
        viewBinding.btnConnection.isVisible = state is ViewState.Connection.Error
        viewBinding.groupData.isVisible = state is ViewState.Data
    }

    private fun showConnection(state: ViewState.Connection) {
        with(viewBinding) {
            connectionImg.setAnimation(state.animationRes)
            connectionImg.playAnimation()
            connectionText.setText(state.descRes)
        }
    }

    private fun showData(state: ViewState.Data) {
        viewBinding.balanceFrom.setupModel(state.balance.fromWalletModel)
        viewBinding.balanceTo.setupModel(state.balance.toWalletModel)
        viewBinding.balanceFrom.isEnabled = !state.inProgress
        viewBinding.balanceTo.isEnabled = !state.inProgress
        viewBinding.btnTransfer.isEnabled = !state.inProgress
        viewBinding.balanceSending.isVisible = state.inProgress
    }

    private fun handleSideEffect(effect: Effect) {
        when (effect) {
            is Navigation -> handleNavigation(effect)
            is ViewEffect -> handleEffect(effect)
        }
    }

    private fun handleNavigation(navigation: Navigation) {}
    private fun handleEffect(effect: ViewEffect) {
        when (effect) {
            is ViewEffect.ShowToast -> showToast(effect)
        }
    }

    private fun showToast(effect: ViewEffect.ShowToast) {
        Toast.makeText(requireContext(), effect.textResId, Toast.LENGTH_LONG).show()
    }
}