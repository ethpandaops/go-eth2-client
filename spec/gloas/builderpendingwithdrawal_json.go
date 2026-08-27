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

package gloas

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

// builderPendingWithdrawalJSON is the spec representation of the struct.
type builderPendingWithdrawalJSON struct {
	FeeRecipient string `json:"fee_recipient"`
	Amount       string `json:"amount"`
	BuilderIndex string `json:"builder_index"`
}

// builderPendingWithdrawalYAML is the spec representation of the struct.
type builderPendingWithdrawalYAML struct {
	FeeRecipient string `yaml:"fee_recipient"`
	Amount       uint64 `yaml:"amount"`
	BuilderIndex uint64 `yaml:"builder_index"`
}

// MarshalJSON implements json.Marshaler.
func (p *BuilderPendingWithdrawal) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderPendingWithdrawalJSON{
		FeeRecipient: p.FeeRecipient.String(),
		Amount:       fmt.Sprintf("%d", p.Amount),
		BuilderIndex: fmt.Sprintf("%d", p.BuilderIndex),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *BuilderPendingWithdrawal) UnmarshalJSON(input []byte) error {
	var data builderPendingWithdrawalJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.FeeRecipient == "" {
		return errors.New("fee recipient missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.FeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid fee recipient")
	}
	if len(feeRecipient) != bellatrix.ExecutionAddressLength {
		return fmt.Errorf("incorrect length %d for fee recipient", len(feeRecipient))
	}
	copy(p.FeeRecipient[:], feeRecipient)

	if data.Amount == "" {
		return errors.New("amount missing")
	}
	amount, err := strconv.ParseUint(data.Amount, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid amount")
	}
	p.Amount = phase0.Gwei(amount)

	if data.BuilderIndex == "" {
		return errors.New("builder index missing")
	}
	builderIndex, err := strconv.ParseUint(data.BuilderIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid builder index")
	}
	p.BuilderIndex = BuilderIndex(builderIndex)

	return nil
}
