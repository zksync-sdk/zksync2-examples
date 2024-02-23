package org.example.balance;

import io.zksync.protocol.ZkSync;
import io.zksync.protocol.account.Wallet;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.http.HttpService;

import java.math.BigInteger;

public class CheckBalance {
    public static void main(String[] args) throws Exception {
        String L1_NODE = "https://rpc.ankr.com/eth_sepolia";
        String L2_NODE = "https://sepolia.era.zksync.dev";
        final String privateKey = "PRIVATE_KEY";

        Web3j l1Web3 = Web3j.build(new HttpService(L1_NODE));
        ZkSync zksync = ZkSync.build(new HttpService(L2_NODE));
        Credentials credentials = Credentials.create(privateKey);

        Wallet wallet = new Wallet(l1Web3, zksync, credentials);

        BigInteger l1Balance = wallet.getBalanceL1().send();
        System.out.println("L1 balance:" + l1Balance);

        BigInteger l2Balance = wallet.getBalance().send();
        System.out.println("L2 balance:" + l2Balance);

        String l1DAI = "0x70a0F165d6f8054d0d0CF8dFd4DD2005f0AF6B55";
        String l2DAI = wallet.l2TokenAddress(l1DAI);
        BigInteger l2Erc20Balance = wallet.getBalance(l2DAI).send();
        System.out.println("L2 erc20 balance:" + l2Erc20Balance);
    }
}
