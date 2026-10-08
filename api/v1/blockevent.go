// Copyright © 2020, 2021 Attestant Limited.
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

// BlockEvent is the data for the block event.
type BlockEvent struct {
	Slot                phase0.Slot
	Block               phase0.Root
	ExecutionOptimistic bool
	// BuilderIndex and BlockHash come from the block's
	// signed_execution_payload_bid.message. They are only sent from Gloas
	// onwards, and are nil when the beacon node does not send them. A pointer
	// is used because 0 is a valid builder index.
	BuilderIndex *uint64
	BlockHash    *phase0.Hash32
}

// blockEventJSON is the spec representation of the struct.
type blockEventJSON struct {
	Slot                string `json:"slot"`
	Block               string `json:"block"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
	BuilderIndex        string `json:"builder_index,omitempty"`
	BlockHash           string `json:"block_hash,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (e *BlockEvent) MarshalJSON() ([]byte, error) {
	data := &blockEventJSON{
		Slot:                fmt.Sprintf("%d", e.Slot),
		Block:               fmt.Sprintf("%#x", e.Block),
		ExecutionOptimistic: e.ExecutionOptimistic,
	}
	if e.BuilderIndex != nil {
		data.BuilderIndex = strconv.FormatUint(*e.BuilderIndex, 10)
	}

	if e.BlockHash != nil {
		data.BlockHash = fmt.Sprintf("%#x", *e.BlockHash)
	}

	return json.Marshal(data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *BlockEvent) UnmarshalJSON(input []byte) error {
	var err error

	var blockEventJSON blockEventJSON
	if err = json.Unmarshal(input, &blockEventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if blockEventJSON.Slot == "" {
		return errors.New("slot missing")
	}

	slot, err := strconv.ParseUint(blockEventJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}

	e.Slot = phase0.Slot(slot)

	if blockEventJSON.Block == "" {
		return errors.New("block missing")
	}

	block, err := hex.DecodeString(strings.TrimPrefix(blockEventJSON.Block, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block")
	}

	if len(block) != rootLength {
		return fmt.Errorf("incorrect length %d for block", len(block))
	}

	copy(e.Block[:], block)
	e.ExecutionOptimistic = blockEventJSON.ExecutionOptimistic

	e.BuilderIndex = nil
	if blockEventJSON.BuilderIndex != "" {
		builderIndex, err := strconv.ParseUint(blockEventJSON.BuilderIndex, 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for builder index")
		}

		e.BuilderIndex = &builderIndex
	}

	e.BlockHash = nil
	if blockEventJSON.BlockHash != "" {
		var blockHash phase0.Hash32
		if err := decodeFixedBytes(blockHash[:], blockEventJSON.BlockHash, phase0.Hash32Length, "block hash"); err != nil {
			return err
		}

		e.BlockHash = &blockHash
	}

	return nil
}

// String returns a string version of the structure.
func (e *BlockEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
