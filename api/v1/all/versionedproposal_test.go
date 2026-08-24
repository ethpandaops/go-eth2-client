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
	"testing"

	apiv1all "github.com/ethpandaops/go-eth2-client/api/v1/all"
	"github.com/ethpandaops/go-eth2-client/spec/all"
	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	"github.com/stretchr/testify/require"
)

// testGloasSignedBeaconBlock builds a minimal gloas signed beacon block.
func testGloasSignedBeaconBlock() *gloas.SignedBeaconBlock {
	return &gloas.SignedBeaconBlock{
		Message: &gloas.BeaconBlock{
			Slot:          12345,
			ProposerIndex: 42,
			ParentRoot:    phase0.Root{0x01},
			StateRoot:     phase0.Root{0x02},
			Body: &gloas.BeaconBlockBody{
				RANDAOReveal: phase0.BLSSignature{0x03},
				Graffiti:     [32]byte{0x04},
			},
		},
		Signature: phase0.BLSSignature{0x05},
	}
}

// TestProposalFromSignedBlock verifies the bare-block-to-proposal mapping for
// post-Gloas forks: the block lands in the field matching its version, with
// no blinded flag and no other fields set.
func TestProposalFromSignedBlock(t *testing.T) {
	tests := []struct {
		name    string
		version version.DataVersion
	}{
		{name: "gloas", version: version.DataVersionGloas},
		{name: "heze", version: version.DataVersionHeze},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			block := &all.SignedBeaconBlock{Version: test.version}
			require.NoError(t, block.FromView(testGloasSignedBeaconBlock()))
			block.Version = test.version

			// Every fork converts to the view of its own package, so the
			// expected value comes from the block rather than from the seed.
			view, err := block.ToView()
			require.NoError(t, err)

			proposal, err := apiv1all.ProposalFromSignedBlock(block)
			require.NoError(t, err)
			require.Equal(t, test.version, proposal.Version)
			require.False(t, proposal.Blinded)

			switch test.version {
			case version.DataVersionGloas:
				require.NotNil(t, proposal.Gloas)
				require.Nil(t, proposal.Heze)
				require.Equal(t, view, proposal.Gloas)
			case version.DataVersionHeze:
				require.NotNil(t, proposal.Heze)
				require.Nil(t, proposal.Gloas)
				require.Equal(t, view, proposal.Heze)
			default:
				t.Fatalf("unexpected version %s", test.version)
			}
		})
	}
}

// TestProposalFromSignedBlockPreDeneb verifies the bare-block-to-proposal
// mapping for the pre-blobs forks: bellatrix and capella proposals are plain
// signed beacon blocks.
func TestProposalFromSignedBlockPreDeneb(t *testing.T) {
	t.Run("bellatrix", func(t *testing.T) {
		view := &bellatrix.SignedBeaconBlock{
			Message: &bellatrix.BeaconBlock{
				Slot:          123,
				ProposerIndex: 7,
				ParentRoot:    phase0.Root{0x01},
				StateRoot:     phase0.Root{0x02},
				Body: &bellatrix.BeaconBlockBody{
					RANDAOReveal:     phase0.BLSSignature{0x03},
					Graffiti:         [32]byte{0x04},
					ExecutionPayload: &bellatrix.ExecutionPayload{BlockNumber: 9},
				},
			},
			Signature: phase0.BLSSignature{0x05},
		}

		block := &all.SignedBeaconBlock{Version: version.DataVersionBellatrix}
		require.NoError(t, block.FromView(view))

		proposal, err := apiv1all.ProposalFromSignedBlock(block)
		require.NoError(t, err)
		require.Equal(t, version.DataVersionBellatrix, proposal.Version)
		require.False(t, proposal.Blinded)
		require.NotNil(t, proposal.Bellatrix)
		require.Nil(t, proposal.Capella)
		require.Equal(t, view, proposal.Bellatrix)
	})

	t.Run("capella", func(t *testing.T) {
		view := &capella.SignedBeaconBlock{
			Message: &capella.BeaconBlock{
				Slot:          456,
				ProposerIndex: 8,
				ParentRoot:    phase0.Root{0x11},
				StateRoot:     phase0.Root{0x12},
				Body: &capella.BeaconBlockBody{
					RANDAOReveal:     phase0.BLSSignature{0x13},
					Graffiti:         [32]byte{0x14},
					ExecutionPayload: &capella.ExecutionPayload{BlockNumber: 10},
				},
			},
			Signature: phase0.BLSSignature{0x15},
		}

		block := &all.SignedBeaconBlock{Version: version.DataVersionCapella}
		require.NoError(t, block.FromView(view))

		proposal, err := apiv1all.ProposalFromSignedBlock(block)
		require.NoError(t, err)
		require.Equal(t, version.DataVersionCapella, proposal.Version)
		require.False(t, proposal.Blinded)
		require.NotNil(t, proposal.Capella)
		require.Nil(t, proposal.Bellatrix)
		require.Equal(t, view, proposal.Capella)
	})
}

// TestProposalFromSignedBlockUnsupportedVersion verifies blob-carrying
// versions (Deneb through Fulu) are rejected: those forks submit
// SignedBlockContents, not bare blocks.
func TestProposalFromSignedBlockUnsupportedVersion(t *testing.T) {
	for _, v := range []version.DataVersion{
		version.DataVersionDeneb,
		version.DataVersionElectra,
		version.DataVersionFulu,
	} {
		block := &all.SignedBeaconBlock{Version: v}
		_, err := apiv1all.ProposalFromSignedBlock(block)
		require.ErrorContains(t, err, "unsupported version", "version %s", v)
	}

	_, err := apiv1all.ProposalFromSignedBlock(nil)
	require.ErrorContains(t, err, "nil block")
}
