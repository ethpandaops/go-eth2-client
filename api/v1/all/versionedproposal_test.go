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

// TestProposalFromSignedBlock verifies the post-Gloas block-to-proposal
// mapping: the bare block lands in the field matching its version, with no
// blinded flag and no other fields set.
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
			view := testGloasSignedBeaconBlock()

			block := &all.SignedBeaconBlock{Version: test.version}
			require.NoError(t, block.FromView(view))
			block.Version = test.version

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

// TestProposalFromSignedBlockUnsupportedVersion verifies pre-Gloas versions
// are rejected: those forks submit SignedBlockContents, not bare blocks.
func TestProposalFromSignedBlockUnsupportedVersion(t *testing.T) {
	block := &all.SignedBeaconBlock{Version: version.DataVersionElectra}
	_, err := apiv1all.ProposalFromSignedBlock(block)
	require.ErrorContains(t, err, "unsupported version")

	_, err = apiv1all.ProposalFromSignedBlock(nil)
	require.ErrorContains(t, err, "nil block")
}
