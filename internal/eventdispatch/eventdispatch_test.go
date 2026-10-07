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

package eventdispatch_test

import (
	"bytes"
	"context"
	"maps"
	"reflect"
	"testing"

	"github.com/ethpandaops/go-eth2-client/api"
	apiv1 "github.com/ethpandaops/go-eth2-client/api/v1"
	"github.com/ethpandaops/go-eth2-client/internal/eventdispatch"
	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestBindsEveryTopicHandler confirms that each specific handler field of
// api.EventsOpts is bound to exactly one supported topic. The two light-client topics have only
// the generic handler. The fields are walked reflectively, so a handler added to api.EventsOpts
// but not bound fails here.
func TestBindsEveryTopicHandler(t *testing.T) {
	optsType := reflect.TypeFor[api.EventsOpts]()

	boundTopics := make(map[string]bool)
	for i := range optsType.NumField() {
		field := optsType.Field(i)
		if field.Type.Kind() != reflect.Func || field.Name == "Handler" {
			continue
		}

		t.Run(field.Name, func(t *testing.T) {
			opts := &api.EventsOpts{}
			reflect.ValueOf(opts).Elem().Field(i).Set(reflect.MakeFunc(field.Type,
				func([]reflect.Value) []reflect.Value { return nil }))

			var topics []string
			for topic := range apiv1.SupportedEventTopics {
				if eventdispatch.HasTopicHandler(opts, topic) {
					topics = append(topics, topic)
				}
			}
			require.Len(t, topics, 1, "handler is not bound to exactly one topic")
			require.NotContains(t, boundTopics, topics[0], "topic bound to more than one handler")
			boundTopics[topics[0]] = true
		})
	}

	withSpecificHandlers := maps.Clone(apiv1.SupportedEventTopics)
	delete(withSpecificHandlers, "light_client_finality_update")
	delete(withSpecificHandlers, "light_client_optimistic_update")
	delete(withSpecificHandlers, "inclusion_list")
	require.Equal(t, withSpecificHandlers, boundTopics)
}

func TestHandle(t *testing.T) {
	ctx := context.Background()
	data := []byte(`{"slot":"10","block_root":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}`)

	t.Run("Specific", func(t *testing.T) {
		var received *apiv1.ExecutionPayloadAvailableEvent
		opts := &api.EventsOpts{
			Handler: func(*apiv1.Event) { require.Fail(t, "generic handler called") },
			ExecutionPayloadAvailableHandler: func(_ context.Context, event *apiv1.ExecutionPayloadAvailableEvent) {
				received = event
			},
		}
		require.NoError(t, eventdispatch.Handle(ctx, opts, "execution_payload_available", data))
		require.NotNil(t, received)
	})

	t.Run("Generic", func(t *testing.T) {
		var received *apiv1.Event
		opts := &api.EventsOpts{Handler: func(event *apiv1.Event) { received = event }}
		require.NoError(t, eventdispatch.Handle(ctx, opts, "execution_payload_available", data))
		require.Equal(t, "execution_payload_available", received.Topic)
		require.IsType(t, &apiv1.ExecutionPayloadAvailableEvent{}, received.Data)
	})

	t.Run("NoHandler", func(t *testing.T) {
		var output bytes.Buffer
		logCtx := zerolog.New(&output).Level(zerolog.DebugLevel).WithContext(ctx)

		require.NoError(t, eventdispatch.Handle(logCtx, &api.EventsOpts{}, "execution_payload_available", data))
		require.JSONEq(t, `{"level":"debug","topic":"execution_payload_available","message":"No specific or generic handler supplied; ignoring"}`,
			output.String())
	})

	t.Run("UnsupportedTopic", func(t *testing.T) {
		err := eventdispatch.Handle(ctx, &api.EventsOpts{}, "unknown", data)
		require.ErrorIs(t, err, eventdispatch.ErrUnsupportedTopic)
		require.EqualError(t, err, "unsupported event topic unknown")
	})

	t.Run("Malformed", func(t *testing.T) {
		opts := &api.EventsOpts{Handler: func(*apiv1.Event) { require.Fail(t, "handler called") }}
		require.Error(t, eventdispatch.Handle(ctx, opts, "execution_payload_available", []byte(`invalid`)))
	})
}

func TestForwardingGuarded(t *testing.T) {
	ctx := context.Background()
	data := []byte(`{"slot":"10","block_root":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}`)

	specific := 0
	generic := 0
	opts := &api.EventsOpts{
		Topics:  []string{"execution_payload_available", "head"},
		Handler: func(*apiv1.Event) { generic++ },
		ExecutionPayloadAvailableHandler: func(context.Context, *apiv1.ExecutionPayloadAvailableEvent) {
			specific++
		},
	}

	forward := false
	var forwardedTopics []string
	filtered := eventdispatch.ForwardingGuarded(opts, func(topic string) bool {
		forwardedTopics = append(forwardedTopics, topic)

		return forward
	})

	require.NotSame(t, opts, filtered)
	require.Equal(t, opts.Topics, filtered.Topics)
	require.Nil(t, filtered.HeadHandler, "unsupplied handler was filled in")

	require.NoError(t, eventdispatch.Handle(ctx, filtered, "execution_payload_available", data))
	require.Zero(t, specific, "event forwarded although forward refused it")

	forward = true
	require.NoError(t, eventdispatch.Handle(ctx, filtered, "execution_payload_available", data))
	require.Equal(t, 1, specific)
	filtered.Handler(&apiv1.Event{Topic: "head"})
	require.Equal(t, 1, generic)
	require.Equal(t, []string{"execution_payload_available", "execution_payload_available", "head"}, forwardedTopics)

	filtered.Topics[0] = "changed"
	require.Equal(t, "execution_payload_available", opts.Topics[0], "topics shared with the original options")
}

func TestSupports(t *testing.T) {
	for topic := range apiv1.SupportedEventTopics {
		require.True(t, eventdispatch.Supports(topic), topic)
	}
	require.False(t, eventdispatch.Supports("unknown"))
	require.False(t, eventdispatch.Supports(""))
}

// TestHandleReroutesSingleAttestation pins the Grandine shape: a post-Electra
// SingleAttestation published on the attestation topic reaches the single_attestation
// handler instead of failing to decode as an aggregate, while aggregates still reach the
// attestation handler.
func TestHandleReroutesSingleAttestation(t *testing.T) {
	ctx := context.Background()
	single := []byte(`{"committee_index":"32","attester_index":"1861792","data":{"slot":"14326127","index":"0","beacon_block_root":"0x62ab3e7b1366549bd040a3523f0c2d5a2743164355ac5adbc07d51c1dfd05f88","source":{"epoch":"447690","root":"0x5d419ad7acb52b038214f39171670e4af6b877850fcb1ebf30e18d8f5bb6a372"},"target":{"epoch":"447691","root":"0xdf078bb15f0f0532eb480da59eb63a05025807b77a3d4df5cb598810974a4304"}},"signature":"0x8f5bd89876a78f836a19494d36c10e0c0d6197117f20c4354d74d06022d7e95619889a54c9feadf2c4160f93aee02bb41606b8e70e296d2e39abd1fb67bca01d019ffe6a9714ac7f99273a8acaa91c7f7ce518091051f95134b6017449503472"}`)
	aggregate := []byte(`{"aggregation_bits":"0x00002840403040000000020008800040008042800000020220","data":{"slot":"4095945","index":"12","beacon_block_root":"0xff27c7551bf4cfe4dc4cce00920e7a5c5074860d1dbd8aa8b3b5f888523f51ff","source":{"epoch":"127997","root":"0x38758fb180459583bd5e8e1a31711eb09e63eb92be974485397e9a2c57de2783"},"target":{"epoch":"127998","root":"0x46d4629861bd81cfc94007501b4edb1b3ca9444b41d7a98681b6c2f4bdb978bd"}},"signature":"0xacb9f562a28c4ef5b60b88678068ea51573a3237d3331dda3b2d845a0d03bc56ab2994d2deb90d9f074a8bdab59945150d0a7717e74b1bf2627f8971c81091f724c211dfce8fa16fb839c6a1bfd341ddec5e7eb88472682fd1a170e373660534"}`)

	t.Run("SingleToSpecificHandler", func(t *testing.T) {
		var received *electra.SingleAttestation
		opts := &api.EventsOpts{
			AttestationHandler: func(context.Context, *spec.VersionedAttestation) {
				require.Fail(t, "attestation handler called")
			},
			SingleAttestationHandler: func(_ context.Context, attestation *electra.SingleAttestation) {
				received = attestation
			},
		}
		require.NoError(t, eventdispatch.Handle(ctx, opts, "attestation", single))
		require.NotNil(t, received)
		require.Equal(t, phase0.ValidatorIndex(1861792), received.AttesterIndex)
	})

	t.Run("SingleToGenericHandler", func(t *testing.T) {
		var received *apiv1.Event
		opts := &api.EventsOpts{Handler: func(event *apiv1.Event) { received = event }}
		require.NoError(t, eventdispatch.Handle(ctx, opts, "attestation", single))
		require.Equal(t, "single_attestation", received.Topic)
		require.IsType(t, &electra.SingleAttestation{}, received.Data)
	})

	t.Run("AggregateNotRerouted", func(t *testing.T) {
		var received *spec.VersionedAttestation
		opts := &api.EventsOpts{
			AttestationHandler: func(_ context.Context, attestation *spec.VersionedAttestation) {
				received = attestation
			},
			SingleAttestationHandler: func(context.Context, *electra.SingleAttestation) {
				require.Fail(t, "single attestation handler called")
			},
		}
		require.NoError(t, eventdispatch.Handle(ctx, opts, "attestation", aggregate))
		require.NotNil(t, received)
	})
}
