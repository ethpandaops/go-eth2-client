package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/ethpandaops/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestPayloadAttributesEventJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.payloadAttributesEventJSON",
		},
		{
			name:  "VersionMissing",
			input: []byte(`{"data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "unsupported data version",
		},
		{
			name:  "VersionInvalid",
			input: []byte(`{"version":"invalid","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: unrecognised data version \"invalid\"",
		},
		{
			name:  "ProposerIndexMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "proposer index missing",
		},
		{
			name:  "ProposerIndexWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":123,"proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field payloadAttributesDataJSON.data.proposer_index of type string",
		},
		{
			name:  "ProposerIndexInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"invalid","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for proposer index: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ProposerSlotMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "proposal slot missing",
		},
		{
			name:  "ProposerSlotWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":10,"parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field payloadAttributesDataJSON.data.proposal_slot of type string",
		},
		{
			name:  "ProposerSlotInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"invalid","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for proposal slot: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ParentBlockNumberMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "parent block number missing",
		},
		{
			name:  "ParentBlockNumberWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":9,"parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field payloadAttributesDataJSON.data.parent_block_number of type string",
		},
		{
			name:  "ParentBlockNumberInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"invalid","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for parent block number: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ParentBlockRootMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "parent block root missing",
		},
		{
			name:  "ParentBlockRootWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":true,"parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesDataJSON.data.parent_block_root of type string",
		},
		{
			name:  "ParentBlockRootInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"invalid","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for parent block root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "ParentBlockHashMissing",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "parent block hash missing",
		},
		{
			name:  "ParentBlockHashWrongType",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":true,"payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesDataJSON.data.parent_block_hash of type string",
		},
		{
			name:  "ParentBlockHashInvalid",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"invalid","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "invalid value for parent block hash: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BadPayloadAttributesV1Data",
			input: []byte(`{"version":"bellatrix","data":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesEventJSON.data of type v1.payloadAttributesDataJSON",
		},
		{
			name:  "BadPayloadAttributesV2Data",
			input: []byte(`{"version":"capella","data":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesEventJSON.data of type v1.payloadAttributesDataJSON",
		},
		{
			name:  "BadPayloadAttributesV3Data",
			input: []byte(`{"version":"deneb","data":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field payloadAttributesEventJSON.data of type v1.payloadAttributesDataJSON",
		},
		{
			name:  "BadPayloadAttributesV1",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":true}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go value of type v1.payloadAttributesV1JSON",
		},
		{
			name:  "BadPayloadAttributesV2",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":true}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go value of type v1.payloadAttributesV2JSON",
		},
		{
			name:  "BadPayloadAttributesV3",
			input: []byte(`{"version":"deneb","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":true}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go value of type v1.payloadAttributesV3JSON",
		},
		{
			name:  "BadPayloadAttributesV1",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{}}}`),
			err:   "payload attributes timestamp missing",
		},
		{
			name:  "BadPayloadAttributesV2",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{}}}`),
			err:   "payload attributes timestamp missing",
		},
		{
			name:  "BadPayloadAttributesV3",
			input: []byte(`{"version":"deneb","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{}}}`),
			err:   "payload attributes timestamp missing",
		},
		{
			name:  "MissingPayloadAttributes",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}}`),
			err:   "payload attributes missing",
		},
		{
			name:  "MissingWithdrawals",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			err:   "payload attributes withdrawals missing",
		},
		{
			name:  "NullWithdrawalV2",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"},null]}}}`),
			err:   "withdrawals entry 1 missing",
		},
		{
			name:  "NullWithdrawalV3",
			input: []byte(`{"version":"deneb","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"},null],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df"}}}`),
			err:   "withdrawals entry 1 missing",
		},
		{
			name:  "GoodPayloadAttributesV1",
			input: []byte(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
		},
		{
			name:  "GoodPayloadAttributesV2",
			input: []byte(`{"version":"capella","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}]}}}`),
		},
		{
			name:  "GoodPayloadAttributesV3",
			input: []byte(`{"version":"deneb","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df"}}}`),
		},
		{
			name:  "GoodPayloadAttributesV3Electra",
			input: []byte(`{"version":"electra","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df"}}}`),
		},
		{
			name:  "GoodPayloadAttributesV3Fulu",
			input: []byte(`{"version":"fulu","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df"}}}`),
		},
		{
			name:  "ParentBlockNumberMissingFulu",
			input: []byte(`{"version":"fulu","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df"}}}`),
			err:   "parent block number missing",
		},
		{
			name:  "GoodPayloadAttributesV4",
			input: []byte(`{"version":"gloas","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","safe_block_hash":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505","finalized_block_hash":"0x2c77bd2fc774dacd6a6a957e7fd16456ce019feb84f781bf999eb52bf2826616","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10","target_gas_limit":"60000000"}}}`),
		},
		{
			name:  "SlotNumberMissingV4",
			input: []byte(`{"version":"gloas","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","target_gas_limit":"60000000"}}}`),
			err:   "payload attributes slot number missing",
		},
		{
			name:  "SlotNumberInvalidV4",
			input: []byte(`{"version":"gloas","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"invalid","target_gas_limit":"60000000"}}}`),
			err:   "invalid value for payload attributes slot number: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "TargetGasLimitMissingV4",
			input: []byte(`{"version":"gloas","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10"}}}`),
			err:   "payload attributes target gas limit missing",
		},
		{
			name:  "TargetGasLimitInvalidV4",
			input: []byte(`{"version":"gloas","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10","target_gas_limit":"invalid"}}}`),
			err:   "invalid value for payload attributes target gas limit: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ParentBlockNumberInvalidV4",
			input: []byte(`{"version":"gloas","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"invalid","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10","target_gas_limit":"60000000"}}}`),
			err:   "invalid value for parent block number: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "GoodPayloadAttributesV5",
			input: []byte(`{"version":"heze","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","safe_block_hash":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505","finalized_block_hash":"0x2c77bd2fc774dacd6a6a957e7fd16456ce019feb84f781bf999eb52bf2826616","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10","target_gas_limit":"60000000","inclusion_list_transactions":["0x02f870","0x01"]}}}`),
		},
		{
			name:  "GoodPayloadAttributesV5EmptyInclusionList",
			input: []byte(`{"version":"heze","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","safe_block_hash":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505","finalized_block_hash":"0x2c77bd2fc774dacd6a6a957e7fd16456ce019feb84f781bf999eb52bf2826616","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10","target_gas_limit":"60000000","inclusion_list_transactions":[]}}}`),
		},
		{
			name:  "InclusionListTransactionsMissingV5",
			input: []byte(`{"version":"heze","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10","target_gas_limit":"60000000"}}}`),
			err:   "payload attributes inclusion list transactions missing",
		},
		{
			name:  "InclusionListTransactionInvalidV5",
			input: []byte(`{"version":"heze","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10","target_gas_limit":"60000000","inclusion_list_transactions":["0xzz"]}}}`),
			err:   "invalid value for inclusion list transactions entry 0: encoding/hex: invalid byte: U+007A 'z'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.PayloadAttributesEvent
			err := json.Unmarshal(test.input, &res)

			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}

// TestPayloadAttributesEventForkchoiceHashesGloas checks that the safe and
// finalized block hashes are read from Gloas onwards, and that beacon nodes
// which do not send them yet are tolerated.
func TestPayloadAttributesEventForkchoiceHashesGloas(t *testing.T) {
	const attrs = `"payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10","target_gas_limit":"60000000"}`
	const prefix = `{"version":"gloas","data":{"proposer_index":"123","proposal_slot":"10","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf",`

	var res api.PayloadAttributesEvent
	require.NoError(t, json.Unmarshal([]byte(prefix+`"safe_block_hash":"0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505","finalized_block_hash":"0x2c77bd2fc774dacd6a6a957e7fd16456ce019feb84f781bf999eb52bf2826616",`+attrs+`}}`), &res))
	require.NotNil(t, res.Data.SafeBlockHash)
	require.NotNil(t, res.Data.FinalizedBlockHash)
	assert.Equal(t, "0x1b66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505", res.Data.SafeBlockHash.String())
	assert.Equal(t, "0x2c77bd2fc774dacd6a6a957e7fd16456ce019feb84f781bf999eb52bf2826616", res.Data.FinalizedBlockHash.String())

	var missing api.PayloadAttributesEvent
	missingInput := []byte(prefix + attrs + `}}`)
	require.NoError(t, json.Unmarshal(missingInput, &missing))
	assert.Nil(t, missing.Data.SafeBlockHash)
	assert.Nil(t, missing.Data.FinalizedBlockHash)

	// Absent hashes must not be re-emitted as zero hashes.
	rt, err := json.Marshal(&missing)
	require.NoError(t, err)
	assert.Equal(t, string(missingInput), string(rt))

	var short api.PayloadAttributesEvent
	require.EqualError(t, json.Unmarshal([]byte(prefix+`"safe_block_hash":"0x66ac1fb663c9bc59509846d6ec05345bd908eda73e670af888da41af171505",`+attrs+`}}`), &short), "incorrect length for safe block hash")

	var invalid api.PayloadAttributesEvent
	require.EqualError(t, json.Unmarshal([]byte(prefix+`"finalized_block_hash":"invalid",`+attrs+`}}`), &invalid), "invalid value for finalized block hash: encoding/hex: invalid byte: U+0069 'i'")
}

// TestPayloadAttributesEventParentBlockNumberGloas checks that a parent block number
// still sent by a beacon node from Gloas onwards is read but not marshaled.
func TestPayloadAttributesEventParentBlockNumberGloas(t *testing.T) {
	input := []byte(`{"version":"gloas","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000","withdrawals":[{"index":"5","validator_index":"10","address":"0x0000000000000000000000000000000000000000","amount":"15640"}],"parent_beacon_block_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","slot_number":"10","target_gas_limit":"60000000"}}}`)

	var res api.PayloadAttributesEvent
	require.NoError(t, json.Unmarshal(input, &res))
	assert.Equal(t, uint64(9), res.Data.ParentBlockNumber)
	require.NotNil(t, res.Data.V4)
	assert.Equal(t, uint64(10), res.Data.V4.SlotNumber)
	assert.Equal(t, uint64(60000000), res.Data.V4.TargetGasLimit)

	rt, err := json.Marshal(&res)
	require.NoError(t, err)
	assert.NotContains(t, string(rt), "parent_block_number")
}
