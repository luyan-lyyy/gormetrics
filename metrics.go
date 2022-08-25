// Copyright 2019 Profects Group B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gormetrics

const (
	labelStatus   = "status"
	labelDatabase = "database"
	labelDriver   = "driver"

	// Statuses for metrics (values of labelStatus).
	metricStatusFail    = "fail"
	metricStatusSuccess = "success"

	metricOpenConnections              = "connections_open"
	metricIdleConnections              = "connections_idle"
	metricInUseConnections             = "connections_in_use"
	metricMaxOpenConnections           = "connections_max_open"
	metricWaitedForConnections         = "connections_wait_for"
	metricBlockedSecondsConnections    = "connections_blocked_seconds"
	metricClosedMaxLifetimeConnections = "connections_closed_max_lifetime"

	helpOpenConnections              = `Currently open connections to the database`
	helpIdleConnections              = `Currently idle connections to the database`
	helpInUseConnections             = `Currently in use connections`
	helpMaxOpenConnections           = `Currently Maximum number of open connections to the database`
	helpWaitedForConnections         = `Currently waited for connections to the database`
	helpBlockedSecondsConnections    = `The total time blocked waiting for a new connection`
	helpClosedMaxLifetimeConnections = `The total number of connections closed due to SetConnMaxLifetime`

	metricAllTotal     = "all_total"
	metricCreatesTotal = "creates_total"
	metricDeletesTotal = "deletes_total"
	metricQueriesTotal = "queries_total"
	metricUpdatesTotal = "updates_total"

	helpAllTotal     = `All queries requested`
	helpCreatesTotal = `All create queries requested`
	helpDeletesTotal = `All delete queries requested`
	helpQueriesTotal = `All select queries requested`
	helpUpdatesTotal = `All update queries requested`
)
