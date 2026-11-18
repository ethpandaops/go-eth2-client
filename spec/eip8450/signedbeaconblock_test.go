package eip8450_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec/eip8450"
	"github.com/stretchr/testify/require"
)

func TestSignedBeaconBlockJSON(t *testing.T) {
	input, err := os.ReadFile("testdata/signedbeaconblock.json")
	require.NoError(t, err)

	var block eip8450.SignedBeaconBlock
	require.NoError(t, json.Unmarshal(input, &block))
	require.NotNil(t, block.Message.Body.PQKeyRegistrations)

	output, err := json.Marshal(&block)
	require.NoError(t, err)
	require.JSONEq(t, string(input), string(output))
}

func TestBeaconStateRegistryJSON(t *testing.T) {
	material := &eip8450.PQKeyMaterial{}
	material.XMSSPublicKey[0] = 0x11
	material.VRFRoot[31] = 0x22

	encoded, err := json.Marshal(material)
	require.NoError(t, err)
	require.JSONEq(t, `{"xmss_public_key":"0x1100000000000000000000000000000000000000000000000000000000000000","vrf_root":"0x0000000000000000000000000000000000000000000000000000000000000022"}`, string(encoded))

	var decoded eip8450.PQKeyMaterial
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, *material, decoded)

	root, err := material.HashTreeRoot()
	require.NoError(t, err)
	require.NotEqual(t, [32]byte{}, root)

	ssz, err := material.MarshalSSZ()
	require.NoError(t, err)
	require.Len(t, ssz, 64)

	registration := &eip8450.SignedPQKeyRegistration{Message: &eip8450.PQKeyRegistration{KeyMaterial: material}}
	ssz, err = registration.MarshalSSZ()
	require.NoError(t, err)
	require.Len(t, ssz, 8+64+8+2404+96)
}
