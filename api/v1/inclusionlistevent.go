// Copyright © 2025 Attestant Limited.
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
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// InclusionListEvent represents the inclusion list event.
type InclusionListEvent struct {
	// Version is the fork version of the beacon chain.
	Version spec.DataVersion
	// Data is the data of the event.
	Data *SignedInclusionList
}

// SignedInclusionList represents the data of the inclusion list event.
type SignedInclusionList struct {
	Message   *InclusionList
	Signature phase0.BLSSignature `ssz-size:"96"`
}

// InclusionList represents the inclusion list.
type InclusionList struct {
	Slot           phase0.Slot
	ValidatorIndex phase0.ValidatorIndex
	DependentRoot  phase0.Root             `ssz-size:"32"`
	Transactions   []bellatrix.Transaction `ssz-max:"1048576,1073741824" ssz-size:"?,?"`
}

// inclusionListEventJSON is the spec representation of the event.
type inclusionListEventJSON struct {
	Version spec.DataVersion         `json:"version"`
	Data    *signedInclusionListJSON `json:"data"`
}

// signedInclusionListJSON is the spec representation of the signed inclusion list.
type signedInclusionListJSON struct {
	Message   json.RawMessage `json:"message"`
	Signature string          `json:"signature"`
}

// inclusionListJSON is the spec representation of the inclusion list.
type inclusionListJSON struct {
	Slot           string   `json:"slot"`
	ValidatorIndex string   `json:"validator_index"`
	DependentRoot  string   `json:"dependent_root"`
	Transactions   []string `json:"transactions"`
}

func (e *InclusionList) UnmarshalJSON(input []byte) error {
	var message inclusionListJSON
	if err := json.Unmarshal(input, &message); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	// Parse slot.
	if message.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(message.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	e.Slot = phase0.Slot(slot)

	// Parse validator index.
	if message.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(message.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	e.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	// Parse dependent root.
	if message.DependentRoot == "" {
		return errors.New("dependent root missing")
	}
	dependentRoot, err := hex.DecodeString(strings.TrimPrefix(message.DependentRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for dependent root")
	}
	if len(dependentRoot) != phase0.RootLength {
		return fmt.Errorf("incorrect length %d for dependent root", len(dependentRoot))
	}
	copy(e.DependentRoot[:], dependentRoot)

	// Parse transactions.
	if message.Transactions == nil {
		return errors.New("transactions missing")
	}
	if len(message.Transactions) > bellatrix.MaxTransactionsPerPayload {
		return errors.Wrap(err, "incorrect length for transactions")
	}
	transactions := make([]bellatrix.Transaction, len(message.Transactions))
	for i := range message.Transactions {
		if message.Transactions[i] == "" {
			return errors.New("transaction missing")
		}
		transactions[i], err = hex.DecodeString(strings.TrimPrefix(message.Transactions[i], "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for transaction")
		}
		if len(transactions[i]) > bellatrix.MaxBytesPerTransaction {
			return errors.Wrap(err, "incorrect length for transaction")
		}
	}
	e.Transactions = transactions

	return nil
}

// MarshalJSON implements json.Marshaler.
func (e *InclusionListEvent) MarshalJSON() ([]byte, error) {
	var inclusionListMessage []byte
	var err error

	switch e.Version {
	case spec.DataVersionHeze:
		if e.Data.Message == nil {
			return nil, errors.New("inclusion list message missing")
		}
		transactions := make([]string, len(e.Data.Message.Transactions))
		for i := range e.Data.Message.Transactions {
			transactions[i] = fmt.Sprintf("%#x", e.Data.Message.Transactions[i])
		}
		inclusionListMessage, err = json.Marshal(&inclusionListJSON{
			Slot:           fmt.Sprintf("%d", e.Data.Message.Slot),
			ValidatorIndex: fmt.Sprintf("%d", e.Data.Message.ValidatorIndex),
			DependentRoot:  fmt.Sprintf("%#x", e.Data.Message.DependentRoot),
			Transactions:   transactions,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal inclusion list message")
		}
	default:
		return nil, errors.New("unsupported inclusion list version")
	}

	data := signedInclusionListJSON{
		Message:   inclusionListMessage,
		Signature: fmt.Sprintf("%#x", e.Data.Signature),
	}

	return json.Marshal(&inclusionListEventJSON{
		Version: e.Version,
		Data:    &data,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *InclusionListEvent) UnmarshalJSON(input []byte) error {
	var event inclusionListEventJSON
	if err := json.Unmarshal(input, &event); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	// Parse version.
	e.Version = event.Version

	// Parse data.
	if event.Data == nil {
		return errors.New("inclusion list data missing")
	}
	e.Data = &SignedInclusionList{}

	// Parse message of data.
	if event.Data.Message == nil {
		return errors.New("inclusion list message missing")
	}
	switch event.Version {
	case spec.DataVersionHeze:
		var message InclusionList
		if err := json.Unmarshal(event.Data.Message, &message); err != nil {
			return err
		}
		e.Data.Message = &message
	default:
		return errors.New("unsupported inclusion list version")
	}

	// Parse signature of data.
	if event.Data.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(event.Data.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != phase0.SignatureLength {
		return fmt.Errorf("incorrect length %d for signature", len(signature))
	}
	copy(e.Data.Signature[:], signature)

	return nil
}

// String returns a string version of the structure.
func (e *InclusionListEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
