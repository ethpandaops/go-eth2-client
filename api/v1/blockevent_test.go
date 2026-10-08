// Copyright © 2020 Attestant Limited.
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

func TestBlockEventJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.blockEventJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blockEventJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BlockMissing",
			input: []byte(`{"slot":"525277","execution_optimistic":false}`),
			err:   "block missing",
		},
		{
			name:  "BlockWrongType",
			input: []byte(`{"slot":"525277","block":true,"execution_optimistic":false}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blockEventJSON.block of type string",
		},
		{
			name:  "BlockInvalid",
			input: []byte(`{"slot":"525277","block":"invalid","execution_optimistic":false}`),
			err:   "invalid value for block: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockShort",
			input: []byte(`{"slot":"525277","block":"0xe3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false}`),
			err:   "incorrect length 31 for block",
		},
		{
			name:  "BlockLong",
			input: []byte(`{"slot":"525277","block":"0x9999e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false}`),
			err:   "incorrect length 33 for block",
		},
		{
			name:  "OptimisticWrongType",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":"false"}`),
			err:   "invalid JSON: json: cannot unmarshal string into Go struct field blockEventJSON.execution_optimistic of type bool",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false}`),
		},
		{
			name:  "Optimistic",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":true}`),
		},
		{
			name:  "BuilderIndexInvalid",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false,"builder_index":"-1","block_hash":"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}`),
			err:   "invalid value for builder index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BlockHashInvalid",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false,"builder_index":"42","block_hash":"invalid"}`),
			err:   "invalid value for block hash: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockHashShort",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false,"builder_index":"42","block_hash":"0x34567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}`),
			err:   "incorrect length 31 for block hash",
		},
		{
			name:  "GoodGloas",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false,"builder_index":"42","block_hash":"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}`),
		},
		{
			name:  "GoodGloasSelfBuild",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false,"builder_index":"18446744073709551615","block_hash":"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.BlockEvent
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

func TestBlockEventGloasFields(t *testing.T) {
	var res api.BlockEvent
	require.NoError(t, json.Unmarshal([]byte(`{"slot":"10","block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","execution_optimistic":false,"builder_index":"42","block_hash":"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}`), &res))
	require.NotNil(t, res.BuilderIndex)
	assert.Equal(t, uint64(42), *res.BuilderIndex)
	require.NotNil(t, res.BlockHash)
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", res.BlockHash.String())

	// Decoding a pre-Gloas event into the same struct must not keep the
	// previous Gloas values.
	require.NoError(t, json.Unmarshal([]byte(`{"slot":"11","block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","execution_optimistic":false}`), &res))
	assert.Nil(t, res.BuilderIndex)
	assert.Nil(t, res.BlockHash)
}
