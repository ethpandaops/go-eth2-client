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

	// The spec wraps the response in a data container, but some clients return it unwrapped.
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(httpResponse.body, &fields); err != nil {
		return nil, errors.Join(errors.New("failed to parse fork choice"), err)
	}

	if _, wrapped := fields["data"]; wrapped {
		data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), new(apiv1.ForkChoiceV2))
		if err != nil {
			return nil, errors.Join(errors.New("failed to parse fork choice"), err)
		}

		return &api.Response[*apiv1.ForkChoiceV2]{
			Data:     data,
			Metadata: metadata,
		}, nil
	}

	var data apiv1.ForkChoiceV2
	if err := json.Unmarshal(httpResponse.body, &data); err != nil {
		return nil, errors.Join(errors.New("failed to parse fork choice"), err)
	}

	return &api.Response[*apiv1.ForkChoiceV2]{
		Data:     &data,
		Metadata: make(map[string]any),
	}, nil
}
