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
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"time"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/gsdenys/healthcheck/checks/db"
	"github.com/gsdenys/healthcheck/checks/dns"
	"github.com/gsdenys/healthcheck/checks/goroutine"
)

func Example() {
	// Create a Handler that we can use to register liveness and readiness checks.
	health := NewHandler()

	// Add a readiness check to make sure an upstream dependency resolves in DNS.
	// If this fails we don't want to receive requests, but we shouldn't be
	// restarted or rescheduled.
	upstreamHost := "upstream.example.com"
	health.AddReadinessCheck(
		"upstream-dep-dns",
		dns.Resolve(upstreamHost, 50*time.Millisecond))

	// Add a liveness check to detect Goroutine leaks. If this fails we want
	// to be restarted/rescheduled.
	health.AddLivenessCheck("goroutine-threshold", goroutine.Count(100))

	// Serve http://0.0.0.0:8080/live and http://0.0.0.0:8080/ready endpoints.
	// go http.ListenAndServe("0.0.0.0:8080", health)

	// Make a request to the readiness endpoint and print the response.
	fmt.Print(dumpRequest(health, "GET", "/ready"))

	// Output:
	// HTTP/1.1 503 Service Unavailable
	// Connection: close
	// Content-Type: application/json; charset=utf-8
	//
	// {}
}

func Example_database() {
	// Connect to a database/sql database
	database := connectToDatabase()

	// Create a Handler that we can use to register liveness and readiness checks.
	health := NewHandler()

	// Add a readiness check to we don't receive requests unless we can reach
	// the database with a ping in <1 second.
	health.AddReadinessCheck("database", db.Ping(database, 1*time.Second))

	// Serve http://0.0.0.0:8080/live and http://0.0.0.0:8080/ready endpoints.
	// go http.ListenAndServe("0.0.0.0:8080", health)

	// Make a request to the readiness endpoint and print the response.
	fmt.Print(dumpRequest(health, "GET", "/ready?full=1"))

	// Output:
	// HTTP/1.1 200 OK
	// Connection: close
	// Content-Type: application/json; charset=utf-8
	//
	// {
	//     "database": "OK"
	// }
}

func dumpRequest(handler http.Handler, method string, path string) string {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	dump, err := httputil.DumpResponse(rr.Result(), true)
	if err != nil {
		panic(err)
	}
	return strings.Replace(string(dump), "\r\n", "\n", -1)
}

func connectToDatabase() *sql.DB {
	db, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return db
}
