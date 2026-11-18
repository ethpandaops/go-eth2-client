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

	"github.com/ethpandaops/go-eth2-client/spec/daisugi"
	"github.com/stretchr/testify/require"
)

// TestSignedBeaconBlockDevnetFixture decodes the first Daisugi block of a local
// Lighthouse/Nethermind devnet and checks it against the SSZ encoding and block root
// reported by the beacon node.
func TestSignedBeaconBlockDevnetFixture(t *testing.T) {
	jsonData, err := os.ReadFile("testdata/signedbeaconblock.json")
	require.NoError(t, err)

	sszData, err := os.ReadFile("testdata/signedbeaconblock.ssz")
	require.NoError(t, err)

	expectedRoot, err := os.ReadFile("testdata/signedbeaconblock.root")
	require.NoError(t, err)

	var fromJSON daisugi.SignedBeaconBlock
	require.NoError(t, json.Unmarshal(jsonData, &fromJSON))

	var fromSSZ daisugi.SignedBeaconBlock
	require.NoError(t, fromSSZ.UnmarshalSSZ(sszData))

	for name, block := range map[string]*daisugi.SignedBeaconBlock{"json": &fromJSON, "ssz": &fromSSZ} {
		t.Run(name, func(t *testing.T) {
			root, err := block.Message.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, strings.TrimSpace(string(expectedRoot)), fmt.Sprintf("%#x", root))

			encoded, err := block.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, sszData, encoded)

			// The library prints fee recipients checksummed, the beacon node in lower case.
			reencoded, err := json.Marshal(block)
			require.NoError(t, err)
			require.JSONEq(t, strings.ToLower(string(jsonData)), strings.ToLower(string(reencoded)))
		})
	}
}
