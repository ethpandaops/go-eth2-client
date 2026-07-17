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

package gloas_test

import (
	"encoding/json"
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

func TestBuilderJSON(t *testing.T) {
	// Beacon API JSON encodes all integers as strings.
	input := `{"pubkey":"0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":"1","execution_address":"0x000102030405060708090a0b0c0d0e0f10111213","balance":"32000000000","deposit_epoch":"12","withdrawable_epoch":"18446744073709551615"}`

	var builder gloas.Builder
	require.NoError(t, json.Unmarshal([]byte(input), &builder))
	require.Equal(t, uint8(1), builder.Version)

	// Round-trip.
	output, err := json.Marshal(&builder)
	require.NoError(t, err)
	require.JSONEq(t, input, string(output))
}

func TestBuilderJSONInvalidVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		errorMsg string
	}{
		{
			name:     "missing",
			input:    `{"pubkey":"0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","execution_address":"0x000102030405060708090a0b0c0d0e0f10111213","balance":"32000000000","deposit_epoch":"12","withdrawable_epoch":"13"}`,
			errorMsg: "version missing",
		},
		{
			name:     "not a number",
			input:    `{"pubkey":"0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":"banana","execution_address":"0x000102030405060708090a0b0c0d0e0f10111213","balance":"32000000000","deposit_epoch":"12","withdrawable_epoch":"13"}`,
			errorMsg: "invalid value for version",
		},
		{
			name:     "out of range for uint8",
			input:    `{"pubkey":"0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":"256","execution_address":"0x000102030405060708090a0b0c0d0e0f10111213","balance":"32000000000","deposit_epoch":"12","withdrawable_epoch":"13"}`,
			errorMsg: "invalid value for version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var builder gloas.Builder
			err := json.Unmarshal([]byte(test.input), &builder)
			require.ErrorContains(t, err, test.errorMsg)
		})
	}
}

func TestBuilderYAML(t *testing.T) {
	input := `{"pubkey":"0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c","version":"1","execution_address":"0x000102030405060708090a0b0c0d0e0f10111213","balance":"32000000000","deposit_epoch":"12","withdrawable_epoch":"13"}`

	var builder gloas.Builder
	require.NoError(t, json.Unmarshal([]byte(input), &builder))

	yamlBytes, err := builder.MarshalYAML()
	require.NoError(t, err)
	require.Contains(t, string(yamlBytes), "version: '1'")

	var roundTripped gloas.Builder
	require.NoError(t, roundTripped.UnmarshalYAML(yamlBytes))
	require.Equal(t, builder, roundTripped)
}
