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
	apiv1deneb "github.com/ethpandaops/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/ethpandaops/go-eth2-client/api/v1/electra"
	apiv1fulu "github.com/ethpandaops/go-eth2-client/api/v1/fulu"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	"github.com/stretchr/testify/require"
)

// blockContentsTests enumerates the per-fork views the agnostic BlockContents
// must be wire-compatible with.
func blockContentsTests() []struct {
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
			name:    "Deneb",
			version: version.DataVersionDeneb,
			view: &apiv1deneb.BlockContents{
				Block:     testDenebSignedBeaconBlock().Message,
				KZGProofs: testKZGProofs(),
				Blobs:     testBlobs(),
			},
		},
		{
			name:    "Electra",
			version: version.DataVersionElectra,
			view: &apiv1electra.BlockContents{
				Block:     testElectraSignedBeaconBlock().Message,
				KZGProofs: testKZGProofs(),
				Blobs:     testBlobs(),
			},
		},
		{
			name:    "Fulu",
			version: version.DataVersionFulu,
			view: &apiv1fulu.BlockContents{
				Block:     testElectraSignedBeaconBlock().Message,
				KZGProofs: testKZGProofs(),
				Blobs:     testBlobs(),
			},
		},
	}
}

// TestBlockContentsJSONWireCompat verifies that the agnostic type marshals to
// byte-identical JSON versus the per-fork type, and that JSON round-trips
// losslessly through the agnostic type.
func TestBlockContentsJSONWireCompat(t *testing.T) {
	for _, test := range blockContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			expected, err := json.Marshal(test.view)
			require.NoError(t, err)

			agnostic := &apiv1all.BlockContents{}
			require.NoError(t, agnostic.FromView(test.view))
			require.Equal(t, test.version, agnostic.Version)

			got, err := json.Marshal(agnostic)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(got),
				"agnostic JSON differs from per-fork JSON")

			// Round-trip through UnmarshalJSON with Version pre-set.
			rt := &apiv1all.BlockContents{Version: test.version}
			require.NoError(t, json.Unmarshal(expected, rt))

			rtJSON, err := json.Marshal(rt)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(rtJSON),
				"round-tripped JSON differs from per-fork JSON")

			// Version must propagate to nested versionable children.
			require.NotNil(t, rt.Block)
			require.Equal(t, test.version, rt.Block.Version)
			require.NotNil(t, rt.Block.Body)
			require.Equal(t, test.version, rt.Block.Body.Version)
		})
	}
}

// TestBlockContentsSSZWireCompat verifies that the agnostic type produces
// byte-identical SSZ and the same hash tree root as the per-fork type, and
// that SSZ round-trips losslessly through the agnostic type.
func TestBlockContentsSSZWireCompat(t *testing.T) {
	for _, test := range blockContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			codec, ok := test.view.(sszCodec)
			require.True(t, ok)

			expected, err := codec.MarshalSSZ()
			require.NoError(t, err)

			agnostic := &apiv1all.BlockContents{}
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
			rt := &apiv1all.BlockContents{Version: test.version}
			require.NoError(t, rt.UnmarshalSSZ(expected))

			rtSSZ, err := rt.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expected, rtSSZ,
				"round-tripped SSZ differs from per-fork SSZ")

			// Version must propagate to nested versionable children.
			require.NotNil(t, rt.Block)
			require.Equal(t, test.version, rt.Block.Version)
		})
	}
}

// TestBlockContentsViewRoundtrip verifies ToView reproduces the per-fork view
// the agnostic instance was built from.
func TestBlockContentsViewRoundtrip(t *testing.T) {
	for _, test := range blockContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			agnostic := &apiv1all.BlockContents{}
			require.NoError(t, agnostic.FromView(test.view))

			view, err := agnostic.ToView()
			require.NoError(t, err)
			require.IsType(t, test.view, view)
			require.Equal(t, test.view, view)
		})
	}
}

// TestBlockContentsToVersioned verifies conversion to and from
// api.VersionedProposal.
func TestBlockContentsToVersioned(t *testing.T) {
	for _, test := range blockContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			agnostic := &apiv1all.BlockContents{}
			require.NoError(t, agnostic.FromView(test.view))

			versioned, err := agnostic.ToVersioned()
			require.NoError(t, err)
			require.Equal(t, test.version, versioned.Version)

			rt := &apiv1all.BlockContents{}
			require.NoError(t, rt.FromVersioned(versioned))
			require.Equal(t, test.version, rt.Version)
			require.Equal(t, agnostic, rt)
		})
	}
}

// TestBlockContentsUnsupportedVersion verifies unsupported versions are
// rejected with a clear error.
func TestBlockContentsUnsupportedVersion(t *testing.T) {
	for _, v := range []version.DataVersion{
		version.DataVersionUnknown,
		version.DataVersionPhase0,
		version.DataVersionCapella,
		version.DataVersionGloas,
	} {
		b := &apiv1all.BlockContents{Version: v}

		_, err := b.ToView()
		require.ErrorContains(t, err, "unsupported version")

		_, err = json.Marshal(b)
		require.ErrorContains(t, err, "unsupported version")

		_, err = b.MarshalSSZ()
		require.ErrorContains(t, err, "unsupported version")
	}
}
