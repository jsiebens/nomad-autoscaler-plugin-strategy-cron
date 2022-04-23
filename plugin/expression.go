package plugin

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/antonmedv/expr"
	"github.com/hashicorp/nomad-autoscaler/sdk"
)

type env struct {
	Count      int64 // current count
	Metrics    sdk.TimestampedMetrics
	MetricsAvg float64
	MetricsMin float64
	MetricsMax float64
}

func sum(metrics sdk.TimestampedMetrics) float64 {
	s := 0.0
	for _, v := range metrics {
		s += v.Value
	}
	return s
}

func avg(metrics sdk.TimestampedMetrics) float64 {
	if len(metrics) == 0 {
		return 0
	}
	return sum(metrics) / float64(len(metrics))
}

func max(metrics sdk.TimestampedMetrics) float64 {
	if len(metrics) == 0 {
		return 0
	}
	sortedMetrics := make(sdk.TimestampedMetrics, len(metrics))
	copy(sortedMetrics, metrics)
	sort.Slice(sortedMetrics, func(i, j int) bool {
		return sortedMetrics[i].Value > sortedMetrics[j].Value
	})
	return sortedMetrics[0].Value
}

func min(metrics sdk.TimestampedMetrics) float64 {
	if len(metrics) == 0 {
		return 0
	}
	sortedMetrics := make(sdk.TimestampedMetrics, len(metrics))
	copy(sortedMetrics, metrics)
	sort.Slice(sortedMetrics, func(i, j int) bool {
		return sortedMetrics[i].Value < sortedMetrics[j].Value
	})
	return sortedMetrics[0].Value
}

// evaluateExpression parses the given expression and evaluates the resulting program. The expected result is
// the target count.
// Available variables are:
// "Count" - the current count
// "Metrics" - sdk.TimestampedMetrics
// "MetricsAvg" - average of Metrics
// "MetricsMin" - min of Metrics
// "MetricsMax" - max of Metrics
func evaluateExpression(expression string, count int64, metrics sdk.TimestampedMetrics) (int64, error) {
	env := env{
		Count:      count,
		Metrics:    metrics,
		MetricsAvg: avg(metrics),
		MetricsMax: max(metrics),
		MetricsMin: min(metrics),
	}
	res, err := expr.Eval(expression, env)
	if err != nil {
		return 0, err
	}
	switch res.(type) {
	case string:
		return strconv.ParseInt(res.(string), 10, 64)

	case bool:
		if res.(bool) {
			return count, nil
		} else {
			return 0, nil
		}

	case uint8:
		return int64(res.(uint8)), nil

	case uint16:
		return int64(res.(uint16)), nil

	case uint32:
		return int64(res.(uint32)), nil

	case uint64:
		return int64(res.(uint64)), nil

	case int8:
		return int64(res.(int8)), nil

	case int16:
		return int64(res.(int16)), nil

	case int32:
		return int64(res.(int32)), nil

	case int64:
		return res.(int64), nil

	case int:
		return int64(res.(int)), nil

	case uint:
		return int64(res.(uint)), nil
	}
	return 0, fmt.Errorf("could not parse expression result")
}
