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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/strikesecurity/strikememongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongodbPingNilClient(t *testing.T) {
	assert.Error(t, MongodbPing(nil, 1*time.Second)(), "nil Mongo Client should fail")
}

func TestMongodbPing(t *testing.T) {
	//start mongo server
	mongoServer, err := strikememongo.StartWithOptions(
		&strikememongo.Options{
			MongoVersion:     "4.0.5",
			ShouldUseReplica: false,
		},
	)

	assert.NoError(t, err, "error when try to start MongoDB server")

	//get client connection
	clientOptions := options.Client().ApplyURI(mongoServer.URI())
	cli, errCli := mongo.Connect(context.TODO(), clientOptions)

	assert.NoError(t, errCli, "error when try to create MongoDB client.")
	assert.NoError(t, MongodbPing(cli, 1*time.Second)())

	mongoServer.Stop()
}
