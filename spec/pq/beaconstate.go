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

package pq

import (
	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/ethpandaops/go-eth2-client/spec/altair"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
)

type BeaconState struct {
	GenesisTime                   uint64                              `ssz-index:"0" json:"genesis_time"`
	GenesisValidatorsRoot         phase0.Root                         `ssz-index:"1" json:"genesis_validators_root" ssz-size:"32"`
	Slot                          phase0.Slot                         `ssz-index:"2" json:"slot"`
	Fork                          *phase0.Fork                        `ssz-index:"3" json:"fork"`
	LatestBlockHeader             *phase0.BeaconBlockHeader           `ssz-index:"4" json:"latest_block_header"`
	BlockRoots                    []phase0.Root                       `ssz-index:"5" dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" json:"block_roots" ssz-size:"8192,32"`
	StateRoots                    []phase0.Root                       `ssz-index:"6" dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" json:"state_roots" ssz-size:"8192,32"`
	HistoricalRoots               []phase0.Root                       `ssz-index:"7" json:"historical_roots" ssz-max:"16777216" ssz-size:"?,32"`
	ETH1Data                      *phase0.ETH1Data                    `ssz-index:"8" json:"eth1_data"`
	ETH1DataVotes                 []*phase0.ETH1Data                  `ssz-index:"9" dynssz-max:"EPOCHS_PER_ETH1_VOTING_PERIOD*SLOTS_PER_EPOCH" json:"eth1_data_votes" ssz-max:"2048"`
	ETH1DepositIndex              uint64                              `ssz-index:"10" json:"eth1_deposit_index"`
	XMSSValidators                []*Validator                        `ssz-index:"11" json:"validators" ssz-type:"progressive-list"`
	Balances                      []phase0.Gwei                       `ssz-index:"12" json:"balances" ssz-type:"progressive-list"`
	RANDAOMixes                   []phase0.Root                       `ssz-index:"13" dynssz-size:"EPOCHS_PER_HISTORICAL_VECTOR,32" json:"randao_mixes" ssz-size:"65536,32"`
	Slashings                     []phase0.Gwei                       `ssz-index:"14" dynssz-size:"EPOCHS_PER_SLASHINGS_VECTOR" json:"slashings" ssz-size:"8192"`
	PreviousEpochParticipation    []altair.ParticipationFlags         `ssz-index:"15" json:"previous_epoch_participation" ssz-type:"progressive-list"`
	CurrentEpochParticipation     []altair.ParticipationFlags         `ssz-index:"16" json:"current_epoch_participation" ssz-type:"progressive-list"`
	JustificationBits             bitfield.Bitvector4                 `ssz-index:"17" json:"justification_bits" ssz-size:"1"`
	PreviousJustifiedCheckpoint   *phase0.Checkpoint                  `ssz-index:"18" json:"previous_justified_checkpoint"`
	CurrentJustifiedCheckpoint    *phase0.Checkpoint                  `ssz-index:"19" json:"current_justified_checkpoint"`
	FinalizedCheckpoint           *phase0.Checkpoint                  `ssz-index:"20" json:"finalized_checkpoint"`
	InactivityScores              []uint64                            `ssz-index:"21" json:"inactivity_scores" ssz-type:"progressive-list"`
	LatestBlockHash               phase0.Hash32                       `ssz-index:"24" json:"latest_block_hash" ssz-size:"32"`
	NextWithdrawalIndex           capella.WithdrawalIndex             `ssz-index:"25" json:"next_withdrawal_index"`
	NextWithdrawalValidatorIndex  phase0.ValidatorIndex               `ssz-index:"26" json:"next_withdrawal_validator_index"`
	HistoricalSummaries           []*capella.HistoricalSummary        `ssz-index:"27" json:"historical_summaries" ssz-max:"16777216"`
	DepositRequestsStartIndex     uint64                              `ssz-index:"28" json:"deposit_requests_start_index"`
	DepositBalanceToConsume       phase0.Gwei                         `ssz-index:"29" json:"deposit_balance_to_consume"`
	ExitBalanceToConsume          phase0.Gwei                         `ssz-index:"30" json:"exit_balance_to_consume"`
	EarliestExitEpoch             phase0.Epoch                        `ssz-index:"31" json:"earliest_exit_epoch"`
	ConsolidationBalanceToConsume phase0.Gwei                         `ssz-index:"32" json:"consolidation_balance_to_consume"`
	EarliestConsolidationEpoch    phase0.Epoch                        `ssz-index:"33" json:"earliest_consolidation_epoch"`
	PendingDeposits               []*electra.PendingDeposit           `ssz-index:"34" json:"pending_deposits" ssz-type:"progressive-list"`
	PendingPartialWithdrawals     []*electra.PendingPartialWithdrawal `ssz-index:"35" json:"pending_partial_withdrawals" ssz-type:"progressive-list"`
	PendingConsolidations         []*electra.PendingConsolidation     `ssz-index:"36" json:"pending_consolidations" ssz-type:"progressive-list"`
	ProposerLookahead             []phase0.ValidatorIndex             `ssz-index:"37" dynssz-size:"(MIN_SEED_LOOKAHEAD+1)*SLOTS_PER_EPOCH" json:"proposer_lookahead" ssz-size:"64"`
	XMSSBuilders                  []*Builder                          `ssz-index:"38" json:"builders" ssz-type:"progressive-list"`
	NextWithdrawalBuilderIndex    gloas.BuilderIndex                  `ssz-index:"39" json:"next_withdrawal_builder_index"`
	ExecutionPayloadAvailability  []uint8                             `ssz-index:"40" dynssz-size:"SLOTS_PER_HISTORICAL_ROOT/8" json:"execution_payload_availability" ssz-size:"1024"`
	BuilderPendingPayments        []*gloas.BuilderPendingPayment      `ssz-index:"41" dynssz-size:"SLOTS_PER_EPOCH*2" json:"builder_pending_payments" ssz-size:"64"`
	BuilderPendingWithdrawals     []*gloas.BuilderPendingWithdrawal   `ssz-index:"42" json:"builder_pending_withdrawals" ssz-type:"progressive-list"`
	LatestExecutionPayloadBid     *gloas.ExecutionPayloadBid          `ssz-index:"43" json:"latest_execution_payload_bid"`
	PayloadExpectedWithdrawals    []*capella.Withdrawal               `ssz-index:"44" json:"payload_expected_withdrawals" ssz-type:"progressive-list"`
	PTCWindow                     [][]phase0.ValidatorIndex           `ssz-index:"45" dynssz-size:"(2+MIN_SEED_LOOKAHEAD)*SLOTS_PER_EPOCH,PTC_SIZE" json:"ptc_window" ssz-size:"96,512"`
	ValidatorPreregistrations     []*ValidatorPreregistration         `ssz-index:"48" json:"validator_preregistrations" ssz-type:"progressive-list"`
	BuilderPreregistrations       []*BuilderPreregistration           `ssz-index:"49" json:"builder_preregistrations" ssz-type:"progressive-list"`
	AggregationCommitteeWindow    [][][]phase0.ValidatorIndex         `ssz-index:"50" dynssz-size:"(2+MIN_SEED_LOOKAHEAD)*SLOTS_PER_EPOCH" json:"aggregation_committee_window" ssz-size:"96" ssz-type:"vector,progressive-list,progressive-list"` //nolint:revive
}
