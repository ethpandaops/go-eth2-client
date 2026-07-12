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

// HeadEventV2 is the data for the head_v2 event (Gloas / EIP-7732).
//
// Relative to HeadEvent it drops previous_duty_dependent_root, renames
// current_duty_dependent_root to current_epoch_dependent_root, and adds
// next_epoch_dependent_root plus payload_status. The event can be emitted
// twice for the same block when payload_status transitions from empty to
// full.
type HeadEventV2 struct {
	Slot                      phase0.Slot
	Block                     phase0.Root
	State                     phase0.Root
	PayloadStatus             string
	EpochTransition           bool
	CurrentEpochDependentRoot phase0.Root
	NextEpochDependentRoot    phase0.Root
	ExecutionOptimistic       bool
}

// headEventV2JSON is the spec representation of the struct.
type headEventV2JSON struct {
	Slot                      string `json:"slot"`
	Block                     string `json:"block"`
	State                     string `json:"state"`
	PayloadStatus             string `json:"payload_status,omitempty"`
	EpochTransition           bool   `json:"epoch_transition"`
	CurrentEpochDependentRoot string `json:"current_epoch_dependent_root,omitempty"`
	NextEpochDependentRoot    string `json:"next_epoch_dependent_root,omitempty"`
	ExecutionOptimistic       bool   `json:"execution_optimistic"`
}

// MarshalJSON implements json.Marshaler.
func (e *HeadEventV2) MarshalJSON() ([]byte, error) {
	data := &headEventV2JSON{
		Slot:                fmt.Sprintf("%d", e.Slot),
		Block:               fmt.Sprintf("%#x", e.Block),
		State:               fmt.Sprintf("%#x", e.State),
		PayloadStatus:       e.PayloadStatus,
		EpochTransition:     e.EpochTransition,
		ExecutionOptimistic: e.ExecutionOptimistic,
	}

	var zeroRoot phase0.Root
	if !bytes.Equal(zeroRoot[:], e.CurrentEpochDependentRoot[:]) {
		data.CurrentEpochDependentRoot = fmt.Sprintf("%#x", e.CurrentEpochDependentRoot)
	}

	if !bytes.Equal(zeroRoot[:], e.NextEpochDependentRoot[:]) {
		data.NextEpochDependentRoot = fmt.Sprintf("%#x", e.NextEpochDependentRoot)
	}

	return json.Marshal(data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *HeadEventV2) UnmarshalJSON(input []byte) error {
	var headEventV2JSON headEventV2JSON
	if err := json.Unmarshal(input, &headEventV2JSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if headEventV2JSON.Slot == "" {
		return errors.New("slot missing")
	}

	slot, err := strconv.ParseUint(headEventV2JSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}

	e.Slot = phase0.Slot(slot)

	if headEventV2JSON.Block == "" {
		return errors.New("block missing")
	}

	if err := decodeFixedBytes(e.Block[:], headEventV2JSON.Block, rootLength, "block"); err != nil {
		return err
	}

	if headEventV2JSON.State == "" {
		return errors.New("state missing")
	}

	if err := decodeFixedBytes(e.State[:], headEventV2JSON.State, rootLength, "state"); err != nil {
		return err
	}

	e.PayloadStatus = headEventV2JSON.PayloadStatus
	e.EpochTransition = headEventV2JSON.EpochTransition
	e.ExecutionOptimistic = headEventV2JSON.ExecutionOptimistic

	// Dependent roots only have partial client coverage so do not complain if not present.
	if headEventV2JSON.CurrentEpochDependentRoot != "" {
		if err := decodeFixedBytes(e.CurrentEpochDependentRoot[:], headEventV2JSON.CurrentEpochDependentRoot, rootLength, "current epoch dependent root"); err != nil {
			return err
		}
	}

	if headEventV2JSON.NextEpochDependentRoot != "" {
		if err := decodeFixedBytes(e.NextEpochDependentRoot[:], headEventV2JSON.NextEpochDependentRoot, rootLength, "next epoch dependent root"); err != nil {
			return err
		}
	}

	return nil
}

// String returns a string version of the structure.
func (e *HeadEventV2) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
