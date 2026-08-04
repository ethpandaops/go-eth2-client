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

// Package apitest holds hermetic httptest-based tests for the http service.
// They live in their own package because package http's TestMain only runs
// tests when HTTP_ADDRESS points at a live node, which would leave these
// tests unexecuted in CI.
package apitest_test

import (
	"context"
	"encoding/json"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	consensusclient "github.com/ethpandaops/go-eth2-client"
	"github.com/ethpandaops/go-eth2-client/api"
	"github.com/ethpandaops/go-eth2-client/http"
	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestSubmitAttesterSlashingV2(t *testing.T) {
	var (
		postCount int
		gotPath   string
		gotHeader string
		gotBody   []byte
	)
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"8504736","sync_distance":"0"}}`))
		case "/eth/v2/beacon/pool/attester_slashings":
			postCount++
			gotPath = r.URL.Path
			gotHeader = r.Header.Get("Eth-Consensus-Version")
			gotBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(nethttp.StatusOK)
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc, err := http.New(ctx, http.WithAddress(srv.URL))
	require.NoError(t, err)

	submitter, ok := svc.(consensusclient.AttesterSlashingSubmitterV2)
	require.True(t, ok, "service does not implement AttesterSlashingSubmitterV2")

	slashing := &phase0.AttesterSlashing{
		Attestation1: &phase0.IndexedAttestation{
			AttestingIndices: []uint64{42},
			Data: &phase0.AttestationData{
				Slot:  100,
				Index: 1,
			},
		},
		Attestation2: &phase0.IndexedAttestation{
			AttestingIndices: []uint64{42},
			Data: &phase0.AttestationData{
				Slot:  100,
				Index: 1,
			},
		},
	}

	// Nil options: must error without submitting.
	err = submitter.SubmitAttesterSlashingV2(ctx, nil)
	require.Error(t, err)
	require.Equal(t, 0, postCount)

	// Wrong arm populated for the version: must error without submitting.
	err = submitter.SubmitAttesterSlashingV2(ctx, &api.SubmitAttesterSlashingOpts{
		AttesterSlashing: &spec.VersionedAttesterSlashing{
			Version: spec.DataVersionElectra,
			Phase0:  slashing,
		},
	})
	require.Error(t, err)
	require.Equal(t, 0, postCount)

	// Correct submission: V2 endpoint, version header, fork-specific body.
	err = submitter.SubmitAttesterSlashingV2(ctx, &api.SubmitAttesterSlashingOpts{
		AttesterSlashing: &spec.VersionedAttesterSlashing{
			Version: spec.DataVersionDeneb,
			Deneb:   slashing,
		},
	})
	require.NoError(t, err)
	require.Equal(t, 1, postCount)
	require.Equal(t, "/eth/v2/beacon/pool/attester_slashings", gotPath)
	require.Equal(t, "deneb", gotHeader)

	var body map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(gotBody, &body))
	require.Contains(t, body, "attestation_1")
	require.Contains(t, body, "attestation_2")
	require.NotEqual(t, "null", string(body["attestation_1"]))
}
