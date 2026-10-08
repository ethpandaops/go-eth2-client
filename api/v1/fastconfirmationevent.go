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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// FastConfirmationEvent is the data for the fast confirmation event.
type FastConfirmationEvent struct {
	Slot  phase0.Slot
	Block phase0.Root
	// CurrentSlot is the wall-clock slot at which the algorithm was executed.
	// It is nil when the beacon node does not send it.
	CurrentSlot *phase0.Slot
}

// fastConfirmationEventJSON is the spec representation of the struct.
type fastConfirmationEventJSON struct {
	Block       string `json:"block"`
	Slot        string `json:"slot"`
	CurrentSlot string `json:"current_slot,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (e *FastConfirmationEvent) MarshalJSON() ([]byte, error) {
	data := &fastConfirmationEventJSON{
		Block: fmt.Sprintf("%#x", e.Block),
		Slot:  fmt.Sprintf("%d", e.Slot),
	}
	if e.CurrentSlot != nil {
		data.CurrentSlot = fmt.Sprintf("%d", *e.CurrentSlot)
	}

	return json.Marshal(data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *FastConfirmationEvent) UnmarshalJSON(input []byte) error {
	var data fastConfirmationEventJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.Slot == "" {
		return errors.New("slot missing")
	}

	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}

	e.Slot = phase0.Slot(slot)

	if data.Block == "" {
		return errors.New("block missing")
	}

	block, err := hex.DecodeString(strings.TrimPrefix(data.Block, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block")
	}

	if len(block) != rootLength {
		return fmt.Errorf("incorrect length %d for block", len(block))
	}

	copy(e.Block[:], block)

	e.CurrentSlot = nil
	if data.CurrentSlot != "" {
		currentSlot, err := strconv.ParseUint(data.CurrentSlot, 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for current slot")
		}

		slot := phase0.Slot(currentSlot)
		e.CurrentSlot = &slot
	}

	return nil
}

// String returns a string version of the structure.
func (e *FastConfirmationEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
