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

	"github.com/ethpandaops/go-eth2-client/codecs"
	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// setSweepThresholdRequestJSON is the spec representation of the struct.
type setSweepThresholdRequestJSON struct {
	SourceAddress   bellatrix.ExecutionAddress `json:"source_address"`
	ValidatorPubkey phase0.BLSPubKey           `json:"validator_pubkey"`
	Threshold       phase0.Gwei                `json:"threshold"`
}

// MarshalJSON implements json.Marshaler.
func (e *SetSweepThresholdRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&setSweepThresholdRequestJSON{
		SourceAddress:   e.SourceAddress,
		ValidatorPubkey: e.ValidatorPubkey,
		Threshold:       e.Threshold,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *SetSweepThresholdRequest) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&setSweepThresholdRequestJSON{}, input)
	if err != nil {
		return err
	}

	if err := e.SourceAddress.UnmarshalJSON(raw["source_address"]); err != nil {
		return errors.Wrap(err, "source_address")
	}

	if err := e.ValidatorPubkey.UnmarshalJSON(raw["validator_pubkey"]); err != nil {
		return errors.Wrap(err, "validator_pubkey")
	}

	if err := e.Threshold.UnmarshalJSON(raw["threshold"]); err != nil {
		return errors.Wrap(err, "threshold")
	}

	return nil
}
