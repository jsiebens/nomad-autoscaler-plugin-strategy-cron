package plugin

import (
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-autoscaler/plugins/base"
	"github.com/hashicorp/nomad-autoscaler/sdk"
	"github.com/stretchr/testify/assert"
)

func TestStrategyPlugin_PluginInfo(t *testing.T) {
	s := &StrategyPlugin{logger: hclog.NewNullLogger()}
	expectedOutput := &base.PluginInfo{Name: "cron", PluginType: "strategy"}
	actualOutput, err := s.PluginInfo()
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, actualOutput)
}

func TestStrategyPlugin_SetConfig(t *testing.T) {
	testCases := []struct {
		config            map[string]string
		expectedSeparator string
	}{
		{config: map[string]string{}, expectedSeparator: defaultSeparator},
		{config: map[string]string{configKeySeparator: ";"}, expectedSeparator: ";"},
	}

	for _, tc := range testCases {
		s := &StrategyPlugin{logger: hclog.NewNullLogger()}
		err := s.SetConfig(tc.config)
		assert.Nil(t, err)
		assert.Equal(t, tc.expectedSeparator, s.separator)
	}
}

func TestStrategyPlugin_calculateDirection(t *testing.T) {
	testCases := []struct {
		inputCount     int64
		fixedCount     int64
		expectedOutput sdk.ScaleDirection
	}{
		{inputCount: 0, fixedCount: 1, expectedOutput: sdk.ScaleDirectionUp},
		{inputCount: 5, fixedCount: 5, expectedOutput: sdk.ScaleDirectionNone},
		{inputCount: 4, fixedCount: 0, expectedOutput: sdk.ScaleDirectionDown},
	}

	for _, tc := range testCases {
		s := &StrategyPlugin{logger: hclog.NewNullLogger()}
		assert.Equal(t, tc.expectedOutput, s.calculateDirection(tc.inputCount, tc.fixedCount))
	}
}

func TestStrategyPlugin_calculateTargetCount(t *testing.T) {
	location, _ := time.LoadLocation("Local")

	config := map[string]string{
		"count":           "1",
		"period_business": "* * 9-17 * * mon-fri * -> 10",
		"period_mon_100":  "* * 9-17 * * mon * -> 7",
		"period_sat":      "* * * * * sat * -> 5",
	}

	testCases := []struct {
		now           time.Time
		expectedCount int64
	}{
		{
			now:           time.Date(2021, time.March, 30, 8, 0, 0, 0, location),
			expectedCount: 1,
		},
		{
			now:           time.Date(2021, time.March, 30, 10, 0, 0, 0, location),
			expectedCount: 10,
		},
		{
			now:           time.Date(2021, time.March, 29, 10, 0, 0, 0, location),
			expectedCount: 7,
		},
		{
			now:           time.Date(2021, time.March, 28, 10, 0, 0, 0, location),
			expectedCount: 1,
		},
		{
			now:           time.Date(2021, time.March, 27, 10, 0, 0, 0, location),
			expectedCount: 5,
		},
	}

	for _, tc := range testCases {
		s := &StrategyPlugin{
			separator: defaultSeparator,
			logger:    hclog.NewNullLogger(),
		}
		count, _ := s.calculateTargetCount(config, fromTime(tc.now))
		assert.Equal(t, tc.expectedCount, count)
	}
}

func fromTime(t time.Time) func() time.Time {
	return func() time.Time { return t }
}
