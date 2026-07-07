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
	apiv1bellatrix "github.com/ethpandaops/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/ethpandaops/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/ethpandaops/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/ethpandaops/go-eth2-client/api/v1/electra"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	"github.com/stretchr/testify/require"
)

// blindedBeaconBlockTests enumerates the per-fork views the agnostic
// BlindedBeaconBlock must be wire-compatible with.
func blindedBeaconBlockTests() []struct {
	name    string
	version version.DataVersion
	view    any
} {
	return []struct {
		name    string
		version version.DataVersion
		view    any
	}{
		{
			name:    "Bellatrix",
			version: version.DataVersionBellatrix,
			view: &apiv1bellatrix.BlindedBeaconBlock{
				Slot:          12345,
				ProposerIndex: 42,
				ParentRoot:    phase0.Root{0x11},
				StateRoot:     phase0.Root{0x12},
				Body:          testBellatrixBlindedBody(),
			},
		},
		{
			name:    "Capella",
			version: version.DataVersionCapella,
			view: &apiv1capella.BlindedBeaconBlock{
				Slot:          12345,
				ProposerIndex: 42,
				ParentRoot:    phase0.Root{0x11},
				StateRoot:     phase0.Root{0x12},
				Body:          testCapellaBlindedBody(),
			},
		},
		{
			name:    "Deneb",
			version: version.DataVersionDeneb,
			view: &apiv1deneb.BlindedBeaconBlock{
				Slot:          12345,
				ProposerIndex: 42,
				ParentRoot:    phase0.Root{0x11},
				StateRoot:     phase0.Root{0x12},
				Body:          testDenebBlindedBody(),
			},
		},
		{
			name:    "Electra",
			version: version.DataVersionElectra,
			view: &apiv1electra.BlindedBeaconBlock{
				Slot:          12345,
				ProposerIndex: 42,
				ParentRoot:    phase0.Root{0x11},
				StateRoot:     phase0.Root{0x12},
				Body:          testElectraBlindedBody(),
			},
		},
		{
			// Fulu reuses the Electra wire schema; the agnostic type must
			// accept the Electra view when pinned to the Fulu version.
			name:    "Fulu",
			version: version.DataVersionFulu,
			view: &apiv1electra.BlindedBeaconBlock{
				Slot:          12345,
				ProposerIndex: 42,
				ParentRoot:    phase0.Root{0x11},
				StateRoot:     phase0.Root{0x12},
				Body:          testElectraBlindedBody(),
			},
		},
	}
}

// TestBlindedBeaconBlockJSONWireCompat verifies that the agnostic type
// marshals to byte-identical JSON versus the per-fork type, and that JSON
// round-trips losslessly through the agnostic type.
func TestBlindedBeaconBlockJSONWireCompat(t *testing.T) {
	for _, test := range blindedBeaconBlockTests() {
		t.Run(test.name, func(t *testing.T) {
			expected, err := json.Marshal(test.view)
			require.NoError(t, err)

			agnostic := &apiv1all.BlindedBeaconBlock{Version: test.version}
			require.NoError(t, agnostic.FromView(test.view))
			require.Equal(t, test.version, agnostic.Version)

			got, err := json.Marshal(agnostic)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(got),
				"agnostic JSON differs from per-fork JSON")

			// Round-trip through UnmarshalJSON with Version pre-set.
			rt := &apiv1all.BlindedBeaconBlock{Version: test.version}
			require.NoError(t, json.Unmarshal(expected, rt))

			rtJSON, err := json.Marshal(rt)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(rtJSON),
				"round-tripped JSON differs from per-fork JSON")

			// Version must propagate to nested versionable children.
			require.NotNil(t, rt.Body)
			require.Equal(t, test.version, rt.Body.Version)
		})
	}
}

// TestBlindedBeaconBlockSSZWireCompat verifies that the agnostic type
// produces byte-identical SSZ and the same hash tree root as the per-fork
// type, and that SSZ round-trips losslessly through the agnostic type.
func TestBlindedBeaconBlockSSZWireCompat(t *testing.T) {
	for _, test := range blindedBeaconBlockTests() {
		t.Run(test.name, func(t *testing.T) {
			codec, ok := test.view.(sszCodec)
			require.True(t, ok)

			expected, err := codec.MarshalSSZ()
			require.NoError(t, err)

			agnostic := &apiv1all.BlindedBeaconBlock{Version: test.version}
			require.NoError(t, agnostic.FromView(test.view))

			got, err := agnostic.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expected, got, "agnostic SSZ differs from per-fork SSZ")

			expectedRoot, err := codec.HashTreeRoot()
			require.NoError(t, err)

			gotRoot, err := agnostic.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, expectedRoot, gotRoot,
				"agnostic hash tree root differs from per-fork root")

			// Round-trip through UnmarshalSSZ with Version pre-set.
			rt := &apiv1all.BlindedBeaconBlock{Version: test.version}
			require.NoError(t, rt.UnmarshalSSZ(expected))

			rtSSZ, err := rt.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expected, rtSSZ,
				"round-tripped SSZ differs from per-fork SSZ")
		})
	}
}

// TestBlindedBeaconBlockViewRoundtrip verifies ToView reproduces the
// per-fork view the agnostic instance was built from, and that ToVersioned
// wires the block into api.VersionedBlindedBeaconBlock.
func TestBlindedBeaconBlockViewRoundtrip(t *testing.T) {
	for _, test := range blindedBeaconBlockTests() {
		t.Run(test.name, func(t *testing.T) {
			agnostic := &apiv1all.BlindedBeaconBlock{Version: test.version}
			require.NoError(t, agnostic.FromView(test.view))

			view, err := agnostic.ToView()
			require.NoError(t, err)
			require.IsType(t, test.view, view)
			require.Equal(t, test.view, view)

			versioned, err := agnostic.ToVersioned()
			require.NoError(t, err)
			require.Equal(t, test.version, versioned.Version)

			rt := &apiv1all.BlindedBeaconBlock{}
			require.NoError(t, rt.FromVersioned(versioned))
			require.Equal(t, agnostic, rt)
		})
	}
}

// TestBlindedBeaconBlockBodyWireCompat verifies the standalone agnostic
// blinded body type against the per-fork bodies.
func TestBlindedBeaconBlockBodyWireCompat(t *testing.T) {
	tests := []struct {
		name    string
		version version.DataVersion
		view    any
	}{
		{name: "Bellatrix", version: version.DataVersionBellatrix, view: testBellatrixBlindedBody()},
		{name: "Capella", version: version.DataVersionCapella, view: testCapellaBlindedBody()},
		{name: "Deneb", version: version.DataVersionDeneb, view: testDenebBlindedBody()},
		{name: "Electra", version: version.DataVersionElectra, view: testElectraBlindedBody()},
		{name: "Fulu", version: version.DataVersionFulu, view: testElectraBlindedBody()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expected, err := json.Marshal(test.view)
			require.NoError(t, err)

			agnostic := &apiv1all.BlindedBeaconBlockBody{Version: test.version}
			require.NoError(t, agnostic.FromView(test.view))

			got, err := json.Marshal(agnostic)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(got),
				"agnostic JSON differs from per-fork JSON")

			view, err := agnostic.ToView()
			require.NoError(t, err)
			require.IsType(t, test.view, view)
			require.Equal(t, test.view, view)

			codec, ok := test.view.(sszCodec)
			require.True(t, ok)

			expectedSSZ, err := codec.MarshalSSZ()
			require.NoError(t, err)

			gotSSZ, err := agnostic.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expectedSSZ, gotSSZ,
				"agnostic SSZ differs from per-fork SSZ")

			rt := &apiv1all.BlindedBeaconBlockBody{Version: test.version}
			require.NoError(t, rt.UnmarshalSSZ(expectedSSZ))

			rtSSZ, err := rt.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expectedSSZ, rtSSZ,
				"round-tripped SSZ differs from per-fork SSZ")
		})
	}
}
