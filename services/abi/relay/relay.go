// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package relay

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// RelayABI is the input ABI used to generate the binding from.
const RelayABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_payload\",\"type\":\"bytes\"}],\"name\":\"commitWithdrawHash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newValidators\",\"type\":\"address[]\"},{\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"updateValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"latestWithdrawHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorInfo\",\"outputs\":[{\"name\":\"index\",\"type\":\"uint8\"},{\"name\":\"power\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_withdrawInfo\",\"type\":\"bytes\"},{\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tokenContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lastValidatorUpdateTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalVotingPower\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_withdrawInfo\",\"type\":\"bytes\"},{\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"withdrawRootHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"reserved\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logicContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"consumedIds\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"latestBlockHeight\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_validators\",\"type\":\"address[]\"},{\"name\":\"_votingPowers\",\"type\":\"uint32[]\"},{\"name\":\"_tokenContract\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// Relay is an auto generated Go binding around an Ethereum contract.
type Relay struct {
	RelayCaller     // Read-only binding to the contract
	RelayTransactor // Write-only binding to the contract
	RelayFilterer   // Log filterer for contract events
}

// RelayCaller is an auto generated read-only Go binding around an Ethereum contract.
type RelayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RelayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RelayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RelaySession struct {
	Contract     *Relay            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RelayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RelayCallerSession struct {
	Contract *RelayCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RelayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RelayTransactorSession struct {
	Contract     *RelayTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RelayRaw is an auto generated low-level Go binding around an Ethereum contract.
type RelayRaw struct {
	Contract *Relay // Generic contract binding to access the raw methods on
}

// RelayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RelayCallerRaw struct {
	Contract *RelayCaller // Generic read-only contract binding to access the raw methods on
}

// RelayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RelayTransactorRaw struct {
	Contract *RelayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRelay creates a new instance of Relay, bound to a specific deployed contract.
func NewRelay(address common.Address, backend bind.ContractBackend) (*Relay, error) {
	contract, err := bindRelay(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Relay{RelayCaller: RelayCaller{contract: contract}, RelayTransactor: RelayTransactor{contract: contract}, RelayFilterer: RelayFilterer{contract: contract}}, nil
}

// NewRelayCaller creates a new read-only instance of Relay, bound to a specific deployed contract.
func NewRelayCaller(address common.Address, caller bind.ContractCaller) (*RelayCaller, error) {
	contract, err := bindRelay(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RelayCaller{contract: contract}, nil
}

// NewRelayTransactor creates a new write-only instance of Relay, bound to a specific deployed contract.
func NewRelayTransactor(address common.Address, transactor bind.ContractTransactor) (*RelayTransactor, error) {
	contract, err := bindRelay(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RelayTransactor{contract: contract}, nil
}

// NewRelayFilterer creates a new log filterer instance of Relay, bound to a specific deployed contract.
func NewRelayFilterer(address common.Address, filterer bind.ContractFilterer) (*RelayFilterer, error) {
	contract, err := bindRelay(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RelayFilterer{contract: contract}, nil
}

// bindRelay binds a generic wrapper to an already deployed contract.
func bindRelay(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RelayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Relay *RelayRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Relay.Contract.RelayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Relay *RelayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Relay.Contract.RelayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Relay *RelayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Relay.Contract.RelayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Relay *RelayCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Relay.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Relay *RelayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Relay.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Relay *RelayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Relay.Contract.contract.Transact(opts, method, params...)
}

// ConsumedIds is a free data retrieval call binding the contract method 0xf21a2116.
//
// Solidity: function consumedIds( bytes32) constant returns(bool)
func (_Relay *RelayCaller) ConsumedIds(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "consumedIds", arg0)
	return *ret0, err
}

// ConsumedIds is a free data retrieval call binding the contract method 0xf21a2116.
//
// Solidity: function consumedIds( bytes32) constant returns(bool)
func (_Relay *RelaySession) ConsumedIds(arg0 [32]byte) (bool, error) {
	return _Relay.Contract.ConsumedIds(&_Relay.CallOpts, arg0)
}

// ConsumedIds is a free data retrieval call binding the contract method 0xf21a2116.
//
// Solidity: function consumedIds( bytes32) constant returns(bool)
func (_Relay *RelayCallerSession) ConsumedIds(arg0 [32]byte) (bool, error) {
	return _Relay.Contract.ConsumedIds(&_Relay.CallOpts, arg0)
}

// LastValidatorUpdateTime is a free data retrieval call binding the contract method 0x568873ad.
//
// Solidity: function lastValidatorUpdateTime() constant returns(uint256)
func (_Relay *RelayCaller) LastValidatorUpdateTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "lastValidatorUpdateTime")
	return *ret0, err
}

// LastValidatorUpdateTime is a free data retrieval call binding the contract method 0x568873ad.
//
// Solidity: function lastValidatorUpdateTime() constant returns(uint256)
func (_Relay *RelaySession) LastValidatorUpdateTime() (*big.Int, error) {
	return _Relay.Contract.LastValidatorUpdateTime(&_Relay.CallOpts)
}

// LastValidatorUpdateTime is a free data retrieval call binding the contract method 0x568873ad.
//
// Solidity: function lastValidatorUpdateTime() constant returns(uint256)
func (_Relay *RelayCallerSession) LastValidatorUpdateTime() (*big.Int, error) {
	return _Relay.Contract.LastValidatorUpdateTime(&_Relay.CallOpts)
}

// LatestBlockHeight is a free data retrieval call binding the contract method 0xf3f39ee5.
//
// Solidity: function latestBlockHeight() constant returns(uint256)
func (_Relay *RelayCaller) LatestBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "latestBlockHeight")
	return *ret0, err
}

// LatestBlockHeight is a free data retrieval call binding the contract method 0xf3f39ee5.
//
// Solidity: function latestBlockHeight() constant returns(uint256)
func (_Relay *RelaySession) LatestBlockHeight() (*big.Int, error) {
	return _Relay.Contract.LatestBlockHeight(&_Relay.CallOpts)
}

// LatestBlockHeight is a free data retrieval call binding the contract method 0xf3f39ee5.
//
// Solidity: function latestBlockHeight() constant returns(uint256)
func (_Relay *RelayCallerSession) LatestBlockHeight() (*big.Int, error) {
	return _Relay.Contract.LatestBlockHeight(&_Relay.CallOpts)
}

// LatestWithdrawHash is a free data retrieval call binding the contract method 0x3cd3f6a7.
//
// Solidity: function latestWithdrawHash() constant returns(bytes32)
func (_Relay *RelayCaller) LatestWithdrawHash(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "latestWithdrawHash")
	return *ret0, err
}

// LatestWithdrawHash is a free data retrieval call binding the contract method 0x3cd3f6a7.
//
// Solidity: function latestWithdrawHash() constant returns(bytes32)
func (_Relay *RelaySession) LatestWithdrawHash() ([32]byte, error) {
	return _Relay.Contract.LatestWithdrawHash(&_Relay.CallOpts)
}

// LatestWithdrawHash is a free data retrieval call binding the contract method 0x3cd3f6a7.
//
// Solidity: function latestWithdrawHash() constant returns(bytes32)
func (_Relay *RelayCallerSession) LatestWithdrawHash() ([32]byte, error) {
	return _Relay.Contract.LatestWithdrawHash(&_Relay.CallOpts)
}

// LogicContract is a free data retrieval call binding the contract method 0xcc0e97c9.
//
// Solidity: function logicContract() constant returns(address)
func (_Relay *RelayCaller) LogicContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "logicContract")
	return *ret0, err
}

// LogicContract is a free data retrieval call binding the contract method 0xcc0e97c9.
//
// Solidity: function logicContract() constant returns(address)
func (_Relay *RelaySession) LogicContract() (common.Address, error) {
	return _Relay.Contract.LogicContract(&_Relay.CallOpts)
}

// LogicContract is a free data retrieval call binding the contract method 0xcc0e97c9.
//
// Solidity: function logicContract() constant returns(address)
func (_Relay *RelayCallerSession) LogicContract() (common.Address, error) {
	return _Relay.Contract.LogicContract(&_Relay.CallOpts)
}

// Reserved is a free data retrieval call binding the contract method 0x92698814.
//
// Solidity: function reserved( bytes32) constant returns(bytes32)
func (_Relay *RelayCaller) Reserved(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "reserved", arg0)
	return *ret0, err
}

// Reserved is a free data retrieval call binding the contract method 0x92698814.
//
// Solidity: function reserved( bytes32) constant returns(bytes32)
func (_Relay *RelaySession) Reserved(arg0 [32]byte) ([32]byte, error) {
	return _Relay.Contract.Reserved(&_Relay.CallOpts, arg0)
}

// Reserved is a free data retrieval call binding the contract method 0x92698814.
//
// Solidity: function reserved( bytes32) constant returns(bytes32)
func (_Relay *RelayCallerSession) Reserved(arg0 [32]byte) ([32]byte, error) {
	return _Relay.Contract.Reserved(&_Relay.CallOpts, arg0)
}

// TokenContract is a free data retrieval call binding the contract method 0x55a373d6.
//
// Solidity: function tokenContract() constant returns(address)
func (_Relay *RelayCaller) TokenContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "tokenContract")
	return *ret0, err
}

// TokenContract is a free data retrieval call binding the contract method 0x55a373d6.
//
// Solidity: function tokenContract() constant returns(address)
func (_Relay *RelaySession) TokenContract() (common.Address, error) {
	return _Relay.Contract.TokenContract(&_Relay.CallOpts)
}

// TokenContract is a free data retrieval call binding the contract method 0x55a373d6.
//
// Solidity: function tokenContract() constant returns(address)
func (_Relay *RelayCallerSession) TokenContract() (common.Address, error) {
	return _Relay.Contract.TokenContract(&_Relay.CallOpts)
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() constant returns(uint256)
func (_Relay *RelayCaller) TotalVotingPower(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "totalVotingPower")
	return *ret0, err
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() constant returns(uint256)
func (_Relay *RelaySession) TotalVotingPower() (*big.Int, error) {
	return _Relay.Contract.TotalVotingPower(&_Relay.CallOpts)
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() constant returns(uint256)
func (_Relay *RelayCallerSession) TotalVotingPower() (*big.Int, error) {
	return _Relay.Contract.TotalVotingPower(&_Relay.CallOpts)
}

// ValidatorInfo is a free data retrieval call binding the contract method 0x4f1811dd.
//
// Solidity: function validatorInfo( address) constant returns(index uint8, power uint32)
func (_Relay *RelayCaller) ValidatorInfo(opts *bind.CallOpts, arg0 common.Address) (struct {
	Index uint8
	Power uint32
}, error) {
	ret := new(struct {
		Index uint8
		Power uint32
	})
	out := ret
	err := _Relay.contract.Call(opts, out, "validatorInfo", arg0)
	return *ret, err
}

// ValidatorInfo is a free data retrieval call binding the contract method 0x4f1811dd.
//
// Solidity: function validatorInfo( address) constant returns(index uint8, power uint32)
func (_Relay *RelaySession) ValidatorInfo(arg0 common.Address) (struct {
	Index uint8
	Power uint32
}, error) {
	return _Relay.Contract.ValidatorInfo(&_Relay.CallOpts, arg0)
}

// ValidatorInfo is a free data retrieval call binding the contract method 0x4f1811dd.
//
// Solidity: function validatorInfo( address) constant returns(index uint8, power uint32)
func (_Relay *RelayCallerSession) ValidatorInfo(arg0 common.Address) (struct {
	Index uint8
	Power uint32
}, error) {
	return _Relay.Contract.ValidatorInfo(&_Relay.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(address)
func (_Relay *RelayCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(address)
func (_Relay *RelaySession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Relay.Contract.Validators(&_Relay.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(address)
func (_Relay *RelayCallerSession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Relay.Contract.Validators(&_Relay.CallOpts, arg0)
}

// WithdrawRootHash is a free data retrieval call binding the contract method 0x6b12b635.
//
// Solidity: function withdrawRootHash(_withdrawInfo bytes, _proof bytes) constant returns(bytes20)
func (_Relay *RelayCaller) WithdrawRootHash(opts *bind.CallOpts, _withdrawInfo []byte, _proof []byte) ([20]byte, error) {
	var (
		ret0 = new([20]byte)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "withdrawRootHash", _withdrawInfo, _proof)
	return *ret0, err
}

// WithdrawRootHash is a free data retrieval call binding the contract method 0x6b12b635.
//
// Solidity: function withdrawRootHash(_withdrawInfo bytes, _proof bytes) constant returns(bytes20)
func (_Relay *RelaySession) WithdrawRootHash(_withdrawInfo []byte, _proof []byte) ([20]byte, error) {
	return _Relay.Contract.WithdrawRootHash(&_Relay.CallOpts, _withdrawInfo, _proof)
}

// WithdrawRootHash is a free data retrieval call binding the contract method 0x6b12b635.
//
// Solidity: function withdrawRootHash(_withdrawInfo bytes, _proof bytes) constant returns(bytes20)
func (_Relay *RelayCallerSession) WithdrawRootHash(_withdrawInfo []byte, _proof []byte) ([20]byte, error) {
	return _Relay.Contract.WithdrawRootHash(&_Relay.CallOpts, _withdrawInfo, _proof)
}

// CommitWithdrawHash is a paid mutator transaction binding the contract method 0x11f7ee0d.
//
// Solidity: function commitWithdrawHash(_payload bytes) returns()
func (_Relay *RelayTransactor) CommitWithdrawHash(opts *bind.TransactOpts, _payload []byte) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "commitWithdrawHash", _payload)
}

// CommitWithdrawHash is a paid mutator transaction binding the contract method 0x11f7ee0d.
//
// Solidity: function commitWithdrawHash(_payload bytes) returns()
func (_Relay *RelaySession) CommitWithdrawHash(_payload []byte) (*types.Transaction, error) {
	return _Relay.Contract.CommitWithdrawHash(&_Relay.TransactOpts, _payload)
}

// CommitWithdrawHash is a paid mutator transaction binding the contract method 0x11f7ee0d.
//
// Solidity: function commitWithdrawHash(_payload bytes) returns()
func (_Relay *RelayTransactorSession) CommitWithdrawHash(_payload []byte) (*types.Transaction, error) {
	return _Relay.Contract.CommitWithdrawHash(&_Relay.TransactOpts, _payload)
}

// UpdateValidator is a paid mutator transaction binding the contract method 0x2277e53a.
//
// Solidity: function updateValidator(_newValidators address[], _proof bytes) returns()
func (_Relay *RelayTransactor) UpdateValidator(opts *bind.TransactOpts, _newValidators []common.Address, _proof []byte) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "updateValidator", _newValidators, _proof)
}

// UpdateValidator is a paid mutator transaction binding the contract method 0x2277e53a.
//
// Solidity: function updateValidator(_newValidators address[], _proof bytes) returns()
func (_Relay *RelaySession) UpdateValidator(_newValidators []common.Address, _proof []byte) (*types.Transaction, error) {
	return _Relay.Contract.UpdateValidator(&_Relay.TransactOpts, _newValidators, _proof)
}

// UpdateValidator is a paid mutator transaction binding the contract method 0x2277e53a.
//
// Solidity: function updateValidator(_newValidators address[], _proof bytes) returns()
func (_Relay *RelayTransactorSession) UpdateValidator(_newValidators []common.Address, _proof []byte) (*types.Transaction, error) {
	return _Relay.Contract.UpdateValidator(&_Relay.TransactOpts, _newValidators, _proof)
}

// Withdraw is a paid mutator transaction binding the contract method 0x50fbe2d9.
//
// Solidity: function withdraw(_withdrawInfo bytes, _proof bytes) returns()
func (_Relay *RelayTransactor) Withdraw(opts *bind.TransactOpts, _withdrawInfo []byte, _proof []byte) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "withdraw", _withdrawInfo, _proof)
}

// Withdraw is a paid mutator transaction binding the contract method 0x50fbe2d9.
//
// Solidity: function withdraw(_withdrawInfo bytes, _proof bytes) returns()
func (_Relay *RelaySession) Withdraw(_withdrawInfo []byte, _proof []byte) (*types.Transaction, error) {
	return _Relay.Contract.Withdraw(&_Relay.TransactOpts, _withdrawInfo, _proof)
}

// Withdraw is a paid mutator transaction binding the contract method 0x50fbe2d9.
//
// Solidity: function withdraw(_withdrawInfo bytes, _proof bytes) returns()
func (_Relay *RelayTransactorSession) Withdraw(_withdrawInfo []byte, _proof []byte) (*types.Transaction, error) {
	return _Relay.Contract.Withdraw(&_Relay.TransactOpts, _withdrawInfo, _proof)
}
