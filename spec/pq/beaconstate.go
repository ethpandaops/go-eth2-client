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
	GenesisTime                   uint64                              `json:"genesis_time" ssz-index:"0"`
	GenesisValidatorsRoot         phase0.Root                         `json:"genesis_validators_root" ssz-index:"1" ssz-size:"32"`
	Slot                          phase0.Slot                         `json:"slot" ssz-index:"2"`
	Fork                          *phase0.Fork                        `json:"fork" ssz-index:"3"`
	LatestBlockHeader             *phase0.BeaconBlockHeader           `json:"latest_block_header" ssz-index:"4"`
	BlockRoots                    []phase0.Root                       `json:"block_roots" ssz-index:"5" dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" ssz-size:"8192,32"`
	StateRoots                    []phase0.Root                       `json:"state_roots" ssz-index:"6" dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" ssz-size:"8192,32"`
	HistoricalRoots               []phase0.Root                       `json:"historical_roots" ssz-index:"7" ssz-max:"16777216" ssz-size:"?,32"`
	ETH1Data                      *phase0.ETH1Data                    `json:"eth1_data" ssz-index:"8"`
	ETH1DataVotes                 []*phase0.ETH1Data                  `json:"eth1_data_votes" ssz-index:"9" dynssz-max:"EPOCHS_PER_ETH1_VOTING_PERIOD*SLOTS_PER_EPOCH" ssz-max:"2048"`
	ETH1DepositIndex              uint64                              `json:"eth1_deposit_index" ssz-index:"10"`
	XMSSValidators                []*Validator                        `json:"validators" ssz-index:"11" ssz-type:"progressive-list"`
	Balances                      []phase0.Gwei                       `json:"balances" ssz-index:"12" ssz-type:"progressive-list"`
	RANDAOMixes                   []phase0.Root                       `json:"randao_mixes" ssz-index:"13" dynssz-size:"EPOCHS_PER_HISTORICAL_VECTOR,32" ssz-size:"65536,32"`
	Slashings                     []phase0.Gwei                       `json:"slashings" ssz-index:"14" dynssz-size:"EPOCHS_PER_SLASHINGS_VECTOR" ssz-size:"8192"`
	PreviousEpochParticipation    []altair.ParticipationFlags         `json:"previous_epoch_participation" ssz-index:"15" ssz-type:"progressive-list"`
	CurrentEpochParticipation     []altair.ParticipationFlags         `json:"current_epoch_participation" ssz-index:"16" ssz-type:"progressive-list"`
	JustificationBits             bitfield.Bitvector4                 `json:"justification_bits" ssz-index:"17" ssz-size:"1"`
	PreviousJustifiedCheckpoint   *phase0.Checkpoint                  `json:"previous_justified_checkpoint" ssz-index:"18"`
	CurrentJustifiedCheckpoint    *phase0.Checkpoint                  `json:"current_justified_checkpoint" ssz-index:"19"`
	FinalizedCheckpoint           *phase0.Checkpoint                  `json:"finalized_checkpoint" ssz-index:"20"`
	InactivityScores              []uint64                            `json:"inactivity_scores" ssz-index:"21" ssz-type:"progressive-list"`
	LatestBlockHash               phase0.Hash32                       `json:"latest_block_hash" ssz-index:"24" ssz-size:"32"`
	NextWithdrawalIndex           capella.WithdrawalIndex             `json:"next_withdrawal_index" ssz-index:"25"`
	NextWithdrawalValidatorIndex  phase0.ValidatorIndex               `json:"next_withdrawal_validator_index" ssz-index:"26"`
	HistoricalSummaries           []*capella.HistoricalSummary        `json:"historical_summaries" ssz-index:"27" ssz-max:"16777216"`
	DepositRequestsStartIndex     uint64                              `json:"deposit_requests_start_index" ssz-index:"28"`
	DepositBalanceToConsume       phase0.Gwei                         `json:"deposit_balance_to_consume" ssz-index:"29"`
	ExitBalanceToConsume          phase0.Gwei                         `json:"exit_balance_to_consume" ssz-index:"30"`
	EarliestExitEpoch             phase0.Epoch                        `json:"earliest_exit_epoch" ssz-index:"31"`
	ConsolidationBalanceToConsume phase0.Gwei                         `json:"consolidation_balance_to_consume" ssz-index:"32"`
	EarliestConsolidationEpoch    phase0.Epoch                        `json:"earliest_consolidation_epoch" ssz-index:"33"`
	PendingDeposits               []*electra.PendingDeposit           `json:"pending_deposits" ssz-index:"34" ssz-type:"progressive-list"`
	PendingPartialWithdrawals     []*electra.PendingPartialWithdrawal `json:"pending_partial_withdrawals" ssz-index:"35" ssz-type:"progressive-list"`
	PendingConsolidations         []*electra.PendingConsolidation     `json:"pending_consolidations" ssz-index:"36" ssz-type:"progressive-list"`
	ProposerLookahead             []phase0.ValidatorIndex             `json:"proposer_lookahead" ssz-index:"37" dynssz-size:"(MIN_SEED_LOOKAHEAD+1)*SLOTS_PER_EPOCH" ssz-size:"64"`
	XMSSBuilders                  []*Builder                          `json:"builders" ssz-index:"38" ssz-type:"progressive-list"`
	NextWithdrawalBuilderIndex    gloas.BuilderIndex                  `json:"next_withdrawal_builder_index" ssz-index:"39"`
	ExecutionPayloadAvailability  []uint8                             `json:"execution_payload_availability" ssz-index:"40" dynssz-size:"SLOTS_PER_HISTORICAL_ROOT/8" ssz-size:"1024"`
	BuilderPendingPayments        []*gloas.BuilderPendingPayment      `json:"builder_pending_payments" ssz-index:"41" dynssz-size:"SLOTS_PER_EPOCH*2" ssz-size:"64"`
	BuilderPendingWithdrawals     []*gloas.BuilderPendingWithdrawal   `json:"builder_pending_withdrawals" ssz-index:"42" ssz-type:"progressive-list"`
	LatestExecutionPayloadBid     *gloas.ExecutionPayloadBid          `json:"latest_execution_payload_bid" ssz-index:"43"`
	PayloadExpectedWithdrawals    []*capella.Withdrawal               `json:"payload_expected_withdrawals" ssz-index:"44" ssz-type:"progressive-list"`
	PTCWindow                     [][]phase0.ValidatorIndex           `json:"ptc_window" ssz-index:"45" dynssz-size:"(2+MIN_SEED_LOOKAHEAD)*SLOTS_PER_EPOCH,PTC_SIZE" ssz-size:"96,512"`
	ValidatorPreregistrations     []*ValidatorPreregistration         `json:"validator_preregistrations" ssz-index:"48" ssz-type:"progressive-list"`
	BuilderPreregistrations       []*BuilderPreregistration           `json:"builder_preregistrations" ssz-index:"49" ssz-type:"progressive-list"`
	AggregationCommitteeWindow    [][][]phase0.ValidatorIndex         `json:"aggregation_committee_window" ssz-index:"50" ssz-type:"vector,progressive-list,progressive-list" dynssz-size:"(2+MIN_SEED_LOOKAHEAD)*SLOTS_PER_EPOCH" ssz-size:"96"`
}
