// Copyright © 2020 Attestant Limited.
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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/ethpandaops/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestChainReorgEventJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		// expected is the re-marshaled form when it differs from input.
		expected []byte
		err      string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.chainReorgEventJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field chainReorgEventJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "DepthMissing",
			input: []byte(`{"slot":"524986","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "depth missing",
		},
		{
			name:  "DepthWrongType",
			input: []byte(`{"slot":"524986","depth":true,"old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field chainReorgEventJSON.depth of type string",
		},
		{
			name:  "DepthInvalid",
			input: []byte(`{"slot":"524986","depth":"-1","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid value for depth: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "OldHeadBlockMissing",
			input: []byte(`{"slot":"524986","depth":"2","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "old head block missing",
		},
		{
			name:  "OldHeadBlockWrongType",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":true,"new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field chainReorgEventJSON.old_head_block of type string",
		},
		{
			name:  "OldHeadBlockInvalid",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"invalid","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid value for old head block: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "OldHeadBlockShort",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0xfc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "incorrect length 31 for old head block",
		},
		{
			name:  "OldHeadBlock",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2f2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "incorrect length 33 for old head block",
		},
		{
			name:  "NewHeadBlockMissing",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "new head block missing",
		},
		{
			name:  "NewHeadBlockWrongType",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":true,"old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field chainReorgEventJSON.new_head_block of type string",
		},
		{
			name:  "NewHeadBlockInvalid",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"invalid","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid value for new head block: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "NewHeadBlockShort",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xfe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "incorrect length 31 for new head block",
		},
		{
			name:  "NewHeadBlockLong",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3a3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "incorrect length 33 for new head block",
		},
		{
			name:  "OldHeadMissing",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "old head state missing",
		},
		{
			name:  "OldHeadStateWrongType",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":true,"new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field chainReorgEventJSON.old_head_state of type string",
		},
		{
			name:  "OldHeadStateInvalid",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"invalid","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "invalid value for old head state: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "OldHeadStateShort",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0xcc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "incorrect length 31 for old head state",
		},
		{
			name:  "OldHeadStateLong",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x9797cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "incorrect length 33 for old head state",
		},
		{
			name:  "NewHeadStateMissing",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","epoch":"16405"}`),
			err:   "new head state missing",
		},
		{
			name:  "NewHeadStateWrongType",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":true,"epoch":"16405"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field chainReorgEventJSON.new_head_state of type string",
		},
		{
			name:  "NewHeadStateInvalid",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"invalid","epoch":"16405"}`),
			err:   "invalid value for new head state: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "NewHeadStateShort",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0xb800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "incorrect length 31 for new head state",
		},
		{
			name:  "NewHeadStateLong",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4a4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
			err:   "incorrect length 33 for new head state",
		},
		{
			// epoch is optional (Nimbus omits it) and left at zero.
			name:     "EpochMissing",
			input:    []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522"}`),
			expected: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"0"}`),
		},
		{
			// Verbatim Nimbus v26.9.x payload from a glamsterdam-devnet-8 sentry
			// log ("Failed to parse chain reorg event ... epoch missing"): no epoch,
			// plus execution_optimistic which the event does not carry.
			name:     "NimbusNoEpoch",
			input:    []byte(`{"slot":"343871","depth":"2","old_head_block":"0x9af8248c506f8635c5eb734d45508854ef8a1668fcece9813c0a8e2a1ace4b73","new_head_block":"0x0fec5a8fb2e9a13b5393db10bd9df1caaef2774e0b3083a19e5784402c8baa25","old_head_state":"0x405653c3ebba8072ef70a6ab7588f64f1b6a32495810139928c0ac1ecb6ecbc3","new_head_state":"0xc7b8f3b19ba14bdbae92a322c695583ee19dba5cb3ff466114c9c5ca8011ab56","execution_optimistic":false}`),
			expected: []byte(`{"slot":"343871","depth":"2","old_head_block":"0x9af8248c506f8635c5eb734d45508854ef8a1668fcece9813c0a8e2a1ace4b73","new_head_block":"0x0fec5a8fb2e9a13b5393db10bd9df1caaef2774e0b3083a19e5784402c8baa25","old_head_state":"0x405653c3ebba8072ef70a6ab7588f64f1b6a32495810139928c0ac1ecb6ecbc3","new_head_state":"0xc7b8f3b19ba14bdbae92a322c695583ee19dba5cb3ff466114c9c5ca8011ab56","epoch":"0"}`),
		},
		{
			// Spec example shape (every non-Nimbus CL): epoch present.
			name:     "SpecExample",
			input:    []byte(`{"slot":"200", "depth":"50", "old_head_block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf", "new_head_block":"0x76262e91970d375a19bfe8a867288d7b9cde43c8635f598d93d39d041706fc76", "old_head_state":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf", "new_head_state":"0x600e852a08c1200654ddf11025f1ceacb3c2e74bdd5c630cde0838b2591b69f9", "epoch":"2", "execution_optimistic": false}`),
			expected: []byte(`{"slot":"200","depth":"50","old_head_block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","new_head_block":"0x76262e91970d375a19bfe8a867288d7b9cde43c8635f598d93d39d041706fc76","old_head_state":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","new_head_state":"0x600e852a08c1200654ddf11025f1ceacb3c2e74bdd5c630cde0838b2591b69f9","epoch":"2"}`),
		},
		{
			name:  "EpochWrongType",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field chainReorgEventJSON.epoch of type string",
		},
		{
			name:  "EpochInvalid",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"-1"}`),
			err:   "invalid value for epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ChainReorgEvent
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				expected := test.expected
				if expected == nil {
					expected = test.input
				}
				assert.Equal(t, string(expected), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}

// TestChainReorgEventEpochReset checks that decoding a payload without epoch
// into a previously populated event does not leak the old epoch.
func TestChainReorgEventEpochReset(t *testing.T) {
	var res api.ChainReorgEvent
	require.NoError(t, json.Unmarshal([]byte(`{"slot":"524986","depth":"2","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","epoch":"16405"}`), &res))
	require.EqualValues(t, 16405, res.Epoch)

	require.NoError(t, json.Unmarshal([]byte(`{"slot":"343871","depth":"2","old_head_block":"0x9af8248c506f8635c5eb734d45508854ef8a1668fcece9813c0a8e2a1ace4b73","new_head_block":"0x0fec5a8fb2e9a13b5393db10bd9df1caaef2774e0b3083a19e5784402c8baa25","old_head_state":"0x405653c3ebba8072ef70a6ab7588f64f1b6a32495810139928c0ac1ecb6ecbc3","new_head_state":"0xc7b8f3b19ba14bdbae92a322c695583ee19dba5cb3ff466114c9c5ca8011ab56","execution_optimistic":false}`), &res))
	require.EqualValues(t, 343871, res.Slot)
	require.EqualValues(t, 0, res.Epoch)
}
