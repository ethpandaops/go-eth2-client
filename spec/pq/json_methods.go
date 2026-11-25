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

package pq

import (
	"encoding/json"
	"fmt"
)

func stringify(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

func (v *Validator) MarshalJSON() ([]byte, error) {
	return marshalStruct(v)
}

func (v *Validator) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, v)
}

func (b *Builder) MarshalJSON() ([]byte, error) {
	return marshalStruct(b)
}

func (b *Builder) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, b)
}

func (v *ValidatorPreregistrationRequest) MarshalJSON() ([]byte, error) {
	return marshalStruct(v)
}

func (v *ValidatorPreregistrationRequest) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, v)
}

func (v *ValidatorPreregistration) MarshalJSON() ([]byte, error) {
	return marshalStruct(v)
}

func (v *ValidatorPreregistration) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, v)
}

func (b *BuilderPreregistrationRequest) MarshalJSON() ([]byte, error) {
	return marshalStruct(b)
}

func (b *BuilderPreregistrationRequest) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, b)
}

func (b *BuilderPreregistration) MarshalJSON() ([]byte, error) {
	return marshalStruct(b)
}

func (b *BuilderPreregistration) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, b)
}

func (v *VoluntaryExitRequest) MarshalJSON() ([]byte, error) {
	return marshalStruct(v)
}

func (v *VoluntaryExitRequest) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, v)
}

func (e *ExecutionRequests) MarshalJSON() ([]byte, error) {
	return marshalStruct(e)
}

func (e *ExecutionRequests) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, e)
}

func (s *SignedBeaconBlockHeader) MarshalJSON() ([]byte, error) {
	return marshalStruct(s)
}

func (s *SignedBeaconBlockHeader) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, s)
}

func (p *ProposerSlashing) MarshalJSON() ([]byte, error) {
	return marshalStruct(p)
}

func (p *ProposerSlashing) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, p)
}

func (a *Attestation) MarshalJSON() ([]byte, error) {
	return marshalStruct(a)
}

func (a *Attestation) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, a)
}

func (i *IndexedAttestation) MarshalJSON() ([]byte, error) {
	return marshalStruct(i)
}

func (i *IndexedAttestation) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, i)
}

func (a *AttesterSlashing) MarshalJSON() ([]byte, error) {
	return marshalStruct(a)
}

func (a *AttesterSlashing) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, a)
}

func (s *SingleAttestation) MarshalJSON() ([]byte, error) {
	return marshalStruct(s)
}

func (s *SingleAttestation) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, s)
}

func (a *AggregateAndProof) MarshalJSON() ([]byte, error) {
	return marshalStruct(a)
}

func (a *AggregateAndProof) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, a)
}

func (s *SignedAggregateAndProof) MarshalJSON() ([]byte, error) {
	return marshalStruct(s)
}

func (s *SignedAggregateAndProof) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, s)
}

func (p *PayloadAttestationMessage) MarshalJSON() ([]byte, error) {
	return marshalStruct(p)
}

func (p *PayloadAttestationMessage) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, p)
}

func (s *SignedExecutionPayloadBid) MarshalJSON() ([]byte, error) {
	return marshalStruct(s)
}

func (s *SignedExecutionPayloadBid) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, s)
}

func (e *ExecutionPayloadEnvelope) MarshalJSON() ([]byte, error) {
	return marshalStruct(e)
}

func (e *ExecutionPayloadEnvelope) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, e)
}

func (s *SignedExecutionPayloadEnvelope) MarshalJSON() ([]byte, error) {
	return marshalStruct(s)
}

func (s *SignedExecutionPayloadEnvelope) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, s)
}

func (b *BeaconBlockBody) MarshalJSON() ([]byte, error) {
	return marshalStruct(b)
}

func (b *BeaconBlockBody) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, b)
}

func (b *BeaconBlock) MarshalJSON() ([]byte, error) {
	return marshalStruct(b)
}

func (b *BeaconBlock) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, b)
}

func (s *SignedBeaconBlock) MarshalJSON() ([]byte, error) {
	return marshalStruct(s)
}

func (s *SignedBeaconBlock) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, s)
}

func (b *BeaconState) MarshalJSON() ([]byte, error) {
	return marshalStruct(b)
}

func (b *BeaconState) UnmarshalJSON(input []byte) error {
	return unmarshalStruct(input, b)
}

func (v *Validator) String() string {
	return stringify(v)
}

func (b *Builder) String() string {
	return stringify(b)
}

func (v *ValidatorPreregistrationRequest) String() string {
	return stringify(v)
}

func (v *ValidatorPreregistration) String() string {
	return stringify(v)
}

func (b *BuilderPreregistrationRequest) String() string {
	return stringify(b)
}

func (b *BuilderPreregistration) String() string {
	return stringify(b)
}

func (v *VoluntaryExitRequest) String() string {
	return stringify(v)
}

func (e *ExecutionRequests) String() string {
	return stringify(e)
}

func (s *SignedBeaconBlockHeader) String() string {
	return stringify(s)
}

func (p *ProposerSlashing) String() string {
	return stringify(p)
}

func (a *Attestation) String() string {
	return stringify(a)
}

func (i *IndexedAttestation) String() string {
	return stringify(i)
}

func (a *AttesterSlashing) String() string {
	return stringify(a)
}

func (s *SingleAttestation) String() string {
	return stringify(s)
}

func (a *AggregateAndProof) String() string {
	return stringify(a)
}

func (s *SignedAggregateAndProof) String() string {
	return stringify(s)
}

func (p *PayloadAttestationMessage) String() string {
	return stringify(p)
}

func (s *SignedExecutionPayloadBid) String() string {
	return stringify(s)
}

func (e *ExecutionPayloadEnvelope) String() string {
	return stringify(e)
}

func (s *SignedExecutionPayloadEnvelope) String() string {
	return stringify(s)
}

func (b *BeaconBlockBody) String() string {
	return stringify(b)
}

func (b *BeaconBlock) String() string {
	return stringify(b)
}

func (s *SignedBeaconBlock) String() string {
	return stringify(s)
}

func (b *BeaconState) String() string {
	return stringify(b)
}
