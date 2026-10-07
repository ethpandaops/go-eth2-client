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

package multi_test

import (
	"context"
	"errors"
	nethttp "net/http"
	"testing"

	consensusclient "github.com/ethpandaops/go-eth2-client"
	"github.com/ethpandaops/go-eth2-client/api"
	apiv1 "github.com/ethpandaops/go-eth2-client/api/v1"
	"github.com/ethpandaops/go-eth2-client/mock"
	"github.com/ethpandaops/go-eth2-client/multi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// unsupportedForkChoiceV2Client returns a client that answers fork choice v2 requests with the given status code, as
// clients without the endpoint do.
func unsupportedForkChoiceV2Client(t *testing.T, name string, statusCode int) consensusclient.Service {
	t.Helper()

	client, err := mock.New(context.Background(), mock.WithName(name))
	require.NoError(t, err)
	client.ForkChoiceV2Func = func(context.Context, *api.ForkChoiceOpts) (*api.Response[*apiv1.ForkChoiceV2], error) {
		return nil, &api.Error{Method: nethttp.MethodGet, Endpoint: "/eth/v2/debug/fork_choice", StatusCode: statusCode}
	}

	return client
}

func TestForkChoiceV2(t *testing.T) {
	ctx := context.Background()

	supporting, err := mock.New(ctx, mock.WithName("mock 3"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			unsupportedForkChoiceV2Client(t, "mock 1", nethttp.StatusBadRequest),
			unsupportedForkChoiceV2Client(t, "mock 2", nethttp.StatusNotFound),
			supporting,
		}),
	)
	require.NoError(t, err)

	res, err := multiClient.(consensusclient.ForkChoiceV2Provider).ForkChoiceV2(ctx, &api.ForkChoiceOpts{})
	require.NoError(t, err)
	require.NotNil(t, res)

	// Clients without the endpoint are skipped, not deactivated.
	require.Equal(t, "mock 1", multiClient.Address())
}

func TestForkChoiceV2Unsupported(t *testing.T) {
	ctx := context.Background()

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			unsupportedForkChoiceV2Client(t, "mock 1", nethttp.StatusBadRequest),
			unsupportedForkChoiceV2Client(t, "mock 2", nethttp.StatusNotFound),
		}),
	)
	require.NoError(t, err)

	_, err = multiClient.(consensusclient.ForkChoiceV2Provider).ForkChoiceV2(ctx, &api.ForkChoiceOpts{})
	require.Error(t, err)

	var apiErr *api.Error
	require.True(t, errors.As(err, &apiErr), "error is not an api.Error: %v", err)
	require.Equal(t, nethttp.StatusNotFound, apiErr.StatusCode)
}

// failingForkChoiceV2Client returns a client whose fork choice v2 requests fail with a server error.
func failingForkChoiceV2Client(t *testing.T, name string) consensusclient.Service {
	t.Helper()

	client, err := mock.New(context.Background(), mock.WithName(name))
	require.NoError(t, err)
	client.ForkChoiceV2Func = func(context.Context, *api.ForkChoiceOpts) (*api.Response[*apiv1.ForkChoiceV2], error) {
		return nil, &api.Error{Method: nethttp.MethodGet, Endpoint: "/eth/v2/debug/fork_choice", StatusCode: nethttp.StatusInternalServerError}
	}

	return client
}

// TestForkChoiceV2UnsupportedAndFailing ensures that a pool in which some clients lack the endpoint and the others
// fail returns the unsupported error, so that callers fall back to v1, which the former can serve, along with the
// other clients' error, whichever order they are in.
func TestForkChoiceV2UnsupportedAndFailing(t *testing.T) {
	tests := []struct {
		name    string
		clients func(t *testing.T) []consensusclient.Service
	}{
		{
			name: "UnsupportedFirst",
			clients: func(t *testing.T) []consensusclient.Service {
				t.Helper()

				return []consensusclient.Service{
					unsupportedForkChoiceV2Client(t, "mock 1", nethttp.StatusBadRequest),
					failingForkChoiceV2Client(t, "mock 2"),
				}
			},
		},
		{
			name: "FailingFirst",
			clients: func(t *testing.T) []consensusclient.Service {
				t.Helper()

				return []consensusclient.Service{
					failingForkChoiceV2Client(t, "mock 1"),
					unsupportedForkChoiceV2Client(t, "mock 2", nethttp.StatusBadRequest),
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			multiClient, err := multi.New(ctx,
				multi.WithLogLevel(zerolog.Disabled),
				multi.WithClients(test.clients(t)),
			)
			require.NoError(t, err)

			_, err = multiClient.(consensusclient.ForkChoiceV2Provider).ForkChoiceV2(ctx, &api.ForkChoiceOpts{})
			require.Error(t, err)

			var apiErr *api.Error
			require.True(t, errors.As(err, &apiErr), "error is not an api.Error: %v", err)
			require.Equal(t, nethttp.StatusBadRequest, apiErr.StatusCode)
			require.ErrorContains(t, err, "500")
			require.NotContains(t, err.Error(), "empty response")
		})
	}
}
