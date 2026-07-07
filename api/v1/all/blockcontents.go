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
	apiv1deneb "github.com/ethpandaops/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/ethpandaops/go-eth2-client/api/v1/electra"
	apiv1fulu "github.com/ethpandaops/go-eth2-client/api/v1/fulu"
	specall "github.com/ethpandaops/go-eth2-client/spec/all"
	"github.com/ethpandaops/go-eth2-client/spec/deneb"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/pk910/dynamic-ssz/sszutils"
)

// BlockContents is a fork-agnostic block contents: the beacon block together
// with the blob sidecar data (KZG proofs and blobs). Block contents exist as
// a wire type from Deneb onwards; earlier forks use bare beacon blocks and
// are not supported by this type.
type BlockContents struct {
	Version   version.DataVersion
	Block     *specall.BeaconBlock
	KZGProofs []deneb.KZGProof
	Blobs     []deneb.Blob
}

// viewType returns the fork-specific schema type pointer used as the view
// descriptor for the active Version.
func (b *BlockContents) viewType() (any, error) {
	switch b.Version {
	case version.DataVersionDeneb:
		return (*apiv1deneb.BlockContents)(nil), nil
	case version.DataVersionElectra:
		return (*apiv1electra.BlockContents)(nil), nil
	case version.DataVersionFulu:
		// Fulu reuses the Electra block schema; only the KZG proof list
		// limit differs (cell proofs).
		return (*apiv1fulu.BlockContents)(nil), nil
	default:
		return nil, fmt.Errorf("BlockContents: unsupported version %d", b.Version)
	}
}

// MarshalSSZDyn marshals the block contents using the view that matches
// Version.
func (b *BlockContents) MarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) ([]byte, error) {
	view, err := b.viewType()
	if err != nil {
		return nil, err
	}

	m, ok := any(b).(sszutils.DynamicViewMarshaler)
	if !ok {
		return nil, errors.New("BlockContents: generated SSZ code missing")
	}

	fn := m.MarshalSSZDynView(view)
	if fn == nil {
		return nil, fmt.Errorf("BlockContents: no view marshaler for version %d", b.Version)
	}

	return fn(ds, buf)
}

// SizeSSZDyn returns the SSZ size of the block contents for the active
// Version.
func (b *BlockContents) SizeSSZDyn(ds sszutils.DynamicSpecs) int {
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

// UnmarshalSSZDyn decodes the block contents into the view that matches
// Version.
func (b *BlockContents) UnmarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) error {
	view, err := b.viewType()
	if err != nil {
		return err
	}

	u, ok := any(b).(sszutils.DynamicViewUnmarshaler)
	if !ok {
		return errors.New("BlockContents: generated SSZ code missing")
	}

	fn := u.UnmarshalSSZDynView(view)
	if fn == nil {
		return fmt.Errorf("BlockContents: no view unmarshaler for version %d", b.Version)
	}

	if err := fn(ds, buf); err != nil {
		return err
	}

	b.populateVersion(b.Version)

	return nil
}

// populateVersion sets Version and propagates it to any nested versionable
// children allocated by the SSZ unmarshal.
func (b *BlockContents) populateVersion(v version.DataVersion) {
	propagateVersion(b, v)
}

// ToView returns a fresh fork-specific BlockContents populated with b's
// fields, recursing into Block via copyByName.
func (b *BlockContents) ToView() (any, error) {
	return toViewByCopy(b)
}

// FromView populates b from a fork-specific BlockContents.
func (b *BlockContents) FromView(view any) error {
	v, err := blockContentsVersion(view)
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

// blockContentsVersion maps a BlockContents view type to its DataVersion.
func blockContentsVersion(view any) (version.DataVersion, error) {
	switch view.(type) {
	case *apiv1deneb.BlockContents:
		return version.DataVersionDeneb, nil
	case *apiv1electra.BlockContents:
		return version.DataVersionElectra, nil
	case *apiv1fulu.BlockContents:
		return version.DataVersionFulu, nil
	default:
		return version.DataVersionUnknown, fmt.Errorf("BlockContents: unsupported view type %T", view)
	}
}

// HashTreeRootWithDyn computes the SSZ hash tree root using the active
// Version's view.
func (b *BlockContents) HashTreeRootWithDyn(ds sszutils.DynamicSpecs, hh sszutils.HashWalker) error {
	view, err := b.viewType()
	if err != nil {
		return err
	}

	h, ok := any(b).(sszutils.DynamicViewHashRoot)
	if !ok {
		return errors.New("BlockContents: generated SSZ code missing")
	}

	fn := h.HashTreeRootWithDynView(view)
	if fn == nil {
		return fmt.Errorf("BlockContents: no view hasher for version %d", b.Version)
	}

	return fn(ds, hh)
}

// MarshalSSZ implements the fastssz.Marshaler interface.
func (b *BlockContents) MarshalSSZ() ([]byte, error) {
	ds := dynssz.GetGlobalDynSsz()

	return b.MarshalSSZDyn(ds, make([]byte, 0, b.SizeSSZDyn(ds)))
}

// MarshalSSZTo implements the fastssz.Marshaler interface.
func (b *BlockContents) MarshalSSZTo(dst []byte) ([]byte, error) {
	return b.MarshalSSZDyn(dynssz.GetGlobalDynSsz(), dst)
}

// UnmarshalSSZ implements the fastssz.Unmarshaler interface.
func (b *BlockContents) UnmarshalSSZ(buf []byte) error {
	return b.UnmarshalSSZDyn(dynssz.GetGlobalDynSsz(), buf)
}

// SizeSSZ implements the fastssz.Marshaler interface.
func (b *BlockContents) SizeSSZ() int {
	return b.SizeSSZDyn(dynssz.GetGlobalDynSsz())
}

// HashTreeRoot implements the fastssz.HashRoot interface.
func (b *BlockContents) HashTreeRoot() ([32]byte, error) {
	return dynssz.GetGlobalDynSsz().HashTreeRoot(b)
}

// HashTreeRootWith implements the fastssz.HashRoot interface.
func (b *BlockContents) HashTreeRootWith(hh sszutils.HashWalker) error {
	return b.HashTreeRootWithDyn(dynssz.GetGlobalDynSsz(), hh)
}

// ToVersioned converts b into a *api.VersionedProposal.
func (b *BlockContents) ToVersioned() (*api.VersionedProposal, error) {
	out := &api.VersionedProposal{}
	if err := toVersioned(b.Version, b, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates b from src. Only unblinded block contents
// proposals (Deneb onwards) are supported.
func (b *BlockContents) FromVersioned(src *api.VersionedProposal) error {
	return fromVersioned(b, src)
}

// MarshalJSON delegates to the per-fork BlockContents that matches Version.
func (b *BlockContents) MarshalJSON() ([]byte, error) {
	return marshalAsView(b)
}

// UnmarshalJSON delegates to the per-fork BlockContents that matches
// Version. Caller must set Version before calling.
func (b *BlockContents) UnmarshalJSON(data []byte) error {
	if err := unmarshalAsView(b, data); err != nil {
		return err
	}

	b.populateVersion(b.Version)

	return nil
}

// MarshalYAML delegates to the per-fork BlockContents that matches Version.
func (b *BlockContents) MarshalYAML() ([]byte, error) {
	return marshalAsViewYAML(b)
}

// UnmarshalYAML delegates to the per-fork BlockContents that matches
// Version. Caller must set Version before calling.
func (b *BlockContents) UnmarshalYAML(data []byte) error {
	if err := unmarshalAsViewYAML(b, data); err != nil {
		return err
	}

	b.populateVersion(b.Version)

	return nil
}
