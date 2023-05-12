# How to Compile Solidity Smart Contracts

Use `zksolc` compiler to compile Solidity smart contracts.
`zksolc` compiler requires `solc` to be installed. Specific version of
`zksolc` compiler is compatible with specific versions `solc` so make
sure to make correct versions of your compilers.

There are 3 solidity smart contracts:

- `Storage`: contract without constructor.
- `Incrementer`: contract with constructor.
- `Demo`: contract that has dependency to `Foo` contract.

In the following examples `Docker` is used to create containers with already
`solc` installed.

## Compile Smart Contracts

Run the container has `solc` tool already installed:
```shell
# create container with installed solc tool
SOLC_VERSION="0.8.19-alpine"
docker create -it --name zksolc --entrypoint ash  ethereum/solc:${SOLC_VERSION}

# copy smart contracts source files to container
docker cp solidity zksolc:/solidity

# run and attach to the container
docker start -i zksolc
```
Run commands in container:
```shell
# download zksolc
ZKSOLC_VERSION="v1.3.9"
wget https://github.com/matter-labs/zksolc-bin/raw/main/linux-amd64/zksolc-linux-amd64-musl-${ZKSOLC_VERSION} -O /bin/zksolc; chmod +x /bin/zksolc
```

**Compile Storage Smart Contract**
```shell
# compile smart contract
zksolc --bin -O3 \
	-o /solidity/storage/build  \
	/solidity/storage/Incrementer.sol


# create combined-json with abi and binary
zksolc -O3 -o /solidity/storage/build \
  --combined-json abi,bin \
  /solidity/storage/Storage.sol
```

**Compile Incrementer Smart Contract**
```shell
# compile smart contract
zksolc --bin -O3 \
	-o /solidity/incrementer/build  \
	/solidity/incrementer/Incrementer.sol

# create combined-json with abi and binary
zksolc -O3 -o /solidity/incrementer/build \
  --combined-json abi,bin \
  /solidity/incrementer/Incrementer.sol
```

**Compile Demo Smart Contract**
```shell
# compile smart contract
zksolc --bin -O3 \
	-o /solidity/demo/build/  \
	/solidity/demo/Demo.sol \
	/solidity/demo/Foo.sol

# create combined-json with abi
zksolc -O3 -o /solidity/demo/build \
  --combined-json abi \
  /solidity/demo/Demo.sol \
  /solidity/demo/Foo.sol
```
Exit from container
```shell
exit 
```

Copy generated files from container to host machine
```shell
# copy generated files from container to host
docker cp zksolc:/solidity .

# remove container
docker rm zksolc
```

On host machine, for each smart contract there in `build` folder there are binaries
and `combinned.json` files.

## Generate bindings

Next step is to use `abigen` tool along with `combined.json` file to generate
smart contract bindings for golang.

```shell
# create container with installed abigen tool
docker create -it --name abigen \
	--entrypoint ash \
	ethereum/client-go:alltools-v1.11.6

# copy required combine.json configuration files
docker cp solidity abigen:/solidity

# run and attach to the container
docker start -i abigen
```
Run commands in the container

**Generate ABI for Storage Smart Contract**
```shell
# generate bindings based on ABI in combined.json file
abigen \
	--combined-json /solidity/storage/build/combined.json \
	--out /solidity/storage/build/storage.go \
	--pkg storage
```

**Generate ABI for Incrementer Smart Contract**
```shell
# generate bindings based on ABI in combined.json file
abigen \
	--combined-json /solidity/incrementer/build/combined.json \
	--out /solidity/incrementer/build/demo.go \
	--pkg incrementer
```

**Generate ABI for Demo Smart Contract**
```shell
# generate bindings based on ABI in combined.json file
abigen \
	--combined-json /solidity/demo/build/combined.json \
	--out /solidity/demo/build/demo.go \
	--pkg demo
```

Exit from container
```shell
exit 
```

Copy generated smart contract bindings from container to host:
```shell
# copy generated bindings from container to host
docker cp abigen:/solidity .

# remove the container
docker rm abigen
```

For smart contract there is go file in `build` that can be used to
interact with smart contract.