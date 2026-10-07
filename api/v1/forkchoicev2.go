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

package v1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// ForkChoicePayloadStatus is the payload status of a Gloas fork choice node.
type ForkChoicePayloadStatus uint64

const (
	// ForkChoicePayloadStatusUnknown is an unknown payload status.
	ForkChoicePayloadStatusUnknown ForkChoicePayloadStatus = iota
	// ForkChoicePayloadStatusEmpty is the node representing a block without its execution payload.
	ForkChoicePayloadStatusEmpty
	// ForkChoicePayloadStatusFull is the node representing a block with its execution payload.
	ForkChoicePayloadStatusFull
	// ForkChoicePayloadStatusPending is the node that is the parent of a block's empty and full nodes.
	ForkChoicePayloadStatusPending
)

// ForkChoicePayloadStatusStrings are the strings for fork choice payload status names.
var ForkChoicePayloadStatusStrings = [...]string{
	"unknown",
	"empty",
	"full",
	"pending",
}

// ForkChoicePayloadStatusFromString converts a string input to a fork choice payload status.
func ForkChoicePayloadStatusFromString(input string) (ForkChoicePayloadStatus, error) {
	switch strings.ToLower(input) {
	case "unknown":
		return ForkChoicePayloadStatusUnknown, nil
	case "empty":
		return ForkChoicePayloadStatusEmpty, nil
	case "full":
		return ForkChoicePayloadStatusFull, nil
	case "pending":
		return ForkChoicePayloadStatusPending, nil
	default:
		return ForkChoicePayloadStatusUnknown, fmt.Errorf("unrecognised fork choice payload status: %s", input)
	}
}

// String returns a string representation of the ForkChoicePayloadStatus.
func (s ForkChoicePayloadStatus) String() string {
	if uint64(s) >= uint64(len(ForkChoicePayloadStatusStrings)) {
		return "unknown"
	}

	return ForkChoicePayloadStatusStrings[s]
}

// ForkChoiceV2 is the data regarding the node's current fork choice context, as returned by
// GET /eth/v2/debug/fork_choice. It contains one node per (block root, payload status) pair.
//
// The endpoint is not yet finalised (ethereum/beacon-APIs#615) and client implementations differ,
// so decoding is lenient: fields that not every client provides are optional, and unrecognised
// fields are folded into ExtraData rather than discarded.
type ForkChoiceV2 struct {
	// JustifiedCheckpoint is the current justified checkpoint.
	JustifiedCheckpoint phase0.Checkpoint
	// FinalizedCheckpoint is the current finalized checkpoint.
	FinalizedCheckpoint phase0.Checkpoint
	// ForkChoiceNodes contains the fork choice nodes.
	ForkChoiceNodes []*ForkChoiceNodeV2
	// ExtraData is the client-specific extra data of the fork choice store.
	ExtraData map[string]any
}

// forkChoiceV2JSON is the json representation of the struct.
type forkChoiceV2JSON struct {
	JustifiedCheckpoint *phase0.Checkpoint  `json:"justified_checkpoint"`
	FinalizedCheckpoint *phase0.Checkpoint  `json:"finalized_checkpoint"`
	ForkChoiceNodes     []*ForkChoiceNodeV2 `json:"fork_choice_nodes"`
	ExtraData           map[string]any      `json:"extra_data"`
}

// forkChoiceV2Fields are the fields of forkChoiceV2JSON.
var forkChoiceV2Fields = map[string]struct{}{
	"justified_checkpoint": {},
	"finalized_checkpoint": {},
	"fork_choice_nodes":    {},
	"extra_data":           {},
}

// MarshalJSON implements json.Marshaler.
func (f ForkChoiceV2) MarshalJSON() ([]byte, error) {
	extraData := f.ExtraData
	if extraData == nil {
		extraData = map[string]any{}
	}

	nodes := f.ForkChoiceNodes
	if nodes == nil {
		nodes = []*ForkChoiceNodeV2{}
	}

	return json.Marshal(&forkChoiceV2JSON{
		JustifiedCheckpoint: &f.JustifiedCheckpoint,
		FinalizedCheckpoint: &f.FinalizedCheckpoint,
		ForkChoiceNodes:     nodes,
		ExtraData:           extraData,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *ForkChoiceV2) UnmarshalJSON(input []byte) error {
	var forkChoiceJSON forkChoiceV2JSON
	if err := json.Unmarshal(input, &forkChoiceJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if forkChoiceJSON.JustifiedCheckpoint == nil {
		return errors.New("justified checkpoint missing")
	}
	f.JustifiedCheckpoint = *forkChoiceJSON.JustifiedCheckpoint

	if forkChoiceJSON.FinalizedCheckpoint == nil {
		return errors.New("finalized checkpoint missing")
	}
	f.FinalizedCheckpoint = *forkChoiceJSON.FinalizedCheckpoint

	if forkChoiceJSON.ForkChoiceNodes == nil {
		return errors.New("fork choice nodes missing")
	}
	for i := range forkChoiceJSON.ForkChoiceNodes {
		if forkChoiceJSON.ForkChoiceNodes[i] == nil {
			return fmt.Errorf("fork choice node entry %d missing", i)
		}
	}
	f.ForkChoiceNodes = forkChoiceJSON.ForkChoiceNodes

	extraData, err := foldUnknownFields(input, forkChoiceV2Fields, forkChoiceJSON.ExtraData)
	if err != nil {
		return err
	}
	f.ExtraData = extraData

	return nil
}

// String returns a string version of the structure.
func (f *ForkChoiceV2) String() string {
	data, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

// ForkChoiceNodeV2 is a node in the fork choice tree, identified by its block root and payload status.
type ForkChoiceNodeV2 struct {
	// Slot is the slot of the beacon block.
	Slot phase0.Slot
	// BlockRoot is the block root of the beacon block.
	BlockRoot phase0.Root
	// PayloadStatus is the payload status of the node.
	PayloadStatus ForkChoicePayloadStatus
	// ParentRoot is the block root of the parent fork choice node.
	// For Gloas empty and full nodes this is the node's own block root, pointing at its pending node,
	// although not every client implements this yet.
	ParentRoot phase0.Root
	// ParentPayloadStatus is the payload status of the parent fork choice node.
	// Nil if the parent is not retained in the fork choice tree, or if the client does not provide it.
	ParentPayloadStatus *ForkChoicePayloadStatus
	// JustifiedEpoch is the justified epoch of the node, if provided.
	JustifiedEpoch *phase0.Epoch
	// FinalizedEpoch is the finalized epoch of the node, if provided.
	FinalizedEpoch *phase0.Epoch
	// Weight is the raw stored weight of the node in Gwei.
	Weight uint64
	// Validity is the execution validity of the chain ending at ExecutionBlockHash.
	Validity ForkChoiceNodeValidity
	// ExecutionBlockHash is the execution block hash of the node.
	ExecutionBlockHash phase0.Hash32
	// PayloadAttesterCount is the number of PTC positions with a recorded vote, if provided.
	PayloadAttesterCount *uint64
	// PayloadAvailabilityYesCount is the number of PTC positions voting that the payload was received on time, if provided.
	PayloadAvailabilityYesCount *uint64
	// PayloadDataAvailabilityYesCount is the number of PTC positions voting that the blob data is available, if provided.
	PayloadDataAvailabilityYesCount *uint64
	// ExtraData is the client-specific extra data of the node.
	ExtraData map[string]any
}

// forkChoiceNodeV2JSON is the json representation of the struct.
type forkChoiceNodeV2JSON struct {
	Slot                            string         `json:"slot"`
	BlockRoot                       string         `json:"block_root"`
	PayloadStatus                   string         `json:"payload_status"`
	ParentRoot                      string         `json:"parent_root"`
	ParentPayloadStatus             *string        `json:"parent_payload_status"`
	JustifiedEpoch                  string         `json:"justified_epoch,omitempty"`
	FinalizedEpoch                  string         `json:"finalized_epoch,omitempty"`
	Weight                          string         `json:"weight"`
	Validity                        string         `json:"validity"`
	ExecutionBlockHash              string         `json:"execution_block_hash"`
	PayloadAttesterCount            string         `json:"payload_attester_count,omitempty"`
	PayloadAvailabilityYesCount     string         `json:"payload_availability_yes_count,omitempty"`
	PayloadDataAvailabilityYesCount string         `json:"payload_data_availability_yes_count,omitempty"`
	ExtraData                       map[string]any `json:"extra_data"`
}

// forkChoiceNodeV2Fields are the fields of forkChoiceNodeV2JSON.
var forkChoiceNodeV2Fields = map[string]struct{}{
	"slot":                                {},
	"block_root":                          {},
	"payload_status":                      {},
	"parent_root":                         {},
	"parent_payload_status":               {},
	"justified_epoch":                     {},
	"finalized_epoch":                     {},
	"weight":                              {},
	"validity":                            {},
	"execution_block_hash":                {},
	"payload_attester_count":              {},
	"payload_availability_yes_count":      {},
	"payload_data_availability_yes_count": {},
	"extra_data":                          {},
}

// MarshalJSON implements json.Marshaler.
func (f ForkChoiceNodeV2) MarshalJSON() ([]byte, error) {
	data := &forkChoiceNodeV2JSON{
		Slot:                            fmt.Sprintf("%d", f.Slot),
		BlockRoot:                       fmt.Sprintf("%#x", f.BlockRoot),
		PayloadStatus:                   f.PayloadStatus.String(),
		ParentRoot:                      fmt.Sprintf("%#x", f.ParentRoot),
		ParentPayloadStatus:             nil,
		JustifiedEpoch:                  formatOptionalUint64((*uint64)(f.JustifiedEpoch)),
		FinalizedEpoch:                  formatOptionalUint64((*uint64)(f.FinalizedEpoch)),
		Weight:                          strconv.FormatUint(f.Weight, 10),
		Validity:                        f.Validity.String(),
		ExecutionBlockHash:              fmt.Sprintf("%#x", f.ExecutionBlockHash),
		PayloadAttesterCount:            formatOptionalUint64(f.PayloadAttesterCount),
		PayloadAvailabilityYesCount:     formatOptionalUint64(f.PayloadAvailabilityYesCount),
		PayloadDataAvailabilityYesCount: formatOptionalUint64(f.PayloadDataAvailabilityYesCount),
		ExtraData:                       f.ExtraData,
	}

	if f.ParentPayloadStatus != nil {
		parentPayloadStatus := f.ParentPayloadStatus.String()
		data.ParentPayloadStatus = &parentPayloadStatus
	}

	if data.ExtraData == nil {
		data.ExtraData = map[string]any{}
	}

	return json.Marshal(data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *ForkChoiceNodeV2) UnmarshalJSON(input []byte) error {
	var nodeJSON forkChoiceNodeV2JSON
	if err := json.Unmarshal(input, &nodeJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if nodeJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(nodeJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for slot: %s", nodeJSON.Slot))
	}
	f.Slot = phase0.Slot(slot)

	if nodeJSON.BlockRoot == "" {
		return errors.New("block root missing")
	}
	if err := decodeFixedBytes(f.BlockRoot[:], nodeJSON.BlockRoot, rootLength, "block root"); err != nil {
		return err
	}

	if nodeJSON.PayloadStatus == "" {
		return errors.New("payload status missing")
	}
	if f.PayloadStatus, err = ForkChoicePayloadStatusFromString(nodeJSON.PayloadStatus); err != nil {
		return err
	}

	// The parent root is null for the oldest retained node on some clients.
	if nodeJSON.ParentRoot != "" {
		if err := decodeFixedBytes(f.ParentRoot[:], nodeJSON.ParentRoot, rootLength, "parent root"); err != nil {
			return err
		}
	}

	if nodeJSON.ParentPayloadStatus != nil {
		parentPayloadStatus, err := ForkChoicePayloadStatusFromString(*nodeJSON.ParentPayloadStatus)
		if err != nil {
			return errors.Wrap(err, "invalid value for parent payload status")
		}
		f.ParentPayloadStatus = &parentPayloadStatus
	}

	var justifiedEpoch, finalizedEpoch *uint64
	if err := parseOptionalUint64(&justifiedEpoch, nodeJSON.JustifiedEpoch, "justified epoch"); err != nil {
		return err
	}
	f.JustifiedEpoch = (*phase0.Epoch)(justifiedEpoch)

	if err := parseOptionalUint64(&finalizedEpoch, nodeJSON.FinalizedEpoch, "finalized epoch"); err != nil {
		return err
	}
	f.FinalizedEpoch = (*phase0.Epoch)(finalizedEpoch)

	if nodeJSON.Weight == "" {
		return errors.New("weight missing")
	}
	if f.Weight, err = strconv.ParseUint(nodeJSON.Weight, 10, 64); err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for weight: %s", nodeJSON.Weight))
	}

	if nodeJSON.Validity == "" {
		return errors.New("validity missing")
	}
	// Validities beyond the spec's decode as unknown, keeping the original in
	// ExtraData, as for v1 nodes.
	unrecognisedValidity := false
	if f.Validity, err = ForkChoiceNodeValidityFromString(nodeJSON.Validity); err != nil {
		f.Validity = ForkChoiceNodeValidityUnknown
		unrecognisedValidity = true
	}

	if nodeJSON.ExecutionBlockHash == "" {
		return errors.New("execution block hash missing")
	}
	if err := decodeFixedBytes(f.ExecutionBlockHash[:], nodeJSON.ExecutionBlockHash, rootLength, "execution block hash"); err != nil {
		return err
	}

	if err := parseOptionalUint64(&f.PayloadAttesterCount, nodeJSON.PayloadAttesterCount, "payload attester count"); err != nil {
		return err
	}
	if err := parseOptionalUint64(&f.PayloadAvailabilityYesCount, nodeJSON.PayloadAvailabilityYesCount, "payload availability yes count"); err != nil {
		return err
	}
	if err := parseOptionalUint64(&f.PayloadDataAvailabilityYesCount, nodeJSON.PayloadDataAvailabilityYesCount, "payload data availability yes count"); err != nil {
		return err
	}

	if f.ExtraData, err = foldUnknownFields(input, forkChoiceNodeV2Fields, nodeJSON.ExtraData); err != nil {
		return err
	}

	if unrecognisedValidity {
		f.ExtraData = keepUnrecognisedValidity(f.ExtraData, nodeJSON.Validity)
	}

	return nil
}

// String returns a string version of the structure.
func (f *ForkChoiceNodeV2) String() string {
	data, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

// parseOptionalUint64 parses a decimal string into dst, leaving dst nil if the string is empty.
func parseOptionalUint64(dst **uint64, input string, name string) error {
	*dst = nil
	if input == "" {
		return nil
	}

	val, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for %s: %s", name, input))
	}
	*dst = &val

	return nil
}

// formatOptionalUint64 formats a value as a decimal string, returning an empty string if it is nil.
func formatOptionalUint64(input *uint64) string {
	if input == nil {
		return ""
	}

	return strconv.FormatUint(*input, 10)
}

// foldUnknownFields returns extraData with any fields of the input object that are not in known added to it.
// Fields already present in extraData take precedence. Only the values of unknown fields are decoded, so the
// (large) known fields such as fork_choice_nodes are not decoded a second time.
func foldUnknownFields(input []byte, known map[string]struct{}, extraData map[string]any) (map[string]any, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(input, &fields); err != nil {
		return nil, errors.Wrap(err, "invalid JSON")
	}

	for k, raw := range fields {
		if _, isKnown := known[k]; isKnown {
			continue
		}

		if _, exists := extraData[k]; exists {
			continue
		}

		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("invalid value for %s", k))
		}

		if extraData == nil {
			extraData = make(map[string]any)
		}

		extraData[k] = value
	}

	return extraData, nil
}
