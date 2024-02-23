package org.example.withdraw;

import io.zksync.protocol.ZkSync;
import io.zksync.protocol.account.Wallet;
import io.zksync.transaction.type.WithdrawTransaction;
import io.zksync.utils.ZkSyncAddresses;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import org.web3j.protocol.http.HttpService;

import java.io.IOException;
import java.math.BigInteger;

public class Withdraw {
    public static void main(String[] args) throws TransactionException, IOException {
        String L1_NODE = "https://rpc.ankr.com/eth_sepolia";
        String L2_NODE = "https://sepolia.era.zksync.dev";
        final String privateKey = "PRIVATE_KEY";

        Web3j l1Web3 = Web3j.build(new HttpService(L1_NODE));
        ZkSync zksync = ZkSync.build(new HttpService(L2_NODE));
        Credentials credentials = Credentials.create(privateKey);

        Wallet wallet = new Wallet(l1Web3, zksync, credentials);

        //ETH withdraw

        BigInteger amount = BigInteger.valueOf(7_000_000_000L);

        BigInteger balanceBefore = wallet.getBalance().sendAsync().join();
        System.out.println("Balance before: " + balanceBefore);


        WithdrawTransaction transaction = new WithdrawTransaction(ZkSyncAddresses.ETH_ADDRESS, amount, wallet.getAddress());
        TransactionReceipt result = wallet.withdraw(transaction).sendAsync().join();
        TransactionReceipt receipt = wallet.getTransactionReceiptProcessor().waitForTransactionReceipt(result.getTransactionHash());

        BigInteger balanceAfter = wallet.getBalance().sendAsync().join();
        System.out.println("Balance after: " + balanceAfter);

        //ERC20 withdraw

        BigInteger erc20Amount = BigInteger.valueOf(5);
        String l1DAI = "0x70a0F165d6f8054d0d0CF8dFd4DD2005f0AF6B55";
        String l2DAI = wallet.l2TokenAddress(l1DAI);

        BigInteger erc20BalanceBefore = wallet.getBalance(l2DAI).sendAsync().join();
        System.out.println("Balance before: " + erc20BalanceBefore);


        WithdrawTransaction erc20Transaction = new WithdrawTransaction(l2DAI, erc20Amount, wallet.getAddress());
        TransactionReceipt erc20Result = wallet.withdraw(erc20Transaction).sendAsync().join();
        TransactionReceipt erc20Receipt = wallet.getTransactionReceiptProcessor().waitForTransactionReceipt(erc20Result.getTransactionHash());

        BigInteger erc20BalanceAfter = wallet.getBalance(l2DAI).sendAsync().join();
        System.out.println("Balance after: " + erc20BalanceAfter);

    }
}
