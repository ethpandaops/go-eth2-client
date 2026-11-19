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
	"os"
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec/pq"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/stretchr/testify/require"
)

func minimalSpecs() map[string]any {
	return map[string]any{
		"SLOTS_PER_HISTORICAL_ROOT":     uint64(64),
		"EPOCHS_PER_HISTORICAL_VECTOR":  uint64(64),
		"EPOCHS_PER_SLASHINGS_VECTOR":   uint64(64),
		"EPOCHS_PER_ETH1_VOTING_PERIOD": uint64(4),
		"SLOTS_PER_EPOCH":               uint64(8),
		"MIN_SEED_LOOKAHEAD":            uint64(1),
		"PTC_SIZE":                      uint64(16),
		"MAX_COMMITTEES_PER_SLOT_PQ":    uint64(4),
		"MAX_STARK_PROOF_SIZE":          uint64(1048576),
	}
}

func TestGenesisStateFixture(t *testing.T) {
	ds := dynssz.NewDynSsz(minimalSpecs())

	encoded, err := os.ReadFile("testdata/genesis_minimal_ssz.bin")
	require.NoError(t, err)

	state := &pq.BeaconState{}
	require.NoError(t, ds.UnmarshalSSZ(state, encoded))
	require.Len(t, state.XMSSValidators, 8)
	require.Len(t, state.XMSSBuilders, 1)
	require.Equal(t, "0a0f4d59ab5eef65a086a5277b96ceb154f476c5e1bea5bbd49b3a7031d67e41", hex.EncodeToString(state.XMSSValidators[0].PublicKey[:]))

	reencoded, err := ds.MarshalSSZ(state)
	require.NoError(t, err)
	require.True(t, bytes.Equal(encoded, reencoded))

	root, err := ds.HashTreeRoot(state)
	require.NoError(t, err)
	require.Equal(t, "cc0e5639191cf4d1b7d459b8e1d3a7d6aa86eaa37cc63b6f02a9ea6145c15736", hex.EncodeToString(root[:]))

	fixture, err := os.ReadFile("testdata/genesis_minimal.json")
	require.NoError(t, err)

	fromJSON := &pq.BeaconState{}
	require.NoError(t, json.Unmarshal(fixture, fromJSON))

	jsonRoot, err := ds.HashTreeRoot(fromJSON)
	require.NoError(t, err)
	require.Equal(t, root, jsonRoot)

	marshaled, err := json.Marshal(state)
	require.NoError(t, err)

	roundTripped := &pq.BeaconState{}
	require.NoError(t, json.Unmarshal(marshaled, roundTripped))

	roundTrippedRoot, err := ds.HashTreeRoot(roundTripped)
	require.NoError(t, err)
	require.Equal(t, root, roundTrippedRoot)
}
