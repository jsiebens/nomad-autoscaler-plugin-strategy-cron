package plugin

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/cronexpr"
)

func parsePeriodRule(key, value, separator string) (*Rule, error) {
	var count int64 = 1
	var priority int64 = 0

	entries := strings.Split(value, separator)

	period := strings.TrimSpace(entries[0])
	expr, err := cronexpr.Parse(period)
	if err != nil {
		return nil, err
	}

	if len(entries) > 1 {
		v, err := strconv.ParseInt(strings.TrimSpace(entries[1]), 10, 64)
		if err != nil {
			return nil, err
		}
		count = v
	}

	index := strings.LastIndex(key, "_")
	if index != -1 {
		v, err := strconv.ParseInt(strings.TrimSpace(key[index+1:]), 10, 64)
		if err == nil {
			priority = v
		}
	}

	return &Rule{
		expr:     expr,
		period:   period,
		count:    count,
		priority: priority,
	}, nil
}

type Rule struct {
	expr     *cronexpr.Expression
	period   string
	key      string
	count    int64
	priority int64
}

func (t *Rule) InPeriod(now time.Time) bool {
	nextIn := t.expr.Next(now)
	timeSince := now.Sub(nextIn)
	if -time.Second <= timeSince && timeSince <= time.Second {
		return true
	}

	return false
}

type RuleSorter []*Rule

func (r RuleSorter) Len() int      { return len(r) }
func (r RuleSorter) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r RuleSorter) Less(i, j int) bool {
	if r[i].priority == r[j].priority {
		return r[j].count < r[i].count
	} else {
		return r[j].priority < r[i].priority
	}
}
