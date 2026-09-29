// Copyright © 2023, 2024 Attestant Limited.
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
	"context"
	"sort"
	"time"

	"github.com/ethpandaops/go-eth2-client/api"
	"github.com/ethpandaops/go-eth2-client/spec/phase0"
)

// SlotDuration provides the duration of a slot for the chain.
func (s *Service) SlotDuration(ctx context.Context) (time.Duration, error) {
	if err := s.assertIsActive(ctx); err != nil {
		return 0, err
	}

	response, err := s.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return 0, err
	}

	if res, isCorrectType := response.Data["SECONDS_PER_SLOT"].(time.Duration); isCorrectType {
		return res, nil
	}

	// SECONDS_PER_SLOT is deprecated in favour of SLOT_DURATION_MS.
	if ms, isCorrectType := response.Data["SLOT_DURATION_MS"].(uint64); isCorrectType && ms > 0 {
		return time.Duration(ms) * time.Millisecond, nil
	}

	return 0, ErrIncorrectType
}

// SlotDurationSchedule provides the EIP-8198 slot duration schedule of the
// chain, sorted by epoch and starting at epoch 0. Chains without a
// SLOT_DURATION_SCHEDULE get a single entry with the genesis slot duration.
func (s *Service) SlotDurationSchedule(ctx context.Context) ([]*api.SlotDurationScheduleEntry, error) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	response, err := s.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return nil, err
	}

	schedule := make([]*api.SlotDurationScheduleEntry, 0)

	if entries, isList := response.Data["SLOT_DURATION_SCHEDULE"].([]any); isList {
		for _, entry := range entries {
			entryMap, isMap := entry.(map[string]any)
			if !isMap {
				return nil, ErrIncorrectType
			}

			epoch, isEpoch := entryMap["EPOCH"].(uint64)
			ms, isMs := entryMap["SLOT_DURATION_MS"].(uint64)

			if !isEpoch || !isMs || ms == 0 {
				return nil, ErrIncorrectType
			}

			schedule = append(schedule, &api.SlotDurationScheduleEntry{
				Epoch:    phase0.Epoch(epoch),
				Duration: time.Duration(ms) * time.Millisecond,
			})
		}
	}

	sort.Slice(schedule, func(i, j int) bool {
		return schedule[i].Epoch < schedule[j].Epoch
	})

	if len(schedule) == 0 || schedule[0].Epoch != 0 {
		genesisDuration, err := s.SlotDuration(ctx)
		if err != nil {
			return nil, err
		}

		schedule = append([]*api.SlotDurationScheduleEntry{{Epoch: 0, Duration: genesisDuration}}, schedule...)
	}

	return schedule, nil
}
