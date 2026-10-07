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
	"errors"
	"fmt"

	client "github.com/ethpandaops/go-eth2-client"
	"github.com/ethpandaops/go-eth2-client/api"
	apiv1 "github.com/ethpandaops/go-eth2-client/api/v1"
)

// forkChoiceV2MetadataKeys are the beacon API response metadata fields.
var forkChoiceV2MetadataKeys = []string{"execution_optimistic", "finalized", "version", "dependent_root"}

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

	if raw, wrapped := fields["data"]; wrapped {
		var data *apiv1.ForkChoiceV2
		if err := json.Unmarshal(raw, &data); err != nil {
			return nil, errors.Join(errors.New("failed to parse fork choice"), err)
		}

		if data == nil {
			return nil, errors.New("fork choice data missing")
		}

		metadata := make(map[string]any)

		for k, v := range fields {
			if k == "data" {
				continue
			}

			var value any
			if err := json.Unmarshal(v, &value); err != nil {
				return nil, errors.Join(fmt.Errorf("failed to parse fork choice metadata %s", k), err)
			}

			metadata[k] = value
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

	metadata, err := takeUnwrappedMetadata(fields, &data)
	if err != nil {
		return nil, err
	}

	return &api.Response[*apiv1.ForkChoiceV2]{
		Data:     &data,
		Metadata: metadata,
	}, nil
}

// takeUnwrappedMetadata returns the response metadata of an unwrapped response. Decoding folds it into
// the fork choice's ExtraData alongside any other unknown top-level fields, so it is moved out, keeping
// it in Metadata whether or not a client wraps the response.
func takeUnwrappedMetadata(fields map[string]json.RawMessage, data *apiv1.ForkChoiceV2) (map[string]any, error) {
	metadata := make(map[string]any)

	// Values in the store's own extra_data take precedence when folding, so must stay.
	var storeExtraData map[string]json.RawMessage
	if raw, exists := fields["extra_data"]; exists {
		if err := json.Unmarshal(raw, &storeExtraData); err != nil {
			return nil, errors.Join(errors.New("failed to parse fork choice extra data"), err)
		}
	}

	for _, k := range forkChoiceV2MetadataKeys {
		raw, exists := fields[k]
		if !exists {
			continue
		}

		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, errors.Join(fmt.Errorf("failed to parse fork choice metadata %s", k), err)
		}

		metadata[k] = value

		if _, inStore := storeExtraData[k]; !inStore {
			delete(data.ExtraData, k)
		}
	}

	return metadata, nil
}
