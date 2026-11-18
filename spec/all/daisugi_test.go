// Copyright © 2026 Attestant Limited.
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

package all_test

import (
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/all"
	"github.com/ethpandaops/go-eth2-client/spec/altair"
	"github.com/ethpandaops/go-eth2-client/spec/daisugi"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func daisugiSignedBlock(proof []byte) *daisugi.SignedBeaconBlock {
	return &daisugi.SignedBeaconBlock{
		Message: &daisugi.BeaconBlock{
			Slot:          96,
			ProposerIndex: 5,
			ParentRoot:    phase0.Root{0x01},
			StateRoot:     phase0.Root{0x02},
			Body: &daisugi.BeaconBlockBody{
				RANDAOReveal: phase0.BLSSignature{0x03},
				ETH1Data:     &phase0.ETH1Data{BlockHash: make([]byte, 32)},
				Graffiti:     [32]byte{0x04},
				Attestations: []*electra.Attestation{
					{
						AggregationBits: bitfield.NewBitlist(8),
						Data: &phase0.AttestationData{
							Slot:   95,
							Source: &phase0.Checkpoint{},
							Target: &phase0.Checkpoint{},
						},
						CommitteeBits: bitfield.NewBitvector64(),
					},
				},
				SyncAggregate: &altair.SyncAggregate{
					SyncCommitteeBits: bitfield.NewBitvector512(),
				},
				ExecutionPayload: &daisugi.ExecutionPayload{
					BlockNumber:         42,
					Timestamp:           1700000212,
					ExtraData:           []byte("daisugi"),
					BaseFeePerGas:       uint256.NewInt(7),
					BlockHash:           phase0.Hash32{0x05},
					RecursiveStarkProof: proof,
				},
				ExecutionRequests: &electra.ExecutionRequests{},
			},
		},
		Signature: phase0.BLSSignature{0x06},
	}
}

func TestDaisugiSignedBeaconBlockViews(t *testing.T) {
	tests := []struct {
		name  string
		proof []byte
	}{
		{name: "EmptyProof", proof: []byte{}},
		{name: "WithProof", proof: []byte{0xde, 0xad, 0xbe, 0xef, 0x01}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			block := daisugiSignedBlock(test.proof)

			expectedSSZ, err := block.MarshalSSZ()
			require.NoError(t, err)

			expectedRoot, err := block.Message.HashTreeRoot()
			require.NoError(t, err)

			agnostic := &all.SignedBeaconBlock{}
			require.NoError(t, agnostic.FromVersioned(&spec.VersionedSignedBeaconBlock{
				Version: spec.DataVersionDaisugi,
				Daisugi: block,
			}))

			require.Equal(t, version.DataVersionDaisugi, agnostic.Version)
			require.Equal(t, version.DataVersionDaisugi, agnostic.Message.Body.ExecutionPayload.Version)
			require.Equal(t, version.DataVersionDaisugi, agnostic.Message.Body.Attestations[0].Version)
			require.Equal(t, test.proof, agnostic.Message.Body.ExecutionPayload.RecursiveStarkProof)

			agnosticSSZ, err := agnostic.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expectedSSZ, agnosticSSZ)

			agnosticRoot, err := agnostic.Message.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, expectedRoot, agnosticRoot)

			decoded := &all.SignedBeaconBlock{Version: version.DataVersionDaisugi}
			require.NoError(t, decoded.UnmarshalSSZ(expectedSSZ))
			require.Equal(t, test.proof, decoded.Message.Body.ExecutionPayload.RecursiveStarkProof)

			versioned, err := decoded.ToVersioned()
			require.NoError(t, err)
			require.Equal(t, spec.DataVersionDaisugi, versioned.Version)
			require.NotNil(t, versioned.Daisugi)
			require.Nil(t, versioned.Electra)

			roundTripSSZ, err := versioned.Daisugi.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expectedSSZ, roundTripSSZ)
		})
	}
}

func TestDaisugiExecutionPayloadHeaderView(t *testing.T) {
	header := &daisugi.ExecutionPayloadHeader{
		BlockNumber:             42,
		BaseFeePerGas:           uint256.NewInt(7),
		RecursiveStarkProofRoot: phase0.Root{0x0b},
	}

	expectedRoot, err := header.HashTreeRoot()
	require.NoError(t, err)

	agnostic := &all.ExecutionPayloadHeader{}
	require.NoError(t, agnostic.FromView(header))
	require.Equal(t, version.DataVersionDaisugi, agnostic.Version)
	require.Equal(t, header.RecursiveStarkProofRoot, agnostic.RecursiveStarkProofRoot)

	agnosticRoot, err := agnostic.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, expectedRoot, agnosticRoot)

	view, err := agnostic.ToView()
	require.NoError(t, err)
	require.IsType(t, &daisugi.ExecutionPayloadHeader{}, view)
}
