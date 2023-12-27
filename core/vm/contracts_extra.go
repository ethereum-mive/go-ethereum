package vm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ExtraPrecompiledContract is similar to PrecompiledContract,
// but it is used for extra added precompiles.
type ExtraPrecompiledContract interface {
	ContractRef
	RequiredGas(input []byte) uint64                                                               // RequiredGas calculates the contract gas use
	Run(evm *EVM, caller ContractRef, input []byte, value *big.Int, readOnly bool) ([]byte, error) // Run runs the precompiled contract
}

var (
	ExtraPrecompiledContracts = map[common.Address]PrecompiledContract{}
	ExtraPrecompiledAddresses []common.Address
)

// AddPrecompiledContracts adds extra precompiled contracts.
func AddPrecompiledContracts(precompiles ...ExtraPrecompiledContract) {
	for _, precompile := range precompiles {
		ExtraPrecompiledContracts[precompile.Address()] = &extraPrecompiledContract{precompile}
		ExtraPrecompiledAddresses = append(ExtraPrecompiledAddresses, precompile.Address())
	}
}

// RunPrecompiledContract2 runs and evaluates the output of a precompiled contract.
// It is similar to RunPrecompiledContract, but it considers the ExtraPrecompiledContract interface.
func RunPrecompiledContract2(
	evm *EVM,
	p PrecompiledContract,
	caller ContractRef,
	input []byte,
	value *big.Int,
	suppliedGas uint64,
	readOnly bool,
) (ret []byte, remainingGas uint64, err error) {
	if c, ok := p.(*extraPrecompiledContract); ok {
		gasCost := c.RequiredGas(input)
		if suppliedGas < gasCost {
			return nil, 0, ErrOutOfGas
		}
		suppliedGas -= gasCost
		output, err := c.precompile.Run(evm, caller, input, value, readOnly)
		return output, suppliedGas, err
	}

	return RunPrecompiledContract(p, input, suppliedGas)
}

type extraPrecompiledContract struct {
	precompile ExtraPrecompiledContract
}

func (c *extraPrecompiledContract) RequiredGas(input []byte) uint64 {
	return c.precompile.RequiredGas(input)
}

func (c *extraPrecompiledContract) Run(input []byte) ([]byte, error) {
	panic("the Run method is not implemented")
}
