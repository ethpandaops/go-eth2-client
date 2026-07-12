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

func TestHeadEventV2JSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.headEventV2JSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","payload_status":"full","epoch_transition":false,"current_epoch_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","next_epoch_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11","execution_optimistic":false}`),
			err:   "slot missing",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","payload_status":"full","epoch_transition":false,"execution_optimistic":false}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BlockMissing",
			input: []byte(`{"slot":"525277","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","payload_status":"full","epoch_transition":false,"execution_optimistic":false}`),
			err:   "block missing",
		},
		{
			name:  "BlockWrongLength",
			input: []byte(`{"slot":"525277","block":"0xe3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","payload_status":"full","epoch_transition":false,"execution_optimistic":false}`),
			err:   "incorrect length 31 for block",
		},
		{
			name:  "StateMissing",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","payload_status":"full","epoch_transition":false,"execution_optimistic":false}`),
			err:   "state missing",
		},
		{
			name:  "DependentRootsMissing",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","payload_status":"empty","epoch_transition":false,"execution_optimistic":false}`),
		},
		{
			name:  "CurrentEpochDependentRootWrongLength",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","payload_status":"full","epoch_transition":false,"current_epoch_dependent_root":"0x7a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","next_epoch_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11","execution_optimistic":false}`),
			err:   "incorrect length 31 for current epoch dependent root",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","payload_status":"full","epoch_transition":false,"current_epoch_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","next_epoch_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11","execution_optimistic":false}`),
		},
		{
			name:  "GoodEpochTransition",
			input: []byte(`{"slot":"525280","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","payload_status":"empty","epoch_transition":true,"current_epoch_dependent_root":"0x907a3462a2905e3df2624869aa7f9a8635eb35bdcf9ce68a26fab691f9dada61","next_epoch_dependent_root":"0x935569bdc1aaad65dbeb532a125390d039058924ea81799238ed53e4e4639a11","execution_optimistic":true}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.HeadEventV2
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
