package org.example.deposit;

import io.zksync.protocol.ZkSync;
import io.zksync.protocol.account.Wallet;
import io.zksync.transaction.type.DepositTransaction;
import io.zksync.utils.ZkSyncAddresses;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;

import java.math.BigInteger;

public class Deposit {
    public static void main(String[] args) throws Exception {
        String L1_NODE = "https://rpc.ankr.com/eth_sepolia";
        String L2_NODE = "https://sepolia.era.zksync.dev";
        final String privateKey = "PRIVATE_KEY";
        BigInteger amount = BigInteger.valueOf(7000000000L);

        Web3j l1Web3 = Web3j.build(new HttpService(L1_NODE));
        ZkSync zksync = ZkSync.build(new HttpService(L2_NODE));
        Credentials credentials = Credentials.create(privateKey);

        Wallet wallet = new Wallet(l1Web3, zksync, credentials);

        BigInteger balance_before = wallet.getBalanceL1().send();
        System.out.println("Balance before: " + balance_before);

        DepositTransaction transaction = new DepositTransaction(ZkSyncAddresses.ETH_ADDRESS, amount);
        String hash = wallet.deposit(transaction).sendAsync().join().getResult();
        TransactionReceipt l1Receipt = wallet.getTransactionReceiptProcessorL1().waitForTransactionReceipt(hash);

        String l2Hash = zksync.getL2HashFromPriorityOp(l1Receipt, zksync.zksMainContract().sendAsync().join().getResult());
        wallet.getTransactionReceiptProcessor().waitForTransactionReceipt(l2Hash);

        BigInteger balance_after = wallet.getBalanceL1().send();
        System.out.println("Balance before: " + balance_after);

    }
}

