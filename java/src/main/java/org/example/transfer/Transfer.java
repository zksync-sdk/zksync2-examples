package org.example.transfer;

import io.zksync.protocol.ZkSync;
import io.zksync.protocol.account.Wallet;
import io.zksync.transaction.type.TransferTransaction;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;

import java.math.BigInteger;

public class Transfer {
    public static void main(String[] args) throws Exception {
        String L1_NODE = "https://rpc.ankr.com/eth_sepolia";
        String L2_NODE = "https://sepolia.era.zksync.dev";
        final String privateKey = "PRIVATE_KEY";

        BigInteger amount = BigInteger.valueOf(7000000000L);

        Web3j l1Web3 = Web3j.build(new HttpService(L1_NODE));
        ZkSync zksync = ZkSync.build(new HttpService(L2_NODE));
        Credentials credentials = Credentials.create(privateKey);

        Wallet wallet = new Wallet(l1Web3, zksync, credentials);

        // Transfer ETH

        BigInteger balance_before = wallet.getBalance().send();
        System.out.println("Balance before: " + balance_before);

        TransferTransaction transaction = new TransferTransaction("0x7221A759fd029b5B12792bb690839A1292BB21F8", amount, credentials.getAddress());
        TransactionReceipt receipt = wallet.transfer(transaction).send();
        wallet.getTransactionReceiptProcessor().waitForTransactionReceipt(receipt.getTransactionHash());

        BigInteger balance_after = wallet.getBalance().send();
        System.out.println("Balance after: " + balance_after);

        // Transfer ERC20

        BigInteger erc20Amount = BigInteger.valueOf(5);
        String l1DAI = "0x70a0F165d6f8054d0d0CF8dFd4DD2005f0AF6B55";
        String l2DAI = wallet.l2TokenAddress(l1DAI);

        BigInteger erc_balance_before = wallet.getBalance(l2DAI).sendAsync().join();
        System.out.println("ERC20 balance before: " + erc_balance_before);

        TransferTransaction erc_transaction = new TransferTransaction("0xa61464658AfeAf65CccaaFD3a512b69A83B77618", erc20Amount, credentials.getAddress(), l2DAI);
        TransactionReceipt erc_receipt = wallet.transfer(erc_transaction).sendAsync().join();
        wallet.getTransactionReceiptProcessor().waitForTransactionReceipt(erc_receipt.getTransactionHash());

        BigInteger erc_balance_after = wallet.getBalance(l2DAI).send();
        System.out.println("ERC20 balance after: " + erc_balance_after);
    }
}
