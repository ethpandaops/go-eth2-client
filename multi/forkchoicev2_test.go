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
	"fmt"
	nethttp "net/http"
	"sync/atomic"
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

// countingForkChoiceV2Client returns a client whose fork choice v2 requests fail with err, counting calls.
func countingForkChoiceV2Client(t *testing.T, name string, calls *atomic.Int32, err error) consensusclient.Service {
	t.Helper()

	client, newErr := mock.New(context.Background(), mock.WithName(name))
	require.NoError(t, newErr)
	client.ForkChoiceV2Func = func(context.Context, *api.ForkChoiceOpts) (*api.Response[*apiv1.ForkChoiceV2], error) {
		calls.Add(1)

		return nil, err
	}

	return client
}

func unsupportedErr(statusCode int) error {
	return &api.Error{Method: nethttp.MethodGet, Endpoint: "/eth/v2/debug/fork_choice", StatusCode: statusCode}
}

// TestForkChoiceV2 ensures that clients without the endpoint, and clients whose response does not follow the spec, are
// skipped in favour of the next client without being deactivated: they are tried again on the next call.
func TestForkChoiceV2(t *testing.T) {
	ctx := context.Background()

	var badRequest, notFound, notAllowed, notImplemented, invalid atomic.Int32
	supporting, err := mock.New(ctx, mock.WithName("mock 6"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			countingForkChoiceV2Client(t, "mock 1", &badRequest, unsupportedErr(nethttp.StatusBadRequest)),
			countingForkChoiceV2Client(t, "mock 2", &notFound, unsupportedErr(nethttp.StatusNotFound)),
			countingForkChoiceV2Client(t, "mock 3", &notAllowed, unsupportedErr(nethttp.StatusMethodNotAllowed)),
			countingForkChoiceV2Client(t, "mock 4", &notImplemented, unsupportedErr(nethttp.StatusNotImplemented)),
			countingForkChoiceV2Client(t, "mock 5", &invalid,
				fmt.Errorf("%w: fork choice data missing", consensusclient.ErrInvalidResponse)),
			supporting,
		}),
	)
	require.NoError(t, err)

	for range 2 {
		res, err := multiClient.(consensusclient.ForkChoiceV2Provider).ForkChoiceV2(ctx, &api.ForkChoiceOpts{})
		require.NoError(t, err)
		require.NotNil(t, res)
	}

	// Each skipped client was tried on both calls, so none was deactivated.
	for name, calls := range map[string]*atomic.Int32{
		"400": &badRequest, "404": &notFound, "405": &notAllowed, "501": &notImplemented, "invalid response": &invalid,
	} {
		require.Equal(t, int32(2), calls.Load(), name)
	}
	require.Equal(t, "mock 1", multiClient.Address())
}

// TestForkChoiceV2InvalidResponses ensures that when every client returns a response that does not follow the spec,
// the error says so, rather than "empty response".
func TestForkChoiceV2InvalidResponses(t *testing.T) {
	ctx := context.Background()

	var calls atomic.Int32
	invalidErr := fmt.Errorf("%w: fork choice data missing", consensusclient.ErrInvalidResponse)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			countingForkChoiceV2Client(t, "mock 1", &calls, invalidErr),
			countingForkChoiceV2Client(t, "mock 2", &calls, invalidErr),
		}),
	)
	require.NoError(t, err)

	_, err = multiClient.(consensusclient.ForkChoiceV2Provider).ForkChoiceV2(ctx, &api.ForkChoiceOpts{})
	require.ErrorIs(t, err, consensusclient.ErrInvalidResponse)
	require.Equal(t, int32(2), calls.Load())
	require.Equal(t, "mock 1", multiClient.Address(), "clients with invalid responses stay active")
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
