package version_test

import (
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec/version"
	"github.com/stretchr/testify/require"
)

func TestDataVersionFromStringEIP8198(t *testing.T) {
	// EIP-8198 changes no containers, so its data decodes as Heze.
	v, err := version.DataVersionFromString("eip8198")
	require.NoError(t, err)
	require.Equal(t, version.DataVersionHeze, v)

	v, err = version.DataVersionFromString("heze")
	require.NoError(t, err)
	require.Equal(t, version.DataVersionHeze, v)

	_, err = version.DataVersionFromString("eip9999")
	require.Error(t, err)
}
