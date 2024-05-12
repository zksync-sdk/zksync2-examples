package com.zk.android.app.utils

import androidx.fragment.app.Fragment
import androidx.lifecycle.DefaultLifecycleObserver
import androidx.lifecycle.LifecycleOwner
import kotlin.properties.ReadWriteProperty
import kotlin.reflect.KProperty

fun <T> Fragment.viewLifecycle(onDestroyFinalize: ((T?) -> Unit)? = null): ReadWriteProperty<Fragment, T> =
    object : ReadWriteProperty<Fragment, T>, DefaultLifecycleObserver {

        private var field: T? = null

        private var lifecycleOwner: LifecycleOwner? = null

        init {
            viewLifecycleOwnerLiveData
                .observe(this@viewLifecycle) { newViewLifecycleOwner ->
                    lifecycleOwner
                        ?.lifecycle
                        ?.removeObserver(this)

                    lifecycleOwner = newViewLifecycleOwner.also {
                        it.lifecycle.addObserver(this)
                    }
                }
        }

        override fun onDestroy(owner: LifecycleOwner) {
            super.onDestroy(owner)
            onDestroyFinalize?.invoke(field)
            field = null
            lifecycleOwner = null
        }

        override fun getValue(thisRef: Fragment, property: KProperty<*>): T {
            return this.field!!
        }

        override fun setValue(thisRef: Fragment, property: KProperty<*>, value: T) {
            this.field = value
        }
    }

fun <T> Fragment.viewLifecycleNullable(): ReadWriteProperty<Fragment, T?> =
    object : ReadWriteProperty<Fragment, T?>, DefaultLifecycleObserver {

        private var field: T? = null

        private var lifecycleOwner: LifecycleOwner? = null

        init {
            viewLifecycleOwnerLiveData
                .observe(this@viewLifecycleNullable) { newViewLifecycleOwner ->
                    lifecycleOwner
                        ?.lifecycle
                        ?.removeObserver(this)

                    lifecycleOwner = newViewLifecycleOwner.also {
                        it.lifecycle.addObserver(this)
                    }
                }
        }

        override fun onDestroy(owner: LifecycleOwner) {
            super.onDestroy(owner)
            field = null
            lifecycleOwner = null
        }

        override fun getValue(thisRef: Fragment, property: KProperty<*>): T? {
            return this.field
        }

        override fun setValue(thisRef: Fragment, property: KProperty<*>, value: T?) {
            this.field = value
        }
    }