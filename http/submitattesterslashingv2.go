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

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"

	client "github.com/ethpandaops/go-eth2-client"
	"github.com/ethpandaops/go-eth2-client/api"
	"github.com/ethpandaops/go-eth2-client/spec"
)

// SubmitAttesterSlashingV2 submits a versioned attester slashing using the V2
// endpoint. The V1 endpoint was deprecated in the Electra release of the
// beacon APIs and later removed from the specification; V2 carries the fork of
// the submission in the Eth-Consensus-Version header.
func (s *Service) SubmitAttesterSlashingV2(ctx context.Context, opts *api.SubmitAttesterSlashingOpts) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	if opts == nil {
		return client.ErrNoOptions
	}

	unversionedSlashing, err := unversionAttesterSlashing(opts.AttesterSlashing)
	if err != nil {
		return err
	}

	specJSON, err := json.Marshal(unversionedSlashing)
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}

	endpoint := "/eth/v2/beacon/pool/attester_slashings"
	query := ""

	headers := map[string]string{
		"Eth-Consensus-Version": strings.ToLower(opts.AttesterSlashing.Version.String()),
	}
	if _, err := s.post(ctx,
		endpoint,
		query,
		&opts.Common,
		bytes.NewReader(specJSON),
		ContentTypeJSON,
		headers,
	); err != nil {
		return errors.Join(errors.New("failed to submit attester slashing"), err)
	}

	return nil
}

// unversionAttesterSlashing returns the fork-specific slashing for the
// versioned slashing's active version.
func unversionAttesterSlashing(slashing *spec.VersionedAttesterSlashing) (any, error) {
	if slashing == nil {
		return nil, errors.Join(errors.New("no attester slashing supplied"), client.ErrInvalidOptions)
	}

	var unversionedSlashing any
	switch slashing.Version {
	case spec.DataVersionPhase0:
		if slashing.Phase0 != nil {
			unversionedSlashing = slashing.Phase0
		}
	case spec.DataVersionAltair:
		if slashing.Altair != nil {
			unversionedSlashing = slashing.Altair
		}
	case spec.DataVersionBellatrix:
		if slashing.Bellatrix != nil {
			unversionedSlashing = slashing.Bellatrix
		}
	case spec.DataVersionCapella:
		if slashing.Capella != nil {
			unversionedSlashing = slashing.Capella
		}
	case spec.DataVersionDeneb:
		if slashing.Deneb != nil {
			unversionedSlashing = slashing.Deneb
		}
	case spec.DataVersionElectra:
		if slashing.Electra != nil {
			unversionedSlashing = slashing.Electra
		}
	case spec.DataVersionFulu:
		if slashing.Fulu != nil {
			unversionedSlashing = slashing.Fulu
		}
	case spec.DataVersionGloas:
		if slashing.Gloas != nil {
			unversionedSlashing = slashing.Gloas
		}
	case spec.DataVersionHeze:
		if slashing.Heze != nil {
			unversionedSlashing = slashing.Heze
		}
	default:
		return nil, errors.Join(errors.New("unhandled attester slashing version"), client.ErrInvalidOptions)
	}

	if unversionedSlashing == nil {
		return nil, errors.Join(errors.New("attester slashing not present for its version"), client.ErrInvalidOptions)
	}

	return unversionedSlashing, nil
}
