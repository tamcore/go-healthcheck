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

package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/gsdenys/healthcheck/checks"
	"github.com/pkg/errors"
)

//Ping returns a Check function that validates Redis connection.
func Ping(client *redis.Client) checks.Check {
	return func() error {
		if client == nil {
			return fmt.Errorf("redis client is nil")
		}

		err := (*client).Ping().Err()
		if err != nil {
			err = errors.Wrap(err, "Redis healthcheck failed")
		}

		return err
	}
}
