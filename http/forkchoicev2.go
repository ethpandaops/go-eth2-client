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
	"context"
	"encoding/json"
	"fmt"

	client "github.com/ethpandaops/go-eth2-client"
	"github.com/ethpandaops/go-eth2-client/api"
	apiv1 "github.com/ethpandaops/go-eth2-client/api/v1"
)

// ForkChoiceV2 fetches all current fork choice context, with one node per (block root, payload status) pair.
func (s *Service) ForkChoiceV2(ctx context.Context,
	opts *api.ForkChoiceOpts,
) (
	*api.Response[*apiv1.ForkChoiceV2],
	error,
) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	endpoint := "/eth/v2/debug/fork_choice"

	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, false)
	if err != nil {
		return nil, err
	}

	// Responses that do not follow the spec are client.ErrInvalidResponse, so that callers can
	// tell them apart from a failing beacon node.
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(httpResponse.body, &fields); err != nil {
		return nil, fmt.Errorf("%w: failed to parse fork choice: %w", client.ErrInvalidResponse, err)
	}

	raw, wrapped := fields["data"]
	if !wrapped {
		return nil, fmt.Errorf("%w: fork choice data missing", client.ErrInvalidResponse)
	}

	var data *apiv1.ForkChoiceV2
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("%w: failed to parse fork choice: %w", client.ErrInvalidResponse, err)
	}

	if data == nil {
		return nil, fmt.Errorf("%w: fork choice data missing", client.ErrInvalidResponse)
	}

	metadata := make(map[string]any)

	for k, v := range fields {
		if k == "data" {
			continue
		}

		var value any
		if err := json.Unmarshal(v, &value); err != nil {
			return nil, fmt.Errorf("%w: failed to parse fork choice metadata %s: %w", client.ErrInvalidResponse, k, err)
		}

		metadata[k] = value
	}

	return &api.Response[*apiv1.ForkChoiceV2]{
		Data:     data,
		Metadata: metadata,
	}, nil
}
