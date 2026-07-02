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

package all

import (
	"errors"
	"fmt"

	"github.com/ethpandaops/go-eth2-client/api"
	"github.com/ethpandaops/go-eth2-client/spec/all"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/version"
)

// ProposalFromSignedBlock wraps a post-Gloas signed beacon block into the
// versioned signed proposal consumed by the block submission API. From Gloas
// onwards blocks are submitted bare: EIP-7732 separates the execution payload
// (and its blobs) from the block, so there is no block-contents wrapper.
// Pre-Gloas forks submit SignedBlockContents instead — use
// SignedBlockContents.ToVersioned for those.
func ProposalFromSignedBlock(block *all.SignedBeaconBlock) (*api.VersionedSignedProposal, error) {
	if block == nil {
		return nil, errors.New("ProposalFromSignedBlock: nil block")
	}

	view, err := block.ToView()
	if err != nil {
		return nil, err
	}

	gloasBlock, ok := view.(*gloas.SignedBeaconBlock)
	if !ok {
		return nil, fmt.Errorf("ProposalFromSignedBlock: unsupported version %s", block.Version)
	}

	proposal := &api.VersionedSignedProposal{
		Version: block.Version,
	}

	switch block.Version {
	case version.DataVersionGloas:
		proposal.Gloas = gloasBlock
	case version.DataVersionHeze:
		proposal.Heze = gloasBlock
	default:
		return nil, fmt.Errorf("ProposalFromSignedBlock: unsupported version %s", block.Version)
	}

	return proposal, nil
}
