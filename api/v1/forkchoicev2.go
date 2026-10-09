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
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

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
// Only the values the beacon API defines (pending, empty, full) are accepted.
func ForkChoicePayloadStatusFromString(input string) (ForkChoicePayloadStatus, error) {
	switch input {
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

// MarshalJSON implements json.Marshaler.
func (s ForkChoicePayloadStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *ForkChoicePayloadStatus) UnmarshalJSON(input []byte) error {
	var inputString string
	if err := json.Unmarshal(input, &inputString); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	status, err := ForkChoicePayloadStatusFromString(inputString)
	if err != nil {
		return err
	}
	*s = status

	return nil
}

// String returns a string representation of the ForkChoicePayloadStatus.
func (s ForkChoicePayloadStatus) String() string {
	if uint64(s) >= uint64(len(ForkChoicePayloadStatusStrings)) {
		return "unknown"
	}

	return ForkChoicePayloadStatusStrings[s]
}

// ForkChoiceV2 is the data regarding the node's current fork choice context, as returned by
// GET /eth/v2/debug/fork_choice (ethereum/beacon-APIs#615). It contains one node per
// (block root, payload status) pair.
//
// Decoding follows the spec: every field it requires must be present, and enumerations only
// accept the values it defines.
type ForkChoiceV2 struct {
	// JustifiedCheckpoint is the current justified checkpoint.
	JustifiedCheckpoint phase0.Checkpoint
	// FinalizedCheckpoint is the current finalized checkpoint.
	FinalizedCheckpoint phase0.Checkpoint
	// ForkChoiceNodes contains the fork choice nodes.
	ForkChoiceNodes []*ForkChoiceNodeV2
	// ExtraData is the client-specific extra data of the fork choice store; empty if the
	// client provides none.
	ExtraData map[string]any
}

// forkChoiceV2JSON is the json representation of the struct.
type forkChoiceV2JSON struct {
	JustifiedCheckpoint *phase0.Checkpoint  `json:"justified_checkpoint"`
	FinalizedCheckpoint *phase0.Checkpoint  `json:"finalized_checkpoint"`
	ForkChoiceNodes     []*ForkChoiceNodeV2 `json:"fork_choice_nodes"`
	ExtraData           map[string]any      `json:"extra_data"`
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
	fields, err := decodeObject(input)
	if err != nil {
		return err
	}

	var forkChoiceJSON forkChoiceV2JSON
	for _, field := range [...]struct {
		key string
		dst any
	}{
		{"justified_checkpoint", &forkChoiceJSON.JustifiedCheckpoint},
		{"finalized_checkpoint", &forkChoiceJSON.FinalizedCheckpoint},
		{"fork_choice_nodes", &forkChoiceJSON.ForkChoiceNodes},
		{"extra_data", &forkChoiceJSON.ExtraData},
	} {
		if err := takeField(fields, field.key, field.dst); err != nil {
			return err
		}
	}

	if forkChoiceJSON.JustifiedCheckpoint == nil {
		return errors.New("justified checkpoint missing")
	}
	f.JustifiedCheckpoint = *forkChoiceJSON.JustifiedCheckpoint

	if forkChoiceJSON.FinalizedCheckpoint == nil {
		return errors.New("finalized checkpoint missing")
	}
	f.FinalizedCheckpoint = *forkChoiceJSON.FinalizedCheckpoint

	if len(forkChoiceJSON.ForkChoiceNodes) == 0 {
		return errors.New("fork choice nodes missing")
	}
	for i := range forkChoiceJSON.ForkChoiceNodes {
		if forkChoiceJSON.ForkChoiceNodes[i] == nil {
			return fmt.Errorf("fork choice node entry %d missing", i)
		}
	}
	f.ForkChoiceNodes = forkChoiceJSON.ForkChoiceNodes

	if forkChoiceJSON.ExtraData == nil {
		return errors.New("extra data missing")
	}
	f.ExtraData = forkChoiceJSON.ExtraData

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
	// ParentRoot is the block root of the parent fork choice node. For Gloas empty and full nodes
	// this is the node's own block root, pointing at its pending node; otherwise it is the beacon
	// block's parent root.
	ParentRoot phase0.Root
	// ParentPayloadStatus is the payload status of the parent fork choice node, nil if the parent is
	// not retained in the fork choice tree.
	ParentPayloadStatus *ForkChoicePayloadStatus
	// JustifiedCheckpoint is the justified checkpoint recorded for the node.
	JustifiedCheckpoint phase0.Checkpoint
	// FinalizedCheckpoint is the finalized checkpoint recorded for the node.
	FinalizedCheckpoint phase0.Checkpoint
	// Weight is the raw stored weight of the node in Gwei.
	Weight phase0.Gwei
	// Validity is the execution validity of the chain ending at ExecutionBlockHash.
	Validity ForkChoiceNodeValidity
	// ExecutionBlockHash is the execution payload hash for full nodes, and the bid's parent block hash
	// for Gloas pending and empty nodes.
	ExecutionBlockHash phase0.Hash32
	// PayloadAttesterCount is the number of PTC positions with a recorded vote.
	PayloadAttesterCount uint64
	// PayloadAvailabilityYesCount is the number of PTC positions voting that the payload was received on time.
	PayloadAvailabilityYesCount uint64
	// PayloadDataAvailabilityYesCount is the number of PTC positions voting that the blob data is available.
	PayloadDataAvailabilityYesCount uint64
	// ExtraData is the client-specific extra data of the node; empty if the client provides none.
	ExtraData map[string]any
}

// forkChoiceNodeV2JSON is the json representation of the struct.
type forkChoiceNodeV2JSON struct {
	Slot                            string             `json:"slot"`
	BlockRoot                       string             `json:"block_root"`
	PayloadStatus                   string             `json:"payload_status"`
	ParentRoot                      string             `json:"parent_root"`
	ParentPayloadStatus             *string            `json:"parent_payload_status"`
	JustifiedCheckpoint             *phase0.Checkpoint `json:"justified_checkpoint"`
	FinalizedCheckpoint             *phase0.Checkpoint `json:"finalized_checkpoint"`
	Weight                          string             `json:"weight"`
	Validity                        string             `json:"validity"`
	ExecutionBlockHash              string             `json:"execution_block_hash"`
	PayloadAttesterCount            string             `json:"payload_attester_count"`
	PayloadAvailabilityYesCount     string             `json:"payload_availability_yes_count"`
	PayloadDataAvailabilityYesCount string             `json:"payload_data_availability_yes_count"`
	ExtraData                       map[string]any     `json:"extra_data"`
}

// forkChoiceNodeV2Validities are the execution validities the v2 node defines, matched exactly.
var forkChoiceNodeV2Validities = map[string]ForkChoiceNodeValidity{
	"valid":      ForkChoiceNodeValidityValid,
	"invalid":    ForkChoiceNodeValidityInvalid,
	"optimistic": ForkChoiceNodeValidityOptimistic,
}

// MarshalJSON implements json.Marshaler.
func (f ForkChoiceNodeV2) MarshalJSON() ([]byte, error) {
	data := &forkChoiceNodeV2JSON{
		Slot:                            fmt.Sprintf("%d", f.Slot),
		BlockRoot:                       fmt.Sprintf("%#x", f.BlockRoot),
		PayloadStatus:                   f.PayloadStatus.String(),
		ParentRoot:                      fmt.Sprintf("%#x", f.ParentRoot),
		JustifiedCheckpoint:             &f.JustifiedCheckpoint,
		FinalizedCheckpoint:             &f.FinalizedCheckpoint,
		Weight:                          fmt.Sprintf("%d", f.Weight),
		Validity:                        f.Validity.String(),
		ExecutionBlockHash:              fmt.Sprintf("%#x", f.ExecutionBlockHash),
		PayloadAttesterCount:            strconv.FormatUint(f.PayloadAttesterCount, 10),
		PayloadAvailabilityYesCount:     strconv.FormatUint(f.PayloadAvailabilityYesCount, 10),
		PayloadDataAvailabilityYesCount: strconv.FormatUint(f.PayloadDataAvailabilityYesCount, 10),
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
	fields, err := decodeObject(input)
	if err != nil {
		return err
	}

	// parent_payload_status is required but nullable, so its presence is checked separately from
	// its value.
	_, hasParentPayloadStatus := fields["parent_payload_status"]

	// Fields are taken in a fixed order, so that the first of several invalid fields is always the one
	// reported.
	var nodeJSON forkChoiceNodeV2JSON
	for _, field := range [...]struct {
		key string
		dst any
	}{
		{"slot", &nodeJSON.Slot},
		{"block_root", &nodeJSON.BlockRoot},
		{"payload_status", &nodeJSON.PayloadStatus},
		{"parent_root", &nodeJSON.ParentRoot},
		{"parent_payload_status", &nodeJSON.ParentPayloadStatus},
		{"justified_checkpoint", &nodeJSON.JustifiedCheckpoint},
		{"finalized_checkpoint", &nodeJSON.FinalizedCheckpoint},
		{"weight", &nodeJSON.Weight},
		{"validity", &nodeJSON.Validity},
		{"execution_block_hash", &nodeJSON.ExecutionBlockHash},
		{"payload_attester_count", &nodeJSON.PayloadAttesterCount},
		{"payload_availability_yes_count", &nodeJSON.PayloadAvailabilityYesCount},
		{"payload_data_availability_yes_count", &nodeJSON.PayloadDataAvailabilityYesCount},
		{"extra_data", &nodeJSON.ExtraData},
	} {
		if err := takeField(fields, field.key, field.dst); err != nil {
			return err
		}
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
	if err := decodeFixedBytes(f.BlockRoot[:], nodeJSON.BlockRoot, "block root"); err != nil {
		return err
	}

	if nodeJSON.PayloadStatus == "" {
		return errors.New("payload status missing")
	}
	if f.PayloadStatus, err = ForkChoicePayloadStatusFromString(nodeJSON.PayloadStatus); err != nil {
		return err
	}

	if nodeJSON.ParentRoot == "" {
		return errors.New("parent root missing")
	}
	if err := decodeFixedBytes(f.ParentRoot[:], nodeJSON.ParentRoot, "parent root"); err != nil {
		return err
	}

	if !hasParentPayloadStatus {
		return errors.New("parent payload status missing")
	}
	f.ParentPayloadStatus = nil
	if nodeJSON.ParentPayloadStatus != nil {
		parentPayloadStatus, err := ForkChoicePayloadStatusFromString(*nodeJSON.ParentPayloadStatus)
		if err != nil {
			return errors.Wrap(err, "invalid value for parent payload status")
		}
		f.ParentPayloadStatus = &parentPayloadStatus
	}

	if nodeJSON.JustifiedCheckpoint == nil {
		return errors.New("justified checkpoint missing")
	}
	f.JustifiedCheckpoint = *nodeJSON.JustifiedCheckpoint

	if nodeJSON.FinalizedCheckpoint == nil {
		return errors.New("finalized checkpoint missing")
	}
	f.FinalizedCheckpoint = *nodeJSON.FinalizedCheckpoint

	if nodeJSON.Weight == "" {
		return errors.New("weight missing")
	}
	weight, err := strconv.ParseUint(nodeJSON.Weight, 10, 64)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for weight: %s", nodeJSON.Weight))
	}
	f.Weight = phase0.Gwei(weight)

	if nodeJSON.Validity == "" {
		return errors.New("validity missing")
	}
	validity, known := forkChoiceNodeV2Validities[nodeJSON.Validity]
	if !known {
		return fmt.Errorf("invalid value for validity: %s", nodeJSON.Validity)
	}
	f.Validity = validity

	if nodeJSON.ExecutionBlockHash == "" {
		return errors.New("execution block hash missing")
	}
	if err := decodeFixedBytes(f.ExecutionBlockHash[:], nodeJSON.ExecutionBlockHash, "execution block hash"); err != nil {
		return err
	}

	for _, count := range [...]struct {
		input string
		dst   *uint64
		name  string
	}{
		{nodeJSON.PayloadAttesterCount, &f.PayloadAttesterCount, "payload attester count"},
		{nodeJSON.PayloadAvailabilityYesCount, &f.PayloadAvailabilityYesCount, "payload availability yes count"},
		{nodeJSON.PayloadDataAvailabilityYesCount, &f.PayloadDataAvailabilityYesCount, "payload data availability yes count"},
	} {
		if count.input == "" {
			return fmt.Errorf("%s missing", count.name)
		}
		if *count.dst, err = strconv.ParseUint(count.input, 10, 64); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for %s: %s", count.name, count.input))
		}
	}

	if nodeJSON.ExtraData == nil {
		return errors.New("extra data missing")
	}
	f.ExtraData = nodeJSON.ExtraData

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

// decodeObject decodes a JSON object into its fields, leaving their values undecoded.
func decodeObject(input []byte) (map[string]json.RawMessage, error) {
	if trimmed := bytes.TrimSpace(input); len(trimmed) == 0 || trimmed[0] != '{' {
		return nil, errors.New("invalid JSON: not an object")
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(input, &fields); err != nil {
		return nil, errors.Wrap(err, "invalid JSON")
	}

	return fields, nil
}

// takeField decodes the field key, if present, into dst. Unlike encoding/json's struct decoding,
// keys match exactly rather than ignoring case: the API's keys are snake_case.
func takeField(fields map[string]json.RawMessage, key string, dst any) error {
	raw, exists := fields[key]
	if !exists {
		return nil
	}

	if err := json.Unmarshal(raw, dst); err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for %s", key))
	}

	return nil
}
