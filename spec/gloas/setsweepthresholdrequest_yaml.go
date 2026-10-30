// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gloas

import (
	"bytes"
	"encoding/json"

	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// setSweepThresholdRequestYAML is the spec representation of the struct.
type setSweepThresholdRequestYAML struct {
	SourceAddress   bellatrix.ExecutionAddress `yaml:"source_address"`
	ValidatorPubkey phase0.BLSPubKey           `yaml:"validator_pubkey"`
	Threshold       phase0.Gwei                `yaml:"threshold"`
}

// MarshalYAML implements yaml.Marshaler.
func (e *SetSweepThresholdRequest) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&setSweepThresholdRequestYAML{
		SourceAddress:   e.SourceAddress,
		ValidatorPubkey: e.ValidatorPubkey,
		Threshold:       e.Threshold,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *SetSweepThresholdRequest) UnmarshalYAML(input []byte) error {
	// This is very inefficient, but YAML is only used for spec tests so we do this
	// rather than maintain a custom YAML unmarshaller.
	var unmarshaled setSweepThresholdRequestJSON
	if err := yaml.Unmarshal(input, &unmarshaled); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}

	marshaled, err := json.Marshal(&unmarshaled)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return e.UnmarshalJSON(marshaled)
}
