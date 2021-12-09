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

package http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPGet(t *testing.T) {
	assert.NoError(t, Get("https://gsdenys.github.io", 5*time.Second)())
	assert.Error(t, Get("http://gsdenys.github.io", 5*time.Second)(), "redirect should fail")
	assert.Error(t, Get("https://gsdenys.github.io/nonexistent", 5*time.Second)(), "404 should fail")
}
