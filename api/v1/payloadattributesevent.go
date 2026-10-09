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

package v1

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// PayloadAttributesEvent represents the data of a payload_attributes event.
type PayloadAttributesEvent struct {
	// Version is the fork version of the beacon chain.
	Version spec.DataVersion
	// Data is the data of the event.
	Data *PayloadAttributesData
}

// PayloadAttributesData represents the data of a payload_attributes event.
type PayloadAttributesData struct {
	// ProposerIndex is the index of the proposer.
	ProposerIndex phase0.ValidatorIndex
	// ProposalSlot is the slot of the proposal.
	ProposalSlot phase0.Slot
	// ParentBlockNumber is the number of the parent block.
	// The field is not part of the event from Gloas onwards: it is 0 there unless
	// the beacon node still sends it, and is never marshaled.
	ParentBlockNumber uint64
	// ParentBlockRoot is the root of the parent block.
	ParentBlockRoot phase0.Root
	// ParentBlockHash is the hash of the parent block.
	ParentBlockHash phase0.Hash32
	// SafeBlockHash is the execution block hash the node would pass as
	// safeBlockHash in engine_forkchoiceUpdated. Part of the event from Gloas
	// onwards; nil when the beacon node does not send it.
	SafeBlockHash *phase0.Hash32
	// FinalizedBlockHash is the execution block hash the node would pass as
	// finalizedBlockHash in engine_forkchoiceUpdated. Part of the event from
	// Gloas onwards; nil when the beacon node does not send it.
	FinalizedBlockHash *phase0.Hash32
	// V1 is the v1 payload attributes (Bellatrix).
	V1 *PayloadAttributesV1
	// V2 is the v2 payload attributes (Capella).
	V2 *PayloadAttributesV2
	// V3 is the v3 payload attributes (Deneb, Electra, Fulu).
	V3 *PayloadAttributesV3
	// V4 is the v4 payload attributes (Gloas).
	V4 *PayloadAttributesV4
	// V5 is the v5 payload attributes (Heze).
	V5 *PayloadAttributesV5
}

// PayloadAttributesV1 represents the payload attributes.
type PayloadAttributesV1 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient bellatrix.ExecutionAddress
}

// PayloadAttributesV2 represents the payload attributes v2.
type PayloadAttributesV2 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient bellatrix.ExecutionAddress
	// Withdrawals is the list of withdrawals.
	Withdrawals []*capella.Withdrawal
}

// PayloadAttributesV3 represents the payload attributes v3.
type PayloadAttributesV3 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient bellatrix.ExecutionAddress
	// Withdrawals is the list of withdrawals.
	Withdrawals []*capella.Withdrawal
	// ParentBeaconBlockRoot is the parent beacon block root.
	ParentBeaconBlockRoot phase0.Root
}

// PayloadAttributesV4 represents the payload attributes v4.
type PayloadAttributesV4 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient bellatrix.ExecutionAddress
	// Withdrawals is the list of withdrawals.
	Withdrawals []*capella.Withdrawal
	// ParentBeaconBlockRoot is the parent beacon block root.
	ParentBeaconBlockRoot phase0.Root
	// SlotNumber is the slot number of the payload.
	SlotNumber uint64
	// TargetGasLimit is the target gas limit of the payload.
	TargetGasLimit uint64
}

// PayloadAttributesV5 represents the payload attributes v5.
type PayloadAttributesV5 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient bellatrix.ExecutionAddress
	// Withdrawals is the list of withdrawals.
	Withdrawals []*capella.Withdrawal
	// ParentBeaconBlockRoot is the parent beacon block root.
	ParentBeaconBlockRoot phase0.Root
	// SlotNumber is the slot number of the payload.
	SlotNumber uint64
	// TargetGasLimit is the target gas limit of the payload.
	TargetGasLimit uint64
	// InclusionListTransactions is the list of inclusion list transactions.
	InclusionListTransactions []bellatrix.Transaction
}

// payloadAttributesEventJSON is the spec representation of the event.
type payloadAttributesEventJSON struct {
	Version spec.DataVersion           `json:"version"`
	Data    *payloadAttributesDataJSON `json:"data"`
}

// payloadAttributesDataJSON is the spec representation of the payload attributes data.
type payloadAttributesDataJSON struct {
	ProposerIndex      string          `json:"proposer_index"`
	ProposalSlot       string          `json:"proposal_slot"`
	ParentBlockNumber  string          `json:"parent_block_number,omitempty"`
	ParentBlockRoot    string          `json:"parent_block_root"`
	ParentBlockHash    string          `json:"parent_block_hash"`
	SafeBlockHash      string          `json:"safe_block_hash,omitempty"`
	FinalizedBlockHash string          `json:"finalized_block_hash,omitempty"`
	PayloadAttributes  json.RawMessage `json:"payload_attributes"`
}

// payloadAttributesV1JSON is the spec representation of the payload attributes.
type payloadAttributesV1JSON struct {
	Timestamp             string `json:"timestamp"`
	PrevRandao            string `json:"prev_randao"`
	SuggestedFeeRecipient string `json:"suggested_fee_recipient"`
}

// payloadAttributesV2JSON is the spec representation of the payload attributes v2.
type payloadAttributesV2JSON struct {
	Timestamp             string                `json:"timestamp"`
	PrevRandao            string                `json:"prev_randao"`
	SuggestedFeeRecipient string                `json:"suggested_fee_recipient"`
	Withdrawals           []*capella.Withdrawal `json:"withdrawals"`
}

// payloadAttributesV3JSON is the spec representation of the payload attributes v3.
type payloadAttributesV3JSON struct {
	Timestamp             string                `json:"timestamp"`
	PrevRandao            string                `json:"prev_randao"`
	SuggestedFeeRecipient string                `json:"suggested_fee_recipient"`
	Withdrawals           []*capella.Withdrawal `json:"withdrawals"`
	ParentBeaconBlockRoot string                `json:"parent_beacon_block_root"`
}

// payloadAttributesV4JSON is the spec representation of the payload attributes v4.
type payloadAttributesV4JSON struct {
	Timestamp             string                `json:"timestamp"`
	PrevRandao            string                `json:"prev_randao"`
	SuggestedFeeRecipient string                `json:"suggested_fee_recipient"`
	Withdrawals           []*capella.Withdrawal `json:"withdrawals"`
	ParentBeaconBlockRoot string                `json:"parent_beacon_block_root"`
	SlotNumber            string                `json:"slot_number"`
	TargetGasLimit        string                `json:"target_gas_limit"`
}

// payloadAttributesV5JSON is the spec representation of the payload attributes v5.
type payloadAttributesV5JSON struct {
	Timestamp                 string                `json:"timestamp"`
	PrevRandao                string                `json:"prev_randao"`
	SuggestedFeeRecipient     string                `json:"suggested_fee_recipient"`
	Withdrawals               []*capella.Withdrawal `json:"withdrawals"`
	ParentBeaconBlockRoot     string                `json:"parent_beacon_block_root"`
	SlotNumber                string                `json:"slot_number"`
	TargetGasLimit            string                `json:"target_gas_limit"`
	InclusionListTransactions []string              `json:"inclusion_list_transactions"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV1) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV1JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV1) unpack(data *payloadAttributesV1JSON) error {
	var err error

	if data.Timestamp == "" {
		return errors.New("payload attributes timestamp missing")
	}

	p.Timestamp, err = strconv.ParseUint(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes timestamp")
	}

	if data.PrevRandao == "" {
		return errors.New("payload attributes prev randao missing")
	}

	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes prev randao")
	}

	if len(prevRandao) != 32 {
		return errors.New("incorrect length for payload attributes prev randao")
	}

	copy(p.PrevRandao[:], prevRandao)

	if data.SuggestedFeeRecipient == "" {
		return errors.New("payload attributes suggested fee recipient missing")
	}

	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.SuggestedFeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes suggested fee recipient")
	}

	if len(feeRecipient) != bellatrix.FeeRecipientLength {
		return errors.New("incorrect length for payload attributes suggested fee recipient")
	}

	copy(p.SuggestedFeeRecipient[:], feeRecipient)

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV2) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV2JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV2) unpack(data *payloadAttributesV2JSON) error {
	var err error

	if data.Timestamp == "" {
		return errors.New("payload attributes timestamp missing")
	}

	p.Timestamp, err = strconv.ParseUint(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes timestamp")
	}

	if data.PrevRandao == "" {
		return errors.New("payload attributes prev randao missing")
	}

	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes prev randao")
	}

	if len(prevRandao) != 32 {
		return errors.New("incorrect length for payload attributes prev randao")
	}

	copy(p.PrevRandao[:], prevRandao)

	if data.SuggestedFeeRecipient == "" {
		return errors.New("payload attributes suggested fee recipient missing")
	}

	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.SuggestedFeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes suggested fee recipient")
	}

	if len(feeRecipient) != bellatrix.FeeRecipientLength {
		return errors.New("incorrect length for payload attributes suggested fee recipient")
	}

	copy(p.SuggestedFeeRecipient[:], feeRecipient)

	if data.Withdrawals == nil {
		return errors.New("payload attributes withdrawals missing")
	}

	for i := range data.Withdrawals {
		if data.Withdrawals[i] == nil {
			return fmt.Errorf("withdrawals entry %d missing", i)
		}
	}

	p.Withdrawals = data.Withdrawals

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV3) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV3JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV3) unpack(data *payloadAttributesV3JSON) error {
	var err error

	if data.Timestamp == "" {
		return errors.New("payload attributes timestamp missing")
	}

	p.Timestamp, err = strconv.ParseUint(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes timestamp")
	}

	if data.PrevRandao == "" {
		return errors.New("payload attributes prev randao missing")
	}

	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes prev randao")
	}

	if len(prevRandao) != 32 {
		return errors.New("incorrect length for payload attributes prev randao")
	}

	copy(p.PrevRandao[:], prevRandao)

	if data.SuggestedFeeRecipient == "" {
		return errors.New("payload attributes suggested fee recipient missing")
	}

	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.SuggestedFeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes suggested fee recipient")
	}

	if len(feeRecipient) != bellatrix.FeeRecipientLength {
		return errors.New("incorrect length for payload attributes suggested fee recipient")
	}

	copy(p.SuggestedFeeRecipient[:], feeRecipient)

	if data.Withdrawals == nil {
		return errors.New("payload attributes withdrawals missing")
	}

	for i := range data.Withdrawals {
		if data.Withdrawals[i] == nil {
			return fmt.Errorf("withdrawals entry %d missing", i)
		}
	}

	p.Withdrawals = data.Withdrawals

	if data.ParentBeaconBlockRoot == "" {
		return errors.New("payload attributes parent beacon block root missing")
	}

	parentBeaconBlockRoot, err := hex.DecodeString(strings.TrimPrefix(data.ParentBeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes parent beacon block root")
	}

	if len(parentBeaconBlockRoot) != phase0.RootLength {
		return errors.New("incorrect length for payload attributes parent beacon block root")
	}

	copy(p.ParentBeaconBlockRoot[:], parentBeaconBlockRoot)

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV4) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV4JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV4) unpack(data *payloadAttributesV4JSON) error {
	var v3 PayloadAttributesV3

	err := v3.unpack(&payloadAttributesV3JSON{
		Timestamp:             data.Timestamp,
		PrevRandao:            data.PrevRandao,
		SuggestedFeeRecipient: data.SuggestedFeeRecipient,
		Withdrawals:           data.Withdrawals,
		ParentBeaconBlockRoot: data.ParentBeaconBlockRoot,
	})
	if err != nil {
		return err
	}

	p.Timestamp = v3.Timestamp
	p.PrevRandao = v3.PrevRandao
	p.SuggestedFeeRecipient = v3.SuggestedFeeRecipient
	p.Withdrawals = v3.Withdrawals
	p.ParentBeaconBlockRoot = v3.ParentBeaconBlockRoot

	if data.SlotNumber == "" {
		return errors.New("payload attributes slot number missing")
	}

	p.SlotNumber, err = strconv.ParseUint(data.SlotNumber, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes slot number")
	}

	if data.TargetGasLimit == "" {
		return errors.New("payload attributes target gas limit missing")
	}

	p.TargetGasLimit, err = strconv.ParseUint(data.TargetGasLimit, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes target gas limit")
	}

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV5) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV5JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV5) unpack(data *payloadAttributesV5JSON) error {
	var v4 PayloadAttributesV4

	err := v4.unpack(&payloadAttributesV4JSON{
		Timestamp:             data.Timestamp,
		PrevRandao:            data.PrevRandao,
		SuggestedFeeRecipient: data.SuggestedFeeRecipient,
		Withdrawals:           data.Withdrawals,
		ParentBeaconBlockRoot: data.ParentBeaconBlockRoot,
		SlotNumber:            data.SlotNumber,
		TargetGasLimit:        data.TargetGasLimit,
	})
	if err != nil {
		return err
	}

	p.Timestamp = v4.Timestamp
	p.PrevRandao = v4.PrevRandao
	p.SuggestedFeeRecipient = v4.SuggestedFeeRecipient
	p.Withdrawals = v4.Withdrawals
	p.ParentBeaconBlockRoot = v4.ParentBeaconBlockRoot
	p.SlotNumber = v4.SlotNumber
	p.TargetGasLimit = v4.TargetGasLimit

	if data.InclusionListTransactions == nil {
		return errors.New("payload attributes inclusion list transactions missing")
	}

	p.InclusionListTransactions = make([]bellatrix.Transaction, len(data.InclusionListTransactions))

	for i, transaction := range data.InclusionListTransactions {
		if transaction == "" {
			return fmt.Errorf("inclusion list transactions entry %d missing", i)
		}

		p.InclusionListTransactions[i], err = hex.DecodeString(strings.TrimPrefix(transaction, "0x"))
		if err != nil {
			return errors.Wrapf(err, "invalid value for inclusion list transactions entry %d", i)
		}
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (e *PayloadAttributesEvent) MarshalJSON() ([]byte, error) {
	var (
		payloadAttributes []byte
		err               error
	)

	switch e.Version {
	case spec.DataVersionBellatrix:
		if e.Data.V1 == nil {
			return nil, errors.New("no payload attributes v1 data")
		}

		payloadAttributes, err = json.Marshal(&payloadAttributesV1JSON{
			Timestamp:             strconv.FormatUint(e.Data.V1.Timestamp, 10),
			PrevRandao:            fmt.Sprintf("%#x", e.Data.V1.PrevRandao),
			SuggestedFeeRecipient: e.Data.V1.SuggestedFeeRecipient.String(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v1")
		}
	case spec.DataVersionCapella:
		if e.Data.V2 == nil {
			return nil, errors.New("no payload attributes v2 data")
		}

		payloadAttributes, err = json.Marshal(&payloadAttributesV2JSON{
			Timestamp:             strconv.FormatUint(e.Data.V2.Timestamp, 10),
			PrevRandao:            fmt.Sprintf("%#x", e.Data.V2.PrevRandao),
			SuggestedFeeRecipient: e.Data.V2.SuggestedFeeRecipient.String(),
			Withdrawals:           e.Data.V2.Withdrawals,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v2")
		}
	case spec.DataVersionDeneb, spec.DataVersionElectra, spec.DataVersionFulu:
		if e.Data.V3 == nil {
			return nil, errors.New("no payload attributes v3 data")
		}

		payloadAttributes, err = json.Marshal(&payloadAttributesV3JSON{
			Timestamp:             strconv.FormatUint(e.Data.V3.Timestamp, 10),
			PrevRandao:            fmt.Sprintf("%#x", e.Data.V3.PrevRandao),
			SuggestedFeeRecipient: e.Data.V3.SuggestedFeeRecipient.String(),
			Withdrawals:           e.Data.V3.Withdrawals,
			ParentBeaconBlockRoot: fmt.Sprintf("%#x", e.Data.V3.ParentBeaconBlockRoot),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v3")
		}
	case spec.DataVersionGloas:
		if e.Data.V4 == nil {
			return nil, errors.New("no payload attributes v4 data")
		}

		payloadAttributes, err = json.Marshal(&payloadAttributesV4JSON{
			Timestamp:             strconv.FormatUint(e.Data.V4.Timestamp, 10),
			PrevRandao:            fmt.Sprintf("%#x", e.Data.V4.PrevRandao),
			SuggestedFeeRecipient: e.Data.V4.SuggestedFeeRecipient.String(),
			Withdrawals:           e.Data.V4.Withdrawals,
			ParentBeaconBlockRoot: fmt.Sprintf("%#x", e.Data.V4.ParentBeaconBlockRoot),
			SlotNumber:            strconv.FormatUint(e.Data.V4.SlotNumber, 10),
			TargetGasLimit:        strconv.FormatUint(e.Data.V4.TargetGasLimit, 10),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v4")
		}
	case spec.DataVersionHeze:
		if e.Data.V5 == nil {
			return nil, errors.New("no payload attributes v5 data")
		}

		inclusionListTransactions := make([]string, len(e.Data.V5.InclusionListTransactions))
		for i, transaction := range e.Data.V5.InclusionListTransactions {
			inclusionListTransactions[i] = fmt.Sprintf("%#x", []byte(transaction))
		}

		payloadAttributes, err = json.Marshal(&payloadAttributesV5JSON{
			Timestamp:                 strconv.FormatUint(e.Data.V5.Timestamp, 10),
			PrevRandao:                fmt.Sprintf("%#x", e.Data.V5.PrevRandao),
			SuggestedFeeRecipient:     e.Data.V5.SuggestedFeeRecipient.String(),
			Withdrawals:               e.Data.V5.Withdrawals,
			ParentBeaconBlockRoot:     fmt.Sprintf("%#x", e.Data.V5.ParentBeaconBlockRoot),
			SlotNumber:                strconv.FormatUint(e.Data.V5.SlotNumber, 10),
			TargetGasLimit:            strconv.FormatUint(e.Data.V5.TargetGasLimit, 10),
			InclusionListTransactions: inclusionListTransactions,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v5")
		}
	default:
		return nil, fmt.Errorf("unsupported payload attributes version: %s", e.Version)
	}

	data := payloadAttributesDataJSON{
		ProposerIndex:     fmt.Sprintf("%d", e.Data.ProposerIndex),
		ProposalSlot:      fmt.Sprintf("%d", e.Data.ProposalSlot),
		ParentBlockRoot:   fmt.Sprintf("%#x", e.Data.ParentBlockRoot),
		ParentBlockHash:   fmt.Sprintf("%#x", e.Data.ParentBlockHash),
		PayloadAttributes: payloadAttributes,
	}

	// The parent block number is not part of the event from Gloas onwards.
	if e.Version < spec.DataVersionGloas {
		data.ParentBlockNumber = strconv.FormatUint(e.Data.ParentBlockNumber, 10)
	}

	if e.Data.SafeBlockHash != nil {
		data.SafeBlockHash = fmt.Sprintf("%#x", *e.Data.SafeBlockHash)
	}

	if e.Data.FinalizedBlockHash != nil {
		data.FinalizedBlockHash = fmt.Sprintf("%#x", *e.Data.FinalizedBlockHash)
	}

	return json.Marshal(&payloadAttributesEventJSON{
		Version: e.Version,
		Data:    &data,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *PayloadAttributesEvent) UnmarshalJSON(input []byte) error {
	var event payloadAttributesEventJSON
	if err := json.Unmarshal(input, &event); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return e.unpack(&event)
}

// String returns a string version of the structure.
func (e *PayloadAttributesEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

func (e *PayloadAttributesEvent) unpack(data *payloadAttributesEventJSON) error {
	var err error

	if data.Data == nil {
		return errors.New("payload attributes data missing")
	}

	e.Data = &PayloadAttributesData{}

	if data.Data.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}

	proposerIndex, err := strconv.ParseUint(data.Data.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}

	e.Data.ProposerIndex = phase0.ValidatorIndex(proposerIndex)

	if data.Data.ProposalSlot == "" {
		return errors.New("proposal slot missing")
	}

	proposalSlot, err := strconv.ParseUint(data.Data.ProposalSlot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposal slot")
	}

	e.Data.ProposalSlot = phase0.Slot(proposalSlot)

	// The parent block number is not part of the event from Gloas onwards;
	// beacon nodes that still send it are tolerated.
	if data.Data.ParentBlockNumber != "" {
		parentBlockNumber, err := strconv.ParseUint(data.Data.ParentBlockNumber, 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for parent block number")
		}

		e.Data.ParentBlockNumber = parentBlockNumber
	} else if data.Version < spec.DataVersionGloas {
		return errors.New("parent block number missing")
	}

	if data.Data.ParentBlockRoot == "" {
		return errors.New("parent block root missing")
	}

	parentBlockRoot, err := hex.DecodeString(strings.TrimPrefix(data.Data.ParentBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent block root")
	}

	if len(parentBlockRoot) != phase0.RootLength {
		return errors.New("incorrect length for parent block root")
	}

	copy(e.Data.ParentBlockRoot[:], parentBlockRoot)

	if data.Data.ParentBlockHash == "" {
		return errors.New("parent block hash missing")
	}

	parentBlockHash, err := hex.DecodeString(strings.TrimPrefix(data.Data.ParentBlockHash, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent block hash")
	}

	if len(parentBlockHash) != phase0.Hash32Length {
		return errors.New("incorrect length for parent block hash")
	}

	copy(e.Data.ParentBlockHash[:], parentBlockHash)

	// The safe and finalized block hashes are part of the event from Gloas
	// onwards; beacon nodes that do not send them yet are tolerated.
	if data.Data.SafeBlockHash != "" {
		safeBlockHash, err := hex.DecodeString(strings.TrimPrefix(data.Data.SafeBlockHash, "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for safe block hash")
		}

		if len(safeBlockHash) != phase0.Hash32Length {
			return errors.New("incorrect length for safe block hash")
		}

		var hash phase0.Hash32
		copy(hash[:], safeBlockHash)
		e.Data.SafeBlockHash = &hash
	}

	if data.Data.FinalizedBlockHash != "" {
		finalizedBlockHash, err := hex.DecodeString(strings.TrimPrefix(data.Data.FinalizedBlockHash, "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for finalized block hash")
		}

		if len(finalizedBlockHash) != phase0.Hash32Length {
			return errors.New("incorrect length for finalized block hash")
		}

		var hash phase0.Hash32
		copy(hash[:], finalizedBlockHash)
		e.Data.FinalizedBlockHash = &hash
	}

	if data.Data.PayloadAttributes == nil {
		return errors.New("payload attributes missing")
	}

	switch data.Version {
	case spec.DataVersionBellatrix:
		var payloadAttributes PayloadAttributesV1

		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}

		e.Data.V1 = &payloadAttributes
	case spec.DataVersionCapella:
		var payloadAttributes PayloadAttributesV2

		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}

		e.Data.V2 = &payloadAttributes
	case spec.DataVersionDeneb, spec.DataVersionElectra, spec.DataVersionFulu:
		var payloadAttributes PayloadAttributesV3

		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}

		e.Data.V3 = &payloadAttributes
	case spec.DataVersionGloas:
		var payloadAttributes PayloadAttributesV4

		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}

		e.Data.V4 = &payloadAttributes
	case spec.DataVersionHeze:
		var payloadAttributes PayloadAttributesV5

		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}

		e.Data.V5 = &payloadAttributes
	default:
		return errors.New("unsupported data version")
	}

	e.Version = data.Version

	return nil
}
