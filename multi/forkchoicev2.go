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

package multi

import (
	"context"
	"errors"
	"net/http"

	consensusclient "github.com/ethpandaops/go-eth2-client"
	"github.com/ethpandaops/go-eth2-client/api"
	apiv1 "github.com/ethpandaops/go-eth2-client/api/v1"
)

// ForkChoiceV2 fetches all current fork choice context, with one node per (block root, payload status) pair.
//
// Not every client implements the endpoint yet, and those that do not answer with an error status
// (Lighthouse 400, Nimbus 404) rather than failing over. Such clients, and clients whose response
// does not follow the spec (client.ErrInvalidResponse), are skipped without being deactivated in
// favour of the next client: they still serve other endpoints. If no client returns a fork choice and any of them
// lacked the endpoint, the error returned is (or, alongside the other clients' error, wraps) the
// last of those errors, so callers can fall back to GET /eth/v1/debug/fork_choice, which the
// clients without the v2 endpoint can serve.
func (s *Service) ForkChoiceV2(ctx context.Context,
	opts *api.ForkChoiceOpts,
) (
	*api.Response[*apiv1.ForkChoiceV2],
	error,
) {
	// The last errors of clients without the endpoint and of other failing clients. doCall only
	// keeps the last client's error, which is "empty response" if it lacked the endpoint.
	var unsupportedErr, otherErr error

	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		forkChoice, err := client.(consensusclient.ForkChoiceV2Provider).ForkChoiceV2(ctx, opts)
		if err != nil {
			var apiErr *api.Error
			if errors.As(err, &apiErr) && isUnsupportedEndpoint(apiErr.StatusCode) {
				unsupportedErr = err

				// No response, so that the next client is tried without deactivating this one; an
				// error would deactivate it, or end the call on a 4xx.
				return nil, nil //nolint:nilnil
			}
			otherErr = err

			if errors.Is(err, consensusclient.ErrInvalidResponse) {
				return nil, nil //nolint:nilnil
			}

			return nil, err
		}

		return forkChoice, nil
	}, nil)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return nil, err
		case unsupportedErr == nil && otherErr == nil:
			return nil, err
		case unsupportedErr == nil:
			return nil, otherErr
		case otherErr == nil:
			return nil, unsupportedErr
		default:
			// The unsupported error comes first, so that errors.As finds its status code.
			return nil, errors.Join(unsupportedErr, otherErr)
		}
	}

	response, isResponse := res.(*api.Response[*apiv1.ForkChoiceV2])
	if !isResponse {
		return nil, ErrIncorrectType
	}

	return response, nil
}

// isUnsupportedEndpoint returns true if the status code is one that clients return for an endpoint
// they do not implement.
//
// This includes 400, which Lighthouse returns for endpoint versions it does not have. The endpoint
// takes no parameters, so a client implementing it has no reason to return 400; if one did, it
// would be skipped as unsupported, and callers would fall back to v1.
func isUnsupportedEndpoint(statusCode int) bool {
	switch statusCode {
	case http.StatusBadRequest, http.StatusNotFound, http.StatusMethodNotAllowed, http.StatusNotImplemented:
		return true
	default:
		return false
	}
}
