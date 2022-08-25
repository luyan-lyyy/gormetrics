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
	"fmt"

	"github.com/luyan-lyyy/gormetrics/gormi"
	"github.com/pkg/errors"

	"github.com/prometheus/client_golang/prometheus"
)

// callbackHandler manages gorm query callback handling so query
// statistics are always up to date.
type callbackHandler struct {
	opts          *pluginOpts
	counters      *queryCounters
	defaultLabels map[string]string
}

func (h *callbackHandler) registerCallback(cb gormi.Callback) {
	cb.Create().After("gorm:after_create").Register(
		h.opts.callbackName("after_create"),
		h.afterCreate,
	)

	cb.Delete().After("gorm:after_delete").Register(
		h.opts.callbackName("after_delete"),
		h.afterDelete,
	)

	cb.Query().After("gorm:after_query").Register(
		h.opts.callbackName("after_query"),
		h.afterQuery,
	)

	cb.Update().After("gorm:after_update").Register(
		h.opts.callbackName("after_update"),
		h.afterUpdate,
	)
}

func (h *callbackHandler) afterCreate(scope gormi.Scope) {
	h.updateVectors(scope, h.counters.creates)
}

func (h *callbackHandler) afterDelete(scope gormi.Scope) {
	h.updateVectors(scope, h.counters.deletes)
}

func (h *callbackHandler) afterQuery(scope gormi.Scope) {
	h.updateVectors(scope, h.counters.queries)
}

func (h *callbackHandler) afterUpdate(scope gormi.Scope) {
	h.updateVectors(scope, h.counters.updates)
}

// updateVectors registers one or more of prometheus.CounterVec to increment
// with the status in scope (any type of query). If any errors are in
// scope.DB().GetErrors(), a status "fail" will be assigned to the increment.
// Otherwise, a status "success" will be assigned.
// Increments h.gauges.all (gormetrics_all_total) by default.
func (h *callbackHandler) updateVectors(scope gormi.Scope, vectors ...*prometheus.CounterVec) {
	vectors = append(vectors, h.counters.all)

	hasError := scope.DB().Error() != nil
	status := metricStatusFail
	if !hasError {
		status = metricStatusSuccess
	}

	labels := mergeLabels(prometheus.Labels{
		labelStatus: status,
	}, h.defaultLabels)

	for _, counter := range vectors {
		if counter == nil {
			continue
		}
		counter.With(labels).Inc()
	}
}

// extraInfo contains information for filtering the provided metrics.
type extraInfo struct {
	// The name of the database in use.
	dbName string

	// The name of the driver powering database/sql (underlying database for GORM).
	driverName string
}

// newCallbackHandler creates a new callback handler configured with info and opts.
// info does not contain any mandatory information for the functioning of the
// function, but sets label values which can be useful in the usage of
// the provided metrics (driver, database, connection).
// Automatically registers metrics.
func newCallbackHandler(info extraInfo, opts *pluginOpts) (*callbackHandler, error) {
	counters, err := newQueryCounters(opts.prometheusNamespace)
	if err != nil {
		return nil, errors.Wrap(err, "could not create query gauges")
	}

	return &callbackHandler{
		opts:     opts,
		counters: counters,
		defaultLabels: prometheus.Labels{
			labelDriver:   info.driverName,
			labelDatabase: info.dbName,
		},
	}, nil
}

// callbackName creates a GORM callback name based on the configured plugin
// scope and callback name.
func (c *pluginOpts) callbackName(callback string) string {
	return fmt.Sprintf("%v:%v", c.gormPluginScope, callback)
}

// Merges maps a and b. a is returned with extra values from b. Existing items
// in a with a matching key in b will not get overwritten.
func mergeLabels(a, b prometheus.Labels) prometheus.Labels {
	for k, v := range b {
		if _, exists := a[k]; !exists {
			a[k] = v
		}
	}
	return a
}
