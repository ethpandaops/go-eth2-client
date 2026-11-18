// Copyright © 2023 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package daisugi

import (
	"fmt"

	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/holiman/uint256"
	"github.com/pk910/dynamic-ssz/hasher"
	"github.com/pk910/dynamic-ssz/sszutils"
)

// MaxRecursiveProofBytes is the transport bound of the recursive STARK proof (EIP-8288).
const MaxRecursiveProofBytes = 8388608

// ExecutionPayload represents an execution layer payload.
type ExecutionPayload struct {
	ParentHash          phase0.Hash32              `ssz-size:"32"`
	FeeRecipient        bellatrix.ExecutionAddress `ssz-size:"20"`
	StateRoot           phase0.Root                `ssz-size:"32"`
	ReceiptsRoot        phase0.Root                `ssz-size:"32"`
	LogsBloom           [256]byte                  `ssz-size:"256"`
	PrevRandao          [32]byte                   `ssz-size:"32"`
	BlockNumber         uint64
	GasLimit            uint64
	GasUsed             uint64
	Timestamp           uint64
	ExtraData           []byte                  `dynssz-max:"MAX_EXTRA_DATA_BYTES"                                   ssz-max:"32"`
	BaseFeePerGas       *uint256.Int            `ssz-type:"uint256"`
	BlockHash           phase0.Hash32           `ssz-size:"32"`
	Transactions        []bellatrix.Transaction `dynssz-max:"MAX_TRANSACTIONS_PER_PAYLOAD,MAX_BYTES_PER_TRANSACTION" ssz-max:"1048576,1073741824" ssz-size:"?,?"`
	Withdrawals         []*capella.Withdrawal   `dynssz-max:"MAX_WITHDRAWALS_PER_PAYLOAD"                            ssz-max:"16"`
	BlobGasUsed         uint64
	ExcessBlobGas       uint64
	RecursiveStarkProof []byte `ssz-max:"8388608"`
}

// RecursiveStarkProofRoot returns the SSZ hash tree root of a recursive STARK proof, which
// is the value ExecutionPayloadHeader.RecursiveStarkProofRoot commits to.
func RecursiveStarkProofRoot(proof []byte) (phase0.Root, error) {
	length := uint64(len(proof))
	if length > MaxRecursiveProofBytes {
		return phase0.Root{}, fmt.Errorf("recursive stark proof exceeds %d bytes", MaxRecursiveProofBytes)
	}

	var root phase0.Root

	err := hasher.WithDefaultHasher(func(hh sszutils.HashWalker) error {
		idx := hh.StartTree(sszutils.TreeTypeBinary | sszutils.TreeTypePacked)
		hh.AppendBytes32(proof)
		hh.MerkleizeWithMixin(idx, length, sszutils.CalculateLimit(MaxRecursiveProofBytes, length, 1))

		res, err := hh.HashRoot()
		if err != nil {
			return err
		}

		root = res

		return nil
	})
	if err != nil {
		return phase0.Root{}, fmt.Errorf("failed to hash recursive stark proof: %w", err)
	}

	return root, nil
}

// String returns a string version of the structure.
func (e *ExecutionPayload) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
