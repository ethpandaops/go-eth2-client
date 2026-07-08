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

package v1

import (
	"fmt"
)

// NodeIdentity contains the node's network identity.
type NodeIdentity struct {
	// PeerID is the node's libp2p peer ID.
	PeerID string `json:"peer_id"`
	// ENR is the node's Ethereum node record.
	ENR string `json:"enr"`
	// P2PAddresses are the multiaddrs on which eth2 RPC requests are served.
	P2PAddresses []string `json:"p2p_addresses"`
	// DiscoveryAddresses are the multiaddrs on which the node listens for
	// discv5 requests.
	DiscoveryAddresses []string `json:"discovery_addresses"`
	// Metadata is the node's libp2p metadata.
	Metadata NodeIdentityMetadata `json:"metadata"`
}

// NodeIdentityMetadata is the node's libp2p metadata.
type NodeIdentityMetadata struct {
	// SeqNumber versions the node's metadata (decimal string).
	SeqNumber string `json:"seq_number"`
	// Attnets is the hex bitvector of persistent attestation subnet
	// subscriptions.
	Attnets string `json:"attnets"`
	// Syncnets is the hex bitvector of sync committee subnet subscriptions
	// (Altair onwards).
	Syncnets string `json:"syncnets,omitempty"`
	// CustodyGroupCount is the node's custody group count (Fulu onwards,
	// decimal string).
	CustodyGroupCount string `json:"custody_group_count,omitempty"`
}

// String returns a string version of the structure.
func (n *NodeIdentity) String() string {
	return fmt.Sprintf("%s %v", n.PeerID, n.P2PAddresses)
}
