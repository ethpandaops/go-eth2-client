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

package eip8450

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// XMSSPublicKey is a leanXMSS public key.
type XMSSPublicKey [32]byte

// XMSSSignature is a leanXMSS signature.
type XMSSSignature [2404]byte

// PQKeyMaterial holds the post-quantum keys of one validator or builder.
type PQKeyMaterial struct {
	XMSSPublicKey XMSSPublicKey `ssz-size:"32"`
	VRFRoot       phase0.Root   `ssz-size:"32"`
}

// PQKeyRegistration is the message of a post-quantum key registration.
type PQKeyRegistration struct {
	Index               phase0.ValidatorIndex
	KeyMaterial         *PQKeyMaterial
	PreviousChangeEpoch phase0.Epoch
}

// SignedPQKeyRegistration is a post-quantum key registration as it appears in a beacon block.
type SignedPQKeyRegistration struct {
	Message    *PQKeyRegistration
	Possession XMSSSignature       `ssz-size:"2404"`
	Signature  phase0.BLSSignature `ssz-size:"96"`
}

type pqKeyMaterialJSON struct {
	XMSSPublicKey string `json:"xmss_public_key"`
	VRFRoot       string `json:"vrf_root"`
}

type pqKeyRegistrationJSON struct {
	Index               string         `json:"index"`
	KeyMaterial         *PQKeyMaterial `json:"key_material"`
	PreviousChangeEpoch string         `json:"previous_change_epoch"`
}

type signedPQKeyRegistrationJSON struct {
	Message    *PQKeyRegistration `json:"message"`
	Possession string             `json:"possession"`
	Signature  string             `json:"signature"`
}

func decodeFixedHex(name, value string, target []byte) error {
	if value == "" {
		return fmt.Errorf("%s missing", name)
	}

	decoded, err := hex.DecodeString(strings.TrimPrefix(value, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid "+name)
	}

	if len(decoded) != len(target) {
		return fmt.Errorf("incorrect length for %s", name)
	}

	copy(target, decoded)

	return nil
}

func marshalYAML(data any) ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(data, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

func yamlToJSON(input []byte, data any) ([]byte, error) {
	if err := yaml.Unmarshal(input, data); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal YAML")
	}

	marshaled, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal JSON")
	}

	return marshaled, nil
}

// MarshalJSON implements json.Marshaler.
func (p *PQKeyMaterial) MarshalJSON() ([]byte, error) {
	return json.Marshal(&pqKeyMaterialJSON{
		XMSSPublicKey: fmt.Sprintf("%#x", p.XMSSPublicKey),
		VRFRoot:       fmt.Sprintf("%#x", p.VRFRoot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PQKeyMaterial) UnmarshalJSON(input []byte) error {
	var data pqKeyMaterialJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if err := decodeFixedHex("xmss public key", data.XMSSPublicKey, p.XMSSPublicKey[:]); err != nil {
		return err
	}

	return decodeFixedHex("vrf root", data.VRFRoot, p.VRFRoot[:])
}

// MarshalYAML implements yaml.Marshaler.
func (p *PQKeyMaterial) MarshalYAML() ([]byte, error) {
	return marshalYAML(&pqKeyMaterialJSON{
		XMSSPublicKey: fmt.Sprintf("%#x", p.XMSSPublicKey),
		VRFRoot:       fmt.Sprintf("%#x", p.VRFRoot),
	})
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (p *PQKeyMaterial) UnmarshalYAML(input []byte) error {
	marshaled, err := yamlToJSON(input, &pqKeyMaterialJSON{})
	if err != nil {
		return err
	}

	return p.UnmarshalJSON(marshaled)
}

// MarshalJSON implements json.Marshaler.
func (p *PQKeyRegistration) MarshalJSON() ([]byte, error) {
	return json.Marshal(&pqKeyRegistrationJSON{
		Index:               strconv.FormatUint(uint64(p.Index), 10),
		KeyMaterial:         p.KeyMaterial,
		PreviousChangeEpoch: strconv.FormatUint(uint64(p.PreviousChangeEpoch), 10),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PQKeyRegistration) UnmarshalJSON(input []byte) error {
	var data pqKeyRegistrationJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	index, err := strconv.ParseUint(data.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid index")
	}

	if data.KeyMaterial == nil {
		return errors.New("key material missing")
	}

	previousChangeEpoch, err := strconv.ParseUint(data.PreviousChangeEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid previous change epoch")
	}

	p.Index = phase0.ValidatorIndex(index)
	p.KeyMaterial = data.KeyMaterial
	p.PreviousChangeEpoch = phase0.Epoch(previousChangeEpoch)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (p *PQKeyRegistration) MarshalYAML() ([]byte, error) {
	return marshalYAML(&pqKeyRegistrationJSON{
		Index:               strconv.FormatUint(uint64(p.Index), 10),
		KeyMaterial:         p.KeyMaterial,
		PreviousChangeEpoch: strconv.FormatUint(uint64(p.PreviousChangeEpoch), 10),
	})
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (p *PQKeyRegistration) UnmarshalYAML(input []byte) error {
	marshaled, err := yamlToJSON(input, &pqKeyRegistrationJSON{})
	if err != nil {
		return err
	}

	return p.UnmarshalJSON(marshaled)
}

// MarshalJSON implements json.Marshaler.
func (s *SignedPQKeyRegistration) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedPQKeyRegistrationJSON{
		Message:    s.Message,
		Possession: fmt.Sprintf("%#x", s.Possession),
		Signature:  fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedPQKeyRegistration) UnmarshalJSON(input []byte) error {
	var data signedPQKeyRegistrationJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.Message == nil {
		return errors.New("message missing")
	}

	s.Message = data.Message

	if err := decodeFixedHex("possession", data.Possession, s.Possession[:]); err != nil {
		return err
	}

	return decodeFixedHex("signature", data.Signature, s.Signature[:])
}

// MarshalYAML implements yaml.Marshaler.
func (s *SignedPQKeyRegistration) MarshalYAML() ([]byte, error) {
	return marshalYAML(&signedPQKeyRegistrationJSON{
		Message:    s.Message,
		Possession: fmt.Sprintf("%#x", s.Possession),
		Signature:  fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SignedPQKeyRegistration) UnmarshalYAML(input []byte) error {
	marshaled, err := yamlToJSON(input, &signedPQKeyRegistrationJSON{})
	if err != nil {
		return err
	}

	return s.UnmarshalJSON(marshaled)
}
