package plugin

import (
	"testing"

	"github.com/hashicorp/nomad-autoscaler/sdk"
	"github.com/stretchr/testify/assert"
)

func TestExpression(t *testing.T) {
	testCases := []struct {
		expression    string
		count         int64
		metrics       sdk.TimestampedMetrics
		expectedCount int64
	}{
		{
			expression:    "Count - 1",
			count:         5,
			metrics:       nil,
			expectedCount: 4,
		},
		{
			expression:    "MetricsAvg == 9 ? 15 : 3",
			count:         5,
			metrics:       sdk.TimestampedMetrics{{Value: 17}, {Value: 1}},
			expectedCount: 15,
		},
		{
			expression:    "Metrics.Len() == 2 ? 15 : 3",
			count:         5,
			metrics:       sdk.TimestampedMetrics{{Value: 17}, {Value: 1}},
			expectedCount: 15,
		},
	}
	for _, tc := range testCases {
		res, _ := evaluateExpression(tc.expression, tc.count, tc.metrics)
		assert.Equal(t, tc.expectedCount, res)
	}
}
