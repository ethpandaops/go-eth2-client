// Copyright © 2023 Attestant Limited.
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

package heze

import (
	"fmt"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/ethpandaops/go-eth2-client/spec/altair"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// BeaconState represents a beacon state for Heze, whose latest execution payload
// bid is the one extended by EIP-7805.
type BeaconState struct {
	GenesisTime                   uint64                              `ssz-index:"0"`
	GenesisValidatorsRoot         phase0.Root                         `ssz-index:"1" ssz-size:"32"`
	Slot                          phase0.Slot                         `ssz-index:"2"`
	Fork                          *phase0.Fork                        `ssz-index:"3"`
	LatestBlockHeader             *phase0.BeaconBlockHeader           `ssz-index:"4"`
	BlockRoots                    []phase0.Root                       `ssz-index:"5" dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" ssz-size:"8192,32"`
	StateRoots                    []phase0.Root                       `ssz-index:"6" dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" ssz-size:"8192,32"`
	HistoricalRoots               []phase0.Root                       `ssz-index:"7" ssz-max:"16777216"                         ssz-size:"?,32"`
	ETH1Data                      *phase0.ETH1Data                    `ssz-index:"8"`
	ETH1DataVotes                 []*phase0.ETH1Data                  `ssz-index:"9" dynssz-max:"EPOCHS_PER_ETH1_VOTING_PERIOD*SLOTS_PER_EPOCH" ssz-max:"2048"`
	ETH1DepositIndex              uint64                              `ssz-index:"10"`
	Validators                    []*phase0.Validator                 `ssz-index:"11" ssz-type:"progressive-list"`
	Balances                      []phase0.Gwei                       `ssz-index:"12" ssz-type:"progressive-list"`
	RANDAOMixes                   []phase0.Root                       `ssz-index:"13" dynssz-size:"EPOCHS_PER_HISTORICAL_VECTOR,32" ssz-size:"65536,32"`
	Slashings                     []phase0.Gwei                       `ssz-index:"14" dynssz-size:"EPOCHS_PER_SLASHINGS_VECTOR"     ssz-size:"8192"`
	PreviousEpochParticipation    []altair.ParticipationFlags         `ssz-index:"15" ssz-type:"progressive-list"`
	CurrentEpochParticipation     []altair.ParticipationFlags         `ssz-index:"16" ssz-type:"progressive-list"`
	JustificationBits             bitfield.Bitvector4                 `ssz-index:"17" ssz-size:"1"`
	PreviousJustifiedCheckpoint   *phase0.Checkpoint                  `ssz-index:"18"`
	CurrentJustifiedCheckpoint    *phase0.Checkpoint                  `ssz-index:"19"`
	FinalizedCheckpoint           *phase0.Checkpoint                  `ssz-index:"20"`
	InactivityScores              []uint64                            `ssz-index:"21" ssz-type:"progressive-list"`
	CurrentSyncCommittee          *altair.SyncCommittee               `ssz-index:"22"`
	NextSyncCommittee             *altair.SyncCommittee               `ssz-index:"23"`
	LatestBlockHash               phase0.Hash32                       `ssz-index:"24" ssz-size:"32"`
	NextWithdrawalIndex           capella.WithdrawalIndex             `ssz-index:"25"`
	NextWithdrawalValidatorIndex  phase0.ValidatorIndex               `ssz-index:"26"`
	HistoricalSummaries           []*capella.HistoricalSummary        `ssz-index:"27" ssz-max:"16777216"`
	DepositRequestsStartIndex     uint64                              `ssz-index:"28"`
	DepositBalanceToConsume       phase0.Gwei                         `ssz-index:"29"`
	ExitBalanceToConsume          phase0.Gwei                         `ssz-index:"30"`
	EarliestExitEpoch             phase0.Epoch                        `ssz-index:"31"`
	ConsolidationBalanceToConsume phase0.Gwei                         `ssz-index:"32"`
	EarliestConsolidationEpoch    phase0.Epoch                        `ssz-index:"33"`
	PendingDeposits               []*electra.PendingDeposit           `ssz-index:"34" ssz-type:"progressive-list"`
	PendingPartialWithdrawals     []*electra.PendingPartialWithdrawal `ssz-index:"35" ssz-type:"progressive-list"`
	PendingConsolidations         []*electra.PendingConsolidation     `ssz-index:"36" ssz-type:"progressive-list"`
	ProposerLookahead             []phase0.ValidatorIndex             `ssz-index:"37" dynssz-size:"(MIN_SEED_LOOKAHEAD+1)*SLOTS_PER_EPOCH" ssz-size:"64"`
	Builders                      []*gloas.Builder                    `ssz-index:"38" ssz-type:"progressive-list"`
	NextWithdrawalBuilderIndex    gloas.BuilderIndex                  `ssz-index:"39"`
	ExecutionPayloadAvailability  []uint8                             `ssz-index:"40" dynssz-size:"SLOTS_PER_HISTORICAL_ROOT/8"      ssz-size:"1024"`
	BuilderPendingPayments        []*gloas.BuilderPendingPayment      `ssz-index:"41" dynssz-size:"SLOTS_PER_EPOCH*2"                ssz-size:"64"`
	BuilderPendingWithdrawals     []*gloas.BuilderPendingWithdrawal   `ssz-index:"42" ssz-type:"progressive-list"`
	LatestExecutionPayloadBid     *ExecutionPayloadBid                `ssz-index:"43"`
	PayloadExpectedWithdrawals    []*capella.Withdrawal               `ssz-index:"44" ssz-type:"progressive-list"`
	PTCWindow                     [][]phase0.ValidatorIndex           `ssz-index:"45" dynssz-size:"(2+MIN_SEED_LOOKAHEAD)*SLOTS_PER_EPOCH,PTC_SIZE" ssz-size:"96,512"`
}

// String returns a string version of the structure.
func (b *BeaconState) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
