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

package all

import (
	"errors"
	"fmt"

	"github.com/ethpandaops/go-eth2-client/api"
	apiv1bellatrix "github.com/ethpandaops/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/ethpandaops/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/ethpandaops/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/ethpandaops/go-eth2-client/api/v1/electra"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/pk910/dynamic-ssz/sszutils"
)

// BlindedBeaconBlock is a fork-agnostic blinded beacon block. The Body's
// contents depend on Version. Blinded blocks exist as a wire type from
// Bellatrix up to (excluding) Gloas; other forks are not supported by this
// type.
type BlindedBeaconBlock struct {
	Version       version.DataVersion
	Slot          phase0.Slot
	ProposerIndex phase0.ValidatorIndex
	ParentRoot    phase0.Root
	StateRoot     phase0.Root
	Body          *BlindedBeaconBlockBody
}

// viewType returns the fork-specific schema type pointer used as the view
// descriptor for the active Version.
func (b *BlindedBeaconBlock) viewType() (any, error) {
	switch b.Version {
	case version.DataVersionBellatrix:
		return (*apiv1bellatrix.BlindedBeaconBlock)(nil), nil
	case version.DataVersionCapella:
		return (*apiv1capella.BlindedBeaconBlock)(nil), nil
	case version.DataVersionDeneb:
		return (*apiv1deneb.BlindedBeaconBlock)(nil), nil
	case version.DataVersionElectra,
		version.DataVersionFulu:
		// Fulu reuses the Electra blinded block schema unchanged.
		return (*apiv1electra.BlindedBeaconBlock)(nil), nil
	default:
		return nil, fmt.Errorf("BlindedBeaconBlock: unsupported version %d", b.Version)
	}
}

// MarshalSSZDyn marshals the blinded block using the view that matches
// Version.
func (b *BlindedBeaconBlock) MarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) ([]byte, error) {
	view, err := b.viewType()
	if err != nil {
		return nil, err
	}

	m, ok := any(b).(sszutils.DynamicViewMarshaler)
	if !ok {
		return nil, errors.New("BlindedBeaconBlock: generated SSZ code missing")
	}

	fn := m.MarshalSSZDynView(view)
	if fn == nil {
		return nil, fmt.Errorf("BlindedBeaconBlock: no view marshaler for version %d", b.Version)
	}

	return fn(ds, buf)
}

// SizeSSZDyn returns the SSZ size of the blinded block for the active
// Version.
func (b *BlindedBeaconBlock) SizeSSZDyn(ds sszutils.DynamicSpecs) int {
	view, err := b.viewType()
	if err != nil {
		return 0
	}

	sz, ok := any(b).(sszutils.DynamicViewSizer)
	if !ok {
		return 0
	}

	fn := sz.SizeSSZDynView(view)
	if fn == nil {
		return 0
	}

	return fn(ds)
}

// UnmarshalSSZDyn decodes the blinded block into the view that matches
// Version.
func (b *BlindedBeaconBlock) UnmarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) error {
	view, err := b.viewType()
	if err != nil {
		return err
	}

	u, ok := any(b).(sszutils.DynamicViewUnmarshaler)
	if !ok {
		return errors.New("BlindedBeaconBlock: generated SSZ code missing")
	}

	fn := u.UnmarshalSSZDynView(view)
	if fn == nil {
		return fmt.Errorf("BlindedBeaconBlock: no view unmarshaler for version %d", b.Version)
	}

	if err := fn(ds, buf); err != nil {
		return err
	}

	b.populateVersion(b.Version)

	return nil
}

// populateVersion sets Version and propagates it to any nested versionable
// children.
func (b *BlindedBeaconBlock) populateVersion(v version.DataVersion) {
	propagateVersion(b, v)
}

// ToView returns a fresh fork-specific BlindedBeaconBlock populated with b's
// fields, recursing into Body via copyByName.
func (b *BlindedBeaconBlock) ToView() (any, error) {
	return toViewByCopy(b)
}

// FromView populates b from a fork-specific BlindedBeaconBlock.
func (b *BlindedBeaconBlock) FromView(view any) error {
	v, err := blindedBeaconBlockVersion(view)
	if err != nil {
		return err
	}

	if b.Version == version.DataVersionUnknown {
		b.Version = v
	}

	if err := copyByName(view, b); err != nil {
		return err
	}

	b.populateVersion(b.Version)

	return nil
}

// blindedBeaconBlockVersion maps a BlindedBeaconBlock view type to its
// DataVersion.
func blindedBeaconBlockVersion(view any) (version.DataVersion, error) {
	switch view.(type) {
	case *apiv1bellatrix.BlindedBeaconBlock:
		return version.DataVersionBellatrix, nil
	case *apiv1capella.BlindedBeaconBlock:
		return version.DataVersionCapella, nil
	case *apiv1deneb.BlindedBeaconBlock:
		return version.DataVersionDeneb, nil
	case *apiv1electra.BlindedBeaconBlock:
		return version.DataVersionElectra, nil
	default:
		return version.DataVersionUnknown, fmt.Errorf("BlindedBeaconBlock: unsupported view type %T", view)
	}
}

// HashTreeRootWithDyn computes the SSZ hash tree root using the active
// Version's view.
func (b *BlindedBeaconBlock) HashTreeRootWithDyn(ds sszutils.DynamicSpecs, hh sszutils.HashWalker) error {
	view, err := b.viewType()
	if err != nil {
		return err
	}

	h, ok := any(b).(sszutils.DynamicViewHashRoot)
	if !ok {
		return errors.New("BlindedBeaconBlock: generated SSZ code missing")
	}

	fn := h.HashTreeRootWithDynView(view)
	if fn == nil {
		return fmt.Errorf("BlindedBeaconBlock: no view hasher for version %d", b.Version)
	}

	return fn(ds, hh)
}

// MarshalSSZ implements the fastssz.Marshaler interface.
func (b *BlindedBeaconBlock) MarshalSSZ() ([]byte, error) {
	ds := dynssz.GetGlobalDynSsz()

	return b.MarshalSSZDyn(ds, make([]byte, 0, b.SizeSSZDyn(ds)))
}

// MarshalSSZTo implements the fastssz.Marshaler interface.
func (b *BlindedBeaconBlock) MarshalSSZTo(dst []byte) ([]byte, error) {
	return b.MarshalSSZDyn(dynssz.GetGlobalDynSsz(), dst)
}

// UnmarshalSSZ implements the fastssz.Unmarshaler interface.
func (b *BlindedBeaconBlock) UnmarshalSSZ(buf []byte) error {
	return b.UnmarshalSSZDyn(dynssz.GetGlobalDynSsz(), buf)
}

// SizeSSZ implements the fastssz.Marshaler interface.
func (b *BlindedBeaconBlock) SizeSSZ() int {
	return b.SizeSSZDyn(dynssz.GetGlobalDynSsz())
}

// HashTreeRoot implements the fastssz.HashRoot interface.
func (b *BlindedBeaconBlock) HashTreeRoot() ([32]byte, error) {
	return dynssz.GetGlobalDynSsz().HashTreeRoot(b)
}

// HashTreeRootWith implements the fastssz.HashRoot interface.
func (b *BlindedBeaconBlock) HashTreeRootWith(hh sszutils.HashWalker) error {
	return b.HashTreeRootWithDyn(dynssz.GetGlobalDynSsz(), hh)
}

// ToVersioned converts b into a *api.VersionedBlindedBeaconBlock.
func (b *BlindedBeaconBlock) ToVersioned() (*api.VersionedBlindedBeaconBlock, error) {
	out := &api.VersionedBlindedBeaconBlock{}
	if err := toVersioned(b.Version, b, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates b from src.
func (b *BlindedBeaconBlock) FromVersioned(src *api.VersionedBlindedBeaconBlock) error {
	return fromVersioned(b, src)
}

// MarshalJSON delegates to the per-fork BlindedBeaconBlock that matches
// Version.
func (b *BlindedBeaconBlock) MarshalJSON() ([]byte, error) {
	return marshalAsView(b)
}

// UnmarshalJSON delegates to the per-fork BlindedBeaconBlock that matches
// Version. Caller must set Version before calling.
func (b *BlindedBeaconBlock) UnmarshalJSON(data []byte) error {
	if err := unmarshalAsView(b, data); err != nil {
		return err
	}

	b.populateVersion(b.Version)

	return nil
}

// MarshalYAML delegates to the per-fork BlindedBeaconBlock that matches
// Version.
func (b *BlindedBeaconBlock) MarshalYAML() ([]byte, error) {
	return marshalAsViewYAML(b)
}

// UnmarshalYAML delegates to the per-fork BlindedBeaconBlock that matches
// Version. Caller must set Version before calling.
func (b *BlindedBeaconBlock) UnmarshalYAML(data []byte) error {
	if err := unmarshalAsViewYAML(b, data); err != nil {
		return err
	}

	b.populateVersion(b.Version)

	return nil
}
