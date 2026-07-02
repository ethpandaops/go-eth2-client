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

	bitfield "github.com/OffchainLabs/go-bitfield"
	apiv1all "github.com/ethpandaops/go-eth2-client/api/v1/all"
	apiv1bellatrix "github.com/ethpandaops/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/ethpandaops/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/ethpandaops/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/ethpandaops/go-eth2-client/api/v1/electra"
	"github.com/ethpandaops/go-eth2-client/spec/altair"
	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/deneb"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

// testAttestationData builds a phase0 attestation data with realistic-ish data.
func testAttestationData() *phase0.AttestationData {
	return &phase0.AttestationData{
		Slot:            12344,
		Index:           1,
		BeaconBlockRoot: phase0.Root{0x31},
		Source:          &phase0.Checkpoint{Epoch: 384, Root: phase0.Root{0x32}},
		Target:          &phase0.Checkpoint{Epoch: 385, Root: phase0.Root{0x33}},
	}
}

// testPhase0Attestation builds a phase0 attestation.
func testPhase0Attestation() *phase0.Attestation {
	return &phase0.Attestation{
		AggregationBits: bitfield.NewBitlist(128),
		Data:            testAttestationData(),
		Signature:       phase0.BLSSignature{0x34},
	}
}

// testElectraAttestation builds an electra attestation.
func testElectraAttestation() *electra.Attestation {
	return &electra.Attestation{
		AggregationBits: bitfield.NewBitlist(128),
		Data:            testAttestationData(),
		Signature:       phase0.BLSSignature{0x34},
		CommitteeBits:   bitfield.NewBitvector64(),
	}
}

// testBellatrixBlindedBody builds a Bellatrix blinded beacon block body.
func testBellatrixBlindedBody() *apiv1bellatrix.BlindedBeaconBlockBody {
	return &apiv1bellatrix.BlindedBeaconBlockBody{
		RANDAOReveal: phase0.BLSSignature{0x13},
		ETH1Data: &phase0.ETH1Data{
			DepositRoot:  phase0.Root{0x14},
			DepositCount: 15,
			BlockHash:    make([]byte, 32),
		},
		Graffiti:          [32]byte{'b', 'u', 'i', 'l', 'd'},
		ProposerSlashings: []*phase0.ProposerSlashing{},
		AttesterSlashings: []*phase0.AttesterSlashing{},
		Attestations:      []*phase0.Attestation{testPhase0Attestation()},
		Deposits:          []*phase0.Deposit{},
		VoluntaryExits:    []*phase0.SignedVoluntaryExit{},
		SyncAggregate: &altair.SyncAggregate{
			SyncCommitteeBits:      bitfield.NewBitvector512(),
			SyncCommitteeSignature: phase0.BLSSignature{0x16},
		},
		ExecutionPayloadHeader: &bellatrix.ExecutionPayloadHeader{
			ParentHash:       phase0.Hash32{0x01},
			FeeRecipient:     bellatrix.ExecutionAddress{0x02},
			StateRoot:        [32]byte{0x03},
			ReceiptsRoot:     [32]byte{0x04},
			LogsBloom:        [256]byte{0x05},
			PrevRandao:       [32]byte{0x06},
			BlockNumber:      7,
			GasLimit:         30_000_000,
			GasUsed:          21_000,
			Timestamp:        1_600_000_000,
			ExtraData:        []byte{0x08},
			BaseFeePerGasLE:  [32]byte{0x09},
			BlockHash:        phase0.Hash32{0x0a},
			TransactionsRoot: phase0.Root{0x0b},
		},
	}
}

// testCapellaBlindedBody builds a Capella blinded beacon block body.
func testCapellaBlindedBody() *apiv1capella.BlindedBeaconBlockBody {
	return &apiv1capella.BlindedBeaconBlockBody{
		RANDAOReveal: phase0.BLSSignature{0x13},
		ETH1Data: &phase0.ETH1Data{
			DepositRoot:  phase0.Root{0x14},
			DepositCount: 15,
			BlockHash:    make([]byte, 32),
		},
		Graffiti:          [32]byte{'b', 'u', 'i', 'l', 'd'},
		ProposerSlashings: []*phase0.ProposerSlashing{},
		AttesterSlashings: []*phase0.AttesterSlashing{},
		Attestations:      []*phase0.Attestation{testPhase0Attestation()},
		Deposits:          []*phase0.Deposit{},
		VoluntaryExits:    []*phase0.SignedVoluntaryExit{},
		SyncAggregate: &altair.SyncAggregate{
			SyncCommitteeBits:      bitfield.NewBitvector512(),
			SyncCommitteeSignature: phase0.BLSSignature{0x16},
		},
		ExecutionPayloadHeader: &capella.ExecutionPayloadHeader{
			ParentHash:       phase0.Hash32{0x01},
			FeeRecipient:     bellatrix.ExecutionAddress{0x02},
			StateRoot:        [32]byte{0x03},
			ReceiptsRoot:     [32]byte{0x04},
			LogsBloom:        [256]byte{0x05},
			PrevRandao:       [32]byte{0x06},
			BlockNumber:      7,
			GasLimit:         30_000_000,
			GasUsed:          21_000,
			Timestamp:        1_600_000_000,
			ExtraData:        []byte{0x08},
			BaseFeePerGasLE:  [32]byte{0x09},
			BlockHash:        phase0.Hash32{0x0a},
			TransactionsRoot: phase0.Root{0x0b},
			WithdrawalsRoot:  phase0.Root{0x0c},
		},
		BLSToExecutionChanges: []*capella.SignedBLSToExecutionChange{},
	}
}

// testDenebExecutionPayloadHeader builds a Deneb execution payload header
// (also used by Electra/Fulu blinded bodies).
func testDenebExecutionPayloadHeader() *deneb.ExecutionPayloadHeader {
	return &deneb.ExecutionPayloadHeader{
		ParentHash:       phase0.Hash32{0x01},
		FeeRecipient:     bellatrix.ExecutionAddress{0x02},
		StateRoot:        phase0.Root{0x03},
		ReceiptsRoot:     phase0.Root{0x04},
		LogsBloom:        [256]byte{0x05},
		PrevRandao:       [32]byte{0x06},
		BlockNumber:      7,
		GasLimit:         30_000_000,
		GasUsed:          21_000,
		Timestamp:        1_600_000_000,
		ExtraData:        []byte{0x08},
		BaseFeePerGas:    uint256.NewInt(9),
		BlockHash:        phase0.Hash32{0x0a},
		TransactionsRoot: phase0.Root{0x0b},
		WithdrawalsRoot:  phase0.Root{0x0c},
		BlobGasUsed:      131072,
		ExcessBlobGas:    262144,
	}
}

// testDenebBlindedBody builds a Deneb blinded beacon block body.
func testDenebBlindedBody() *apiv1deneb.BlindedBeaconBlockBody {
	return &apiv1deneb.BlindedBeaconBlockBody{
		RANDAOReveal: phase0.BLSSignature{0x13},
		ETH1Data: &phase0.ETH1Data{
			DepositRoot:  phase0.Root{0x14},
			DepositCount: 15,
			BlockHash:    make([]byte, 32),
		},
		Graffiti:          [32]byte{'b', 'u', 'i', 'l', 'd'},
		ProposerSlashings: []*phase0.ProposerSlashing{},
		AttesterSlashings: []*phase0.AttesterSlashing{},
		Attestations:      []*phase0.Attestation{testPhase0Attestation()},
		Deposits:          []*phase0.Deposit{},
		VoluntaryExits:    []*phase0.SignedVoluntaryExit{},
		SyncAggregate: &altair.SyncAggregate{
			SyncCommitteeBits:      bitfield.NewBitvector512(),
			SyncCommitteeSignature: phase0.BLSSignature{0x16},
		},
		ExecutionPayloadHeader: testDenebExecutionPayloadHeader(),
		BLSToExecutionChanges:  []*capella.SignedBLSToExecutionChange{},
		BlobKZGCommitments:     []deneb.KZGCommitment{{0x17}},
	}
}

// testElectraBlindedBody builds an Electra blinded beacon block body (also
// the Fulu wire schema).
func testElectraBlindedBody() *apiv1electra.BlindedBeaconBlockBody {
	return &apiv1electra.BlindedBeaconBlockBody{
		RANDAOReveal: phase0.BLSSignature{0x13},
		ETH1Data: &phase0.ETH1Data{
			DepositRoot:  phase0.Root{0x14},
			DepositCount: 15,
			BlockHash:    make([]byte, 32),
		},
		Graffiti:          [32]byte{'b', 'u', 'i', 'l', 'd'},
		ProposerSlashings: []*phase0.ProposerSlashing{},
		AttesterSlashings: []*electra.AttesterSlashing{},
		Attestations:      []*electra.Attestation{testElectraAttestation()},
		Deposits:          []*phase0.Deposit{},
		VoluntaryExits:    []*phase0.SignedVoluntaryExit{},
		SyncAggregate: &altair.SyncAggregate{
			SyncCommitteeBits:      bitfield.NewBitvector512(),
			SyncCommitteeSignature: phase0.BLSSignature{0x16},
		},
		ExecutionPayloadHeader: testDenebExecutionPayloadHeader(),
		BLSToExecutionChanges:  []*capella.SignedBLSToExecutionChange{},
		BlobKZGCommitments:     []deneb.KZGCommitment{{0x17}},
		ExecutionRequests: &electra.ExecutionRequests{
			Deposits:       []*electra.DepositRequest{},
			Withdrawals:    []*electra.WithdrawalRequest{},
			Consolidations: []*electra.ConsolidationRequest{},
		},
	}
}

// signedBlindedBeaconBlockTests enumerates the per-fork views the agnostic
// SignedBlindedBeaconBlock must be wire-compatible with.
func signedBlindedBeaconBlockTests() []struct {
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
			name:    "Bellatrix",
			version: version.DataVersionBellatrix,
			view: &apiv1bellatrix.SignedBlindedBeaconBlock{
				Message: &apiv1bellatrix.BlindedBeaconBlock{
					Slot:          12345,
					ProposerIndex: 42,
					ParentRoot:    phase0.Root{0x11},
					StateRoot:     phase0.Root{0x12},
					Body:          testBellatrixBlindedBody(),
				},
				Signature: phase0.BLSSignature{0x18},
			},
		},
		{
			name:    "Capella",
			version: version.DataVersionCapella,
			view: &apiv1capella.SignedBlindedBeaconBlock{
				Message: &apiv1capella.BlindedBeaconBlock{
					Slot:          12345,
					ProposerIndex: 42,
					ParentRoot:    phase0.Root{0x11},
					StateRoot:     phase0.Root{0x12},
					Body:          testCapellaBlindedBody(),
				},
				Signature: phase0.BLSSignature{0x18},
			},
		},
		{
			name:    "Deneb",
			version: version.DataVersionDeneb,
			view: &apiv1deneb.SignedBlindedBeaconBlock{
				Message: &apiv1deneb.BlindedBeaconBlock{
					Slot:          12345,
					ProposerIndex: 42,
					ParentRoot:    phase0.Root{0x11},
					StateRoot:     phase0.Root{0x12},
					Body:          testDenebBlindedBody(),
				},
				Signature: phase0.BLSSignature{0x18},
			},
		},
		{
			name:    "Electra",
			version: version.DataVersionElectra,
			view: &apiv1electra.SignedBlindedBeaconBlock{
				Message: &apiv1electra.BlindedBeaconBlock{
					Slot:          12345,
					ProposerIndex: 42,
					ParentRoot:    phase0.Root{0x11},
					StateRoot:     phase0.Root{0x12},
					Body:          testElectraBlindedBody(),
				},
				Signature: phase0.BLSSignature{0x18},
			},
		},
		{
			// Fulu reuses the Electra wire schema; the agnostic type must
			// accept the Electra view when pinned to the Fulu version.
			name:    "Fulu",
			version: version.DataVersionFulu,
			view: &apiv1electra.SignedBlindedBeaconBlock{
				Message: &apiv1electra.BlindedBeaconBlock{
					Slot:          12345,
					ProposerIndex: 42,
					ParentRoot:    phase0.Root{0x11},
					StateRoot:     phase0.Root{0x12},
					Body:          testElectraBlindedBody(),
				},
				Signature: phase0.BLSSignature{0x18},
			},
		},
	}
}

// TestSignedBlindedBeaconBlockJSONWireCompat verifies that the agnostic type
// marshals to byte-identical JSON versus the per-fork type, and that JSON
// round-trips losslessly through the agnostic type.
func TestSignedBlindedBeaconBlockJSONWireCompat(t *testing.T) {
	for _, test := range signedBlindedBeaconBlockTests() {
		t.Run(test.name, func(t *testing.T) {
			expected, err := json.Marshal(test.view)
			require.NoError(t, err)

			agnostic := &apiv1all.SignedBlindedBeaconBlock{Version: test.version}
			require.NoError(t, agnostic.FromView(test.view))
			require.Equal(t, test.version, agnostic.Version)

			got, err := json.Marshal(agnostic)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(got),
				"agnostic JSON differs from per-fork JSON")

			// Round-trip through UnmarshalJSON with Version pre-set.
			rt := &apiv1all.SignedBlindedBeaconBlock{Version: test.version}
			require.NoError(t, json.Unmarshal(expected, rt))

			rtJSON, err := json.Marshal(rt)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(rtJSON),
				"round-tripped JSON differs from per-fork JSON")

			// Version must propagate to nested versionable children.
			require.NotNil(t, rt.Message)
			require.Equal(t, test.version, rt.Message.Version)
			require.NotNil(t, rt.Message.Body)
			require.Equal(t, test.version, rt.Message.Body.Version)
			require.NotNil(t, rt.Message.Body.ExecutionPayloadHeader)
			require.Equal(t, test.version, rt.Message.Body.ExecutionPayloadHeader.Version)
		})
	}
}

// TestSignedBlindedBeaconBlockSSZWireCompat verifies that the agnostic type
// produces byte-identical SSZ and the same hash tree root as the per-fork
// type, and that SSZ round-trips losslessly through the agnostic type.
func TestSignedBlindedBeaconBlockSSZWireCompat(t *testing.T) {
	for _, test := range signedBlindedBeaconBlockTests() {
		t.Run(test.name, func(t *testing.T) {
			codec, ok := test.view.(sszCodec)
			require.True(t, ok)

			expected, err := codec.MarshalSSZ()
			require.NoError(t, err)

			agnostic := &apiv1all.SignedBlindedBeaconBlock{Version: test.version}
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
			rt := &apiv1all.SignedBlindedBeaconBlock{Version: test.version}
			require.NoError(t, rt.UnmarshalSSZ(expected))

			rtSSZ, err := rt.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expected, rtSSZ,
				"round-tripped SSZ differs from per-fork SSZ")

			// Version must propagate to nested versionable children.
			require.NotNil(t, rt.Message)
			require.Equal(t, test.version, rt.Message.Version)
		})
	}
}

// TestSignedBlindedBeaconBlockViewRoundtrip verifies ToView reproduces the
// per-fork view the agnostic instance was built from.
func TestSignedBlindedBeaconBlockViewRoundtrip(t *testing.T) {
	for _, test := range signedBlindedBeaconBlockTests() {
		t.Run(test.name, func(t *testing.T) {
			agnostic := &apiv1all.SignedBlindedBeaconBlock{Version: test.version}
			require.NoError(t, agnostic.FromView(test.view))

			view, err := agnostic.ToView()
			require.NoError(t, err)
			require.IsType(t, test.view, view)
			require.Equal(t, test.view, view)
		})
	}
}

// TestSignedBlindedBeaconBlockToVersioned verifies conversion to and from
// api.VersionedSignedBlindedBeaconBlock.
func TestSignedBlindedBeaconBlockToVersioned(t *testing.T) {
	for _, test := range signedBlindedBeaconBlockTests() {
		t.Run(test.name, func(t *testing.T) {
			agnostic := &apiv1all.SignedBlindedBeaconBlock{Version: test.version}
			require.NoError(t, agnostic.FromView(test.view))

			versioned, err := agnostic.ToVersioned()
			require.NoError(t, err)
			require.Equal(t, test.version, versioned.Version)

			rt := &apiv1all.SignedBlindedBeaconBlock{}
			require.NoError(t, rt.FromVersioned(versioned))
			require.Equal(t, test.version, rt.Version)
			require.Equal(t, agnostic, rt)
		})
	}
}

// TestSignedBlindedBeaconBlockUnsupportedVersion verifies unsupported
// versions are rejected with a clear error.
func TestSignedBlindedBeaconBlockUnsupportedVersion(t *testing.T) {
	for _, v := range []version.DataVersion{
		version.DataVersionUnknown,
		version.DataVersionPhase0,
		version.DataVersionAltair,
		version.DataVersionGloas,
	} {
		s := &apiv1all.SignedBlindedBeaconBlock{Version: v}

		_, err := s.ToView()
		require.ErrorContains(t, err, "unsupported version")

		_, err = json.Marshal(s)
		require.ErrorContains(t, err, "unsupported version")

		_, err = s.MarshalSSZ()
		require.ErrorContains(t, err, "unsupported version")
	}
}
