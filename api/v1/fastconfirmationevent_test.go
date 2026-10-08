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
	require "github.com/stretchr/testify/require"
)

func TestFastConfirmationEventJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.fastConfirmationEventJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"block":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","current_slot":"2"}`),
			err:   "slot missing",
		},
		{
			name:  "BlockMissing",
			input: []byte(`{"slot":"1","current_slot":"2"}`),
			err:   "block missing",
		},
		{
			name:  "BlockShort",
			input: []byte(`{"block":"0x8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","current_slot":"2"}`),
			err:   "incorrect length 31 for block",
		},
		{
			name:  "CurrentSlotInvalid",
			input: []byte(`{"block":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","current_slot":"-1"}`),
			err:   "invalid value for current slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"block":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1","current_slot":"2"}`),
		},
		{
			name:  "GoodCurrentSlotZero",
			input: []byte(`{"block":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"0","current_slot":"0"}`),
		},
		{
			// Beacon nodes that predate current_slot omit it.
			name:  "GoodCurrentSlotMissing",
			input: []byte(`{"block":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"1"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.FastConfirmationEvent
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}

func TestFastConfirmationEventCurrentSlot(t *testing.T) {
	var res api.FastConfirmationEvent
	require.NoError(t, json.Unmarshal([]byte(`{"block":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"11307428","current_slot":"11307429"}`), &res))
	assert.Equal(t, phase0.Slot(11307428), res.Slot)
	require.NotNil(t, res.CurrentSlot)
	assert.Equal(t, phase0.Slot(11307429), *res.CurrentSlot)

	// Decoding an event without current_slot into the same struct must not keep
	// the previous value.
	require.NoError(t, json.Unmarshal([]byte(`{"block":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","slot":"11307429"}`), &res))
	assert.Nil(t, res.CurrentSlot)
}
