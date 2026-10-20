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

package spec_test

import (
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionedAttesterSlashingAttestations(t *testing.T) {
	data := func(slot phase0.Slot) *phase0.AttestationData {
		return &phase0.AttestationData{
			Slot:   slot,
			Source: &phase0.Checkpoint{},
			Target: &phase0.Checkpoint{},
		}
	}

	phase0Slashing := &phase0.AttesterSlashing{
		Attestation1: &phase0.IndexedAttestation{AttestingIndices: []uint64{1, 2}, Data: data(10)},
		Attestation2: &phase0.IndexedAttestation{AttestingIndices: []uint64{2, 3}, Data: data(11)},
	}
	electraSlashing := &electra.AttesterSlashing{
		Attestation1: &electra.IndexedAttestation{AttestingIndices: []uint64{1, 2}, Data: data(10)},
		Attestation2: &electra.IndexedAttestation{AttestingIndices: []uint64{2, 3}, Data: data(11)},
	}
	gloasSlashing := &gloas.AttesterSlashing{
		Attestation1: &gloas.IndexedAttestation{AttestingIndices: []uint64{1, 2}, Data: data(10)},
		Attestation2: &gloas.IndexedAttestation{AttestingIndices: []uint64{2, 3}, Data: data(11)},
	}

	tests := []spec.VersionedAttesterSlashing{
		{Version: spec.DataVersionPhase0, Phase0: phase0Slashing},
		{Version: spec.DataVersionAltair, Altair: phase0Slashing},
		{Version: spec.DataVersionBellatrix, Bellatrix: phase0Slashing},
		{Version: spec.DataVersionCapella, Capella: phase0Slashing},
		{Version: spec.DataVersionDeneb, Deneb: phase0Slashing},
		{Version: spec.DataVersionElectra, Electra: electraSlashing},
		{Version: spec.DataVersionFulu, Fulu: electraSlashing},
		{Version: spec.DataVersionGloas, Gloas: gloasSlashing},
		{Version: spec.DataVersionHeze, Heze: gloasSlashing},
	}

	for _, slashing := range tests {
		t.Run(slashing.Version.String(), func(t *testing.T) {
			att1, err := slashing.Attestation1()
			require.NoError(t, err)
			assert.Equal(t, slashing.Version, att1.Version)

			indices, err := att1.AttestingIndices()
			require.NoError(t, err)
			assert.Equal(t, []uint64{1, 2}, indices)

			attData, err := att1.Data()
			require.NoError(t, err)
			assert.Equal(t, phase0.Slot(10), attData.Slot)

			_, err = att1.Signature()
			require.NoError(t, err)

			att2, err := slashing.Attestation2()
			require.NoError(t, err)
			assert.Equal(t, slashing.Version, att2.Version)

			indices, err = att2.AttestingIndices()
			require.NoError(t, err)
			assert.Equal(t, []uint64{2, 3}, indices)

			attData, err = att2.Data()
			require.NoError(t, err)
			assert.Equal(t, phase0.Slot(11), attData.Slot)
		})
	}
}
