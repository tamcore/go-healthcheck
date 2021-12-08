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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestDatabasePing(t *testing.T) {
	assert.Error(t, DatabasePing(nil, 1*time.Second)(), "nil DB should fail")

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	assert.NoError(t, DatabasePing(db, 1*time.Second)(), "ping should succeed")
}
