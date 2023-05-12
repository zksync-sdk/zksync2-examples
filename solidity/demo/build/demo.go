// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package demo

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DemoMetaData contains all meta data concerning the Demo contract.
var DemoMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"foo\",\"outputs\":[{\"internalType\":\"contractFoo\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFooName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DemoABI is the input ABI used to generate the binding from.
// Deprecated: Use DemoMetaData.ABI instead.
var DemoABI = DemoMetaData.ABI

// Demo is an auto generated Go binding around an Ethereum contract.
type Demo struct {
	DemoCaller     // Read-only binding to the contract
	DemoTransactor // Write-only binding to the contract
	DemoFilterer   // Log filterer for contract events
}

// DemoCaller is an auto generated read-only Go binding around an Ethereum contract.
type DemoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DemoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DemoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DemoSession struct {
	Contract     *Demo             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DemoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DemoCallerSession struct {
	Contract *DemoCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DemoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DemoTransactorSession struct {
	Contract     *DemoTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DemoRaw is an auto generated low-level Go binding around an Ethereum contract.
type DemoRaw struct {
	Contract *Demo // Generic contract binding to access the raw methods on
}

// DemoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DemoCallerRaw struct {
	Contract *DemoCaller // Generic read-only contract binding to access the raw methods on
}

// DemoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DemoTransactorRaw struct {
	Contract *DemoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDemo creates a new instance of Demo, bound to a specific deployed contract.
func NewDemo(address common.Address, backend bind.ContractBackend) (*Demo, error) {
	contract, err := bindDemo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Demo{DemoCaller: DemoCaller{contract: contract}, DemoTransactor: DemoTransactor{contract: contract}, DemoFilterer: DemoFilterer{contract: contract}}, nil
}

// NewDemoCaller creates a new read-only instance of Demo, bound to a specific deployed contract.
func NewDemoCaller(address common.Address, caller bind.ContractCaller) (*DemoCaller, error) {
	contract, err := bindDemo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DemoCaller{contract: contract}, nil
}

// NewDemoTransactor creates a new write-only instance of Demo, bound to a specific deployed contract.
func NewDemoTransactor(address common.Address, transactor bind.ContractTransactor) (*DemoTransactor, error) {
	contract, err := bindDemo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DemoTransactor{contract: contract}, nil
}

// NewDemoFilterer creates a new log filterer instance of Demo, bound to a specific deployed contract.
func NewDemoFilterer(address common.Address, filterer bind.ContractFilterer) (*DemoFilterer, error) {
	contract, err := bindDemo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DemoFilterer{contract: contract}, nil
}

// bindDemo binds a generic wrapper to an already deployed contract.
func bindDemo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DemoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Demo *DemoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Demo.Contract.DemoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Demo *DemoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Demo.Contract.DemoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Demo *DemoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Demo.Contract.DemoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Demo *DemoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Demo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Demo *DemoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Demo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Demo *DemoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Demo.Contract.contract.Transact(opts, method, params...)
}

// Foo is a free data retrieval call binding the contract method 0xc2985578.
//
// Solidity: function foo() view returns(address)
func (_Demo *DemoCaller) Foo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Demo.contract.Call(opts, &out, "foo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Foo is a free data retrieval call binding the contract method 0xc2985578.
//
// Solidity: function foo() view returns(address)
func (_Demo *DemoSession) Foo() (common.Address, error) {
	return _Demo.Contract.Foo(&_Demo.CallOpts)
}

// Foo is a free data retrieval call binding the contract method 0xc2985578.
//
// Solidity: function foo() view returns(address)
func (_Demo *DemoCallerSession) Foo() (common.Address, error) {
	return _Demo.Contract.Foo(&_Demo.CallOpts)
}

// GetFooName is a free data retrieval call binding the contract method 0xc6261261.
//
// Solidity: function getFooName() view returns(string)
func (_Demo *DemoCaller) GetFooName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Demo.contract.Call(opts, &out, "getFooName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetFooName is a free data retrieval call binding the contract method 0xc6261261.
//
// Solidity: function getFooName() view returns(string)
func (_Demo *DemoSession) GetFooName() (string, error) {
	return _Demo.Contract.GetFooName(&_Demo.CallOpts)
}

// GetFooName is a free data retrieval call binding the contract method 0xc6261261.
//
// Solidity: function getFooName() view returns(string)
func (_Demo *DemoCallerSession) GetFooName() (string, error) {
	return _Demo.Contract.GetFooName(&_Demo.CallOpts)
}

// FooMetaData contains all meta data concerning the Foo contract.
var FooMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// FooABI is the input ABI used to generate the binding from.
// Deprecated: Use FooMetaData.ABI instead.
var FooABI = FooMetaData.ABI

// Foo is an auto generated Go binding around an Ethereum contract.
type Foo struct {
	FooCaller     // Read-only binding to the contract
	FooTransactor // Write-only binding to the contract
	FooFilterer   // Log filterer for contract events
}

// FooCaller is an auto generated read-only Go binding around an Ethereum contract.
type FooCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FooTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FooTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FooFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FooFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FooSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FooSession struct {
	Contract     *Foo              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FooCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FooCallerSession struct {
	Contract *FooCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FooTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FooTransactorSession struct {
	Contract     *FooTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FooRaw is an auto generated low-level Go binding around an Ethereum contract.
type FooRaw struct {
	Contract *Foo // Generic contract binding to access the raw methods on
}

// FooCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FooCallerRaw struct {
	Contract *FooCaller // Generic read-only contract binding to access the raw methods on
}

// FooTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FooTransactorRaw struct {
	Contract *FooTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFoo creates a new instance of Foo, bound to a specific deployed contract.
func NewFoo(address common.Address, backend bind.ContractBackend) (*Foo, error) {
	contract, err := bindFoo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Foo{FooCaller: FooCaller{contract: contract}, FooTransactor: FooTransactor{contract: contract}, FooFilterer: FooFilterer{contract: contract}}, nil
}

// NewFooCaller creates a new read-only instance of Foo, bound to a specific deployed contract.
func NewFooCaller(address common.Address, caller bind.ContractCaller) (*FooCaller, error) {
	contract, err := bindFoo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FooCaller{contract: contract}, nil
}

// NewFooTransactor creates a new write-only instance of Foo, bound to a specific deployed contract.
func NewFooTransactor(address common.Address, transactor bind.ContractTransactor) (*FooTransactor, error) {
	contract, err := bindFoo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FooTransactor{contract: contract}, nil
}

// NewFooFilterer creates a new log filterer instance of Foo, bound to a specific deployed contract.
func NewFooFilterer(address common.Address, filterer bind.ContractFilterer) (*FooFilterer, error) {
	contract, err := bindFoo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FooFilterer{contract: contract}, nil
}

// bindFoo binds a generic wrapper to an already deployed contract.
func bindFoo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FooMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Foo *FooRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Foo.Contract.FooCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Foo *FooRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Foo.Contract.FooTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Foo *FooRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Foo.Contract.FooTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Foo *FooCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Foo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Foo *FooTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Foo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Foo *FooTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Foo.Contract.contract.Transact(opts, method, params...)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Foo *FooCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Foo.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Foo *FooSession) Name() (string, error) {
	return _Foo.Contract.Name(&_Foo.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Foo *FooCallerSession) Name() (string, error) {
	return _Foo.Contract.Name(&_Foo.CallOpts)
}
