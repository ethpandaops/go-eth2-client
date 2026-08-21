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
	"fmt"
	"strings"

	client "github.com/ethpandaops/go-eth2-client"
	"github.com/ethpandaops/go-eth2-client/api"
	apiv1all "github.com/ethpandaops/go-eth2-client/api/v1/all"
	apiv1gloas "github.com/ethpandaops/go-eth2-client/api/v1/gloas"
	apiv2 "github.com/ethpandaops/go-eth2-client/api/v2"
	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/deneb"
)

// SubmitExecutionPayloadEnvelope submits a signed execution payload envelope
// using the stateless request form: a SignedExecutionPayloadEnvelopeContents
// body (signed envelope plus blobs and KZG proofs) with the
// Eth-Blob-Data-Included header set to true. The stateful form (bare
// envelope, header false) only works when the beacon node cached the blobs
// and KZG proofs from its own block production, so builders publishing
// externally-built payloads must use the stateless form.
func (s *Service) SubmitExecutionPayloadEnvelope(ctx context.Context,
	opts *api.SubmitExecutionPayloadEnvelopeOpts,
) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	if opts == nil {
		return client.ErrNoOptions
	}

	if opts.SignedExecutionPayloadEnvelope == nil {
		return errors.Join(errors.New("no envelope supplied"), client.ErrInvalidOptions)
	}

	versioned := opts.SignedExecutionPayloadEnvelope

	var contents any

	switch versioned.Version {
	case spec.DataVersionGloas, spec.DataVersionHeze:
		// All current envelope-bearing forks (Gloas, Heze) reuse the gloas schema.
		if versioned.Gloas == nil {
			return errors.Join(errors.New("no gloas envelope supplied"), client.ErrInvalidOptions)
		}

		contents = &apiv1gloas.SignedExecutionPayloadEnvelopeContents{
			SignedExecutionPayloadEnvelope: versioned.Gloas,
			KZGProofs:                      nonNilKZGProofs(opts.KZGProofs),
			Blobs:                          nonNilBlobs(opts.Blobs),
		}
	default:
		return errors.Join(
			fmt.Errorf("unsupported envelope version %s", versioned.Version),
			client.ErrInvalidOptions,
		)
	}

	body, contentType, err := s.submitExecutionPayloadEnvelopeData(ctx, contents)
	if err != nil {
		return err
	}

	return s.postExecutionPayloadEnvelope(ctx, &opts.Common, versioned.Version,
		opts.BroadcastValidation, body, contentType)
}

// SubmitAgnosticExecutionPayloadEnvelope submits a signed execution payload
// envelope supplied as a fork-agnostic *all.SignedExecutionPayloadEnvelope,
// using the same stateless request form as SubmitExecutionPayloadEnvelope.
func (s *Service) SubmitAgnosticExecutionPayloadEnvelope(ctx context.Context,
	opts *api.SubmitAgnosticExecutionPayloadEnvelopeOpts,
) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}

	if opts == nil {
		return client.ErrNoOptions
	}

	if opts.SignedExecutionPayloadEnvelope == nil {
		return errors.Join(errors.New("no envelope supplied"), client.ErrInvalidOptions)
	}

	contents := &apiv1all.SignedExecutionPayloadEnvelopeContents{
		Version:                        opts.SignedExecutionPayloadEnvelope.Version,
		SignedExecutionPayloadEnvelope: opts.SignedExecutionPayloadEnvelope,
		KZGProofs:                      nonNilKZGProofs(opts.KZGProofs),
		Blobs:                          nonNilBlobs(opts.Blobs),
	}

	body, contentType, err := s.submitExecutionPayloadEnvelopeData(ctx, contents)
	if err != nil {
		return err
	}

	return s.postExecutionPayloadEnvelope(ctx, &opts.Common, contents.Version,
		opts.BroadcastValidation, body, contentType)
}

// postExecutionPayloadEnvelope performs the POST shared by both submit
// variants, setting the consensus version and stateless-form headers.
func (s *Service) postExecutionPayloadEnvelope(ctx context.Context,
	common *api.CommonOpts,
	consensusVersion spec.DataVersion,
	broadcastValidation *apiv2.BroadcastValidation,
	body []byte,
	contentType ContentType,
) error {
	endpoint := "/eth/v1/beacon/execution_payload_envelopes"

	query := ""
	if broadcastValidation != nil {
		query = "broadcast_validation=" + broadcastValidation.String()
	}

	headers := make(map[string]string)
	headers["Eth-Consensus-Version"] = strings.ToLower(consensusVersion.String())
	// Always the stateless form: Eth-Blob-Data-Included "true" selects the
	// Contents body schema (beacon-APIs#624); strict consensus clients reject
	// the request when the header is missing. The pre-#624 discriminator
	// (Eth-Execution-Payload-Blinded "false" = Contents) is still sent for
	// beacon nodes that have not adopted the rename yet.
	headers["Eth-Blob-Data-Included"] = "true"
	headers["Eth-Execution-Payload-Blinded"] = "false"

	if _, err := s.post(ctx, endpoint, query, common, bytes.NewBuffer(body), contentType, headers); err != nil {
		return errors.Join(errors.New("failed to submit execution payload envelope"), err)
	}

	return nil
}

// submitExecutionPayloadEnvelopeData marshals the envelope contents to the
// negotiated content type (SSZ unless JSON is enforced).
func (s *Service) submitExecutionPayloadEnvelopeData(ctx context.Context,
	contents any,
) (
	[]byte,
	ContentType,
	error,
) {
	if s.enforceJSON {
		body, err := json.Marshal(contents)
		if err != nil {
			return nil, ContentTypeUnknown, errors.Join(errors.New("failed to marshal JSON"), err)
		}

		return body, ContentTypeJSON, nil
	}

	ds, err := s.dynSSZForRequest(ctx)
	if err != nil {
		return nil, ContentTypeUnknown, err
	}

	body, err := ds.MarshalSSZ(contents)
	if err != nil {
		return nil, ContentTypeUnknown, errors.Join(errors.New("failed to marshal SSZ"), err)
	}

	return body, ContentTypeSSZ, nil
}

// nonNilKZGProofs normalises a nil KZG proof slice to an empty one: the
// beacon-API schema requires kzg_proofs to be present (as an empty array
// when the payload carries no blobs).
func nonNilKZGProofs(proofs []deneb.KZGProof) []deneb.KZGProof {
	if proofs == nil {
		return []deneb.KZGProof{}
	}

	return proofs
}

// nonNilBlobs normalises a nil blob slice to an empty one: the beacon-API
// schema requires blobs to be present (as an empty array when the payload
// carries no blobs).
func nonNilBlobs(blobs []deneb.Blob) []deneb.Blob {
	if blobs == nil {
		return []deneb.Blob{}
	}

	return blobs
}
