package plugin

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRule_New(t *testing.T) {
	testCases := []struct {
		name          string
		inputValue    string
		expectedError bool
		expectedCount int64
	}{
		{
			name:          "period_default",
			inputValue:    "* * 9-17 * * mon-fri *",
			expectedCount: 1,
		},
		{
			name:          "period_valid_count",
			inputValue:    "* * 9-17 * * mon-fri *;5",
			expectedCount: 5,
		},
		{
			name:          "period_invalid_expression",
			inputValue:    "invalid",
			expectedError: true,
		},
		{
			name:          "period_invalid_count",
			inputValue:    "* * 9-17 * * mon-fri *;invalid",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		rule, err := parsePeriodRule(tc.name, tc.inputValue, ";")
		if tc.expectedError {
			assert.NotNil(t, err)
			assert.Nil(t, rule)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedCount, rule.count)
		}
	}
}

func TestRule_Sort(t *testing.T) {
	rule1, _ := parsePeriodRule("a", "* 1 * * *;1", ";")
	rule2, _ := parsePeriodRule("b", "* 1 * * *;6", ";")
	rule3, _ := parsePeriodRule("c", "* 1 * * *;5", ";")
	rule4, _ := parsePeriodRule("c_100", "* 1 * * *;4", ";")

	rules := []*Rule{
		rule1,
		rule2,
		rule3,
		rule4,
	}

	sort.Sort(RuleSorter(rules))

	assert.Equal(t, rule4, rules[0])
	assert.Equal(t, rule2, rules[1])
	assert.Equal(t, rule3, rules[2])
	assert.Equal(t, rule1, rules[3])
}
