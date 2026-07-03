package http

import (
	"testing"

	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Lighthouse (spec-compliant) wraps gloas SSE payloads in a versioned
// envelope; Prysm sends the bare object. Both must decode.
func TestUnmarshalVersionedEventData(t *testing.T) {
	bare := `{"validator_index":"325","data":{"beacon_block_root":"0x5793a892cd76b015aaaaaf8251a3061bbdb79cb7547fd926208e9beffec67e96","slot":"54715","payload_present":true,"blob_data_available":true},"signature":"0x946643c1b1d601a969fa99680e7d27701ad7c98ab4569caa6d76ed523cc74c1aa81361737ca37e42a21ae659d390bf580166086ccb6b412ad9463eceb31a1a54ac2c9b1c5238181f2ea7cb85750d332763908746f6035bfad655228a4ce45440"}`
	wrapped := `{"version":"gloas","data":` + bare + `}`

	for name, payload := range map[string]string{"bare": bare, "wrapped": wrapped} {
		t.Run(name, func(t *testing.T) {
			msg := &gloas.PayloadAttestationMessage{}
			require.NoError(t, unmarshalVersionedEventData([]byte(payload), msg))
			assert.Equal(t, uint64(325), uint64(msg.ValidatorIndex))
			assert.Equal(t, uint64(54715), uint64(msg.Data.Slot))
			assert.True(t, msg.Data.PayloadPresent)
		})
	}
}
