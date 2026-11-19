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
	"encoding/binary"

	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
)

const (
	XMSSPubKeySize    = 32
	XMSSSignatureSize = 2404
	VRFProofLength    = 32
	DepositPrefixSize = 16
)

type XMSSPubKey [XMSSPubKeySize]byte

type XMSSSignature [XMSSSignatureSize]byte

type VRFProof [VRFProofLength]phase0.Root

func (s *XMSSSignature) LeafIndex() uint32 {
	return binary.LittleEndian.Uint32(s[:4])
}

func (s *XMSSSignature) IsZero() bool {
	return *s == XMSSSignature{}
}

func (p XMSSPubKey) DepositPubKey() phase0.BLSPubKey {
	var pubkey phase0.BLSPubKey

	copy(pubkey[DepositPrefixSize:], p[:])

	return pubkey
}

func (v *Validator) Legacy() *phase0.Validator {
	return &phase0.Validator{
		PublicKey:                  v.PublicKey.DepositPubKey(),
		WithdrawalCredentials:      v.WithdrawalCredentials,
		EffectiveBalance:           v.EffectiveBalance,
		Slashed:                    v.Slashed,
		ActivationEligibilityEpoch: v.ActivationEligibilityEpoch,
		ActivationEpoch:            v.ActivationEpoch,
		ExitEpoch:                  v.ExitEpoch,
		WithdrawableEpoch:          v.WithdrawableEpoch,
	}
}

func (b *Builder) Legacy() *gloas.Builder {
	return &gloas.Builder{
		PublicKey:         b.PublicKey.DepositPubKey(),
		Version:           b.Version,
		ExecutionAddress:  b.ExecutionAddress,
		Balance:           b.Balance,
		DepositEpoch:      b.DepositEpoch,
		WithdrawableEpoch: b.WithdrawableEpoch,
	}
}

func (p *ProposerSlashing) Legacy() *phase0.ProposerSlashing {
	slashing := &phase0.ProposerSlashing{}

	if p.SignedHeader1 != nil {
		slashing.SignedHeader1 = &phase0.SignedBeaconBlockHeader{Message: p.SignedHeader1.Message}
	}

	if p.SignedHeader2 != nil {
		slashing.SignedHeader2 = &phase0.SignedBeaconBlockHeader{Message: p.SignedHeader2.Message}
	}

	return slashing
}

func LegacyValidators(validators []*Validator) []*phase0.Validator {
	legacy := make([]*phase0.Validator, len(validators))
	for i, validator := range validators {
		legacy[i] = validator.Legacy()
	}

	return legacy
}

func LegacyBuilders(builders []*Builder) []*gloas.Builder {
	legacy := make([]*gloas.Builder, len(builders))
	for i, builder := range builders {
		legacy[i] = builder.Legacy()
	}

	return legacy
}

func LegacyProposerSlashings(slashings []*ProposerSlashing) []*phase0.ProposerSlashing {
	legacy := make([]*phase0.ProposerSlashing, len(slashings))
	for i, slashing := range slashings {
		legacy[i] = slashing.Legacy()
	}

	return legacy
}
