package com.zk.android.app.utils

import android.content.Context
import android.util.TypedValue

fun Context.toPixelF(dpSize: Int) = toPixelF(dpSize.toFloat())

fun Context.toPixelF(dpSizeF: Float) = TypedValue.applyDimension(
    TypedValue.COMPLEX_UNIT_DIP, dpSizeF, resources.displayMetrics
)

fun Context.toPixelI(dpSize: Int) = toPixelF(dpSize).toInt()