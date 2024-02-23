package org.example.paymaster;

import io.zksync.abi.TransactionEncoder;
import io.zksync.crypto.signer.PrivateKeyEthSigner;
import io.zksync.methods.request.Eip712Meta;
import io.zksync.methods.request.PaymasterParams;
import io.zksync.methods.request.Transaction;
import io.zksync.protocol.ZkSync;
import io.zksync.protocol.account.Wallet;
import io.zksync.protocol.core.Token;
import io.zksync.transaction.fee.DefaultTransactionFeeProvider;
import io.zksync.transaction.type.Transaction712;
import io.zksync.transaction.type.TransferTransaction;
import io.zksync.transaction.type.WithdrawTransaction;
import io.zksync.utils.Paymaster;
import io.zksync.utils.ZkSyncAddresses;
import io.zksync.wrappers.IEthToken;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.EthSendTransaction;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;
import org.web3j.protocol.http.HttpService;
import org.web3j.utils.Numeric;

import java.io.IOException;
import java.lang.reflect.Array;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

public class paymaster {
    protected static final String PAYMASTER = "0x13D0D8550769f59aa241a41897D4859c87f7Dd46";
    protected static final String TOKEN = "0x927488F48ffbc32112F1fF721759649A89721F8F";

    public static void main(String[] args) throws Exception {
        String L1_NODE = "https://rpc.ankr.com/eth_sepolia";
        String L2_NODE = "https://sepolia.era.zksync.dev";
        final String privateKey = "PRIVATE_KEY";


        Web3j l1Web3 = Web3j.build(new HttpService(L1_NODE));
        ZkSync zksync = ZkSync.build(new HttpService(L2_NODE));
        DefaultTransactionFeeProvider feeProvider = new DefaultTransactionFeeProvider(zksync, Token.ETH);

        Credentials credentials = Credentials.create(privateKey);
        PrivateKeyEthSigner signer = new PrivateKeyEthSigner(credentials, zksync.ethChainId().send().getChainId().longValue());

        Wallet wallet = new Wallet(l1Web3, zksync, credentials);

        Function function = new Function(
                IEthToken.FUNC_MINT,
                Arrays.asList(new Address(160, credentials.getAddress()), new Uint256(BigInteger.TWO)),
                Collections.emptyList()
        );
        String calldata = FunctionEncoder.encode(function);

        PaymasterParams paymasterParams = new PaymasterParams(PAYMASTER, Numeric.hexStringToByteArray(FunctionEncoder.encode(Paymaster.encodeApprovalBased(TOKEN, BigInteger.ONE, new byte[] {}))));
        Eip712Meta meta = new Eip712Meta(BigInteger.valueOf(50000), null, null, paymasterParams);
        Transaction estimate = new Transaction(
                wallet.getAddress(),
                TOKEN,
                BigInteger.ZERO,
                zksync.ethGasPrice().send().getGasPrice(),
                BigInteger.ZERO,
                calldata,
                meta
        );
        BigInteger gasLimit = feeProvider.getGasLimit(estimate);

        meta = new Eip712Meta(BigInteger.valueOf(50000), null, null, paymasterParams);
        Transaction712 prepared = new Transaction712(
                zksync.ethChainId().send().getChainId().longValue(),
                wallet.getNonce().send(),
                gasLimit,
                TOKEN,
                BigInteger.ZERO,
                calldata,
                BigInteger.valueOf(100000000L),
                zksync.ethGasPrice().send().getGasPrice(),
                credentials.getAddress(),
                meta
        );

        String signature = signer.getDomain().thenCompose(domain -> signer.signTypedData(domain, prepared)).join();
        byte[] signed = TransactionEncoder.encode(prepared, TransactionEncoder.getSignatureData(signature));

        EthSendTransaction receipt = zksync.ethSendRawTransaction(Numeric.toHexString(signed)).sendAsync().join();
        System.out.println(receipt);
    }
}
