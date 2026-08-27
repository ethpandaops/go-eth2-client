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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// builderPendingPaymentJSON is the spec representation of the struct.
type builderPendingPaymentJSON struct {
	Weight        string                    `json:"weight"`
	Withdrawal    *BuilderPendingWithdrawal `json:"withdrawal"`
	ProposerIndex string                    `json:"proposer_index"`
}

// builderPendingPaymentYAML is the spec representation of the struct.
type builderPendingPaymentYAML struct {
	Weight        uint64                    `yaml:"weight"`
	Withdrawal    *BuilderPendingWithdrawal `yaml:"withdrawal"`
	ProposerIndex uint64                    `yaml:"proposer_index"`
}

// MarshalJSON implements json.Marshaler.
func (p *BuilderPendingPayment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderPendingPaymentJSON{
		Weight:        fmt.Sprintf("%d", p.Weight),
		Withdrawal:    p.Withdrawal,
		ProposerIndex: fmt.Sprintf("%d", p.ProposerIndex),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *BuilderPendingPayment) UnmarshalJSON(input []byte) error {
	var data builderPendingPaymentJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.Weight == "" {
		return errors.New("weight missing")
	}
	weight, err := strconv.ParseUint(data.Weight, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid weight")
	}
	p.Weight = phase0.Gwei(weight)

	if data.Withdrawal == nil {
		return errors.New("withdrawal missing")
	}
	p.Withdrawal = data.Withdrawal

	if data.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	proposerIndex, err := strconv.ParseUint(data.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid proposer index")
	}
	p.ProposerIndex = phase0.ValidatorIndex(proposerIndex)

	return nil
}
