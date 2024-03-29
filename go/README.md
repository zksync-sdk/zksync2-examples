# zkSync Era Go Examples

Examples are made to demonstrate how [zksync2-go](https://github.com/zksync-sdk/zksync2-go) 
SDK can be used for development. The examples demonstrate how to:

1. Deposit ETH from Ethereum into zkSync.
2. Transfer ETH on zkSync. 
3. Withdraw ETH from zkSync to Ethereum.
4. Deploy a smart contract using create opcode.
5. Deploy a smart contract with constructor using create method.
6. Deploy a smart contract with dependency using create method.
7. Deploy a smart contract using create2 opcode.
8. Deploy a smart contract with constructor using create2 method.
9. Deploy a smart contract with dependency using create2 method.
10. Deposit token from Ethereum into zkSync.
11. Transfer token on zkSync.
12. Withdraw token from zkSync to Ethereum.
13. Deploy custom token on zkSync.
14. Deploy token using create method.
15. Deploy token using create2 method.
16. Deploy smart account using create method.
17. Deploy smart account using create2 method.
18. Use paymaster to pay fee with token.

Smart contract deployment use already generated binaries and go bindings.
There is a [user guide](../solidity/README.md) on how those artifacts 
are generated. Same approach can be used to generate required artifact
for other smart contracts.

Run example:

```go
PRIVATE_KEY=<...> go run 02_transfer.go
```