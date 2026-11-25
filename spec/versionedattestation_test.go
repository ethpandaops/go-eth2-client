// Copyright © 2025 Attestant Limited.
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
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/OffchainLabs/go-bitfield"
	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
)

func TestVersionedAttestation_CommitteeIndex(t *testing.T) {
	// Test cases
	tests := []struct {
		name            string
		expectedIndices []phase0.CommitteeIndex
		errorMsg        string
	}{
		{
			name:            "Valid index 0",
			expectedIndices: []phase0.CommitteeIndex{0},
		},
		{
			name:            "Valid index 4",
			expectedIndices: []phase0.CommitteeIndex{4},
		},
		{
			name:            "Valid index 40",
			expectedIndices: []phase0.CommitteeIndex{40},
		},
		{
			name:            "Invalid index 64",
			expectedIndices: []phase0.CommitteeIndex{64},
			errorMsg:        "no committee index found in committee bits",
		},
		{
			name:            "Invalid no index set",
			expectedIndices: []phase0.CommitteeIndex{64},
			errorMsg:        "no committee index found in committee bits",
		},
		{
			name:            "Invalid multiple index set",
			expectedIndices: []phase0.CommitteeIndex{4, 40},
			errorMsg:        "multiple committee indices found in committee bits",
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			committeeBits := bitfield.NewBitvector64()
			for _, expectedIndex := range tt.expectedIndices {
				committeeBits.SetBitAt(uint64(expectedIndex), true)
			}
			attestation := spec.VersionedAttestation{
				Version: spec.DataVersionElectra,
				Electra: &electra.Attestation{
					CommitteeBits: committeeBits,
				},
			}
			committeeIndex, err := attestation.CommitteeIndex()
			if tt.errorMsg == "" {
				require.NoError(t, err)
				for _, expectedIndex := range tt.expectedIndices {
					assert.Equal(t, expectedIndex, committeeIndex)
				}
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestVersionedAttestation_UnmarshalJSON(t *testing.T) {
	data := `{"slot":"4070","index":"1","beacon_block_root":"0x1111111111111111111111111111111111111111111111111111111111111111","source":{"epoch":"125","root":"0x2222222222222222222222222222222222222222222222222222222222222222"},"target":{"epoch":"127","root":"0x3333333333333333333333333333333333333333333333333333333333333333"}}`

	tests := []struct {
		name      string
		input     string
		version   spec.DataVersion
		committee phase0.CommitteeIndex
		bits      uint64
	}{
		{
			name:      "PQAggregateWithoutSignature",
			input:     `{"aggregation_bits":"0x0b","data":` + data + `,"committee_bits":"0x0200"}`,
			version:   spec.DataVersionPQ,
			committee: 1,
			bits:      2,
		},
		{
			name:      "ElectraAggregateWithSignature",
			input:     `{"aggregation_bits":"0x0b","data":` + data + `,"signature":"0x` + strings.Repeat("00", 96) + `","committee_bits":"0x0400000000000000"}`,
			version:   spec.DataVersionElectra,
			committee: 2,
			bits:      2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var attestation spec.VersionedAttestation
			require.NoError(t, json.Unmarshal([]byte(test.input), &attestation))
			require.Equal(t, test.version, attestation.Version)

			committee, err := attestation.CommitteeIndex()
			require.NoError(t, err)
			require.Equal(t, test.committee, committee)

			bits, err := attestation.AggregationBits()
			require.NoError(t, err)
			require.Equal(t, test.bits, bits.Count())

			attestationData, err := attestation.Data()
			require.NoError(t, err)
			require.Equal(t, phase0.Slot(4070), attestationData.Slot)
		})
	}
}
