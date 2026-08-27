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

package gloas_test

import (
	"encoding/json"
	"strings"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/ethpandaops/go-eth2-client/spec/altair"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/deneb"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// testBeaconState builds a small but fully-populated state so the JSON
// round-trip exercises every field's encoder and decoder.
func testBeaconState() *gloas.BeaconState {
	var pubkey phase0.BLSPubKey
	pubkey[0] = 0xaa

	syncCommittee := &altair.SyncCommittee{
		Pubkeys:         []phase0.BLSPubKey{pubkey},
		AggregatePubkey: pubkey,
	}

	return &gloas.BeaconState{
		GenesisTime:           1784030400,
		GenesisValidatorsRoot: phase0.Root{0x01},
		Slot:                  25736,
		Fork: &phase0.Fork{
			PreviousVersion: phase0.Version{0x70, 0x66, 0x95, 0x68},
			CurrentVersion:  phase0.Version{0x80, 0x66, 0x95, 0x68},
			Epoch:           38,
		},
		LatestBlockHeader: &phase0.BeaconBlockHeader{
			Slot:          25735,
			ProposerIndex: 3,
			ParentRoot:    phase0.Root{0x02},
			StateRoot:     phase0.Root{0x03},
			BodyRoot:      phase0.Root{0x04},
		},
		BlockRoots:      []phase0.Root{{0x05}},
		StateRoots:      []phase0.Root{{0x06}},
		HistoricalRoots: []phase0.Root{{0x07}},
		ETH1Data: &phase0.ETH1Data{
			DepositRoot:  phase0.Root{0x08},
			DepositCount: 1,
			BlockHash:    make([]byte, 32),
		},
		ETH1DataVotes:    []*phase0.ETH1Data{},
		ETH1DepositIndex: 1,
		Validators: []*phase0.Validator{
			{
				PublicKey:                  pubkey,
				WithdrawalCredentials:      make([]byte, 32),
				EffectiveBalance:           32000000000,
				ActivationEligibilityEpoch: 0,
				ActivationEpoch:            0,
				ExitEpoch:                  phase0.Epoch(0xffffffffffffffff),
				WithdrawableEpoch:          phase0.Epoch(0xffffffffffffffff),
			},
		},
		Balances:                      []phase0.Gwei{32000000000},
		RANDAOMixes:                   []phase0.Root{{0x09}},
		Slashings:                     []phase0.Gwei{0},
		PreviousEpochParticipation:    []altair.ParticipationFlags{7},
		CurrentEpochParticipation:     []altair.ParticipationFlags{3},
		JustificationBits:             bitfield.Bitvector4{0x0f},
		PreviousJustifiedCheckpoint:   &phase0.Checkpoint{Epoch: 802, Root: phase0.Root{0x0a}},
		CurrentJustifiedCheckpoint:    &phase0.Checkpoint{Epoch: 803, Root: phase0.Root{0x0b}},
		FinalizedCheckpoint:           &phase0.Checkpoint{Epoch: 802, Root: phase0.Root{0x0c}},
		InactivityScores:              []uint64{0},
		CurrentSyncCommittee:          syncCommittee,
		NextSyncCommittee:             syncCommittee,
		LatestBlockHash:               phase0.Hash32{0x0d},
		NextWithdrawalIndex:           15638,
		NextWithdrawalValidatorIndex:  470013,
		HistoricalSummaries:           []*capella.HistoricalSummary{},
		DepositRequestsStartIndex:     0,
		DepositBalanceToConsume:       0,
		ExitBalanceToConsume:          256000000000,
		EarliestExitEpoch:             804,
		ConsolidationBalanceToConsume: 0,
		EarliestConsolidationEpoch:    805,
		ProposerLookahead:             []phase0.ValidatorIndex{1039, 503062},
		Builders: []*gloas.Builder{
			{
				PublicKey:         pubkey,
				Version:           1,
				Balance:           32000000000,
				DepositEpoch:      12,
				WithdrawableEpoch: 18,
			},
		},
		NextWithdrawalBuilderIndex: 0,
		// A bitvector: one bit per historical slot, hex-encoded on the wire.
		ExecutionPayloadAvailability: []uint8{0xff, 0x7f, 0x00, 0xbf},
		BuilderPendingPayments: []*gloas.BuilderPendingPayment{
			{
				Weight: 0,
				Withdrawal: &gloas.BuilderPendingWithdrawal{
					Amount:       0,
					BuilderIndex: 0,
				},
				ProposerIndex: 0,
			},
		},
		BuilderPendingWithdrawals: []*gloas.BuilderPendingWithdrawal{},
		LatestExecutionPayloadBid: &gloas.ExecutionPayloadBid{
			ParentBlockHash:    phase0.Hash32{0x0e},
			ParentBlockRoot:    phase0.Root{0x0f},
			BlockHash:          phase0.Hash32{0x10},
			PrevRandao:         phase0.Root{0x11},
			GasLimit:           60000000,
			BuilderIndex:       0,
			Slot:               25736,
			Value:              0,
			ExecutionPayment:   0,
			BlobKZGCommitments: []deneb.KZGCommitment{},
		},
		PayloadExpectedWithdrawals: []*capella.Withdrawal{
			{
				Index:          15637,
				ValidatorIndex: 1099511627778,
				Amount:         128018656,
			},
		},
		PTCWindow: [][]phase0.ValidatorIndex{
			{1261, 1969, 503057},
			{323, 2478},
		},
	}
}

func TestBeaconStateJSONRoundTrip(t *testing.T) {
	state := testBeaconState()

	encoded, err := json.Marshal(state)
	require.NoError(t, err)

	// Bitvectors are hex strings on the wire, not per-byte arrays.
	require.Contains(t, string(encoded), `"execution_payload_availability":"0xff7f00bf"`)
	// uint64 values are quoted per the beacon API convention.
	require.Contains(t, string(encoded), `"ptc_window":[["1261","1969","503057"],["323","2478"]]`)

	decoded := &gloas.BeaconState{}
	require.NoError(t, json.Unmarshal(encoded, decoded))
	require.Equal(t, state.ExecutionPayloadAvailability, decoded.ExecutionPayloadAvailability)
	require.Equal(t, state.PTCWindow, decoded.PTCWindow)
	require.Equal(t, state.Builders, decoded.Builders)
}

// TestBeaconStateJSONLighthouse checks the encodings lighthouse actually
// serves where they diverge from the quoted-integer convention.
func TestBeaconStateJSONLighthouse(t *testing.T) {
	state := testBeaconState()

	encoded, err := json.Marshal(state)
	require.NoError(t, err)

	// Lighthouse serves ptc_window indices as bare JSON numbers.
	lighthouse := strings.Replace(
		string(encoded),
		`"ptc_window":[["1261","1969","503057"],["323","2478"]]`,
		`"ptc_window":[[1261,1969,503057],[323,2478]]`,
		1,
	)
	require.NotEqual(t, string(encoded), lighthouse)

	decoded := &gloas.BeaconState{}
	require.NoError(t, json.Unmarshal([]byte(lighthouse), decoded))
	require.Equal(t, state.PTCWindow, decoded.PTCWindow)
}

func TestBeaconStateJSONInvalidAvailability(t *testing.T) {
	state := testBeaconState()

	encoded, err := json.Marshal(state)
	require.NoError(t, err)

	invalid := strings.Replace(
		string(encoded),
		`"execution_payload_availability":"0xff7f00bf"`,
		`"execution_payload_availability":"0xzz"`,
		1,
	)

	decoded := &gloas.BeaconState{}
	require.ErrorContains(t, json.Unmarshal([]byte(invalid), decoded), "execution_payload_availability")
}
