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

package apitest_test

import (
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"os"
	"testing"

	consensusclient "github.com/ethpandaops/go-eth2-client"
	"github.com/ethpandaops/go-eth2-client/api"
	apiv1 "github.com/ethpandaops/go-eth2-client/api/v1"
	"github.com/ethpandaops/go-eth2-client/http"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// forkChoiceV2Service returns a fork choice v2 provider backed by a server that answers
// GET /eth/v2/debug/fork_choice with the given status code and body.
func forkChoiceV2Service(t *testing.T, statusCode int, body []byte) consensusclient.ForkChoiceV2Provider {
	t.Helper()

	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"11299992","sync_distance":"0"}}`))
		case "/eth/v2/debug/fork_choice":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_, _ = w.Write(body)
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	svc, err := http.New(ctx, http.WithAddress(srv.URL))
	require.NoError(t, err)

	provider, ok := svc.(consensusclient.ForkChoiceV2Provider)
	require.True(t, ok, "service does not implement ForkChoiceV2Provider")

	return provider
}

// TestForkChoiceV2Clients decodes responses captured from Sepolia beacon nodes after the Gloas fork.
func TestForkChoiceV2Clients(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		check   func(t *testing.T, data *apiv1.ForkChoiceV2)
	}{
		{
			// Teku wraps the response in data, has PTC counts on every node and emits non-spec
			// fields at the top level of each node.
			name:    "Teku",
			fixture: "testdata/forkchoicev2_teku.json",
			check: func(t *testing.T, data *apiv1.ForkChoiceV2) {
				t.Helper()

				for _, node := range data.ForkChoiceNodes {
					require.NotNil(t, node.PayloadAttesterCount)
					require.NotNil(t, node.PayloadAvailabilityYesCount)
					require.NotNil(t, node.PayloadDataAvailabilityYesCount)
					require.Contains(t, node.ExtraData, "unrealised_justified_epoch")
					require.Contains(t, node.ExtraData, "unrealized_justified_root")
					require.Contains(t, node.ExtraData, "state_root")
				}
			},
		},
		{
			// Prysm returns the response unwrapped, with checkpoint epochs and PTC counts in the
			// extra data of the pending node only.
			name:    "Prysm",
			fixture: "testdata/forkchoicev2_prysm.json",
			check: func(t *testing.T, data *apiv1.ForkChoiceV2) {
				t.Helper()

				require.Contains(t, data.ExtraData, "unrealized_justified_checkpoint")
				pending := data.ForkChoiceNodes[0]
				require.Nil(t, pending.PayloadAttesterCount)
				require.Contains(t, pending.ExtraData, "payload_attester_count")
				require.Contains(t, pending.ExtraData, "unrealized_justified_epoch")
			},
		},
		{
			// Lodestar returns the response unwrapped, with checkpoint epochs and PTC counts in the
			// extra data of every node.
			name:    "Lodestar",
			fixture: "testdata/forkchoicev2_lodestar.json",
			check: func(t *testing.T, data *apiv1.ForkChoiceV2) {
				t.Helper()

				require.Contains(t, data.ExtraData, "unrealized_justified_checkpoint")
				for _, node := range data.ForkChoiceNodes {
					require.Nil(t, node.PayloadAttesterCount)
					require.Contains(t, node.ExtraData, "payload_attester_count")
					require.Contains(t, node.ExtraData, "unrealized_justified_epoch")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := os.ReadFile(test.fixture)
			require.NoError(t, err)

			provider := forkChoiceV2Service(t, nethttp.StatusOK, body)
			response, err := provider.ForkChoiceV2(context.Background(), &api.ForkChoiceOpts{})
			require.NoError(t, err)
			require.NotNil(t, response.Data)

			data := response.Data
			require.Equal(t, phase0.Epoch(353123), data.JustifiedCheckpoint.Epoch)
			require.Equal(t, phase0.Epoch(353122), data.FinalizedCheckpoint.Epoch)
			require.Len(t, data.ForkChoiceNodes, 3)

			statuses := make([]apiv1.ForkChoicePayloadStatus, 0, len(data.ForkChoiceNodes))
			for _, node := range data.ForkChoiceNodes {
				require.Equal(t, data.ForkChoiceNodes[0].BlockRoot, node.BlockRoot)
				require.Nil(t, node.ParentPayloadStatus)
				require.Nil(t, node.JustifiedEpoch)
				require.Nil(t, node.FinalizedEpoch)
				statuses = append(statuses, node.PayloadStatus)
			}
			require.Equal(t, []apiv1.ForkChoicePayloadStatus{
				apiv1.ForkChoicePayloadStatusPending,
				apiv1.ForkChoicePayloadStatusEmpty,
				apiv1.ForkChoicePayloadStatusFull,
			}, statuses)

			test.check(t, data)
		})
	}
}

// TestForkChoiceV2DataNull ensures a wrapped response without data is an
// error rather than a nil fork choice.
func TestForkChoiceV2DataNull(t *testing.T) {
	provider := forkChoiceV2Service(t, nethttp.StatusOK, []byte(`{"data":null}`))
	_, err := provider.ForkChoiceV2(context.Background(), &api.ForkChoiceOpts{})
	require.EqualError(t, err, "fork choice data missing")
}

// TestForkChoiceV2Metadata ensures fields beside data are kept as metadata.
func TestForkChoiceV2Metadata(t *testing.T) {
	body, err := os.ReadFile("testdata/forkchoicev2_teku.json")
	require.NoError(t, err)

	var wrapped map[string]any
	require.NoError(t, json.Unmarshal(body, &wrapped))
	wrapped["execution_optimistic"] = false
	body, err = json.Marshal(wrapped)
	require.NoError(t, err)

	provider := forkChoiceV2Service(t, nethttp.StatusOK, body)
	response, err := provider.ForkChoiceV2(context.Background(), &api.ForkChoiceOpts{})
	require.NoError(t, err)
	require.Len(t, response.Data.ForkChoiceNodes, 3)
	require.Equal(t, map[string]any{"execution_optimistic": false}, response.Metadata)
}

// TestForkChoiceV2UnwrappedMetadata ensures response metadata of an unwrapped response is kept as
// metadata, as for wrapped responses, rather than as fork choice extra data.
func TestForkChoiceV2UnwrappedMetadata(t *testing.T) {
	body, err := os.ReadFile("testdata/forkchoicev2_prysm.json")
	require.NoError(t, err)

	var unwrapped map[string]any
	require.NoError(t, json.Unmarshal(body, &unwrapped))
	unwrapped["execution_optimistic"] = false
	unwrapped["unknown_field"] = "kept"
	body, err = json.Marshal(unwrapped)
	require.NoError(t, err)

	provider := forkChoiceV2Service(t, nethttp.StatusOK, body)
	response, err := provider.ForkChoiceV2(context.Background(), &api.ForkChoiceOpts{})
	require.NoError(t, err)
	require.Equal(t, map[string]any{"execution_optimistic": false}, response.Metadata)
	require.NotContains(t, response.Data.ExtraData, "execution_optimistic")
	require.Equal(t, "kept", response.Data.ExtraData["unknown_field"])
	require.Contains(t, response.Data.ExtraData, "unrealized_justified_checkpoint")
}

// TestForkChoiceV2Unsupported ensures that clients without the endpoint surface their status code,
// so that callers can fall back to GET /eth/v1/debug/fork_choice.
func TestForkChoiceV2Unsupported(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       []byte
	}{
		{
			name:       "Lighthouse",
			statusCode: nethttp.StatusBadRequest,
			body:       []byte(`{"code":400,"message":"BAD_REQUEST: Unsupported endpoint version: v2","stacktraces":[]}`),
		},
		{
			name:       "Nimbus",
			statusCode: nethttp.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := forkChoiceV2Service(t, test.statusCode, test.body)
			_, err := provider.ForkChoiceV2(context.Background(), &api.ForkChoiceOpts{})
			require.Error(t, err)

			var apiErr *api.Error
			require.True(t, errors.As(err, &apiErr), "error is not an api.Error: %v", err)
			require.Equal(t, test.statusCode, apiErr.StatusCode)
		})
	}
}
