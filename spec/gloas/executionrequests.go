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
	"fmt"

	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/goccy/go-yaml"
)

// ExecutionRequests represents the execution layer triggered requests of a payload (EIP-8282).
type ExecutionRequests struct {
	Deposits        []*electra.DepositRequest       `dynssz-max:"MAX_DEPOSIT_REQUESTS_PER_PAYLOAD"         ssz-max:"8192"`
	Withdrawals     []*electra.WithdrawalRequest    `dynssz-max:"MAX_WITHDRAWAL_REQUESTS_PER_PAYLOAD"      ssz-max:"16"`
	Consolidations  []*electra.ConsolidationRequest `dynssz-max:"MAX_CONSOLIDATION_REQUESTS_PER_PAYLOAD"   ssz-max:"2"`
	BuilderDeposits []*BuilderDepositRequest        `dynssz-max:"MAX_BUILDER_DEPOSIT_REQUESTS_PER_PAYLOAD" ssz-max:"256"`
	BuilderExits    []*BuilderExitRequest           `dynssz-max:"MAX_BUILDER_EXIT_REQUESTS_PER_PAYLOAD"    ssz-max:"16"`
}

// String returns a string version of the structure.
func (e *ExecutionRequests) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
