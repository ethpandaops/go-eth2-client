// Copyright © 2020 - 2024 Attestant Limited.
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

package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	client "github.com/ethpandaops/go-eth2-client"
	"github.com/ethpandaops/go-eth2-client/api"
	"github.com/ethpandaops/go-eth2-client/spec"
	"github.com/ethpandaops/go-eth2-client/spec/all"
	"github.com/ethpandaops/go-eth2-client/spec/altair"
	"github.com/ethpandaops/go-eth2-client/spec/bellatrix"
	"github.com/ethpandaops/go-eth2-client/spec/capella"
	"github.com/ethpandaops/go-eth2-client/spec/deneb"
	"github.com/ethpandaops/go-eth2-client/spec/electra"
	"github.com/ethpandaops/go-eth2-client/spec/fulu"
	"github.com/ethpandaops/go-eth2-client/spec/gloas"
	"github.com/ethpandaops/go-eth2-client/spec/heze"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
)

// BeaconState fetches a beacon state given a state ID and decodes it directly
// into the per-fork view stored on a *spec.VersionedBeaconState. SSZ
// responses are streamed straight from the wire into the SSZ decoder so the
// (potentially multi-GB) payload is never held in memory in full.
func (s *Service) BeaconState(ctx context.Context,
	opts *api.BeaconStateOpts,
) (
	*api.Response[*spec.VersionedBeaconState],
	error,
) {
	httpResponse, err := s.fetchBeaconState(ctx, opts)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Close()

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		return s.beaconStateFromSSZ(ctx, httpResponse)
	case ContentTypeJSON:
		return s.beaconStateFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
}

// AgnosticBeaconState fetches a beacon state and decodes it directly into a
// fork-agnostic *all.BeaconState. The Version is set from the consensus
// version header before unmarshaling so the union type's view-aware codec
// dispatches into the correct fork's schema. SSZ responses are streamed
// from the wire into the SSZ decoder; no intermediate copy.
func (s *Service) AgnosticBeaconState(ctx context.Context,
	opts *api.BeaconStateOpts,
) (
	*api.Response[*all.BeaconState],
	error,
) {
	httpResponse, err := s.fetchBeaconState(ctx, opts)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Close()

	state := &all.BeaconState{Version: httpResponse.consensusVersion}
	metadata := metadataFromHeaders(httpResponse.headers)

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		ds, err := s.dynSSZForRequest(ctx)
		if err != nil {
			return nil, err
		}

		if httpResponse.bodyReader != nil {
			if err := ds.UnmarshalSSZReader(state, httpResponse.bodyReader, int(httpResponse.bodySize)); err != nil {
				return nil, errors.Join(fmt.Errorf("failed to decode %s beacon state", httpResponse.consensusVersion), err)
			}
		} else {
			if err := ds.UnmarshalSSZ(state, httpResponse.body); err != nil {
				return nil, errors.Join(fmt.Errorf("failed to decode %s beacon state", httpResponse.consensusVersion), err)
			}
		}
	case ContentTypeJSON:
		if httpResponse.bodyReader != nil {
			decoded, jsonMetadata, err := decodeJSONResponse(httpResponse.bodyReader, state)
			if err != nil {
				return nil, errors.Join(fmt.Errorf("failed to decode %s beacon state", httpResponse.consensusVersion), err)
			}
			state = decoded
			metadata = jsonMetadata
		} else {
			decoded, jsonMetadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), state)
			if err != nil {
				return nil, errors.Join(fmt.Errorf("failed to decode %s beacon state", httpResponse.consensusVersion), err)
			}
			state = decoded
			metadata = jsonMetadata
		}
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}

	return &api.Response[*all.BeaconState]{
		Data:     state,
		Metadata: metadata,
	}, nil
}

// fetchBeaconState performs the GET request shared by BeaconState and
// AgnosticBeaconState: validates opts and hits the endpoint. The response is
// fetched via getStream so SSZ payloads are not buffered into memory; the
// caller is responsible for invoking httpResponse.Close.
func (s *Service) fetchBeaconState(ctx context.Context,
	opts *api.BeaconStateOpts,
) (*httpResponse, error) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if opts.State == "" {
		return nil, errors.Join(errors.New("no state specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v2/debug/beacon/states/%s", opts.State)

	if opts.Stream {
		return s.getStream(ctx, endpoint, "", &opts.Common, true)
	}

	return s.get(ctx, endpoint, "", &opts.Common, true)
}

func (s *Service) beaconStateFromSSZ(ctx context.Context, res *httpResponse) (*api.Response[*spec.VersionedBeaconState], error) {
	response := &api.Response[*spec.VersionedBeaconState]{
		Data: &spec.VersionedBeaconState{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	dynSSZ, err := s.dynSSZForRequest(ctx)
	if err != nil {
		return nil, err
	}

	size := int(res.bodySize)

	var target any
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0 = &phase0.BeaconState{}
		target = response.Data.Phase0
	case spec.DataVersionAltair:
		response.Data.Altair = &altair.BeaconState{}
		target = response.Data.Altair
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix = &bellatrix.BeaconState{}
		target = response.Data.Bellatrix
	case spec.DataVersionCapella:
		response.Data.Capella = &capella.BeaconState{}
		target = response.Data.Capella
	case spec.DataVersionDeneb:
		response.Data.Deneb = &deneb.BeaconState{}
		target = response.Data.Deneb
	case spec.DataVersionElectra:
		response.Data.Electra = &electra.BeaconState{}
		target = response.Data.Electra
	case spec.DataVersionFulu:
		response.Data.Fulu = &fulu.BeaconState{}
		target = response.Data.Fulu
	case spec.DataVersionGloas:
		response.Data.Gloas = &gloas.BeaconState{}
		target = response.Data.Gloas
	case spec.DataVersionHeze:
		response.Data.Heze = &heze.BeaconState{}
		target = response.Data.Heze
	default:
		return nil, fmt.Errorf("unhandled state version %s", res.consensusVersion)
	}

	if res.bodyReader != nil {
		err = dynSSZ.UnmarshalSSZReader(target, res.bodyReader, size)
	} else {
		err = dynSSZ.UnmarshalSSZ(target, res.body)
	}
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to decode %s beacon state", res.consensusVersion), err)
	}

	return response, nil
}

func (*Service) beaconStateFromJSON(res *httpResponse) (*api.Response[*spec.VersionedBeaconState], error) {
	response := &api.Response[*spec.VersionedBeaconState]{
		Data: &spec.VersionedBeaconState{
			Version: res.consensusVersion,
		},
	}

	var err error

	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		if res.bodyReader != nil {
			response.Data.Phase0, response.Metadata, err = decodeJSONResponse(res.bodyReader, &phase0.BeaconState{})
		} else {
			response.Data.Phase0, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &phase0.BeaconState{})
		}
	case spec.DataVersionAltair:
		if res.bodyReader != nil {
			response.Data.Altair, response.Metadata, err = decodeJSONResponse(res.bodyReader, &altair.BeaconState{})
		} else {
			response.Data.Altair, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &altair.BeaconState{})
		}
	case spec.DataVersionBellatrix:
		if res.bodyReader != nil {
			response.Data.Bellatrix, response.Metadata, err = decodeJSONResponse(res.bodyReader, &bellatrix.BeaconState{})
		} else {
			response.Data.Bellatrix, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &bellatrix.BeaconState{})
		}
	case spec.DataVersionCapella:
		if res.bodyReader != nil {
			response.Data.Capella, response.Metadata, err = decodeJSONResponse(res.bodyReader, &capella.BeaconState{})
		} else {
			response.Data.Capella, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &capella.BeaconState{})
		}
	case spec.DataVersionDeneb:
		if res.bodyReader != nil {
			response.Data.Deneb, response.Metadata, err = decodeJSONResponse(res.bodyReader, &deneb.BeaconState{})
		} else {
			response.Data.Deneb, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &deneb.BeaconState{})
		}
	case spec.DataVersionElectra:
		if res.bodyReader != nil {
			response.Data.Electra, response.Metadata, err = decodeJSONResponse(res.bodyReader, &electra.BeaconState{})
		} else {
			response.Data.Electra, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &electra.BeaconState{})
		}
	case spec.DataVersionFulu:
		if res.bodyReader != nil {
			response.Data.Fulu, response.Metadata, err = decodeJSONResponse(res.bodyReader, &fulu.BeaconState{})
		} else {
			response.Data.Fulu, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &fulu.BeaconState{})
		}
	case spec.DataVersionGloas:
		if res.bodyReader != nil {
			response.Data.Gloas, response.Metadata, err = decodeJSONResponse(res.bodyReader, &gloas.BeaconState{})
		} else {
			response.Data.Gloas, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &gloas.BeaconState{})
		}
	case spec.DataVersionHeze:
		if res.bodyReader != nil {
			response.Data.Heze, response.Metadata, err = decodeJSONResponse(res.bodyReader, &heze.BeaconState{})
		} else {
			response.Data.Heze, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &heze.BeaconState{})
		}
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}
