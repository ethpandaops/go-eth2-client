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

	apiv1gloas "github.com/ethpandaops/go-eth2-client/api/v1/gloas"
	specall "github.com/ethpandaops/go-eth2-client/spec/all"
	"github.com/ethpandaops/go-eth2-client/spec/deneb"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/pk910/dynamic-ssz/sszutils"
)

// SignedExecutionPayloadEnvelopeContents is the fork-agnostic stateless
// request body for publishing an execution payload envelope: the signed
// envelope together with the blob sidecar data (KZG cell proofs and blobs).
// Envelopes exist as a wire type from Gloas onwards; later forks that share
// the schema reuse the gloas view.
type SignedExecutionPayloadEnvelopeContents struct {
	Version                        version.DataVersion
	SignedExecutionPayloadEnvelope *specall.SignedExecutionPayloadEnvelope
	KZGProofs                      []deneb.KZGProof
	Blobs                          []deneb.Blob
}

// viewType returns the fork-specific schema type pointer used as the view
// descriptor for the active Version.
func (s *SignedExecutionPayloadEnvelopeContents) viewType() (any, error) {
	switch s.Version {
	case version.DataVersionGloas, version.DataVersionHeze:
		return (*apiv1gloas.SignedExecutionPayloadEnvelopeContents)(nil), nil
	default:
		return nil, fmt.Errorf("SignedExecutionPayloadEnvelopeContents: unsupported version %d", s.Version)
	}
}

// MarshalSSZDyn marshals the signed envelope contents using the view that
// matches Version.
func (s *SignedExecutionPayloadEnvelopeContents) MarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) ([]byte, error) {
	view, err := s.viewType()
	if err != nil {
		return nil, err
	}

	m, ok := any(s).(sszutils.DynamicViewMarshaler)
	if !ok {
		return nil, errors.New("SignedExecutionPayloadEnvelopeContents: generated SSZ code missing")
	}

	fn := m.MarshalSSZDynView(view)
	if fn == nil {
		return nil, fmt.Errorf("SignedExecutionPayloadEnvelopeContents: no view marshaler for version %d", s.Version)
	}

	return fn(ds, buf)
}

// SizeSSZDyn returns the SSZ size of the signed envelope contents for the
// active Version.
func (s *SignedExecutionPayloadEnvelopeContents) SizeSSZDyn(ds sszutils.DynamicSpecs) int {
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

// UnmarshalSSZDyn decodes the signed envelope contents into the view that
// matches Version.
func (s *SignedExecutionPayloadEnvelopeContents) UnmarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) error {
	view, err := s.viewType()
	if err != nil {
		return err
	}

	u, ok := any(s).(sszutils.DynamicViewUnmarshaler)
	if !ok {
		return errors.New("SignedExecutionPayloadEnvelopeContents: generated SSZ code missing")
	}

	fn := u.UnmarshalSSZDynView(view)
	if fn == nil {
		return fmt.Errorf("SignedExecutionPayloadEnvelopeContents: no view unmarshaler for version %d", s.Version)
	}

	if err := fn(ds, buf); err != nil {
		return err
	}

	s.populateVersion(s.Version)

	return nil
}

// populateVersion sets Version and propagates it to any nested versionable
// children.
func (s *SignedExecutionPayloadEnvelopeContents) populateVersion(v version.DataVersion) {
	propagateVersion(s, v)
}

// ToView returns a fresh fork-specific SignedExecutionPayloadEnvelopeContents
// populated with s's fields, recursing into the envelope via copyByName.
func (s *SignedExecutionPayloadEnvelopeContents) ToView() (any, error) {
	return toViewByCopy(s)
}

// FromView populates s from a fork-specific
// SignedExecutionPayloadEnvelopeContents.
func (s *SignedExecutionPayloadEnvelopeContents) FromView(view any) error {
	v, err := signedExecutionPayloadEnvelopeContentsVersion(view)
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

// signedExecutionPayloadEnvelopeContentsVersion maps a
// SignedExecutionPayloadEnvelopeContents view type to its DataVersion.
func signedExecutionPayloadEnvelopeContentsVersion(view any) (version.DataVersion, error) {
	switch view.(type) {
	case *apiv1gloas.SignedExecutionPayloadEnvelopeContents:
		return version.DataVersionGloas, nil
	default:
		return version.DataVersionUnknown,
			fmt.Errorf("SignedExecutionPayloadEnvelopeContents: unsupported view type %T", view)
	}
}

// HashTreeRootWithDyn computes the SSZ hash tree root using the active
// Version's view.
func (s *SignedExecutionPayloadEnvelopeContents) HashTreeRootWithDyn(ds sszutils.DynamicSpecs, hh sszutils.HashWalker) error {
	view, err := s.viewType()
	if err != nil {
		return err
	}

	h, ok := any(s).(sszutils.DynamicViewHashRoot)
	if !ok {
		return errors.New("SignedExecutionPayloadEnvelopeContents: generated SSZ code missing")
	}

	fn := h.HashTreeRootWithDynView(view)
	if fn == nil {
		return fmt.Errorf("SignedExecutionPayloadEnvelopeContents: no view hasher for version %d", s.Version)
	}

	return fn(ds, hh)
}

// MarshalSSZ implements the fastssz.Marshaler interface.
func (s *SignedExecutionPayloadEnvelopeContents) MarshalSSZ() ([]byte, error) {
	ds := dynssz.GetGlobalDynSsz()

	return s.MarshalSSZDyn(ds, make([]byte, 0, s.SizeSSZDyn(ds)))
}

// MarshalSSZTo implements the fastssz.Marshaler interface.
func (s *SignedExecutionPayloadEnvelopeContents) MarshalSSZTo(dst []byte) ([]byte, error) {
	return s.MarshalSSZDyn(dynssz.GetGlobalDynSsz(), dst)
}

// UnmarshalSSZ implements the fastssz.Unmarshaler interface.
func (s *SignedExecutionPayloadEnvelopeContents) UnmarshalSSZ(buf []byte) error {
	return s.UnmarshalSSZDyn(dynssz.GetGlobalDynSsz(), buf)
}

// SizeSSZ implements the fastssz.Marshaler interface.
func (s *SignedExecutionPayloadEnvelopeContents) SizeSSZ() int {
	return s.SizeSSZDyn(dynssz.GetGlobalDynSsz())
}

// HashTreeRoot implements the fastssz.HashRoot interface.
func (s *SignedExecutionPayloadEnvelopeContents) HashTreeRoot() ([32]byte, error) {
	return dynssz.GetGlobalDynSsz().HashTreeRoot(s)
}

// HashTreeRootWith implements the fastssz.HashRoot interface.
func (s *SignedExecutionPayloadEnvelopeContents) HashTreeRootWith(hh sszutils.HashWalker) error {
	return s.HashTreeRootWithDyn(dynssz.GetGlobalDynSsz(), hh)
}

// MarshalJSON delegates to the per-fork
// SignedExecutionPayloadEnvelopeContents that matches Version.
func (s *SignedExecutionPayloadEnvelopeContents) MarshalJSON() ([]byte, error) {
	return marshalAsView(s)
}

// UnmarshalJSON delegates to the per-fork
// SignedExecutionPayloadEnvelopeContents that matches Version. Caller must
// set Version before calling.
func (s *SignedExecutionPayloadEnvelopeContents) UnmarshalJSON(data []byte) error {
	if err := unmarshalAsView(s, data); err != nil {
		return err
	}

	s.populateVersion(s.Version)

	return nil
}

// MarshalYAML delegates to the per-fork
// SignedExecutionPayloadEnvelopeContents that matches Version.
func (s *SignedExecutionPayloadEnvelopeContents) MarshalYAML() ([]byte, error) {
	return marshalAsViewYAML(s)
}

// UnmarshalYAML delegates to the per-fork
// SignedExecutionPayloadEnvelopeContents that matches Version. Caller must
// set Version before calling.
func (s *SignedExecutionPayloadEnvelopeContents) UnmarshalYAML(data []byte) error {
	if err := unmarshalAsViewYAML(s, data); err != nil {
		return err
	}

	s.populateVersion(s.Version)

	return nil
}
