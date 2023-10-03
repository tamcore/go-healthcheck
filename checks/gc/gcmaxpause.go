// Copyright 2021 by the contributors.
//
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

package gc

import (
	"fmt"
	"runtime"
	"time"

	"github.com/tamcore/go-healthcheck/checks"
)

// MaxPause returns a Check that fails if any recent Go garbage
// collection pause exceeds the provided threshold.
func MaxPause(threshold time.Duration) checks.Check {
	thresholdNanoseconds := uint64(threshold.Nanoseconds())

	return func() error {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		for _, pause := range stats.PauseNs {
			if pause > thresholdNanoseconds {
				return fmt.Errorf("recent GC cycle took %s > %s", time.Duration(pause), threshold)
			}
		}

		return nil
	}
}
