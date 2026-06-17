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

	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/pk910/dynamic-ssz/sszutils"
)

// ExecutionRequests is a fork-agnostic set of execution layer triggered
// requests. The container is introduced in Electra; the builder request fields
// (BuilderDeposits, BuilderExits) are populated from Gloas onwards (EIP-8282).
type ExecutionRequests struct {
	Version         version.DataVersion
	Deposits        []*electra.DepositRequest
	Withdrawals     []*electra.WithdrawalRequest
	Consolidations  []*electra.ConsolidationRequest
	BuilderDeposits []*gloas.BuilderDepositRequest
	BuilderExits    []*gloas.BuilderExitRequest
}

// viewType returns the fork-specific schema type pointer used as the view
// descriptor for the active Version.
func (e *ExecutionRequests) viewType() (any, error) {
	switch e.Version {
	case version.DataVersionElectra,
		version.DataVersionFulu:
		return (*electra.ExecutionRequests)(nil), nil
	case version.DataVersionGloas,
		version.DataVersionHeze:
		return (*gloas.ExecutionRequests)(nil), nil
	default:
		return nil, fmt.Errorf("ExecutionRequests: unsupported version %d", e.Version)
	}
}

// MarshalSSZDyn marshals the requests using the view that matches Version.
func (e *ExecutionRequests) MarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) ([]byte, error) {
	view, err := e.viewType()
	if err != nil {
		return nil, err
	}

	m, ok := any(e).(sszutils.DynamicViewMarshaler)
	if !ok {
		return nil, errors.New("ExecutionRequests: generated SSZ code missing")
	}

	fn := m.MarshalSSZDynView(view)
	if fn == nil {
		return nil, fmt.Errorf("ExecutionRequests: no view marshaler for version %d", e.Version)
	}

	return fn(ds, buf)
}

// SizeSSZDyn returns the SSZ size of the requests for the active Version.
func (e *ExecutionRequests) SizeSSZDyn(ds sszutils.DynamicSpecs) int {
	view, err := e.viewType()
	if err != nil {
		return 0
	}

	s, ok := any(e).(sszutils.DynamicViewSizer)
	if !ok {
		return 0
	}

	fn := s.SizeSSZDynView(view)
	if fn == nil {
		return 0
	}

	return fn(ds)
}

// UnmarshalSSZDyn decodes the requests into the view that matches Version.
func (e *ExecutionRequests) UnmarshalSSZDyn(ds sszutils.DynamicSpecs, buf []byte) error {
	view, err := e.viewType()
	if err != nil {
		return err
	}

	u, ok := any(e).(sszutils.DynamicViewUnmarshaler)
	if !ok {
		return errors.New("ExecutionRequests: generated SSZ code missing")
	}

	fn := u.UnmarshalSSZDynView(view)
	if fn == nil {
		return fmt.Errorf("ExecutionRequests: no view unmarshaler for version %d", e.Version)
	}

	if err := fn(ds, buf); err != nil {
		return err
	}

	e.populateVersion(e.Version)

	return nil
}

// populateVersion sets Version. ExecutionRequests has no versionable children.
func (e *ExecutionRequests) populateVersion(v version.DataVersion) {
	e.Version = v
}

// ToView returns a fresh fork-specific ExecutionRequests populated with e's
// fields, recursing into the request slices via copyByName.
func (e *ExecutionRequests) ToView() (any, error) {
	return toViewByCopy(e)
}

// FromView populates e from a fork-specific ExecutionRequests.
func (e *ExecutionRequests) FromView(view any) error {
	v, err := executionRequestsVersion(view)
	if err != nil {
		return err
	}

	if e.Version == version.DataVersionUnknown {
		e.Version = v
	}

	if err := copyByName(view, e); err != nil {
		return err
	}

	e.populateVersion(e.Version)

	return nil
}

// executionRequestsVersion maps an ExecutionRequests view type to its
// DataVersion. Forks that share a view (Electra/Fulu, Gloas/Heze) resolve to
// the earliest; fromVersioned pins the precise version before FromView runs.
func executionRequestsVersion(view any) (version.DataVersion, error) {
	switch view.(type) {
	case *electra.ExecutionRequests:
		return version.DataVersionElectra, nil
	case *gloas.ExecutionRequests:
		return version.DataVersionGloas, nil
	default:
		return version.DataVersionUnknown, fmt.Errorf("ExecutionRequests: unsupported view type %T", view)
	}
}

// HashTreeRootWithDyn computes the SSZ hash tree root using the active Version's view.
func (e *ExecutionRequests) HashTreeRootWithDyn(ds sszutils.DynamicSpecs, hh sszutils.HashWalker) error {
	view, err := e.viewType()
	if err != nil {
		return err
	}

	h, ok := any(e).(sszutils.DynamicViewHashRoot)
	if !ok {
		return errors.New("ExecutionRequests: generated SSZ code missing")
	}

	fn := h.HashTreeRootWithDynView(view)
	if fn == nil {
		return fmt.Errorf("ExecutionRequests: no view hasher for version %d", e.Version)
	}

	return fn(ds, hh)
}

// MarshalSSZ implements the fastssz.Marshaler interface.
func (e *ExecutionRequests) MarshalSSZ() ([]byte, error) {
	ds := dynssz.GetGlobalDynSsz()

	return e.MarshalSSZDyn(ds, make([]byte, 0, e.SizeSSZDyn(ds)))
}

// MarshalSSZTo implements the fastssz.Marshaler interface.
func (e *ExecutionRequests) MarshalSSZTo(dst []byte) ([]byte, error) {
	return e.MarshalSSZDyn(dynssz.GetGlobalDynSsz(), dst)
}

// UnmarshalSSZ implements the fastssz.Unmarshaler interface.
func (e *ExecutionRequests) UnmarshalSSZ(buf []byte) error {
	return e.UnmarshalSSZDyn(dynssz.GetGlobalDynSsz(), buf)
}

// SizeSSZ implements the fastssz.Marshaler interface.
func (e *ExecutionRequests) SizeSSZ() int {
	return e.SizeSSZDyn(dynssz.GetGlobalDynSsz())
}

// HashTreeRoot implements the fastssz.HashRoot interface.
func (e *ExecutionRequests) HashTreeRoot() ([32]byte, error) {
	return dynssz.GetGlobalDynSsz().HashTreeRoot(e)
}

// HashTreeRootWith implements the fastssz.HashRoot interface.
func (e *ExecutionRequests) HashTreeRootWith(hh sszutils.HashWalker) error {
	return e.HashTreeRootWithDyn(dynssz.GetGlobalDynSsz(), hh)
}

// ToVersioned converts e into a *spec.VersionedExecutionRequests.
func (e *ExecutionRequests) ToVersioned() (*spec.VersionedExecutionRequests, error) {
	out := &spec.VersionedExecutionRequests{}
	if err := toVersioned(e.Version, e, out); err != nil {
		return nil, err
	}

	return out, nil
}

// FromVersioned populates e from src.
func (e *ExecutionRequests) FromVersioned(src *spec.VersionedExecutionRequests) error {
	return fromVersioned(e, src)
}

// MarshalJSON delegates to the per-fork ExecutionRequests that matches Version.
func (e *ExecutionRequests) MarshalJSON() ([]byte, error) {
	return marshalAsView(e)
}

// UnmarshalJSON delegates to the per-fork ExecutionRequests that matches Version.
// Caller must set Version before calling.
func (e *ExecutionRequests) UnmarshalJSON(data []byte) error {
	if err := unmarshalAsView(e, data); err != nil {
		return err
	}

	e.populateVersion(e.Version)

	return nil
}

// MarshalYAML delegates to the per-fork ExecutionRequests that matches Version.
func (e *ExecutionRequests) MarshalYAML() ([]byte, error) {
	return marshalAsViewYAML(e)
}

// UnmarshalYAML delegates to the per-fork ExecutionRequests that matches Version.
// Caller must set Version before calling.
func (e *ExecutionRequests) UnmarshalYAML(data []byte) error {
	if err := unmarshalAsViewYAML(e, data); err != nil {
		return err
	}

	e.populateVersion(e.Version)

	return nil
}
