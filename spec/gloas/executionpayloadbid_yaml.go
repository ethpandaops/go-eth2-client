// Copyright © 2023 Attestant Limited.
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
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// executionPayloadBidYAML is the spec representation of the struct.
type executionPayloadBidYAML struct {
	ParentBlockHash       string   `yaml:"parent_block_hash"`
	ParentBlockRoot       string   `yaml:"parent_block_root"`
	BlockHash             string   `yaml:"block_hash"`
	PrevRandao            string   `yaml:"prev_randao"`
	FeeRecipient          string   `yaml:"fee_recipient"`
	GasLimit              uint64   `yaml:"gas_limit"`
	BuilderIndex          uint64   `yaml:"builder_index"`
	Slot                  uint64   `yaml:"slot"`
	Value                 uint64   `yaml:"value"`
	ExecutionPayment      uint64   `yaml:"execution_payment"`
	BlobKZGCommitments    []string `yaml:"blob_kzg_commitments"`
	ExecutionRequestsRoot string   `yaml:"execution_requests_root"`
}

// MarshalYAML implements yaml.Marshaler.
func (e *ExecutionPayloadBid) MarshalYAML() ([]byte, error) {
	blobKZGCommitments := make([]string, len(e.BlobKZGCommitments))
	for i := range e.BlobKZGCommitments {
		blobKZGCommitments[i] = fmt.Sprintf("%#x", e.BlobKZGCommitments[i])
	}

	yamlBytes, err := yaml.MarshalWithOptions(&executionPayloadBidYAML{
		ParentBlockHash:       fmt.Sprintf("%#x", e.ParentBlockHash),
		ParentBlockRoot:       fmt.Sprintf("%#x", e.ParentBlockRoot),
		BlockHash:             fmt.Sprintf("%#x", e.BlockHash),
		PrevRandao:            fmt.Sprintf("%#x", e.PrevRandao),
		FeeRecipient:          fmt.Sprintf("%#x", e.FeeRecipient),
		GasLimit:              uint64(e.GasLimit),
		BuilderIndex:          uint64(e.BuilderIndex),
		Slot:                  uint64(e.Slot),
		Value:                 uint64(e.Value),
		ExecutionPayment:      uint64(e.ExecutionPayment),
		BlobKZGCommitments:    blobKZGCommitments,
		ExecutionRequestsRoot: fmt.Sprintf("%#x", e.ExecutionRequestsRoot),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *ExecutionPayloadBid) UnmarshalYAML(input []byte) error {
	var data executionPayloadBidJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	marshaled, err := json.Marshal(&data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return e.UnmarshalJSON(marshaled)
}
