package gloas_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/golang/snappy"
	"github.com/pk910/dynamic-ssz/sszutils"
	require "github.com/stretchr/testify/require"
)

// sszHTRProbe validates SSZ round-trip + HashTreeRoot (no YAML) for a type.
func sszHTRProbe(t *testing.T, typeName string, factory func() sszutils.FastsszMarshaler) {
	base := filepath.Join(os.Getenv("CONSENSUS_SPEC_TESTS_DIR"), "tests", "mainnet", "gloas", "ssz_static", typeName, "ssz_random")
	entries, err := os.ReadDir(base)
	require.NoError(t, err)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(base, e.Name())
		t.Run(typeName+"/"+e.Name(), func(t *testing.T) {
			comp, err := os.ReadFile(filepath.Join(path, "serialized.ssz_snappy"))
			require.NoError(t, err)
			specSSZ, err := snappy.Decode(nil, comp)
			require.NoError(t, err)

			s := factory()
			require.NoError(t, s.(sszutils.FastsszUnmarshaler).UnmarshalSSZ(specSSZ))
			out, err := s.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, specSSZ, out, "SSZ bytes mismatch")

			rootYAML, err := os.ReadFile(filepath.Join(path, "roots.yaml"))
			require.NoError(t, err)
			rootBytes, err := s.(sszutils.FastsszHashRoot).HashTreeRoot()
			require.NoError(t, err)
			got := fmt.Sprintf("{root: '%#x'}\n", string(rootBytes[:]))
			require.YAMLEq(t, string(rootYAML), got, "HASH TREE ROOT mismatch")
		})
	}
}

func TestSSZHTRProbe(t *testing.T) {
	if os.Getenv("CONSENSUS_SPEC_TESTS_DIR") == "" {
		t.Skip("no dir")
	}
	sszHTRProbe(t, "Builder", func() sszutils.FastsszMarshaler { return &gloas.Builder{} })
	sszHTRProbe(t, "PayloadAttestationData", func() sszutils.FastsszMarshaler { return &gloas.PayloadAttestationData{} })
	sszHTRProbe(t, "ProposerPreferences", func() sszutils.FastsszMarshaler { return &gloas.ProposerPreferences{} })
	sszHTRProbe(t, "BuilderPendingPayment", func() sszutils.FastsszMarshaler { return &gloas.BuilderPendingPayment{} })
	sszHTRProbe(t, "ExecutionPayloadBid", func() sszutils.FastsszMarshaler { return &gloas.ExecutionPayloadBid{} })
	sszHTRProbe(t, "PayloadAttestation", func() sszutils.FastsszMarshaler { return &gloas.PayloadAttestation{} })
	sszHTRProbe(t, "IndexedPayloadAttestation", func() sszutils.FastsszMarshaler { return &gloas.IndexedPayloadAttestation{} })
	sszHTRProbe(t, "ExecutionPayload", func() sszutils.FastsszMarshaler { return &gloas.ExecutionPayload{} })
	sszHTRProbe(t, "Attestation", func() sszutils.FastsszMarshaler { return &gloas.Attestation{} })
	sszHTRProbe(t, "IndexedAttestation", func() sszutils.FastsszMarshaler { return &gloas.IndexedAttestation{} })
}
