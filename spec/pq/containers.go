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
	"errors"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
)

type Validator struct {
	PublicKey                  XMSSPubKey   `json:"pubkey" ssz-size:"32"`
	VRFRoot                    phase0.Root  `json:"vrf_root" ssz-size:"32"`
	WithdrawalCredentials      []byte       `json:"withdrawal_credentials" ssz-size:"32"`
	EffectiveBalance           phase0.Gwei  `json:"effective_balance"`
	Slashed                    bool         `json:"slashed"`
	ActivationEligibilityEpoch phase0.Epoch `json:"activation_eligibility_epoch"`
	ActivationEpoch            phase0.Epoch `json:"activation_epoch"`
	ExitEpoch                  phase0.Epoch `json:"exit_epoch"`
	WithdrawableEpoch          phase0.Epoch `json:"withdrawable_epoch"`
}

type Builder struct {
	PublicKey         XMSSPubKey                 `json:"pubkey" ssz-size:"32"`
	Version           uint8                      `json:"version"`
	ExecutionAddress  bellatrix.ExecutionAddress `json:"execution_address" ssz-size:"20"`
	Balance           phase0.Gwei                `json:"balance"`
	DepositEpoch      phase0.Epoch               `json:"deposit_epoch"`
	WithdrawableEpoch phase0.Epoch               `json:"withdrawable_epoch"`
}

type ValidatorPreregistrationRequest struct {
	PublicKey             XMSSPubKey  `json:"pubkey" ssz-size:"32"`
	WithdrawalCredentials []byte      `json:"withdrawal_credentials" ssz-size:"32"`
	VRFRoot               phase0.Root `json:"vrf_root" ssz-size:"32"`
}

type ValidatorPreregistration struct {
	PublicKey             XMSSPubKey  `json:"pubkey" ssz-size:"32"`
	WithdrawalCredentials []byte      `json:"withdrawal_credentials" ssz-size:"32"`
	VRFRoot               phase0.Root `json:"vrf_root" ssz-size:"32"`
	ExpirySlot            phase0.Slot `json:"expiry_slot"`
}

type BuilderPreregistrationRequest struct {
	PublicKey             XMSSPubKey `json:"pubkey" ssz-size:"32"`
	WithdrawalCredentials []byte     `json:"withdrawal_credentials" ssz-size:"32"`
}

type BuilderPreregistration struct {
	PublicKey             XMSSPubKey  `json:"pubkey" ssz-size:"32"`
	WithdrawalCredentials []byte      `json:"withdrawal_credentials" ssz-size:"32"`
	ExpirySlot            phase0.Slot `json:"expiry_slot"`
}

type VoluntaryExitRequest struct {
	PublicKey XMSSPubKey `json:"pubkey" ssz-size:"32"`
}

type ExecutionRequests struct {
	Deposits                  []*electra.DepositRequest          `ssz-index:"0" json:"deposits" ssz-type:"progressive-list"`
	Withdrawals               []*electra.WithdrawalRequest       `ssz-index:"1" json:"withdrawals" ssz-type:"progressive-list"`
	Consolidations            []*electra.ConsolidationRequest    `ssz-index:"2" json:"consolidations" ssz-type:"progressive-list"`
	BuilderDeposits           []*gloas.BuilderDepositRequest     `ssz-index:"3" json:"builder_deposits" ssz-type:"progressive-list"`
	BuilderExits              []*gloas.BuilderExitRequest        `ssz-index:"4" json:"builder_exits" ssz-type:"progressive-list"`
	ValidatorPreregistrations []*ValidatorPreregistrationRequest `ssz-index:"7" json:"validator_preregistrations" ssz-type:"progressive-list"`
	BuilderPreregistrations   []*BuilderPreregistrationRequest   `ssz-index:"8" json:"builder_preregistrations" ssz-type:"progressive-list"`
	VoluntaryExits            []*VoluntaryExitRequest            `ssz-index:"9" json:"voluntary_exits" ssz-type:"progressive-list"`
}

type SignedBeaconBlockHeader struct {
	Message       *phase0.BeaconBlockHeader `json:"message"`
	XMSSSignature XMSSSignature             `json:"signature" ssz-size:"2404"`
}

type ProposerSlashing struct {
	SignedHeader1 *SignedBeaconBlockHeader `json:"signed_header_1"`
	SignedHeader2 *SignedBeaconBlockHeader `json:"signed_header_2"`
}

type Attestation struct {
	AggregationBits bitfield.Bitlist        `ssz-index:"0" json:"aggregation_bits" ssz-type:"progressive-bitlist"`
	Data            *phase0.AttestationData `ssz-index:"1" json:"data"`
	CommitteeBits   bitfield.Bitvector64    `ssz-index:"3" dynssz-size:"MAX_COMMITTEES_PER_SLOT_PQ/8" json:"committee_bits" ssz-size:"2"`
}

type IndexedAttestation struct {
	AttestingIndices []uint64                `ssz-index:"0" json:"attesting_indices" ssz-type:"progressive-list"`
	Data             *phase0.AttestationData `ssz-index:"1" json:"data"`
	Proof            []byte                  `ssz-index:"2" dynssz-max:"MAX_STARK_PROOF_SIZE" json:"proof" ssz-max:"1048576"`
}

type AttesterSlashing struct {
	Attestation1 *IndexedAttestation `json:"attestation_1"`
	Attestation2 *IndexedAttestation `json:"attestation_2"`
}

type SingleAttestation struct {
	CommitteeIndex phase0.CommitteeIndex   `json:"committee_index"`
	AttesterIndex  phase0.ValidatorIndex   `json:"attester_index"`
	Data           *phase0.AttestationData `json:"data"`
	XMSSSignature  XMSSSignature           `json:"signature" ssz-size:"2404"`
}

type AggregateAndProof struct {
	AggregatorIndex phase0.ValidatorIndex `json:"aggregator_index"`
	Aggregate       *Attestation          `json:"aggregate"`
	Proof           []byte                `dynssz-max:"MAX_STARK_PROOF_SIZE" json:"proof" ssz-max:"1048576"`
}

type SignedAggregateAndProof struct {
	Message       *AggregateAndProof `json:"message"`
	XMSSSignature XMSSSignature      `json:"signature" ssz-size:"2404"`
}

type PayloadAttestationMessage struct {
	ValidatorIndex phase0.ValidatorIndex         `json:"validator_index"`
	Data           *gloas.PayloadAttestationData `json:"data"`
	XMSSSignature  XMSSSignature                 `json:"signature" ssz-size:"2404"`
}

type SignedExecutionPayloadBid struct {
	Message       *gloas.ExecutionPayloadBid `json:"message"`
	XMSSSignature XMSSSignature              `json:"signature" ssz-size:"2404"`
}

type ExecutionPayloadEnvelope struct {
	Payload               *gloas.ExecutionPayload `ssz-index:"0" json:"payload"`
	ExecutionRequests     *ExecutionRequests      `ssz-index:"1" json:"execution_requests"`
	BuilderIndex          gloas.BuilderIndex      `ssz-index:"2" json:"builder_index"`
	BeaconBlockRoot       phase0.Root             `ssz-index:"3" json:"beacon_block_root"`
	ParentBeaconBlockRoot phase0.Root             `ssz-index:"4" json:"parent_beacon_block_root"`
}

type SignedExecutionPayloadEnvelope struct {
	Message       *ExecutionPayloadEnvelope `json:"message"`
	XMSSSignature XMSSSignature             `json:"signature" ssz-size:"2404"`
}

type BeaconBlockBody struct {
	VRFReveal                 VRFProof                   `ssz-index:"0" json:"randao_reveal" ssz-size:"32,32"`
	ETH1Data                  *phase0.ETH1Data           `ssz-index:"1" json:"eth1_data"`
	Graffiti                  [32]byte                   `ssz-index:"2" json:"graffiti"`
	XMSSProposerSlashings     []*ProposerSlashing        `ssz-index:"3" json:"proposer_slashings" ssz-type:"progressive-list"`
	AttesterSlashings         []*AttesterSlashing        `ssz-index:"4" json:"attester_slashings" ssz-type:"progressive-list"`
	Attestations              []*Attestation             `ssz-index:"5" json:"attestations" ssz-type:"progressive-list"`
	Deposits                  []*phase0.Deposit          `ssz-index:"6" json:"deposits" ssz-type:"progressive-list"`
	SignedExecutionPayloadBid *SignedExecutionPayloadBid `ssz-index:"10" json:"signed_execution_payload_bid"`
	ParentExecutionRequests   *ExecutionRequests         `ssz-index:"12" json:"parent_execution_requests"`
	AttestationProof          []byte                     `ssz-index:"13" dynssz-max:"MAX_STARK_PROOF_SIZE" json:"proof" ssz-max:"1048576"`
}

type BeaconBlock struct {
	Slot          phase0.Slot           `json:"slot"`
	ProposerIndex phase0.ValidatorIndex `json:"proposer_index"`
	ParentRoot    phase0.Root           `json:"parent_root" ssz-size:"32"`
	StateRoot     phase0.Root           `json:"state_root" ssz-size:"32"`
	Body          *BeaconBlockBody      `json:"body"`
}

type SignedBeaconBlock struct {
	Message       *BeaconBlock  `json:"message"`
	XMSSSignature XMSSSignature `json:"signature" ssz-size:"2404"`
}

func (a *Attestation) CommitteeIndex() (phase0.CommitteeIndex, error) {
	indices := a.CommitteeBits.BitIndices()
	if len(indices) == 0 {
		return 0, errors.New("no committee index found in committee bits")
	}

	return phase0.CommitteeIndex(indices[0]), nil
}
