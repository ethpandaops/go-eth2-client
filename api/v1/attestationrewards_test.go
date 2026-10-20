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

func TestValidatorAttestationRewardsJSON(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		head       int64
		target     int64
		source     int64
		inactivity int64
		err        string
	}{
		{
			// Mainnet epoch 200759, during the May 2023 inactivity leak.
			name:       "InactivityLeakPenalty",
			input:      []byte(`{"validator_index":"9","head":"0","target":"-6162","source":"-3318","inactivity":"-5722"}`),
			target:     -6162,
			source:     -3318,
			inactivity: -5722,
		},
		{
			name:   "Rewarded",
			input:  []byte(`{"validator_index":"2","head":"1825","target":"5034","source":"2543","inactivity":"0"}`),
			head:   1825,
			target: 5034,
			source: 2543,
		},
		{
			name:   "HeadPenalty",
			input:  []byte(`{"validator_index":"2","head":"-1825","target":"5034","source":"2543","inactivity":"0"}`),
			head:   -1825,
			target: 5034,
			source: 2543,
		},
		{
			name:  "InactivityInvalid",
			input: []byte(`{"validator_index":"2","head":"0","target":"0","source":"0","inactivity":"x"}`),
			err:   "invalid value for inactivity: strconv.ParseInt: parsing \"x\": invalid syntax",
		},
		{
			name:  "InactivityMissing",
			input: []byte(`{"validator_index":"2","head":"0","target":"0","source":"0"}`),
			err:   "inactivity missing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ValidatorAttestationRewards
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.head, res.Head)
			assert.Equal(t, test.target, res.Target)
			assert.Equal(t, test.source, res.Source)
			assert.Equal(t, test.inactivity, res.Inactivity)

			rt, err := json.Marshal(&res)
			require.NoError(t, err)
			assert.JSONEq(t, string(test.input), string(rt))
		})
	}
}

func TestIdealAttestationRewardsJSON(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		inactivity int64
		err        string
	}{
		{
			name:  "Good",
			input: []byte(`{"effective_balance":"32000000000","head":"1825","target":"5034","source":"2543","inactivity":"0"}`),
		},
		{
			name:       "InactivityPenalty",
			input:      []byte(`{"effective_balance":"32000000000","head":"1825","target":"5034","source":"2543","inactivity":"-4768"}`),
			inactivity: -4768,
		},
		{
			name:  "EffectiveBalanceNegative",
			input: []byte(`{"effective_balance":"-1","head":"1825","target":"5034","source":"2543","inactivity":"0"}`),
			err:   "invalid value for effective balance: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.IdealAttestationRewards
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.inactivity, res.Inactivity)

			rt, err := json.Marshal(&res)
			require.NoError(t, err)
			assert.JSONEq(t, string(test.input), string(rt))
		})
	}
}
