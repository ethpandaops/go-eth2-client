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

package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	consensushttp "github.com/ethpandaops/go-eth2-client/http"
	"github.com/stretchr/testify/require"
)

// dvtMockNode serves the minimum endpoints needed for connection activation, with a version
// string and an up/down state that a test can change at will.
type dvtMockNode struct {
	version atomic.Value // string
	up      atomic.Bool
}

func newDVTMockNode(version string) *dvtMockNode {
	m := &dvtMockNode{}
	m.version.Store(version)
	m.up.Store(true)

	return m
}

func (m *dvtMockNode) server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.up.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)

			return
		}

		switch r.URL.Path {
		case "/eth/v1/node/version":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"version":"` + m.version.Load().(string) + `"}}`))
		case "/eth/v1/node/syncing":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"100","sync_distance":"0"}}`))
		default:
			w.WriteHeader(http.StatusTeapot)
		}
	}))
}

// TestDVTMiddlewareFlagResetsOnDisconnect confirms that connecting to a node whose self-reported
// version string mentions Charon sets the DVT-middleware flag, that the flag is cleared once the
// connection goes inactive, and that reconnecting afterwards to a plain node does not inherit the
// stale flag.
func TestDVTMiddlewareFlagResetsOnDisconnect(t *testing.T) {
	node := newDVTMockNode("charon/v1.0.0")
	srv := node.server()
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, err := consensushttp.New(ctx, consensushttp.WithAddress(srv.URL), consensushttp.WithTimeout(2*time.Second))
	require.NoError(t, err)

	s, ok := svc.(*consensushttp.Service)
	require.True(t, ok)

	require.True(t, s.IsConnectedToDVTMiddleware(), "expected DVT flag to be set after connecting to a Charon-flagged node")

	// Node goes down: connection should become inactive, clearing the flag.
	node.up.Store(false)
	s.CheckConnectionState(ctx)
	require.False(t, s.IsActive())
	require.False(t, s.IsConnectedToDVTMiddleware(), "expected DVT flag to be cleared once the connection went inactive")

	// Reconnect, this time to a plain node with no DVT middleware in its version string.
	node.version.Store("lighthouse/v5.0.0")
	node.up.Store(true)
	s.CheckConnectionState(ctx)
	require.True(t, s.IsActive())
	require.False(t, s.IsConnectedToDVTMiddleware(), "expected DVT flag to stay clear after reconnecting to a plain node")
}

// TestDVTMiddlewareFlagConcurrentAccess exercises repeated connection-state transitions
// concurrently with reads of the DVT flag, matching the read/write pattern used by
// CheckConnectionState and Proposal/BlindedProposal, under the race detector.
func TestDVTMiddlewareFlagConcurrentAccess(t *testing.T) {
	node := newDVTMockNode("charon/v1.0.0")
	srv := node.server()
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, err := consensushttp.New(ctx, consensushttp.WithAddress(srv.URL), consensushttp.WithTimeout(2*time.Second))
	require.NoError(t, err)

	s, ok := svc.(*consensushttp.Service)
	require.True(t, ok)

	var (
		wg   sync.WaitGroup
		done atomic.Bool
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		defer done.Store(true)

		for i := range 100 {
			node.up.Store(i%2 == 0)
			s.CheckConnectionState(ctx)
		}
	}()

	go func() {
		defer wg.Done()

		for !done.Load() {
			_ = s.IsConnectedToDVTMiddleware()
		}
	}()

	wg.Wait()
}
