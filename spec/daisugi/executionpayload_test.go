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

package daisugi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/daisugi"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func testExecutionPayload(proof []byte) *daisugi.ExecutionPayload {
	return &daisugi.ExecutionPayload{
		ParentHash:    phase0.Hash32{0x01},
		FeeRecipient:  bellatrix.ExecutionAddress{0x02},
		StateRoot:     phase0.Root{0x03},
		ReceiptsRoot:  phase0.Root{0x04},
		LogsBloom:     [256]byte{0x05},
		PrevRandao:    [32]byte{0x06},
		BlockNumber:   7,
		GasLimit:      60000000,
		GasUsed:       21000,
		Timestamp:     1700000212,
		ExtraData:     []byte("daisugi"),
		BaseFeePerGas: uint256.NewInt(7),
		BlockHash:     phase0.Hash32{0x08},
		Transactions:  []bellatrix.Transaction{{0x06, 0x01, 0x02}},
		Withdrawals: []*capella.Withdrawal{
			{Index: 1, ValidatorIndex: 2, Address: bellatrix.ExecutionAddress{0x09}, Amount: 3},
		},
		BlobGasUsed:         0,
		ExcessBlobGas:       0,
		RecursiveStarkProof: proof,
	}
}

func TestExecutionPayloadRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		proof []byte
	}{
		{name: "EmptyProof", proof: []byte{}},
		{name: "SmallProof", proof: []byte{0xde, 0xad, 0xbe, 0xef}},
		{name: "LargeProof", proof: bytes.Repeat([]byte{0xab}, 265000)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := testExecutionPayload(test.proof)

			root, err := payload.HashTreeRoot()
			require.NoError(t, err)

			jsonData, err := json.Marshal(payload)
			require.NoError(t, err)

			var raw map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(jsonData, &raw))
			require.JSONEq(t, fmt.Sprintf("%q", fmt.Sprintf("0x%x", test.proof)),
				string(raw["recursive_stark_proof"]))

			var fromJSON daisugi.ExecutionPayload
			require.NoError(t, json.Unmarshal(jsonData, &fromJSON))
			require.Equal(t, test.proof, fromJSON.RecursiveStarkProof)

			jsonRoot, err := fromJSON.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, root, jsonRoot)

			yamlData, err := yaml.Marshal(payload)
			require.NoError(t, err)

			var fromYAML daisugi.ExecutionPayload
			require.NoError(t, yaml.Unmarshal(yamlData, &fromYAML))

			yamlRoot, err := fromYAML.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, root, yamlRoot)

			sszData, err := payload.MarshalSSZ()
			require.NoError(t, err)

			var fromSSZ daisugi.ExecutionPayload
			require.NoError(t, fromSSZ.UnmarshalSSZ(sszData))
			require.Equal(t, test.proof, fromSSZ.RecursiveStarkProof)

			sszRoot, err := fromSSZ.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, root, sszRoot)
		})
	}
}

func TestExecutionPayloadProofChangesRoot(t *testing.T) {
	withoutProof, err := testExecutionPayload([]byte{}).HashTreeRoot()
	require.NoError(t, err)

	withProof, err := testExecutionPayload([]byte{0x01}).HashTreeRoot()
	require.NoError(t, err)

	require.NotEqual(t, withoutProof, withProof)
}

func TestExecutionPayloadJSONProofErrors(t *testing.T) {
	valid, err := json.Marshal(testExecutionPayload([]byte{}))
	require.NoError(t, err)

	tests := []struct {
		name  string
		proof string
		err   string
	}{
		{name: "Missing", proof: "", err: "recursive_stark_proof: missing"},
		{name: "NoPrefix", proof: `"abcd"`, err: "recursive_stark_proof: invalid prefix"},
		{name: "OddLength", proof: `"0xabc"`, err: "recursive_stark_proof: odd length"},
		{name: "NotHex", proof: `"0xzz"`, err: "recursive_stark_proof: encoding/hex: invalid byte: U+007A 'z'"},
		{
			name:  "TooLarge",
			proof: `"0x` + strings.Repeat("00", daisugi.MaxRecursiveProofBytes+1) + `"`,
			err:   "recursive_stark_proof: incorrect length",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var raw map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(valid, &raw))

			if test.proof == "" {
				delete(raw, "recursive_stark_proof")
			} else {
				raw["recursive_stark_proof"] = json.RawMessage(test.proof)
			}

			input, err := json.Marshal(raw)
			require.NoError(t, err)

			var payload daisugi.ExecutionPayload
			require.EqualError(t, json.Unmarshal(input, &payload), test.err)
		})
	}
}

func TestExecutionPayloadSSZProofTooLarge(t *testing.T) {
	payload := testExecutionPayload(make([]byte, daisugi.MaxRecursiveProofBytes+1))

	_, err := payload.MarshalSSZ()
	require.Error(t, err)
}

func TestExecutionPayloadHeaderRoundTrip(t *testing.T) {
	header := &daisugi.ExecutionPayloadHeader{
		ParentHash:              phase0.Hash32{0x01},
		FeeRecipient:            bellatrix.ExecutionAddress{0x02},
		StateRoot:               phase0.Root{0x03},
		ReceiptsRoot:            phase0.Root{0x04},
		LogsBloom:               [256]byte{0x05},
		PrevRandao:              [32]byte{0x06},
		BlockNumber:             7,
		GasLimit:                60000000,
		GasUsed:                 21000,
		Timestamp:               1700000212,
		ExtraData:               []byte("daisugi"),
		BaseFeePerGas:           uint256.NewInt(7),
		BlockHash:               phase0.Hash32{0x08},
		TransactionsRoot:        phase0.Root{0x09},
		WithdrawalsRoot:         phase0.Root{0x0a},
		RecursiveStarkProofRoot: phase0.Root{0x0b},
	}

	root, err := header.HashTreeRoot()
	require.NoError(t, err)

	jsonData, err := json.Marshal(header)
	require.NoError(t, err)
	require.Contains(t, string(jsonData),
		`"recursive_stark_proof_root":"0x0b00000000000000000000000000000000000000000000000000000000000000"`)

	var fromJSON daisugi.ExecutionPayloadHeader
	require.NoError(t, json.Unmarshal(jsonData, &fromJSON))
	require.Equal(t, header.RecursiveStarkProofRoot, fromJSON.RecursiveStarkProofRoot)

	yamlData, err := yaml.Marshal(header)
	require.NoError(t, err)

	var fromYAML daisugi.ExecutionPayloadHeader
	require.NoError(t, yaml.Unmarshal(yamlData, &fromYAML))
	require.Equal(t, header.RecursiveStarkProofRoot, fromYAML.RecursiveStarkProofRoot)

	sszData, err := header.MarshalSSZ()
	require.NoError(t, err)

	var fromSSZ daisugi.ExecutionPayloadHeader
	require.NoError(t, fromSSZ.UnmarshalSSZ(sszData))

	for name, decoded := range map[string]*daisugi.ExecutionPayloadHeader{
		"json": &fromJSON, "yaml": &fromYAML, "ssz": &fromSSZ,
	} {
		decodedRoot, err := decoded.HashTreeRoot()
		require.NoError(t, err, name)
		require.Equal(t, root, decodedRoot, name)
	}

	header.RecursiveStarkProofRoot = phase0.Root{}

	changedRoot, err := header.HashTreeRoot()
	require.NoError(t, err)
	require.NotEqual(t, root, changedRoot)
}
