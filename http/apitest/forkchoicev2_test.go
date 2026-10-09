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

// TestForkChoiceV2Spec decodes a response in the shape ethereum/beacon-APIs#615 defines: one pending,
// empty and full node for a Gloas block, the parent of whose pending node is not retained.
func TestForkChoiceV2Spec(t *testing.T) {
	body, err := os.ReadFile("testdata/forkchoicev2_spec.json")
	require.NoError(t, err)

	provider := forkChoiceV2Service(t, nethttp.StatusOK, body)
	response, err := provider.ForkChoiceV2(context.Background(), &api.ForkChoiceOpts{})
	require.NoError(t, err)
	require.NotNil(t, response.Data)
	require.Empty(t, response.Metadata)

	data := response.Data
	require.Equal(t, phase0.Epoch(353123), data.JustifiedCheckpoint.Epoch)
	require.Equal(t, phase0.Epoch(353122), data.FinalizedCheckpoint.Epoch)
	require.Len(t, data.ForkChoiceNodes, 3)

	pending, empty, full := data.ForkChoiceNodes[0], data.ForkChoiceNodes[1], data.ForkChoiceNodes[2]
	require.Equal(t, apiv1.ForkChoicePayloadStatusPending, pending.PayloadStatus)
	require.Nil(t, pending.ParentPayloadStatus)

	for _, node := range []*apiv1.ForkChoiceNodeV2{empty, full} {
		require.Equal(t, pending.BlockRoot, node.ParentRoot, "empty and full nodes point at their pending node")
		require.NotNil(t, node.ParentPayloadStatus)
		require.Equal(t, apiv1.ForkChoicePayloadStatusPending, *node.ParentPayloadStatus)
	}
	require.Equal(t, apiv1.ForkChoicePayloadStatusEmpty, empty.PayloadStatus)
	require.Equal(t, apiv1.ForkChoicePayloadStatusFull, full.PayloadStatus)

	for _, node := range data.ForkChoiceNodes {
		require.Equal(t, phase0.Epoch(353123), node.JustifiedCheckpoint.Epoch)
		require.Equal(t, phase0.Epoch(353122), node.FinalizedCheckpoint.Epoch)
		require.Equal(t, uint64(512), node.PayloadAttesterCount)
		require.Equal(t, uint64(510), node.PayloadAvailabilityYesCount)
		require.Equal(t, uint64(508), node.PayloadDataAvailabilityYesCount)
	}
}

// TestForkChoiceV2Clients documents how responses captured from Sepolia beacon nodes after the Gloas
// fork fall short of the spec: none decode until the clients implement ethereum/beacon-APIs#615.
func TestForkChoiceV2Clients(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		err     string
	}{
		{
			// Teku wraps the response in data, but its nodes have no parent_payload_status (nor
			// per-node checkpoints).
			name:    "Teku",
			fixture: "testdata/forkchoicev2_teku.json",
			err:     "parent payload status missing",
		},
		{
			// Prysm returns the response without the data wrapper.
			name:    "Prysm",
			fixture: "testdata/forkchoicev2_prysm.json",
			err:     "fork choice data missing",
		},
		{
			// Lodestar returns the response without the data wrapper.
			name:    "Lodestar",
			fixture: "testdata/forkchoicev2_lodestar.json",
			err:     "fork choice data missing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := os.ReadFile(test.fixture)
			require.NoError(t, err)

			provider := forkChoiceV2Service(t, nethttp.StatusOK, body)
			_, err = provider.ForkChoiceV2(context.Background(), &api.ForkChoiceOpts{})
			require.ErrorIs(t, err, consensusclient.ErrInvalidResponse)
			require.ErrorContains(t, err, test.err)
		})
	}
}

// TestForkChoiceV2DataNull ensures a wrapped response without data is an
// error rather than a nil fork choice.
func TestForkChoiceV2DataNull(t *testing.T) {
	provider := forkChoiceV2Service(t, nethttp.StatusOK, []byte(`{"data":null}`))
	_, err := provider.ForkChoiceV2(context.Background(), &api.ForkChoiceOpts{})
	require.ErrorIs(t, err, consensusclient.ErrInvalidResponse)
	require.EqualError(t, err, "invalid response: fork choice data missing")
}

// TestForkChoiceV2Metadata ensures fields beside data are kept as metadata.
func TestForkChoiceV2Metadata(t *testing.T) {
	body, err := os.ReadFile("testdata/forkchoicev2_spec.json")
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
