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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// inclusionListYAML is the spec representation of the struct.
type inclusionListYAML struct {
	Slot           uint64   `yaml:"slot"`
	ValidatorIndex uint64   `yaml:"validator_index"`
	DependentRoot  string   `yaml:"dependent_root"`
	Transactions   []string `yaml:"transactions"`
}

// MarshalYAML implements yaml.Marshaler.
func (i *InclusionList) MarshalYAML() ([]byte, error) {
	transactions := make([]string, len(i.Transactions))
	for idx := range i.Transactions {
		// %#x renders an empty slice as "", while the spec writes "0x".
		transactions[idx] = "0x" + hex.EncodeToString(i.Transactions[idx])
	}

	yamlBytes, err := yaml.MarshalWithOptions(&inclusionListYAML{
		Slot:           uint64(i.Slot),
		ValidatorIndex: uint64(i.ValidatorIndex),
		DependentRoot:  fmt.Sprintf("%#x", i.DependentRoot),
		Transactions:   transactions,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (i *InclusionList) UnmarshalYAML(input []byte) error {
	var data inclusionListJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return i.UnmarshalJSON(marshaled)
}
