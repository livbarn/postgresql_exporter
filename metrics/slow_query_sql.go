package metrics

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type SQSMetrics struct {
	mutex        sync.Mutex
	names        []string
	milliseconds int64
	metrics      map[string]*prometheus.GaugeVec
}

func NewSQSMetrics(dbNames []string, durationToConsiderSlow time.Duration) *SQSMetrics {
	return &SQSMetrics{
		names:        dbNames,
		milliseconds: int64(durationToConsiderSlow / time.Millisecond),
		metrics: map[string]*prometheus.GaugeVec{
			"sql": prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "slow_query_sql",
				Help:      "Sql of slow query",
			}, []string{"db", "sql"}),
		},
	}
}

func (d *SQSMetrics) Scrape(db *sql.DB) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, name := range d.names {
		// var vals []interface{}
		rows, err := db.Query("select query, EXTRACT(milliseconds FROM now() - query_start) as time from pg_stat_activity WHERE datname not in ('rdsadmin', 'postgres') and state='active' and now() - query_start > ($1 || ' milliseconds' )::interval;", d.milliseconds)
		if err != nil {
			return errors.New("failed to get slow query sql: " + err.Error())
		}
		for rows.Next() {
			var query string
			time := new(float64)
			rows.Scan(&query, &time)
			d.metrics["sql"].WithLabelValues(name, query).Set(*time)
		}
	}

	return nil
}

func (d *SQSMetrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range d.metrics {
		m.Describe(ch)
	}
}

func (d *SQSMetrics) Collect(ch chan<- prometheus.Metric) {
	for _, m := range d.metrics {
		m.Collect(ch)
	}
}

// check interface
var _ Collection = new(SQSMetrics)
