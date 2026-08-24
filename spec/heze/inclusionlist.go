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
	"fmt"

	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// InclusionList represents the transactions an inclusion list committee member
// requires the next payload to contain.
// New in Heze (EIP-7805).
type InclusionList struct {
	Slot           phase0.Slot
	ValidatorIndex phase0.ValidatorIndex
	DependentRoot  phase0.Root             `ssz-size:"32"`
	Transactions   []bellatrix.Transaction `ssz-type:"progressive-list,progressive-list"`
}

// String returns a string version of the structure.
func (i *InclusionList) String() string {
	data, err := yaml.Marshal(i)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
