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

package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSingleAttestationPayload(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "GrandineSingleAttestationOnAttestationTopic",
			data: `{"committee_index":"32","attester_index":"1861792","data":{"slot":"14326127","index":"0","beacon_block_root":"0x62ab3e7b1366549bd040a3523f0c2d5a2743164355ac5adbc07d51c1dfd05f88","source":{"epoch":"447690","root":"0x5d419ad7acb52b038214f39171670e4af6b877850fcb1ebf30e18d8f5bb6a372"},"target":{"epoch":"447691","root":"0xdf078bb15f0f0532eb480da59eb63a05025807b77a3d4df5cb598810974a4304"}},"signature":"0x8f5bd89876a78f836a19494d36c10e0c0d6197117f20c4354d74d06022d7e95619889a54c9feadf2c4160f93aee02bb41606b8e70e296d2e39abd1fb67bca01d019ffe6a9714ac7f99273a8acaa91c7f7ce518091051f95134b6017449503472"}`,
			want: true,
		},
		{
			name: "Phase0Attestation",
			data: `{"aggregation_bits":"0x00002840403040000000020008800040008042800000020220","data":{"slot":"4095945","index":"12","beacon_block_root":"0xff27c7551bf4cfe4dc4cce00920e7a5c5074860d1dbd8aa8b3b5f888523f51ff","source":{"epoch":"127997","root":"0x38758fb180459583bd5e8e1a31711eb09e63eb92be974485397e9a2c57de2783"},"target":{"epoch":"127998","root":"0x46d4629861bd81cfc94007501b4edb1b3ca9444b41d7a98681b6c2f4bdb978bd"}},"signature":"0xacb9f562a28c4ef5b60b88678068ea51573a3237d3331dda3b2d845a0d03bc56ab2994d2deb90d9f074a8bdab59945150d0a7717e74b1bf2627f8971c81091f724c211dfce8fa16fb839c6a1bfd341ddec5e7eb88472682fd1a170e373660534"}`,
			want: false,
		},
		{
			name: "ElectraAggregateAttestation",
			data: `{"aggregation_bits":"0xff","committee_bits":"0x0100000000000000","data":{"slot":"1","index":"0","beacon_block_root":"0x0000000000000000000000000000000000000000000000000000000000000000","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"}},"signature":"0x0000"}`,
			want: false,
		},
		{
			name: "InvalidJSON",
			data: `not json`,
			want: false,
		},
		{
			name: "EmptyObject",
			data: `{}`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, isSingleAttestationPayload([]byte(tt.data)))
		})
	}
}
