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

import (
	"database/sql"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type database struct {
	name       string
	driverName string

	db *sql.DB
	sync.Mutex
}

// newDatabase creates a new database wrapper containing the name of the database,
// it's driver and the (sql) database itself.
func newDatabase(info extraInfo, db *sql.DB) *database {
	return &database{
		name:       info.dbName,
		driverName: info.driverName,
		db:         db,
	}
}

// collectConnectionStats collects database connections for Prometheus to scrape.
func (d *database) collectConnectionStats(counters *databaseGauges) {
	d.Lock()
	defer d.Unlock()

	defaultLabels := prometheus.Labels{
		labelDatabase: d.name,
		labelDriver:   d.driverName,
	}

	stats := d.db.Stats()

	counters.idle.
		With(defaultLabels).
		Set(float64(stats.Idle))

	counters.inUse.
		With(defaultLabels).
		Set(float64(stats.InUse))

	counters.open.
		With(defaultLabels).
		Set(float64(stats.OpenConnections))

	counters.maxOpen.
		With(defaultLabels).
		Set(float64(stats.MaxOpenConnections))

	counters.waitedFor.
		With(defaultLabels).
		Set(float64(stats.WaitCount))

	counters.blockedSeconds.
		With(defaultLabels).
		Set(stats.WaitDuration.Seconds())

	counters.closedMaxIdle.
		With(defaultLabels).
		Set(float64(stats.MaxIdleClosed))

	counters.closedMaxLifetime.
		With(defaultLabels).
		Set(float64(stats.MaxLifetimeClosed))

}

// databaseMetrics is a convenience struct for exporting database metrics to Prometheus.
type databaseMetrics struct {
	gauges *databaseGauges
	db     *database
}

// newDatabaseMetrics creates a new databaseMetrics instance with a database backing it
// for statistics. Use maintain to continuously collect statistics.
func newDatabaseMetrics(db *database, opts *pluginOpts) (*databaseMetrics, error) {
	gauges, err := newDatabaseGauges(opts.prometheusNamespace)
	if err != nil {
		return nil, errors.Wrap(err, "could not create database gauges")
	}

	return &databaseMetrics{
		gauges: gauges,
		db:     db,
	}, nil
}

// maintain collects connection statistics every 3 seconds.
func (d *databaseMetrics) maintain() {
	ticker := time.NewTicker(time.Second * 3)

	for range ticker.C {
		d.db.collectConnectionStats(d.gauges)
	}
}
