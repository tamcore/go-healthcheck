// Copyright 2017 by the contributors.
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

package healthcheck

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"log"

	"github.com/gsdenys/healthcheck/checks"
)

// basicHandler is a basic Handler implementation.
type basicHandler struct {
	http.ServeMux
	checksMutex     sync.RWMutex
	livenessChecks  map[string]checks.Check
	readinessChecks map[string]checks.Check
}

const (
	liveness  = "LIVENESS_ENDPOINT"
	readiness = "READINESS_ENDPOINT"

	defaultLivenessEndpoint  = "/live"
	defaultReadinessEndpoint = "/ready"
)

func getLivenessEndpoint() string {
	value, hasValue := os.LookupEnv(liveness)
	if !hasValue {
		return defaultLivenessEndpoint
	}

	return value
}

func getreadinessEndpoint() string {
	value, hasValue := os.LookupEnv(readiness)
	if !hasValue {
		return defaultReadinessEndpoint
	}

	return value
}

// NewHandler creates a new basic Handler
func NewHandler() Handler {
	h := &basicHandler{
		livenessChecks:  make(map[string]checks.Check),
		readinessChecks: make(map[string]checks.Check),
	}

	h.Handle(getLivenessEndpoint(), http.HandlerFunc(h.LiveEndpoint))
	h.Handle(getreadinessEndpoint(), http.HandlerFunc(h.ReadyEndpoint))
	return h
}

func (s *basicHandler) LiveEndpoint(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r, s.livenessChecks)
}

func (s *basicHandler) ReadyEndpoint(w http.ResponseWriter, r *http.Request) {
	s.handle(w, r, s.readinessChecks, s.livenessChecks)
}

func (s *basicHandler) AddLivenessCheck(name string, check checks.Check) {
	s.checksMutex.Lock()
	defer s.checksMutex.Unlock()
	s.livenessChecks[name] = check
}

func (s *basicHandler) AddReadinessCheck(name string, check checks.Check) {
	s.checksMutex.Lock()
	defer s.checksMutex.Unlock()
	s.readinessChecks[name] = check
}

func (s *basicHandler) collectChecks(checks map[string]checks.Check, resultsOut map[string]string, statusOut *int) {
	s.checksMutex.RLock()
	defer s.checksMutex.RUnlock()
	for name, check := range checks {
		if err := check(); err != nil {
			*statusOut = http.StatusServiceUnavailable
			resultsOut[name] = err.Error()
		} else {
			resultsOut[name] = "OK"
		}
	}
}

func (s *basicHandler) handle(w http.ResponseWriter, r *http.Request, checks ...map[string]checks.Check) {
	log.SetPrefix("[ERROR]")

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	checkResults := make(map[string]string)
	status := http.StatusOK
	for _, checks := range checks {
		s.collectChecks(checks, checkResults, &status)
	}

	// write out the response code and content type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// unless ?full=1, return an empty body. Kubernetes only cares about the
	// HTTP status code, so we won't waste bytes on the full body.
	if r.URL.Query().Get("full") != "1" {
		_, errWrite := w.Write([]byte("{}\n"))

		if errWrite == nil {
			return
		}

		log.Println("writing simple response error. Continuing for full respose", errWrite)
	}

	// otherwise, write the JSON body ignoring any encoding errors (which
	// shouldn't really be possible since we're encoding a map[string]string).
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "    ")
	errEncode := encoder.Encode(checkResults)

	if errEncode != nil {
		log.Println("encoding http data error", errEncode)
	}
}
