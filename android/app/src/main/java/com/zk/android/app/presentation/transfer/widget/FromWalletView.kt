package com.zk.android.app.presentation.transfer.widget

import android.content.Context
import android.util.AttributeSet
import android.view.LayoutInflater
import androidx.constraintlayout.widget.ConstraintLayout
import androidx.core.view.setPadding
import com.zk.android.app.R
import com.zk.android.app.databinding.ViewFromBalanceBinding
import com.zk.android.app.utils.toPixelI


class FromWalletView @JvmOverloads constructor(
    context: Context,
    attrs: AttributeSet? = null,
    defStyleAttr: Int = 0
) : ConstraintLayout(context, attrs, defStyleAttr) {

    private val binding = ViewFromBalanceBinding.inflate(LayoutInflater.from(context), this)

    init {
        setBackgroundResource(R.drawable.bg_gray_content)
        setPadding(context.toPixelI(16))
    }

    fun setupModel(model: FromWalletViewModel) {
        binding.amount.setText(model.balance)
        binding.amount.clearFocus()
    }

    fun getAmount(): Float {
        val text = binding.amount.text?.toString() ?: return 0f
        return text.toFloat()
    }
}