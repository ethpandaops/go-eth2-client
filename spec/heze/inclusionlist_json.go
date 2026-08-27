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

package heze

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// inclusionListJSON is the spec representation of the struct.
type inclusionListJSON struct {
	Slot           string   `json:"slot"`
	ValidatorIndex string   `json:"validator_index"`
	DependentRoot  string   `json:"dependent_root"`
	Transactions   []string `json:"transactions"`
}

// MarshalJSON implements json.Marshaler.
func (i *InclusionList) MarshalJSON() ([]byte, error) {
	transactions := make([]string, len(i.Transactions))
	for idx := range i.Transactions {
		// %#x renders an empty slice as "", while the spec writes "0x".
		transactions[idx] = "0x" + hex.EncodeToString(i.Transactions[idx])
	}

	return json.Marshal(&inclusionListJSON{
		Slot:           fmt.Sprintf("%d", i.Slot),
		ValidatorIndex: fmt.Sprintf("%d", i.ValidatorIndex),
		DependentRoot:  fmt.Sprintf("%#x", i.DependentRoot),
		Transactions:   transactions,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *InclusionList) UnmarshalJSON(input []byte) error {
	var data inclusionListJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	// Slot
	if data.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid slot")
	}
	i.Slot = phase0.Slot(slot)

	// Validator index
	if data.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(data.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid validator index")
	}
	i.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	// Dependent root
	if data.DependentRoot == "" {
		return errors.New("dependent root missing")
	}
	dependentRoot, err := hex.DecodeString(strings.TrimPrefix(data.DependentRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid dependent root")
	}
	copy(i.DependentRoot[:], dependentRoot)

	// Transactions
	if data.Transactions == nil {
		data.Transactions = []string{}
	}
	i.Transactions = make([]bellatrix.Transaction, len(data.Transactions))
	for idx, transaction := range data.Transactions {
		transactionBytes, err := hex.DecodeString(strings.TrimPrefix(transaction, "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid transaction %d", idx))
		}
		i.Transactions[idx] = transactionBytes
	}

	return nil
}
