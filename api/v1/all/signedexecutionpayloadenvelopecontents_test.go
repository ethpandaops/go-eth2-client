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

package all_test

import (
	"encoding/json"
	"testing"

	apiv1all "github.com/ethpandaops/go-eth2-client/api/v1/all"
	apiv1gloas "github.com/ethpandaops/go-eth2-client/api/v1/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

// testGloasSignedExecutionPayloadEnvelope builds a gloas signed execution
// payload envelope with realistic-ish data.
func testGloasSignedExecutionPayloadEnvelope() *gloas.SignedExecutionPayloadEnvelope {
	return &gloas.SignedExecutionPayloadEnvelope{
		Message: &gloas.ExecutionPayloadEnvelope{
			Payload: &gloas.ExecutionPayload{
				ParentHash:    phase0.Hash32{0x01},
				FeeRecipient:  bellatrix.ExecutionAddress{0x02},
				StateRoot:     phase0.Root{0x03},
				ReceiptsRoot:  phase0.Root{0x04},
				LogsBloom:     [256]byte{0x05},
				PrevRandao:    [32]byte{0x06},
				BlockNumber:   7,
				GasLimit:      30_000_000,
				GasUsed:       21_000,
				Timestamp:     1_600_000_000,
				ExtraData:     []byte{0x08},
				BaseFeePerGas: uint256.NewInt(9),
				BlockHash:     phase0.Hash32{0x0a},
				Transactions:  []bellatrix.Transaction{{0x0b, 0x0c}},
				Withdrawals: []*capella.Withdrawal{
					{Index: 1, ValidatorIndex: 2, Address: bellatrix.ExecutionAddress{0x0d}, Amount: 3},
				},
				BlobGasUsed:     131072,
				ExcessBlobGas:   262144,
				BlockAccessList: gloas.BlockAccessList{0x0e},
				SlotNumber:      12345,
			},
			ExecutionRequests: &gloas.ExecutionRequests{
				Deposits:        []*electra.DepositRequest{},
				Withdrawals:     []*electra.WithdrawalRequest{},
				Consolidations:  []*electra.ConsolidationRequest{},
				BuilderDeposits: []*gloas.BuilderDepositRequest{},
				BuilderExits:    []*gloas.BuilderExitRequest{},
			},
			BuilderIndex:          42,
			BeaconBlockRoot:       phase0.Root{0x11},
			ParentBeaconBlockRoot: phase0.Root{0x12},
		},
		Signature: phase0.BLSSignature{0x13},
	}
}

// signedExecutionPayloadEnvelopeContentsTests enumerates the versions the
// agnostic SignedExecutionPayloadEnvelopeContents must be wire-compatible
// with. Gloas and Heze share the gloas view.
func signedExecutionPayloadEnvelopeContentsTests() []struct {
	name    string
	version version.DataVersion
	view    *apiv1gloas.SignedExecutionPayloadEnvelopeContents
} {
	return []struct {
		name    string
		version version.DataVersion
		view    *apiv1gloas.SignedExecutionPayloadEnvelopeContents
	}{
		{
			name:    "Gloas",
			version: version.DataVersionGloas,
			view: &apiv1gloas.SignedExecutionPayloadEnvelopeContents{
				SignedExecutionPayloadEnvelope: testGloasSignedExecutionPayloadEnvelope(),
				KZGProofs:                      testKZGProofs(),
				Blobs:                          testBlobs(),
			},
		},
		{
			name:    "Heze",
			version: version.DataVersionHeze,
			view: &apiv1gloas.SignedExecutionPayloadEnvelopeContents{
				SignedExecutionPayloadEnvelope: testGloasSignedExecutionPayloadEnvelope(),
				KZGProofs:                      testKZGProofs(),
				Blobs:                          testBlobs(),
			},
		},
	}
}

// TestSignedExecutionPayloadEnvelopeContentsJSONWireCompat verifies that the
// agnostic type marshals to byte-identical JSON versus the per-fork type,
// and that JSON round-trips losslessly through the agnostic type.
func TestSignedExecutionPayloadEnvelopeContentsJSONWireCompat(t *testing.T) {
	for _, test := range signedExecutionPayloadEnvelopeContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			expected, err := json.Marshal(test.view)
			require.NoError(t, err)

			agnostic := &apiv1all.SignedExecutionPayloadEnvelopeContents{Version: test.version}
			require.NoError(t, agnostic.FromView(test.view))
			require.Equal(t, test.version, agnostic.Version)

			got, err := json.Marshal(agnostic)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(got),
				"agnostic JSON differs from per-fork JSON")

			// Round-trip through UnmarshalJSON with Version pre-set.
			rt := &apiv1all.SignedExecutionPayloadEnvelopeContents{Version: test.version}
			require.NoError(t, json.Unmarshal(expected, rt))

			rtJSON, err := json.Marshal(rt)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(rtJSON),
				"round-tripped JSON differs from per-fork JSON")

			// Version must propagate to the nested versionable envelope.
			require.NotNil(t, rt.SignedExecutionPayloadEnvelope)
			require.Equal(t, test.version, rt.SignedExecutionPayloadEnvelope.Version)
			require.NotNil(t, rt.SignedExecutionPayloadEnvelope.Message)
			require.Equal(t, test.version, rt.SignedExecutionPayloadEnvelope.Message.Version)
		})
	}
}

// TestSignedExecutionPayloadEnvelopeContentsSSZWireCompat verifies that the
// agnostic type produces byte-identical SSZ and the same hash tree root as
// the per-fork type, and that SSZ round-trips losslessly through the
// agnostic type.
func TestSignedExecutionPayloadEnvelopeContentsSSZWireCompat(t *testing.T) {
	for _, test := range signedExecutionPayloadEnvelopeContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			expected, err := test.view.MarshalSSZ()
			require.NoError(t, err)

			agnostic := &apiv1all.SignedExecutionPayloadEnvelopeContents{Version: test.version}
			require.NoError(t, agnostic.FromView(test.view))

			got, err := agnostic.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expected, got, "agnostic SSZ differs from per-fork SSZ")

			expectedRoot, err := test.view.HashTreeRoot()
			require.NoError(t, err)

			gotRoot, err := agnostic.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, expectedRoot, gotRoot,
				"agnostic hash tree root differs from per-fork root")

			// Round-trip through UnmarshalSSZ with Version pre-set.
			rt := &apiv1all.SignedExecutionPayloadEnvelopeContents{Version: test.version}
			require.NoError(t, rt.UnmarshalSSZ(expected))

			rtSSZ, err := rt.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expected, rtSSZ,
				"round-tripped SSZ differs from per-fork SSZ")

			// Version must propagate to the nested versionable envelope.
			require.NotNil(t, rt.SignedExecutionPayloadEnvelope)
			require.Equal(t, test.version, rt.SignedExecutionPayloadEnvelope.Version)
		})
	}
}
