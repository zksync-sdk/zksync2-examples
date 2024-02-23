package org.example.deploy.create.constructor;

import io.zksync.protocol.ZkSync;
import io.zksync.protocol.account.Wallet;
import org.example.contracts.CounterContract;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;
import org.web3j.utils.Numeric;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class DeployCreateConstrucotor {
    public static void main(String[] args) {
        String L1_NODE = "https://rpc.ankr.com/eth_sepolia";
        String L2_NODE = "https://sepolia.era.zksync.dev";
        final String privateKey = "PRIVATE_KEY";

        Web3j l1Web3 = Web3j.build(new HttpService(L1_NODE));
        ZkSync zksync = ZkSync.build(new HttpService(L2_NODE));
        Credentials credentials = Credentials.create(privateKey);

        Wallet wallet = new Wallet(l1Web3, zksync, credentials);

        BigInteger incrementer = BigInteger.TWO;
        List<Type> inputParameter = new ArrayList<>();
        inputParameter.add(new Uint256(incrementer));
        String calldata = FunctionEncoder.encodeConstructor(inputParameter);

        TransactionReceipt result = wallet.deploy(Numeric.hexStringToByteArray(CounterContract.BINARY), Numeric.hexStringToByteArray(calldata)).sendAsync().join();
        System.out.println("Contract address: " + result.getContractAddress());
    }
}
