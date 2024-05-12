package com.zk.android.app.utils

import org.web3j.protocol.core.RemoteCall
import org.web3j.protocol.core.methods.response.TransactionReceipt
import org.web3j.tx.response.TransactionReceiptProcessor
import java.math.BigInteger

fun <T> RemoteCall<T>.sendSafe(): Result<T> {
    return runCatching { send() }
}

fun TransactionReceiptProcessor.waitForTransactionReceiptSafe(transactionHash: String): Result<TransactionReceipt> {
    return runCatching { waitForTransactionReceipt(transactionHash) }
}

fun BigInteger.weiToEther(): Double {
    val etherInWei = BigInteger.TEN.pow(18) // 1 ether = 10^18 wei
    return toDouble() / etherInWei.toDouble()
}

fun Float.toWei(): BigInteger {
    // Ethereum wei is the smallest denomination of Ether (1 Ether = 10^18 wei)
    val weiPerEther = BigInteger.TEN.pow(18)

    // Convert the Float value to BigInteger by multiplying it by the wei per Ether
    return (this * weiPerEther.toFloat()).toBigDecimal().toBigInteger()
}