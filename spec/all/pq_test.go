package all_test

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec/all"
	"github.com/ethpandaops/go-eth2-client/spec/pq"
	"github.com/ethpandaops/go-eth2-client/spec/version"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/stretchr/testify/require"
)

func pqMinimalSpecs() map[string]any {
	return map[string]any{
		"SLOTS_PER_HISTORICAL_ROOT":     uint64(64),
		"EPOCHS_PER_HISTORICAL_VECTOR":  uint64(64),
		"EPOCHS_PER_SLASHINGS_VECTOR":   uint64(64),
		"EPOCHS_PER_ETH1_VOTING_PERIOD": uint64(4),
		"SLOTS_PER_EPOCH":               uint64(8),
		"MIN_SEED_LOOKAHEAD":            uint64(1),
		"PTC_SIZE":                      uint64(16),
		"MAX_COMMITTEES_PER_SLOT_PQ":    uint64(4),
		"MAX_STARK_PROOF_SIZE":          uint64(1048576),
	}
}

func TestPQBeaconStateView(t *testing.T) {
	ds := dynssz.NewDynSsz(pqMinimalSpecs())

	encoded, err := os.ReadFile("../pq/testdata/genesis_minimal.ssz")
	require.NoError(t, err)

	state := &all.BeaconState{Version: version.DataVersionPQ}
	require.NoError(t, ds.UnmarshalSSZ(state, encoded))
	require.Len(t, state.XMSSValidators, 8)
	require.Len(t, state.Validators, 8)
	require.Len(t, state.Builders, 1)
	require.Equal(t, state.XMSSValidators[3].PublicKey[:], state.Validators[3].PublicKey[pq.DepositPrefixSize:])
	require.Equal(t, make([]byte, pq.DepositPrefixSize), state.Validators[3].PublicKey[:pq.DepositPrefixSize])
	require.Len(t, state.AggregationCommitteeWindow, 24)

	root, err := ds.HashTreeRoot(state)
	require.NoError(t, err)
	require.Equal(t, "cc0e5639191cf4d1b7d459b8e1d3a7d6aa86eaa37cc63b6f02a9ea6145c15736", hex.EncodeToString(root[:]))

	reencoded, err := ds.MarshalSSZ(state)
	require.NoError(t, err)
	require.Equal(t, encoded, reencoded)

	marshaled, err := json.Marshal(state)
	require.NoError(t, err)

	fromJSON := &all.BeaconState{Version: version.DataVersionPQ}
	require.NoError(t, json.Unmarshal(marshaled, fromJSON))

	jsonRoot, err := ds.HashTreeRoot(fromJSON)
	require.NoError(t, err)
	require.Equal(t, root, jsonRoot)
	require.Len(t, fromJSON.Validators, 8)
}

func TestPQSignedBeaconBlockView(t *testing.T) {
	ds := dynssz.NewDynSsz(pqMinimalSpecs())

	block := &pq.SignedBeaconBlock{
		Message: &pq.BeaconBlock{
			Slot: 5,
			Body: &pq.BeaconBlockBody{
				Attestations:          []*pq.Attestation{},
				XMSSProposerSlashings: []*pq.ProposerSlashing{},
				AttestationProof:      make([]byte, pq.XMSSSignatureSize*2),
			},
		},
	}
	block.XMSSSignature[0] = 0x2a
	block.Message.Body.VRFReveal[0][0] = 0x01

	encoded, err := ds.MarshalSSZ(block)
	require.NoError(t, err)

	root, err := ds.HashTreeRoot(block.Message)
	require.NoError(t, err)

	union := &all.SignedBeaconBlock{Version: version.DataVersionPQ}
	require.NoError(t, ds.UnmarshalSSZ(union, encoded))
	require.Equal(t, uint32(0x2a), union.XMSSSignature.LeafIndex())
	require.Len(t, union.Message.Body.AttestationProof, 2*pq.XMSSSignatureSize)

	t.Logf("block root %x", root)

	unionRoot, err := ds.HashTreeRoot(union.Message)
	require.NoError(t, err)
	require.Equal(t, root, unionRoot)

	reencoded, err := ds.MarshalSSZ(union)
	require.NoError(t, err)
	require.Equal(t, encoded, reencoded)

	marshaled, err := json.Marshal(union)
	require.NoError(t, err)

	fromJSON := &pq.SignedBeaconBlock{}
	require.NoError(t, json.Unmarshal(marshaled, fromJSON))

	jsonRoot, err := ds.HashTreeRoot(fromJSON.Message)
	require.NoError(t, err)
	require.Equal(t, root, jsonRoot)
}
