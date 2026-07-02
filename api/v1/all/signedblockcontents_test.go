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
	"github.com/ethpandaops/go-eth2-client/api"
	apiv1all "github.com/ethpandaops/go-eth2-client/api/v1/all"
	apiv1deneb "github.com/ethpandaops/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/ethpandaops/go-eth2-client/api/v1/electra"
	apiv1fulu "github.com/ethpandaops/go-eth2-client/api/v1/fulu"
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

// sszCodec is the SSZ surface shared by the per-fork block contents types.
type sszCodec interface {
	MarshalSSZ() ([]byte, error)
	HashTreeRoot() ([32]byte, error)
}

// testExecutionPayload builds a Deneb execution payload with realistic-ish data.
func testExecutionPayload() *deneb.ExecutionPayload {
	return &deneb.ExecutionPayload{
		ParentHash:    phase0.Hash32{0x01},
		FeeRecipient:  bellatrix.ExecutionAddress{0x02},
		StateRoot:     phase0.Root{0x03},
		ReceiptsRoot:  phase0.Root{0x04},
		LogsBloom:     [256]byte{0x05},
		PrevRandao:    [32]byte{0x06},
		BlockNumber:   7,
		GasLimit:      30_000_000,
		GasUsed:       21_000,
		Timestamp:     1_600_000_000,
		ExtraData:     []byte{0x08},
		BaseFeePerGas: uint256.NewInt(9),
		BlockHash:     phase0.Hash32{0x0a},
		Transactions:  []bellatrix.Transaction{{0x0b, 0x0c}},
		Withdrawals: []*capella.Withdrawal{
			{Index: 1, ValidatorIndex: 2, Address: bellatrix.ExecutionAddress{0x0d}, Amount: 3},
		},
		BlobGasUsed:   131072,
		ExcessBlobGas: 262144,
	}
}

// testDenebSignedBeaconBlock builds a Deneb signed beacon block.
func testDenebSignedBeaconBlock() *deneb.SignedBeaconBlock {
	return &deneb.SignedBeaconBlock{
		Message: &deneb.BeaconBlock{
			Slot:          12345,
			ProposerIndex: 42,
			ParentRoot:    phase0.Root{0x11},
			StateRoot:     phase0.Root{0x12},
			Body: &deneb.BeaconBlockBody{
				RANDAOReveal: phase0.BLSSignature{0x13},
				ETH1Data: &phase0.ETH1Data{
					DepositRoot:  phase0.Root{0x14},
					DepositCount: 15,
					BlockHash:    make([]byte, 32),
				},
				Graffiti:          [32]byte{'b', 'u', 'i', 'l', 'd'},
				ProposerSlashings: []*phase0.ProposerSlashing{},
				AttesterSlashings: []*phase0.AttesterSlashing{},
				Attestations:      []*phase0.Attestation{},
				Deposits:          []*phase0.Deposit{},
				VoluntaryExits:    []*phase0.SignedVoluntaryExit{},
				SyncAggregate: &altair.SyncAggregate{
					SyncCommitteeBits:      bitfield.NewBitvector512(),
					SyncCommitteeSignature: phase0.BLSSignature{0x16},
				},
				ExecutionPayload:      testExecutionPayload(),
				BLSToExecutionChanges: []*capella.SignedBLSToExecutionChange{},
				BlobKZGCommitments:    []deneb.KZGCommitment{{0x17}},
			},
		},
		Signature: phase0.BLSSignature{0x18},
	}
}

// testElectraSignedBeaconBlock builds an Electra signed beacon block (also the
// Fulu wire schema).
func testElectraSignedBeaconBlock() *electra.SignedBeaconBlock {
	return &electra.SignedBeaconBlock{
		Message: &electra.BeaconBlock{
			Slot:          12345,
			ProposerIndex: 42,
			ParentRoot:    phase0.Root{0x11},
			StateRoot:     phase0.Root{0x12},
			Body: &electra.BeaconBlockBody{
				RANDAOReveal: phase0.BLSSignature{0x13},
				ETH1Data: &phase0.ETH1Data{
					DepositRoot:  phase0.Root{0x14},
					DepositCount: 15,
					BlockHash:    make([]byte, 32),
				},
				Graffiti:          [32]byte{'b', 'u', 'i', 'l', 'd'},
				ProposerSlashings: []*phase0.ProposerSlashing{},
				AttesterSlashings: []*electra.AttesterSlashing{},
				Attestations:      []*electra.Attestation{},
				Deposits:          []*phase0.Deposit{},
				VoluntaryExits:    []*phase0.SignedVoluntaryExit{},
				SyncAggregate: &altair.SyncAggregate{
					SyncCommitteeBits:      bitfield.NewBitvector512(),
					SyncCommitteeSignature: phase0.BLSSignature{0x16},
				},
				ExecutionPayload:      testExecutionPayload(),
				BLSToExecutionChanges: []*capella.SignedBLSToExecutionChange{},
				BlobKZGCommitments:    []deneb.KZGCommitment{{0x17}},
				ExecutionRequests: &electra.ExecutionRequests{
					Deposits:       []*electra.DepositRequest{},
					Withdrawals:    []*electra.WithdrawalRequest{},
					Consolidations: []*electra.ConsolidationRequest{},
				},
			},
		},
		Signature: phase0.BLSSignature{0x18},
	}
}

// testKZGProofs and testBlobs are the sidecar data shared by every fork.
func testKZGProofs() []deneb.KZGProof {
	return []deneb.KZGProof{{0x21}, {0x22}}
}

func testBlobs() []deneb.Blob {
	return []deneb.Blob{{0x23}}
}

// signedBlockContentsTests enumerates the per-fork views the agnostic
// SignedBlockContents must be wire-compatible with.
func signedBlockContentsTests() []struct {
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
			name:    "Deneb",
			version: version.DataVersionDeneb,
			view: &apiv1deneb.SignedBlockContents{
				SignedBlock: testDenebSignedBeaconBlock(),
				KZGProofs:   testKZGProofs(),
				Blobs:       testBlobs(),
			},
		},
		{
			name:    "Electra",
			version: version.DataVersionElectra,
			view: &apiv1electra.SignedBlockContents{
				SignedBlock: testElectraSignedBeaconBlock(),
				KZGProofs:   testKZGProofs(),
				Blobs:       testBlobs(),
			},
		},
		{
			name:    "Fulu",
			version: version.DataVersionFulu,
			view: &apiv1fulu.SignedBlockContents{
				SignedBlock: testElectraSignedBeaconBlock(),
				KZGProofs:   testKZGProofs(),
				Blobs:       testBlobs(),
			},
		},
	}
}

// TestSignedBlockContentsJSONWireCompat verifies that the agnostic type
// marshals to byte-identical JSON versus the per-fork type, and that JSON
// round-trips losslessly through the agnostic type.
func TestSignedBlockContentsJSONWireCompat(t *testing.T) {
	for _, test := range signedBlockContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			expected, err := json.Marshal(test.view)
			require.NoError(t, err)

			agnostic := &apiv1all.SignedBlockContents{}
			require.NoError(t, agnostic.FromView(test.view))
			require.Equal(t, test.version, agnostic.Version)

			got, err := json.Marshal(agnostic)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(got),
				"agnostic JSON differs from per-fork JSON")

			// Round-trip through UnmarshalJSON with Version pre-set.
			rt := &apiv1all.SignedBlockContents{Version: test.version}
			require.NoError(t, json.Unmarshal(expected, rt))

			rtJSON, err := json.Marshal(rt)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(rtJSON),
				"round-tripped JSON differs from per-fork JSON")

			// Version must propagate to nested versionable children.
			require.NotNil(t, rt.SignedBlock)
			require.Equal(t, test.version, rt.SignedBlock.Version)
			require.NotNil(t, rt.SignedBlock.Message)
			require.Equal(t, test.version, rt.SignedBlock.Message.Version)
			require.NotNil(t, rt.SignedBlock.Message.Body)
			require.Equal(t, test.version, rt.SignedBlock.Message.Body.Version)
		})
	}
}

// TestSignedBlockContentsSSZWireCompat verifies that the agnostic type
// produces byte-identical SSZ and the same hash tree root as the per-fork
// type, and that SSZ round-trips losslessly through the agnostic type.
func TestSignedBlockContentsSSZWireCompat(t *testing.T) {
	for _, test := range signedBlockContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			codec, ok := test.view.(sszCodec)
			require.True(t, ok)

			expected, err := codec.MarshalSSZ()
			require.NoError(t, err)

			agnostic := &apiv1all.SignedBlockContents{}
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
			rt := &apiv1all.SignedBlockContents{Version: test.version}
			require.NoError(t, rt.UnmarshalSSZ(expected))

			rtSSZ, err := rt.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expected, rtSSZ,
				"round-tripped SSZ differs from per-fork SSZ")

			// Version must propagate to nested versionable children.
			require.NotNil(t, rt.SignedBlock)
			require.Equal(t, test.version, rt.SignedBlock.Version)
		})
	}
}

// TestSignedBlockContentsViewRoundtrip verifies ToView reproduces the
// per-fork view the agnostic instance was built from.
func TestSignedBlockContentsViewRoundtrip(t *testing.T) {
	for _, test := range signedBlockContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			agnostic := &apiv1all.SignedBlockContents{}
			require.NoError(t, agnostic.FromView(test.view))

			view, err := agnostic.ToView()
			require.NoError(t, err)
			require.IsType(t, test.view, view)
			require.Equal(t, test.view, view)
		})
	}
}

// TestSignedBlockContentsToVersioned verifies conversion to and from
// api.VersionedSignedProposal: exactly the plain (unblinded) per-fork field
// is populated, the Blinded flag stays false, and the per-fork field carries
// the exact wire bytes the http client submits (it marshals that field for
// both JSON and SSZ after AssertPresent).
func TestSignedBlockContentsToVersioned(t *testing.T) {
	for _, test := range signedBlockContentsTests() {
		t.Run(test.name, func(t *testing.T) {
			agnostic := &apiv1all.SignedBlockContents{}
			require.NoError(t, agnostic.FromView(test.view))

			versioned, err := agnostic.ToVersioned()
			require.NoError(t, err)
			require.Equal(t, test.version, versioned.Version)
			require.False(t, versioned.Blinded, "unblinded contents must not set the Blinded flag")
			require.NoError(t, versioned.AssertPresent())

			// The per-fork field must be the exact view type; the blinded
			// fields must stay untouched.
			var fieldView any

			switch test.version {
			case version.DataVersionDeneb:
				fieldView = versioned.Deneb
				require.Nil(t, versioned.DenebBlinded)
			case version.DataVersionElectra:
				fieldView = versioned.Electra
				require.Nil(t, versioned.ElectraBlinded)
			case version.DataVersionFulu:
				fieldView = versioned.Fulu
				require.Nil(t, versioned.FuluBlinded)
			default:
				t.Fatalf("unexpected version %v", test.version)
			}

			require.IsType(t, test.view, fieldView)
			require.Equal(t, test.view, fieldView)

			// Byte-compare the body the http client would submit
			// (json.Marshal / MarshalSSZ of the per-fork field).
			expectedJSON, err := json.Marshal(test.view)
			require.NoError(t, err)

			gotJSON, err := json.Marshal(fieldView)
			require.NoError(t, err)
			require.Equal(t, string(expectedJSON), string(gotJSON),
				"submitted proposal JSON body differs from per-fork JSON")

			expectedSSZ, err := test.view.(sszCodec).MarshalSSZ()
			require.NoError(t, err)

			gotSSZ, err := fieldView.(sszCodec).MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, expectedSSZ, gotSSZ,
				"submitted proposal SSZ body differs from per-fork SSZ")

			rt := &apiv1all.SignedBlockContents{}
			require.NoError(t, rt.FromVersioned(versioned))
			require.Equal(t, test.version, rt.Version)
			require.Equal(t, agnostic, rt)
		})
	}
}

// TestSignedBlockContentsFromVersionedBlinded verifies that a blinded
// proposal (Blinded flag set, only the *Blinded field populated) is rejected
// rather than silently mis-read.
func TestSignedBlockContentsFromVersionedBlinded(t *testing.T) {
	versioned := &api.VersionedSignedProposal{
		Version:     version.DataVersionFulu,
		Blinded:     true,
		FuluBlinded: &apiv1electra.SignedBlindedBeaconBlock{},
	}

	rt := &apiv1all.SignedBlockContents{}
	require.Error(t, rt.FromVersioned(versioned),
		"blinded proposals must not convert into unblinded block contents")
}

// TestSignedBlockContentsUnsupportedVersion verifies unsupported versions are
// rejected with a clear error.
func TestSignedBlockContentsUnsupportedVersion(t *testing.T) {
	for _, v := range []version.DataVersion{
		version.DataVersionUnknown,
		version.DataVersionPhase0,
		version.DataVersionCapella,
		version.DataVersionGloas,
	} {
		s := &apiv1all.SignedBlockContents{Version: v}

		_, err := s.ToView()
		require.ErrorContains(t, err, "unsupported version")

		_, err = json.Marshal(s)
		require.ErrorContains(t, err, "unsupported version")

		_, err = s.MarshalSSZ()
		require.ErrorContains(t, err, "unsupported version")
	}
}
