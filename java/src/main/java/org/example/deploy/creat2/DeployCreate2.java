package org.example.deploy.creat2;

import io.zksync.abi.TransactionEncoder;
import io.zksync.crypto.signer.PrivateKeyEthSigner;
import io.zksync.methods.request.Eip712Meta;
import io.zksync.methods.response.ZksEstimateFee;
import io.zksync.protocol.ZkSync;
import io.zksync.protocol.account.Wallet;
import io.zksync.protocol.core.ZkBlockParameterName;
import io.zksync.transaction.fee.Fee;
import io.zksync.transaction.response.ZkSyncTransactionReceiptProcessor;
import io.zksync.transaction.type.Transaction712;
import io.zksync.utils.ContractDeployer;
import org.example.contracts.CounterContract;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.datatypes.Address;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.request.Transaction;
import org.web3j.protocol.core.methods.response.EthCall;
import org.web3j.protocol.core.methods.response.EthGasPrice;
import org.web3j.protocol.core.methods.response.EthSendTransaction;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import org.web3j.protocol.http.HttpService;
import org.web3j.utils.Numeric;

import java.io.IOException;
import java.math.BigInteger;
import java.security.SecureRandom;

public class DeployCreate2 {
    public static void main(String[] args) throws IOException, TransactionException {
        String L1_NODE = "https://rpc.ankr.com/eth_sepolia";
        String L2_NODE = "https://sepolia.era.zksync.dev";
        final String privateKey = "PRIVATE_KEY";

        Web3j l1Web3 = Web3j.build(new HttpService(L1_NODE));
        ZkSync zksync = ZkSync.build(new HttpService(L2_NODE));
        Credentials credentials = Credentials.create(privateKey);
        PrivateKeyEthSigner signer = new PrivateKeyEthSigner(credentials, zksync.ethChainId().send().getChainId().longValue());
        ZkSyncTransactionReceiptProcessor processor = new ZkSyncTransactionReceiptProcessor(zksync, 200, 100);

        BigInteger nonce = zksync
                .ethGetTransactionCount(credentials.getAddress(), DefaultBlockParameterName.PENDING).send()
                .getTransactionCount();

        byte[] salt = SecureRandom.getSeed(32);

        String precomputedAddress = ContractDeployer.computeL2Create2Address(new Address(credentials.getAddress()), Numeric.hexStringToByteArray(CounterContract.BINARY), new byte[]{}, salt).getValue();

        io.zksync.methods.request.Transaction estimate = io.zksync.methods.request.Transaction.create2ContractTransaction(
                credentials.getAddress(),
                BigInteger.ZERO,
                BigInteger.ZERO,
                CounterContract.BINARY,
                "0x",
                salt
        );

        ZksEstimateFee estimateFee = zksync.zksEstimateFee(estimate).send();

        EthGasPrice gasPrice = zksync.ethGasPrice().send();

        Fee fee = estimateFee.getResult();

        Eip712Meta meta = estimate.getEip712Meta();
        meta.setGasPerPubdata(fee.getGasPerPubdataLimitNumber());

        Transaction712 transaction = new Transaction712(
                zksync.ethChainId().send().getChainId().longValue(),
                nonce,
                fee.getGasLimitNumber(),
                estimate.getTo(),
                estimate.getValueNumber(),
                estimate.getData(),
                fee.getMaxPriorityFeePerErgNumber(),
                fee.getGasPriceLimitNumber(),
                credentials.getAddress(),
                meta
        );

        String signature = signer.getDomain().thenCompose(domain -> signer.signTypedData(domain, transaction)).join();
        byte[] message = TransactionEncoder.encode(transaction, TransactionEncoder.getSignatureData(signature));

        EthSendTransaction sent = zksync.ethSendRawTransaction(Numeric.toHexString(message)).send();

        TransactionReceipt receipt = processor.waitForTransactionReceipt(sent.getResult());


        String contractAddress = receipt.getContractAddress();
        System.out.println("Deployed `CounterContract as: `" + contractAddress);

        Transaction call = Transaction.createEthCallTransaction(
                credentials.getAddress(),
                contractAddress,
                FunctionEncoder.encode(CounterContract.encodeGet())
        );

        EthCall ethCall = zksync.ethCall(call, ZkBlockParameterName.COMMITTED).send();
    }
}
