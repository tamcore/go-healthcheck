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

package checks

import (
	"fmt"
	"runtime"
)

// GoroutineCountCheck returns a Check that fails if too many goroutines are
// running (which could indicate a resource leak).
func GoroutineCountCheck(threshold int) Check {
	return func() error {
		count := runtime.NumGoroutine()
		if count > threshold {
			return fmt.Errorf("too many goroutines (%d > %d)", count, threshold)
		}
		return nil
	}
}
