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
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestHTTPGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://example.com",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{})
		},
	)
	httpmock.RegisterResponder("GET", "http://example.com/redirect", func(request *http.Request) (*http.Response, error) {
		response := httpmock.NewStringResponse(http.StatusMovedPermanently, "")
		response.Header.Set("Location", "https://example.com")
		return response, nil
	})

	assert.NoError(t, Get("http://example.com", 5*time.Second)())
	assert.Error(t, Get("http://example.com/redirect", 5*time.Second)(), "redirect should fail")
	assert.Error(t, Get("http://example.com/nonexistent", 5*time.Second)(), "404 should fail")
}
