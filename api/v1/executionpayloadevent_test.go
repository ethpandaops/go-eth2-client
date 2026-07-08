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
	require "github.com/stretchr/testify/require"
)

func TestExecutionPayloadEventJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.executionPayloadEventJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"builder_index":"1","block_hash":"0xc8ac520b396aad76b4cf448aa4ff8ca9b4c1fe600e6416234849654e874db714","block_root":"0x2e52507aad249d4dfe1557f9031d5154c1c7b2a2e5819d0639f719725aa538af","execution_optimistic":true}`),
			err:   "slot missing",
		},
		{
			name:  "BlockRootMissing",
			input: []byte(`{"slot":"91000","builder_index":"1","block_hash":"0xc8ac520b396aad76b4cf448aa4ff8ca9b4c1fe600e6416234849654e874db714","execution_optimistic":true}`),
			err:   "block root missing",
		},
		{
			name:  "BlockRootInvalidLength",
			input: []byte(`{"slot":"91000","builder_index":"1","block_hash":"0xc8ac520b396aad76b4cf448aa4ff8ca9b4c1fe600e6416234849654e874db714","block_root":"0x2e52","execution_optimistic":true}`),
			err:   "incorrect length 2 for block root",
		},
		{
			// Real prysm/lighthouse-shape sample: no state_root, execution_optimistic present.
			name:  "GoodBare",
			input: []byte(`{"slot":"91000","builder_index":"1","block_hash":"0xc8ac520b396aad76b4cf448aa4ff8ca9b4c1fe600e6416234849654e874db714","block_root":"0x2e52507aad249d4dfe1557f9031d5154c1c7b2a2e5819d0639f719725aa538af","execution_optimistic":true}`),
		},
		{
			// Real nimbus-shape sample: carries an extra, non-spec state_root that must be ignored.
			name:  "GoodWithNonSpecStateRoot",
			input: []byte(`{"slot":"91894","builder_index":"18446744073709551615","block_hash":"0x716b390040980aa80e891c4881b5f266812ca1fdbf14efc4ef048145b8d0edb6","block_root":"0x957251f4a7a2b9713520a3f44d008a7814efaf5c63ff2cfe2ff908f28d09866f","state_root":"0x0000000000000000000000000000000000000000000000000000000000000000","execution_optimistic":true}`),
		},
		{
			// Gossip-shape sample: no execution_optimistic.
			name:  "GoodGossip",
			input: []byte(`{"slot":"91001","builder_index":"1","block_hash":"0x45aaeada228ea83858022991d8e6e8a2d6492f58e60d5d4a30acc528fcd5063d","block_root":"0xe7b490b90ad5bcc30d933c88f0177df4df2eb32094036421e455737582880268"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ExecutionPayloadEvent
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.NotEmpty(t, rt)
			}
		})
	}
}

func TestExecutionPayloadEventValues(t *testing.T) {
	input := []byte(`{"slot":"91894","builder_index":"18446744073709551615","block_hash":"0x716b390040980aa80e891c4881b5f266812ca1fdbf14efc4ef048145b8d0edb6","block_root":"0x957251f4a7a2b9713520a3f44d008a7814efaf5c63ff2cfe2ff908f28d09866f","state_root":"0x0000000000000000000000000000000000000000000000000000000000000000","execution_optimistic":true}`)

	var res api.ExecutionPayloadEvent
	require.NoError(t, json.Unmarshal(input, &res))

	assert.Equal(t, uint64(91894), uint64(res.Slot))
	assert.Equal(t, uint64(18446744073709551615), res.BuilderIndex)
	assert.Equal(t, "0x716b390040980aa80e891c4881b5f266812ca1fdbf14efc4ef048145b8d0edb6", res.BlockHash.String())
	assert.Equal(t, "0x957251f4a7a2b9713520a3f44d008a7814efaf5c63ff2cfe2ff908f28d09866f", res.BlockRoot.String())
	assert.True(t, res.ExecutionOptimistic)
}
