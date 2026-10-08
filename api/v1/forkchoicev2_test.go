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
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForkChoiceV2JSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
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
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"fork_choice_nodes":[{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}],"extra_data":{}}`),
		},
		{
			name:  "GoodExtraData",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"fork_choice_nodes":[{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}],"extra_data":{"client":"teku"}}`),
		},
		{
			name:  "JustifiedCheckpointMissing",
			input: []byte(`{"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"fork_choice_nodes":[{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}],"extra_data":{}}`),
			err:   "justified checkpoint missing",
		},
		{
			name:  "FinalizedCheckpointMissing",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"fork_choice_nodes":[{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}],"extra_data":{}}`),
			err:   "finalized checkpoint missing",
		},
		{
			name:  "ForkChoiceNodesMissing",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"extra_data":{}}`),
			err:   "fork choice nodes missing",
		},
		{
			name:  "ForkChoiceNodesEmpty",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"fork_choice_nodes":[],"extra_data":{}}`),
			err:   "fork choice nodes missing",
		},
		{
			name:  "ForkChoiceNodeNull",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"fork_choice_nodes":[null],"extra_data":{}}`),
			err:   "fork choice node entry 0 missing",
		},
		{
			name:  "ExtraDataMissing",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"fork_choice_nodes":[{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}]}`),
			err:   "extra data missing",
		},
		{
			name:  "NodeInvalid",
			input: []byte(`{"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"fork_choice_nodes":[{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}],"extra_data":{}}`),
			err:   "invalid value for fork_choice_nodes: justified checkpoint missing",
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
			assert.Equal(t, string(test.input), string(rt))
			assert.Equal(t, string(rt), fc.String())
		})
	}
}

func TestForkChoiceNodeV2JSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
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
			name:  "GoodFull",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
		},
		{
			name:  "GoodPending",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"pending","parent_root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","parent_payload_status":"full","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
		},
		{
			name:  "GoodEmpty",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"empty","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
		},
		{
			name:  "GoodParentNotRetained",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":null,"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
		},
		{
			name:  "GoodOptimistic",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"optimistic","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
		},
		{
			name:  "GoodInvalid",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"invalid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
		},
		{
			name:  "GoodExtraData",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{"balance":"32"}}`),
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "slot missing",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"x","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "invalid value for slot: x: strconv.ParseUint: parsing \"x\": invalid syntax",
		},
		{
			name:  "BlockRootMissing",
			input: []byte(`{"slot":"29","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "block root missing",
		},
		{
			name:  "PayloadStatusMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "payload status missing",
		},
		{
			name:  "PayloadStatusUnknown",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"unknown","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "unrecognised fork choice payload status: unknown",
		},
		{
			name:  "ParentRootMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "parent root missing",
		},
		{
			name:  "ParentRootNull",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":null,"parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "parent root missing",
		},
		{
			name:  "ParentPayloadStatusMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "parent payload status missing",
		},
		{
			name:  "ParentPayloadStatusInvalid",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"bad","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "invalid value for parent payload status: unrecognised fork choice payload status: bad",
		},
		{
			name:  "JustifiedCheckpointMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "justified checkpoint missing",
		},
		{
			name:  "FinalizedCheckpointMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "finalized checkpoint missing",
		},
		{
			name:  "WeightMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "weight missing",
		},
		{
			name:  "ValidityMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "validity missing",
		},
		{
			name:  "ValidityNotYetRevealed",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"not_yet_revealed","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "invalid value for validity: not_yet_revealed",
		},
		{
			name:  "ValidityUnknown",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"unknown","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "invalid value for validity: unknown",
		},
		{
			name:  "ExecutionBlockHashMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "execution block hash missing",
		},
		{
			name:  "PayloadAttesterCountMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "payload attester count missing",
		},
		{
			name:  "PayloadAvailabilityYesCountMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_data_availability_yes_count":"509","extra_data":{}}`),
			err:   "payload availability yes count missing",
		},
		{
			name:  "PayloadDataAvailabilityYesCountMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","extra_data":{}}`),
			err:   "payload data availability yes count missing",
		},
		{
			name:  "ExtraDataMissing",
			input: []byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509"}`),
			err:   "extra data missing",
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
			assert.Equal(t, string(test.input), string(rt))
			assert.Equal(t, string(rt), node.String())
		})
	}
}

func TestForkChoiceNodeV2Values(t *testing.T) {
	var node api.ForkChoiceNodeV2
	require.NoError(t, json.Unmarshal([]byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}`), &node))

	assert.Equal(t, phase0.Slot(29), node.Slot)
	assert.Equal(t, api.ForkChoicePayloadStatusFull, node.PayloadStatus)
	require.NotNil(t, node.ParentPayloadStatus)
	assert.Equal(t, api.ForkChoicePayloadStatusPending, *node.ParentPayloadStatus)
	assert.Equal(t, phase0.Epoch(3), node.JustifiedCheckpoint.Epoch)
	assert.Equal(t, phase0.Root{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa}, node.JustifiedCheckpoint.Root)
	assert.Equal(t, phase0.Epoch(2), node.FinalizedCheckpoint.Epoch)
	assert.Equal(t, phase0.Gwei(32000000000), node.Weight)
	assert.Equal(t, api.ForkChoiceNodeValidityValid, node.Validity)
	assert.Equal(t, uint64(512), node.PayloadAttesterCount)
	assert.Equal(t, uint64(510), node.PayloadAvailabilityYesCount)
	assert.Equal(t, uint64(509), node.PayloadDataAvailabilityYesCount)
	assert.Empty(t, node.ExtraData)
}

// The spec requires at least one fork choice node, so a zero value, as the mock returns, does not decode.
func TestForkChoiceV2ZeroValue(t *testing.T) {
	encoded, err := json.Marshal(&api.ForkChoiceV2{})
	require.NoError(t, err)

	var decoded api.ForkChoiceV2
	require.EqualError(t, json.Unmarshal(encoded, &decoded), "fork choice nodes missing")
}

func TestForkChoicePayloadStatusJSON(t *testing.T) {
	for _, status := range []api.ForkChoicePayloadStatus{
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
	require.EqualError(t, json.Unmarshal([]byte(`"unknown"`), &decoded), "unrecognised fork choice payload status: unknown")
	require.EqualError(t, json.Unmarshal([]byte(`"FULL"`), &decoded), "unrecognised fork choice payload status: FULL")
}

// TestForkChoiceNodeV2JSONReuse ensures that decoding into an existing node, as encoding/json does when decoding into
// a non-empty slice of nodes, leaves nothing from the previous node behind.
func TestForkChoiceNodeV2JSONReuse(t *testing.T) {
	node := new(api.ForkChoiceNodeV2)
	require.NoError(t, json.Unmarshal([]byte(`{"slot":"29","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":"pending","justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{"x":"1"}}`), node))
	require.NotNil(t, node.ParentPayloadStatus)

	nodes := []*api.ForkChoiceNodeV2{node}
	require.NoError(t, json.Unmarshal([]byte(`[{"slot":"30","block_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","payload_status":"full","parent_root":"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","parent_payload_status":null,"justified_checkpoint":{"epoch":"3","root":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"finalized_checkpoint":{"epoch":"2","root":"0x9999999999999999999999999999999999999999999999999999999999999999"},"weight":"32000000000","validity":"valid","execution_block_hash":"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee","payload_attester_count":"512","payload_availability_yes_count":"510","payload_data_availability_yes_count":"509","extra_data":{}}]`), &nodes))
	require.Same(t, node, nodes[0])
	require.Equal(t, phase0.Slot(30), node.Slot)
	require.Nil(t, node.ParentPayloadStatus)
	require.Empty(t, node.ExtraData)
}
