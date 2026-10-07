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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/ethpandaops/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForkChoiceV2JSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		err      string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: not an object",
		},
		{
			name:  "Good",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xa000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x9000000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"pending","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":"full","justified_epoch":"3","finalized_epoch":"2","weight":"128000000000","validity":"valid","execution_block_hash":"0xaa00000000000000000000000000000000000000000000000000000000000000","payload_attester_count":"16","payload_availability_yes_count":"16","payload_data_availability_yes_count":"16","extra_data":{}}],"extra_data":{"head_root":"0xb000000000000000000000000000000000000000000000000000000000000000"}}`),
		},
		{
			name:     "ExtraDataMissing",
			input:    []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xa000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x9000000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[]}`),
			expected: `{"justified_checkpoint":{"epoch":"3","root":"0xa000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x9000000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[],"extra_data":{}}`,
		},
		{
			name:     "UnknownFieldFolded",
			input:    []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xa000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x9000000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[],"head_root":"0xb000000000000000000000000000000000000000000000000000000000000000"}`),
			expected: `{"justified_checkpoint":{"epoch":"3","root":"0xa000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x9000000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[],"extra_data":{"head_root":"0xb000000000000000000000000000000000000000000000000000000000000000"}}`,
		},
		{
			name:  "JustifiedCheckpointMissing",
			input: []byte(`{"finalized_checkpoint":{"epoch":"2","root":"0x9000000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[]}`),
			err:   "justified checkpoint missing",
		},
		{
			name:  "FinalizedCheckpointMissing",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xa000000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[]}`),
			err:   "finalized checkpoint missing",
		},
		{
			name:  "ForkChoiceNodesMissing",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xa000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x9000000000000000000000000000000000000000000000000000000000000000"}}`),
			err:   "fork choice nodes missing",
		},
		{
			name:  "ForkChoiceNodeNull",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xa000000000000000000000000000000000000000000000000000000000000000"},"finalized_checkpoint":{"epoch":"2","root":"0x9000000000000000000000000000000000000000000000000000000000000000"},"fork_choice_nodes":[null]}`),
			err:   "fork choice node entry 0 missing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var fc api.ForkChoiceV2
			err := json.Unmarshal(test.input, &fc)
			if test.err != "" {
				require.EqualError(t, err, test.err)

				return
			}
			require.NoError(t, err)
			rt, err := json.Marshal(&fc)
			require.NoError(t, err)
			expected := test.expected
			if expected == "" {
				expected = string(test.input)
			}
			assert.Equal(t, expected, string(rt))
			assert.Equal(t, string(rt), fc.String())
		})
	}
}

func TestForkChoiceNodeV2JSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		err      string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: not an object",
		},
		{
			name:  "GoodPending",
			input: []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"pending","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":"full","justified_epoch":"3","finalized_epoch":"2","weight":"128000000000","validity":"valid","execution_block_hash":"0xaa00000000000000000000000000000000000000000000000000000000000000","payload_attester_count":"16","payload_availability_yes_count":"15","payload_data_availability_yes_count":"14","extra_data":{"state_root":"0xc000000000000000000000000000000000000000000000000000000000000000"}}`),
		},
		{
			name:  "GoodEmpty",
			input: []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"empty","parent_root":"0xb000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":"pending","justified_epoch":"3","finalized_epoch":"2","weight":"64000000000","validity":"valid","execution_block_hash":"0xaa00000000000000000000000000000000000000000000000000000000000000","payload_attester_count":"16","payload_availability_yes_count":"15","payload_data_availability_yes_count":"14","extra_data":{}}`),
		},
		{
			name:  "ParentPayloadStatusNull",
			input: []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"pending","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":null,"weight":"128000000000","validity":"optimistic","execution_block_hash":"0xaa00000000000000000000000000000000000000000000000000000000000000","extra_data":{}}`),
		},
		{
			name:     "OptionalFieldsMissing",
			input:    []byte(`{"payload_status":"full","slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			expected: `{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":null,"weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{}}`,
		},
		{
			name:     "ParentRootNull",
			input:    []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":null,"weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{}}`),
			expected: `{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0x0000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":null,"weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{}}`,
		},
		{
			name:     "UnknownFieldsFolded",
			input:    []byte(`{"payload_status":"full","slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","unrealised_justified_epoch":"4","state_root":"0xc000000000000000000000000000000000000000000000000000000000000000","extra_data":{"state_root":"0xd000000000000000000000000000000000000000000000000000000000000000"}}`),
			expected: `{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":null,"weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{"state_root":"0xd000000000000000000000000000000000000000000000000000000000000000","unrealised_justified_epoch":"4"}}`,
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			err:   "slot missing",
		},
		{
			name:  "PayloadStatusMissing",
			input: []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			err:   "payload status missing",
		},
		{
			name:     "PayloadStatusUppercase",
			input:    []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"FULL","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":"Pending","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{}}`),
			expected: `{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":"pending","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{}}`,
		},
		{
			name:  "FieldCaseMismatch",
			input: []byte(`{"Slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			err:   "slot missing",
		},
		{
			name:  "SeveralFieldsInvalid",
			input: []byte(`{"slot":1,"block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":2,"validity":3,"execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			err:   "invalid value for slot: json: cannot unmarshal number into Go value of type string",
		},
		{
			name:  "PayloadStatusInvalid",
			input: []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"bad","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			err:   "unrecognised fork choice payload status: bad",
		},
		{
			name:     "ParentPayloadStatusUnrecognised",
			input:    []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":"bad","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			expected: `{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":null,"weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{"parent_payload_status":"bad"}}`,
		},
		{
			name:  "ParentPayloadStatusWrongType",
			input: []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":1,"weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			err:   "invalid value for parent_payload_status: json: cannot unmarshal number into Go value of type string",
		},
		{
			name:     "ParentRootNull",
			input:    []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"pending","parent_root":null,"weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			expected: `{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"pending","parent_root":"0x0000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":null,"weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{}}`,
		},
		{
			name:     "ValidityNotYetRevealed",
			input:    []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"pending","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"not_yet_revealed","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{}}`),
			expected: `{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"pending","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":null,"weight":"0","validity":"unknown","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{"validity":"not_yet_revealed"}}`,
		},
		{
			name:     "ValidityUnrecognisedExtraDataKeyTaken",
			input:    []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"pending","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"not_yet_revealed","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{"validity":"client"}}`),
			expected: `{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"pending","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","parent_payload_status":null,"weight":"0","validity":"unknown","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","extra_data":{"validity":"client"}}`,
		},
		{
			name:  "WeightInvalid",
			input: []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"-1","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000"}`),
			err:   "invalid value for weight: -1: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "PayloadAttesterCountInvalid",
			input: []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"valid","execution_block_hash":"0xbb00000000000000000000000000000000000000000000000000000000000000","payload_attester_count":"x"}`),
			err:   "invalid value for payload attester count: x: strconv.ParseUint: parsing \"x\": invalid syntax",
		},
		{
			name:  "ExecutionBlockHashShort",
			input: []byte(`{"slot":"29","block_root":"0xb000000000000000000000000000000000000000000000000000000000000000","payload_status":"full","parent_root":"0xa000000000000000000000000000000000000000000000000000000000000000","weight":"0","validity":"valid","execution_block_hash":"0xbb"}`),
			err:   "incorrect length 1 for execution block hash",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var node api.ForkChoiceNodeV2
			err := json.Unmarshal(test.input, &node)
			if test.err != "" {
				require.EqualError(t, err, test.err)

				return
			}
			require.NoError(t, err)
			rt, err := json.Marshal(&node)
			require.NoError(t, err)
			expected := test.expected
			if expected == "" {
				expected = string(test.input)
			}
			assert.Equal(t, expected, string(rt))
			assert.Equal(t, string(rt), node.String())
		})
	}
}

// A zero value, as the mock returns, marshals to JSON that decodes again.
func TestForkChoiceV2ZeroValueRoundTrip(t *testing.T) {
	encoded, err := json.Marshal(&api.ForkChoiceV2{})
	require.NoError(t, err)

	var decoded api.ForkChoiceV2
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Empty(t, decoded.ForkChoiceNodes)
}

func TestForkChoicePayloadStatusJSON(t *testing.T) {
	for _, status := range []api.ForkChoicePayloadStatus{
		api.ForkChoicePayloadStatusUnknown,
		api.ForkChoicePayloadStatusEmpty,
		api.ForkChoicePayloadStatusFull,
		api.ForkChoicePayloadStatusPending,
	} {
		data, err := json.Marshal(status)
		require.NoError(t, err)
		require.Equal(t, `"`+status.String()+`"`, string(data))

		var decoded api.ForkChoicePayloadStatus
		require.NoError(t, json.Unmarshal(data, &decoded))
		require.Equal(t, status, decoded)
	}

	var decoded api.ForkChoicePayloadStatus
	require.EqualError(t, json.Unmarshal([]byte(`"bad"`), &decoded), "unrecognised fork choice payload status: bad")
}
