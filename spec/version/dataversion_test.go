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

package version_test

import (
	"encoding/json"
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec/version"
	"github.com/stretchr/testify/require"
)

func TestDataVersionJSON(t *testing.T) {
	tests := []struct {
		name    string
		version version.DataVersion
		str     string
	}{
		{name: "Electra", version: version.DataVersionElectra, str: "electra"},
		{name: "Fulu", version: version.DataVersionFulu, str: "fulu"},
		{name: "Gloas", version: version.DataVersionGloas, str: "gloas"},
		{name: "Heze", version: version.DataVersionHeze, str: "heze"},
		{name: "Daisugi", version: version.DataVersionDaisugi, str: "daisugi"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.str, test.version.String())

			encoded, err := json.Marshal(&test.version)
			require.NoError(t, err)
			require.JSONEq(t, `"`+test.str+`"`, string(encoded))

			var decoded version.DataVersion
			require.NoError(t, json.Unmarshal(encoded, &decoded))
			require.Equal(t, test.version, decoded)

			fromString, err := version.DataVersionFromString(test.str)
			require.NoError(t, err)
			require.Equal(t, test.version, fromString)
		})
	}
}

func TestDataVersionUnknown(t *testing.T) {
	_, err := version.DataVersionFromString("nosuchfork")
	require.EqualError(t, err, `unrecognised data version "nosuchfork"`)

	require.Equal(t, "unknown", (version.DataVersionDaisugi + 1).String())
}

// TestDataVersionValuesStable guards the numeric values: they are persisted and compared by
// downstream users, so a new version must be appended and never inserted.
func TestDataVersionValuesStable(t *testing.T) {
	require.Equal(t, version.DataVersion(6), version.DataVersionElectra)
	require.Equal(t, version.DataVersion(7), version.DataVersionFulu)
	require.Equal(t, version.DataVersion(8), version.DataVersionGloas)
	require.Equal(t, version.DataVersion(9), version.DataVersionHeze)
	require.Equal(t, version.DataVersion(10), version.DataVersionDaisugi)
}
