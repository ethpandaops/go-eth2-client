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

// SignedBlockContents is a fork-agnostic signed block contents: the signed
// beacon block together with the blob sidecar data (KZG proofs and blobs).
// Block contents exist as a wire type from Deneb onwards; earlier forks
// submit bare signed beacon blocks and are not supported by this type.
type SignedBlockContents struct {
	Version     version.DataVersion
	SignedBlock *specall.SignedBeaconBlock
	KZGProofs   []deneb.KZGProof
	Blobs       []deneb.Blob
}

// viewType returns the fork-specific schema type pointer used as the view
// descriptor for the active Version.
func (s *SignedBlockContents) viewType() (any, error) {
	switch s.Version {
	case version.DataVersionDeneb:
		return (*apiv1deneb.SignedBlockContents)(nil), nil
	case version.DataVersionElectra:
		return (*apiv1electra.SignedBlockContents)(nil), nil
	case version.DataVersionFulu:
		// Fulu reuses the Electra signed block schema; only the KZG proof
		// list limit differs (cell proofs).
		return (*apiv1fulu.SignedBlockContents)(nil), nil
	default:
		return nil, fmt.Errorf("SignedBlockContents: unsupported version %d", s.Version)
	}
}

// MarshalSSZDyn marshals the signed block contents using the view that
// matches Version.
func (s *SignedBlockContents) MarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) ([]byte, error) {
	view, err := s.viewType()
	if err != nil {
		return nil, err
	}

	m, ok := any(s).(sszutils.DynamicViewMarshaler)
	if !ok {
		return nil, errors.New("SignedBlockContents: generated SSZ code missing")
	}

	fn := m.MarshalSSZDynView(view)
	if fn == nil {
		return nil, fmt.Errorf("SignedBlockContents: no view marshaler for version %d", s.Version)
	}

	return fn(ds, buf)
}

// SizeSSZDyn returns the SSZ size of the signed block contents for the
// active Version.
func (s *SignedBlockContents) SizeSSZDyn(ds sszutils.DynamicSpecs) int {
	view, err := s.viewType()
	if err != nil {
		return 0
	}

	sz, ok := any(s).(sszutils.DynamicViewSizer)
	if !ok {
		return 0
	}

	fn := sz.SizeSSZDynView(view)
	if fn == nil {
		return 0
	}

	return fn(ds)
}

// UnmarshalSSZDyn decodes the signed block contents into the view that
// matches Version.
func (s *SignedBlockContents) UnmarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) error {
	view, err := s.viewType()
	if err != nil {
		return err
	}

	u, ok := any(s).(sszutils.DynamicViewUnmarshaler)
	if !ok {
		return errors.New("SignedBlockContents: generated SSZ code missing")
	}

	fn := u.UnmarshalSSZDynView(view)
	if fn == nil {
		return fmt.Errorf("SignedBlockContents: no view unmarshaler for version %d", s.Version)
	}

	if err := fn(ds, buf); err != nil {
		return err
	}

	s.populateVersion(s.Version)

	return nil
}

// populateVersion sets Version and propagates it to any nested versionable
// children allocated by the SSZ unmarshal.
func (s *SignedBlockContents) populateVersion(v version.DataVersion) {
	propagateVersion(s, v)
}

// ToView returns a fresh fork-specific SignedBlockContents populated with
// s's fields, recursing into SignedBlock via copyByName.
func (s *SignedBlockContents) ToView() (any, error) {
	return toViewByCopy(s)
}

// FromView populates s from a fork-specific SignedBlockContents.
func (s *SignedBlockContents) FromView(view any) error {
	v, err := signedBlockContentsVersion(view)
	if err != nil {
		return err
	}

	if s.Version == version.DataVersionUnknown {
		s.Version = v
	}

	if err := copyByName(view, s); err != nil {
		return err
	}

	s.populateVersion(s.Version)

	return nil
}

// signedBlockContentsVersion maps a SignedBlockContents view type to its
// DataVersion.
func signedBlockContentsVersion(view any) (version.DataVersion, error) {
	switch view.(type) {
	case *apiv1deneb.SignedBlockContents:
		return version.DataVersionDeneb, nil
	case *apiv1electra.SignedBlockContents:
		return version.DataVersionElectra, nil
	case *apiv1fulu.SignedBlockContents:
		return version.DataVersionFulu, nil
	default:
		return version.DataVersionUnknown, fmt.Errorf("SignedBlockContents: unsupported view type %T", view)
	}
}

// HashTreeRootWithDyn computes the SSZ hash tree root using the active
// Version's view.
func (s *SignedBlockContents) HashTreeRootWithDyn(ds sszutils.DynamicSpecs, hh sszutils.HashWalker) error {
	view, err := s.viewType()
	if err != nil {
		return err
	}

	h, ok := any(s).(sszutils.DynamicViewHashRoot)
	if !ok {
		return errors.New("SignedBlockContents: generated SSZ code missing")
	}

	fn := h.HashTreeRootWithDynView(view)
	if fn == nil {
		return fmt.Errorf("SignedBlockContents: no view hasher for version %d", s.Version)
	}

	return fn(ds, hh)
}

// MarshalSSZ implements the fastssz.Marshaler interface.
func (s *SignedBlockContents) MarshalSSZ() ([]byte, error) {
	ds := dynssz.GetGlobalDynSsz()

	return s.MarshalSSZDyn(ds, make([]byte, 0, s.SizeSSZDyn(ds)))
}

// MarshalSSZTo implements the fastssz.Marshaler interface.
func (s *SignedBlockContents) MarshalSSZTo(dst []byte) ([]byte, error) {
	return s.MarshalSSZDyn(dynssz.GetGlobalDynSsz(), dst)
}

// UnmarshalSSZ implements the fastssz.Unmarshaler interface.
func (s *SignedBlockContents) UnmarshalSSZ(buf []byte) error {
	return s.UnmarshalSSZDyn(dynssz.GetGlobalDynSsz(), buf)
}

// SizeSSZ implements the fastssz.Marshaler interface.
func (s *SignedBlockContents) SizeSSZ() int {
	return s.SizeSSZDyn(dynssz.GetGlobalDynSsz())
}

// HashTreeRoot implements the fastssz.HashRoot interface.
func (s *SignedBlockContents) HashTreeRoot() ([32]byte, error) {
	return dynssz.GetGlobalDynSsz().HashTreeRoot(s)
}

// HashTreeRootWith implements the fastssz.HashRoot interface.
func (s *SignedBlockContents) HashTreeRootWith(hh sszutils.HashWalker) error {
	return s.HashTreeRootWithDyn(dynssz.GetGlobalDynSsz(), hh)
}

// ToVersioned converts s into a *api.VersionedSignedProposal.
func (s *SignedBlockContents) ToVersioned() (*api.VersionedSignedProposal, error) {
	out := &api.VersionedSignedProposal{}
	if err := toVersioned(s.Version, s, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates s from src. Only unblinded block contents
// proposals (Deneb onwards) are supported.
func (s *SignedBlockContents) FromVersioned(src *api.VersionedSignedProposal) error {
	return fromVersioned(s, src)
}

// MarshalJSON delegates to the per-fork SignedBlockContents that matches
// Version.
func (s *SignedBlockContents) MarshalJSON() ([]byte, error) {
	return marshalAsView(s)
}

// UnmarshalJSON delegates to the per-fork SignedBlockContents that matches
// Version. Caller must set Version before calling.
func (s *SignedBlockContents) UnmarshalJSON(data []byte) error {
	if err := unmarshalAsView(s, data); err != nil {
		return err
	}

	s.populateVersion(s.Version)

	return nil
}

// MarshalYAML delegates to the per-fork SignedBlockContents that matches
// Version.
func (s *SignedBlockContents) MarshalYAML() ([]byte, error) {
	return marshalAsViewYAML(s)
}

// UnmarshalYAML delegates to the per-fork SignedBlockContents that matches
// Version. Caller must set Version before calling.
func (s *SignedBlockContents) UnmarshalYAML(data []byte) error {
	if err := unmarshalAsViewYAML(s, data); err != nil {
		return err
	}

	s.populateVersion(s.Version)

	return nil
}
