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

package pq_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/go-eth2-client/spec/pq"
	"github.com/holiman/uint256"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/stretchr/testify/require"
)

func filled(value byte, size int) []byte {
	return bytes.Repeat([]byte{value}, size)
}

func testAttestation() *pq.Attestation {
	bits := bitfield.NewBitlist(3)
	bits.SetBitAt(0, true)
	bits.SetBitAt(2, true)

	return &pq.Attestation{
		AggregationBits: bits,
		Data: &phase0.AttestationData{
			Slot:   9,
			Source: &phase0.Checkpoint{},
			Target: &phase0.Checkpoint{},
		},
		CommitteeBits: bitfield.Bitvector64{0x02},
	}
}

func testExecutionRequests() *pq.ExecutionRequests {
	requests := &pq.ExecutionRequests{
		VoluntaryExits: []*pq.VoluntaryExitRequest{{}},
		ValidatorPreregistrations: []*pq.ValidatorPreregistrationRequest{{
			WithdrawalCredentials: filled(9, 32),
		}},
	}
	copy(requests.VoluntaryExits[0].PublicKey[:], filled(7, 32))
	copy(requests.ValidatorPreregistrations[0].PublicKey[:], filled(8, 32))
	copy(requests.ValidatorPreregistrations[0].VRFRoot[:], filled(10, 32))

	return requests
}

func TestContainerRoots(t *testing.T) {
	ds := dynssz.NewDynSsz(minimalSpecs())

	tests := []struct {
		name  string
		value any
		fresh func() any
		root  string
	}{
		{
			name:  "Attestation",
			value: testAttestation(),
			fresh: func() any { return &pq.Attestation{} },
			root:  "b466724785ecd3e15d9cc9ba91b5e7ce4403e64d16b7a7d4292f173f0c0ae42c",
		},
		{
			name:  "ExecutionRequests",
			value: testExecutionRequests(),
			fresh: func() any { return &pq.ExecutionRequests{} },
			root:  "b6dfa9dc7725d8490dd7bc1b6ffe3b05b62133823ebb86c1f2e856027a7b996b",
		},
		{
			name: "IndexedAttestation",
			value: &pq.IndexedAttestation{
				AttestingIndices: []uint64{1, 2, 3},
				Data:             &phase0.AttestationData{Source: &phase0.Checkpoint{}, Target: &phase0.Checkpoint{}},
				Proof:            filled(5, pq.XMSSSignatureSize),
			},
			fresh: func() any { return &pq.IndexedAttestation{} },
			root:  "ffcf1ce8edee765686dded005254721a8ea0885c1563bccf32083d5e47051962",
		},
		{
			name: "SignedAggregateAndProof",
			value: &pq.SignedAggregateAndProof{
				Message: &pq.AggregateAndProof{
					AggregatorIndex: 4,
					Aggregate:       testAttestation(),
					Proof:           filled(6, 10),
				},
			},
			fresh: func() any { return &pq.SignedAggregateAndProof{} },
			root:  "92cd0d690be592a4dd3179dc8adaef2d0c8eb02f27a3b2be7c78b480c62761b0",
		},
		{
			name: "SignedExecutionPayloadEnvelope",
			value: &pq.SignedExecutionPayloadEnvelope{
				Message: &pq.ExecutionPayloadEnvelope{
					Payload:           &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(0)},
					ExecutionRequests: testExecutionRequests(),
				},
			},
			fresh: func() any { return &pq.SignedExecutionPayloadEnvelope{} },
			root:  "41dc0f1335663ba810a5d2733804aa359987185232522b4dd38a5681009159bd",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := ds.HashTreeRoot(test.value)
			require.NoError(t, err)
			require.Equal(t, test.root, hex.EncodeToString(root[:]))

			encoded, err := ds.MarshalSSZ(test.value)
			require.NoError(t, err)

			decoded := test.fresh()
			require.NoError(t, ds.UnmarshalSSZ(decoded, encoded))

			decodedRoot, err := ds.HashTreeRoot(decoded)
			require.NoError(t, err)
			require.Equal(t, root, decodedRoot)

			marshaled, err := json.Marshal(test.value)
			require.NoError(t, err)

			fromJSON := test.fresh()
			require.NoError(t, json.Unmarshal(marshaled, fromJSON))

			jsonRoot, err := ds.HashTreeRoot(fromJSON)
			require.NoError(t, err)
			require.Equal(t, root, jsonRoot)
		})
	}

	encoded, err := ds.MarshalSSZ(testAttestation())
	require.NoError(t, err)
	require.Equal(t, "850000000900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020d", hex.EncodeToString(encoded))
}
