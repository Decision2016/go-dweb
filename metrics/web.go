/**
  @author: decision
  @date: 2024/6/25
  @note:
**/

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	LoaderTaskCountInQueue = promauto.NewCounter(prometheus.CounterOpts{
		Name: "load_task_count_in_queue",
		Help: "The task count in loader process queue",
	})
	LoaderCurrentTaskProgress = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "loader_current_task_progress",
		Help: "indicates the progress of the current task being processed by the loader",
	})
)
