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

package spec_test

import (
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/daisugi"
	"github.com/ethpandaops/go-eth2-client/spec/deneb"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func daisugiTestBlock(proof []byte) *daisugi.SignedBeaconBlock {
	return &daisugi.SignedBeaconBlock{
		Message: &daisugi.BeaconBlock{
			Slot:          96,
			ProposerIndex: 5,
			ParentRoot:    phase0.Root{0x01},
			StateRoot:     phase0.Root{0x02},
			Body: &daisugi.BeaconBlockBody{
				ETH1Data: &phase0.ETH1Data{},
				ExecutionPayload: &daisugi.ExecutionPayload{
					BlockNumber:         42,
					Timestamp:           1700000212,
					BaseFeePerGas:       uint256.NewInt(7),
					BlockHash:           phase0.Hash32{0x03},
					BlobGasUsed:         1,
					ExcessBlobGas:       2,
					RecursiveStarkProof: proof,
				},
				ExecutionRequests: &electra.ExecutionRequests{},
			},
		},
	}
}

func TestVersionedExecutionPayloadRecursiveStarkProof(t *testing.T) {
	proof := []byte{0x01, 0x02, 0x03}

	tests := []struct {
		name    string
		payload *spec.VersionedExecutionPayload
		proof   []byte
		err     string
	}{
		{
			name:    "Phase0",
			payload: &spec.VersionedExecutionPayload{Version: spec.DataVersionPhase0},
			err:     "no execution payload before bellatrix",
		},
		{
			name: "Electra",
			payload: &spec.VersionedExecutionPayload{
				Version: spec.DataVersionElectra,
				Electra: &deneb.ExecutionPayload{},
			},
			err: "recursive stark proof is only present in daisugi",
		},
		{
			name:    "Heze",
			payload: &spec.VersionedExecutionPayload{Version: spec.DataVersionHeze},
			err:     "recursive stark proof is only present in daisugi",
		},
		{
			name:    "DaisugiMissing",
			payload: &spec.VersionedExecutionPayload{Version: spec.DataVersionDaisugi},
			err:     "no daisugi execution payload",
		},
		{
			name: "Daisugi",
			payload: &spec.VersionedExecutionPayload{
				Version: spec.DataVersionDaisugi,
				Daisugi: daisugiTestBlock(proof).Message.Body.ExecutionPayload,
			},
			proof: proof,
		},
		{
			name:    "Unknown",
			payload: &spec.VersionedExecutionPayload{Version: spec.DataVersionUnknown},
			err:     "unknown version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := test.payload.RecursiveStarkProof()
			if test.err != "" {
				require.EqualError(t, err, test.err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.proof, res)
		})
	}
}

func TestVersionedSignedBeaconBlockDaisugi(t *testing.T) {
	proof := []byte{0xaa, 0xbb}
	block := &spec.VersionedSignedBeaconBlock{
		Version: spec.DataVersionDaisugi,
		Daisugi: daisugiTestBlock(proof),
	}

	slot, err := block.Slot()
	require.NoError(t, err)
	require.Equal(t, phase0.Slot(96), slot)

	proposer, err := block.ProposerIndex()
	require.NoError(t, err)
	require.Equal(t, phase0.ValidatorIndex(5), proposer)

	blockHash, err := block.ExecutionBlockHash()
	require.NoError(t, err)
	require.Equal(t, phase0.Hash32{0x03}, blockHash)

	blockNumber, err := block.ExecutionBlockNumber()
	require.NoError(t, err)
	require.Equal(t, uint64(42), blockNumber)

	gotProof, err := block.ExecutionRecursiveStarkProof()
	require.NoError(t, err)
	require.Equal(t, proof, gotProof)

	payload, err := block.ExecutionPayload()
	require.NoError(t, err)
	require.Equal(t, spec.DataVersionDaisugi, payload.Version)
	require.False(t, payload.IsEmpty())

	payloadProof, err := payload.RecursiveStarkProof()
	require.NoError(t, err)
	require.Equal(t, proof, payloadProof)

	excessBlobGas, err := payload.ExcessBlobGas()
	require.NoError(t, err)
	require.Equal(t, uint64(2), excessBlobGas)

	requests, err := block.ExecutionRequests()
	require.NoError(t, err)
	require.Equal(t, spec.DataVersionDaisugi, requests.Version)

	// Daisugi branches off Electra, so it has none of the Gloas constructs.
	_, err = block.SignedExecutionPayloadBid()
	require.Error(t, err)

	root, err := block.Root()
	require.NoError(t, err)

	expectedRoot, err := block.Daisugi.Message.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, phase0.Root(expectedRoot), root)
}

func TestVersionedSignedBeaconBlockRecursiveStarkProofErrors(t *testing.T) {
	tests := []struct {
		name  string
		block *spec.VersionedSignedBeaconBlock
		err   string
	}{
		{
			name: "Electra",
			block: &spec.VersionedSignedBeaconBlock{
				Version: spec.DataVersionElectra,
				Electra: &electra.SignedBeaconBlock{},
			},
			err: "recursive stark proof is only present in daisugi",
		},
		{
			name:  "Fulu",
			block: &spec.VersionedSignedBeaconBlock{Version: spec.DataVersionFulu},
			err:   "recursive stark proof is only present in daisugi",
		},
		{
			name:  "DaisugiMissing",
			block: &spec.VersionedSignedBeaconBlock{Version: spec.DataVersionDaisugi},
			err:   "no daisugi block",
		},
		{
			name:  "Unknown",
			block: &spec.VersionedSignedBeaconBlock{Version: spec.DataVersionUnknown},
			err:   "unknown version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.block.ExecutionRecursiveStarkProof()
			require.EqualError(t, err, test.err)
		})
	}
}
