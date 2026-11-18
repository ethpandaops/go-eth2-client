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
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	apiv1daisugi "github.com/ethpandaops/go-eth2-client/api/v1/daisugi"
	"github.com/stretchr/testify/require"
)

// TestSignedBlindedBeaconBlockDevnetFixture decodes the blinded form of a devnet block that
// carries a recursive STARK proof. Its root must equal the root of the full block, which
// only holds when the header commits to the proof the same way the payload does.
func TestSignedBlindedBeaconBlockDevnetFixture(t *testing.T) {
	jsonData, err := os.ReadFile("testdata/signedblindedbeaconblock_proof.json")
	require.NoError(t, err)

	expectedRoot, err := os.ReadFile("testdata/signedblindedbeaconblock_proof.root")
	require.NoError(t, err)

	var block apiv1daisugi.SignedBlindedBeaconBlock
	require.NoError(t, json.Unmarshal(jsonData, &block))
	require.False(t, block.Message.Body.ExecutionPayloadHeader.RecursiveStarkProofRoot.IsZero())

	root, err := block.Message.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, strings.TrimSpace(string(expectedRoot)), fmt.Sprintf("%#x", root))

	sszData, err := block.MarshalSSZ()
	require.NoError(t, err)

	var fromSSZ apiv1daisugi.SignedBlindedBeaconBlock
	require.NoError(t, fromSSZ.UnmarshalSSZ(sszData))

	sszRoot, err := fromSSZ.Message.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, root, sszRoot)

	reencoded, err := json.Marshal(&block)
	require.NoError(t, err)
	require.JSONEq(t, strings.ToLower(string(jsonData)), strings.ToLower(string(reencoded)))
}
