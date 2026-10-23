package spec_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
)

func TestVersionedSignedExecutionPayloadEnvelopeHeze(t *testing.T) {
	payload := &gloas.ExecutionPayload{BlockNumber: 42}
	requests := &gloas.ExecutionRequests{}
	root := phase0.Root{0x01}

	envelope := &spec.VersionedSignedExecutionPayloadEnvelope{
		Version: spec.DataVersionHeze,
		Gloas: &gloas.SignedExecutionPayloadEnvelope{
			Message: &gloas.ExecutionPayloadEnvelope{
				Payload:           payload,
				ExecutionRequests: requests,
				BuilderIndex:      7,
				BeaconBlockRoot:   root,
			},
		},
	}

	gotPayload, err := envelope.Payload()
	require.NoError(t, err)
	require.Equal(t, payload, gotPayload)

	gotRequests, err := envelope.ExecutionRequests()
	require.NoError(t, err)
	require.Equal(t, spec.DataVersionHeze, gotRequests.Version)
	require.Equal(t, requests, gotRequests.Heze)

	builderIndex, err := envelope.BuilderIndex()
	require.NoError(t, err)
	require.Equal(t, gloas.BuilderIndex(7), builderIndex)

	blockRoot, err := envelope.BeaconBlockRoot()
	require.NoError(t, err)
	require.Equal(t, root, blockRoot)

	unsigned := &spec.VersionedExecutionPayloadEnvelope{
		Version: spec.DataVersionHeze,
		Gloas:   envelope.Gloas.Message,
	}

	gotPayload, err = unsigned.Payload()
	require.NoError(t, err)
	require.Equal(t, payload, gotPayload)

	gotRequests, err = unsigned.ExecutionRequests()
	require.NoError(t, err)
	require.Equal(t, spec.DataVersionHeze, gotRequests.Version)
	require.Equal(t, requests, gotRequests.Heze)
}
