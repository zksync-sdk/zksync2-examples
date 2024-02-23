package org.example.deploy.create;

import io.zksync.protocol.ZkSync;
import io.zksync.protocol.account.Wallet;
import org.example.contracts.CounterContract;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;
import org.web3j.utils.Numeric;

import java.math.BigInteger;

public class DeployCreate {
    public static void main(String[] args) {
        String L1_NODE = "https://rpc.ankr.com/eth_sepolia";
        String L2_NODE = "https://sepolia.era.zksync.dev";
        final String privateKey = "PRIVATE_KEY";

        Web3j l1Web3 = Web3j.build(new HttpService(L1_NODE));
        ZkSync zksync = ZkSync.build(new HttpService(L2_NODE));
        Credentials credentials = Credentials.create(privateKey);

        Wallet wallet = new Wallet(l1Web3, zksync, credentials);

        TransactionReceipt result = wallet.deploy(Numeric.hexStringToByteArray(CounterContract.BINARY)).sendAsync().join();
        System.out.println("Contract address: " + result.getContractAddress());
    }
}
