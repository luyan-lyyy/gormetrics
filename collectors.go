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
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type globalCollectors struct {
	query    map[string]*queryCounters
	database map[string]*databaseGauges

	sync.Mutex
}

// collectors is used by newQueryCounters and newDatabaseGauges to cache existing
// collectors so none are registered in Prometheus twice (this causes an error).
var collectors = globalCollectors{
	query:    make(map[string]*queryCounters),
	database: make(map[string]*databaseGauges),
}

// queryCounters contains all gauges that are exported.
type queryCounters struct {
	all     *prometheus.CounterVec
	creates *prometheus.CounterVec
	deletes *prometheus.CounterVec
	queries *prometheus.CounterVec
	updates *prometheus.CounterVec
}

func newQueryCounters(namespace string) (*queryCounters, error) {
	collectors.Lock()
	defer collectors.Unlock()

	if gc, exists := collectors.query[namespace]; exists {
		return gc, nil
	}

	cc := counterVecCreator{
		namespace: namespace,
		labels: []string{
			labelDatabase,
			labelDriver,
			labelStatus,
		},
	}

	qc := queryCounters{
		all:     cc.new(metricAllTotal, helpAllTotal),
		creates: cc.new(metricCreatesTotal, helpCreatesTotal),
		deletes: cc.new(metricDeletesTotal, helpDeletesTotal),
		queries: cc.new(metricQueriesTotal, helpQueriesTotal),
		updates: cc.new(metricUpdatesTotal, helpUpdatesTotal),
	}

	if err := registerCollectors(
		qc.all,
		qc.creates,
		qc.deletes,
		qc.queries,
		qc.updates,
	); err != nil {
		return nil, errors.Wrap(err, "could not register collectors")
	}

	collectors.query[namespace] = &qc

	return collectors.query[namespace], nil
}

type databaseGauges struct {
	idle              *prometheus.GaugeVec
	inUse             *prometheus.GaugeVec
	open              *prometheus.GaugeVec
	maxOpen           *prometheus.GaugeVec
	waitedFor         *prometheus.GaugeVec
	blockedSeconds    *prometheus.GaugeVec
	closedMaxIdle     *prometheus.GaugeVec
	closedMaxLifetime *prometheus.GaugeVec
}

func newDatabaseGauges(namespace string) (*databaseGauges, error) {
	collectors.Lock()
	defer collectors.Unlock()

	if gc, exists := collectors.database[namespace]; exists {
		return gc, nil
	}

	vecCreator := gaugeVecCreator{
		namespace: namespace,
		labels: []string{
			labelDatabase,
			labelDriver,
		},
	}

	dg := databaseGauges{
		idle:              vecCreator.new(metricIdleConnections, helpIdleConnections),
		inUse:             vecCreator.new(metricInUseConnections, helpInUseConnections),
		open:              vecCreator.new(metricOpenConnections, helpOpenConnections),
		maxOpen:           vecCreator.new(metricMaxOpenConnections, helpMaxOpenConnections),
		waitedFor:         vecCreator.new(metricWaitedForConnections, helpWaitedForConnections),
		blockedSeconds:    vecCreator.new(metricBlockedSecondsConnections, helpBlockedSecondsConnections),
		closedMaxIdle:     vecCreator.new(metricClosedMaxIdleConnections, helpClosedMaxIdleConnections),
		closedMaxLifetime: vecCreator.new(metricClosedMaxLifetimeConnections, helpClosedMaxLifetimeConnections),
	}

	if err := registerCollectors(
		dg.idle,
		dg.inUse,
		dg.open,
		dg.maxOpen,
		dg.waitedFor,
		dg.blockedSeconds,
		dg.closedMaxIdle,
		dg.closedMaxLifetime,
	); err != nil {
		return nil, err
	}

	collectors.database[namespace] = &dg

	return collectors.database[namespace], nil
}

// registerCollectors registers multiple instances of prometheus.Collector.
func registerCollectors(collectors ...prometheus.Collector) error {
	for _, c := range collectors {
		err := prometheus.Register(c)
		if err != nil {
			return err
		}
	}

	return nil
}
