// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

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

// DwebMetaData contains all meta data concerning the Dweb contract.
var DwebMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"boostrap\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"identity\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"identity\",\"type\":\"string\"}],\"name\":\"initial\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"link\",\"type\":\"string\"}],\"name\":\"join\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newIdent\",\"type\":\"string\"}],\"name\":\"update\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DwebABI is the input ABI used to generate the binding from.
// Deprecated: Use DwebMetaData.ABI instead.
var DwebABI = DwebMetaData.ABI

// Dweb is an auto generated Go binding around an Ethereum contract.
type Dweb struct {
	DwebCaller     // Read-only binding to the contract
	DwebTransactor // Write-only binding to the contract
	DwebFilterer   // Log filterer for contract events
}

// DwebCaller is an auto generated read-only Go binding around an Ethereum contract.
type DwebCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DwebTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DwebTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DwebFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DwebFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DwebSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DwebSession struct {
	Contract     *Dweb             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DwebCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DwebCallerSession struct {
	Contract *DwebCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DwebTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DwebTransactorSession struct {
	Contract     *DwebTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DwebRaw is an auto generated low-level Go binding around an Ethereum contract.
type DwebRaw struct {
	Contract *Dweb // Generic contract binding to access the raw methods on
}

// DwebCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DwebCallerRaw struct {
	Contract *DwebCaller // Generic read-only contract binding to access the raw methods on
}

// DwebTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DwebTransactorRaw struct {
	Contract *DwebTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDweb creates a new instance of Dweb, bound to a specific deployed contract.
func NewDweb(address common.Address, backend bind.ContractBackend) (*Dweb, error) {
	contract, err := bindDweb(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Dweb{DwebCaller: DwebCaller{contract: contract}, DwebTransactor: DwebTransactor{contract: contract}, DwebFilterer: DwebFilterer{contract: contract}}, nil
}

// NewDwebCaller creates a new read-only instance of Dweb, bound to a specific deployed contract.
func NewDwebCaller(address common.Address, caller bind.ContractCaller) (*DwebCaller, error) {
	contract, err := bindDweb(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DwebCaller{contract: contract}, nil
}

// NewDwebTransactor creates a new write-only instance of Dweb, bound to a specific deployed contract.
func NewDwebTransactor(address common.Address, transactor bind.ContractTransactor) (*DwebTransactor, error) {
	contract, err := bindDweb(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DwebTransactor{contract: contract}, nil
}

// NewDwebFilterer creates a new log filterer instance of Dweb, bound to a specific deployed contract.
func NewDwebFilterer(address common.Address, filterer bind.ContractFilterer) (*DwebFilterer, error) {
	contract, err := bindDweb(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DwebFilterer{contract: contract}, nil
}

// bindDweb binds a generic wrapper to an already deployed contract.
func bindDweb(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DwebMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dweb *DwebRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dweb.Contract.DwebCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dweb *DwebRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dweb.Contract.DwebTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dweb *DwebRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dweb.Contract.DwebTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dweb *DwebCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dweb.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dweb *DwebTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dweb.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dweb *DwebTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dweb.Contract.contract.Transact(opts, method, params...)
}

// Boostrap is a free data retrieval call binding the contract method 0xbb332d2c.
//
// Solidity: function boostrap() view returns(string)
func (_Dweb *DwebCaller) Boostrap(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Dweb.contract.Call(opts, &out, "boostrap")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Boostrap is a free data retrieval call binding the contract method 0xbb332d2c.
//
// Solidity: function boostrap() view returns(string)
func (_Dweb *DwebSession) Boostrap() (string, error) {
	return _Dweb.Contract.Boostrap(&_Dweb.CallOpts)
}

// Boostrap is a free data retrieval call binding the contract method 0xbb332d2c.
//
// Solidity: function boostrap() view returns(string)
func (_Dweb *DwebCallerSession) Boostrap() (string, error) {
	return _Dweb.Contract.Boostrap(&_Dweb.CallOpts)
}

// Identity is a free data retrieval call binding the contract method 0x2c159a1a.
//
// Solidity: function identity() view returns(string)
func (_Dweb *DwebCaller) Identity(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Dweb.contract.Call(opts, &out, "identity")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Identity is a free data retrieval call binding the contract method 0x2c159a1a.
//
// Solidity: function identity() view returns(string)
func (_Dweb *DwebSession) Identity() (string, error) {
	return _Dweb.Contract.Identity(&_Dweb.CallOpts)
}

// Identity is a free data retrieval call binding the contract method 0x2c159a1a.
//
// Solidity: function identity() view returns(string)
func (_Dweb *DwebCallerSession) Identity() (string, error) {
	return _Dweb.Contract.Identity(&_Dweb.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Dweb *DwebCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Dweb.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Dweb *DwebSession) Owner() (common.Address, error) {
	return _Dweb.Contract.Owner(&_Dweb.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Dweb *DwebCallerSession) Owner() (common.Address, error) {
	return _Dweb.Contract.Owner(&_Dweb.CallOpts)
}

// Initial is a paid mutator transaction binding the contract method 0xb698b135.
//
// Solidity: function initial(string identity) returns(bool)
func (_Dweb *DwebTransactor) Initial(opts *bind.TransactOpts, identity string) (*types.Transaction, error) {
	return _Dweb.contract.Transact(opts, "initial", identity)
}

// Initial is a paid mutator transaction binding the contract method 0xb698b135.
//
// Solidity: function initial(string identity) returns(bool)
func (_Dweb *DwebSession) Initial(identity string) (*types.Transaction, error) {
	return _Dweb.Contract.Initial(&_Dweb.TransactOpts, identity)
}

// Initial is a paid mutator transaction binding the contract method 0xb698b135.
//
// Solidity: function initial(string identity) returns(bool)
func (_Dweb *DwebTransactorSession) Initial(identity string) (*types.Transaction, error) {
	return _Dweb.Contract.Initial(&_Dweb.TransactOpts, identity)
}

// Join is a paid mutator transaction binding the contract method 0x6a786b07.
//
// Solidity: function join(string link) returns(bool)
func (_Dweb *DwebTransactor) Join(opts *bind.TransactOpts, link string) (*types.Transaction, error) {
	return _Dweb.contract.Transact(opts, "join", link)
}

// Join is a paid mutator transaction binding the contract method 0x6a786b07.
//
// Solidity: function join(string link) returns(bool)
func (_Dweb *DwebSession) Join(link string) (*types.Transaction, error) {
	return _Dweb.Contract.Join(&_Dweb.TransactOpts, link)
}

// Join is a paid mutator transaction binding the contract method 0x6a786b07.
//
// Solidity: function join(string link) returns(bool)
func (_Dweb *DwebTransactorSession) Join(link string) (*types.Transaction, error) {
	return _Dweb.Contract.Join(&_Dweb.TransactOpts, link)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Dweb *DwebTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dweb.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Dweb *DwebSession) RenounceOwnership() (*types.Transaction, error) {
	return _Dweb.Contract.RenounceOwnership(&_Dweb.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Dweb *DwebTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Dweb.Contract.RenounceOwnership(&_Dweb.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Dweb *DwebTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Dweb.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Dweb *DwebSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Dweb.Contract.TransferOwnership(&_Dweb.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Dweb *DwebTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Dweb.Contract.TransferOwnership(&_Dweb.TransactOpts, newOwner)
}

// Update is a paid mutator transaction binding the contract method 0x3d7403a3.
//
// Solidity: function update(string newIdent) returns(bool)
func (_Dweb *DwebTransactor) Update(opts *bind.TransactOpts, newIdent string) (*types.Transaction, error) {
	return _Dweb.contract.Transact(opts, "update", newIdent)
}

// Update is a paid mutator transaction binding the contract method 0x3d7403a3.
//
// Solidity: function update(string newIdent) returns(bool)
func (_Dweb *DwebSession) Update(newIdent string) (*types.Transaction, error) {
	return _Dweb.Contract.Update(&_Dweb.TransactOpts, newIdent)
}

// Update is a paid mutator transaction binding the contract method 0x3d7403a3.
//
// Solidity: function update(string newIdent) returns(bool)
func (_Dweb *DwebTransactorSession) Update(newIdent string) (*types.Transaction, error) {
	return _Dweb.Contract.Update(&_Dweb.TransactOpts, newIdent)
}

// DwebOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Dweb contract.
type DwebOwnershipTransferredIterator struct {
	Event *DwebOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DwebOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DwebOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DwebOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DwebOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DwebOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DwebOwnershipTransferred represents a OwnershipTransferred event raised by the Dweb contract.
type DwebOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Dweb *DwebFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DwebOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Dweb.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DwebOwnershipTransferredIterator{contract: _Dweb.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Dweb *DwebFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DwebOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Dweb.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DwebOwnershipTransferred)
				if err := _Dweb.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Dweb *DwebFilterer) ParseOwnershipTransferred(log types.Log) (*DwebOwnershipTransferred, error) {
	event := new(DwebOwnershipTransferred)
	if err := _Dweb.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
