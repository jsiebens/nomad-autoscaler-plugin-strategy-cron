package plugin

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-autoscaler/plugins"
	"github.com/hashicorp/nomad-autoscaler/plugins/base"
	"github.com/hashicorp/nomad-autoscaler/plugins/strategy"
	"github.com/hashicorp/nomad-autoscaler/sdk"
)

const (
	// pluginName is the unique name of the this plugin amongst strategy
	// plugins.
	pluginName = "cron"

	defaultSeparator = "->"

	configKeySeparator = "separator"

	// These are the keys read from the RunRequest.Config map.
	runConfigKeyCount        = "count"
	runConfigKeyPeriodPrefix = "period_"
)

var (
	PluginID = plugins.PluginID{
		Name:       pluginName,
		PluginType: sdk.PluginTypeStrategy,
	}

	PluginConfig = &plugins.InternalPluginConfig{
		Factory: func(l hclog.Logger) interface{} { return NewCronPlugin(l) },
	}

	pluginInfo = &base.PluginInfo{
		Name:       pluginName,
		PluginType: sdk.PluginTypeStrategy,
	}
)

// Assert that StrategyPlugin meets the strategy.Strategy interface.
var _ strategy.Strategy = (*StrategyPlugin)(nil)

// StrategyPlugin is the Periods implementation of the strategy.Strategy
// interface.
type StrategyPlugin struct {
	separator string
	logger    hclog.Logger
}

// NewCronPlugin returns the Periods implementation of the
// strategy.Strategy interface.
func NewCronPlugin(log hclog.Logger) strategy.Strategy {
	return &StrategyPlugin{
		logger: log,
	}
}

// PluginInfo satisfies the PluginInfo function on the base.Base interface.
func (s *StrategyPlugin) PluginInfo() (*base.PluginInfo, error) {
	return pluginInfo, nil
}

// SetConfig satisfies the SetConfig function on the base.Base interface.
func (s *StrategyPlugin) SetConfig(config map[string]string) error {
	s.separator = defaultSeparator

	sep, ok := config[configKeySeparator]
	if ok {
		s.separator = sep
	}

	return nil
}

// Run satisfies the Run function on the strategy.Strategy interface.
func (s *StrategyPlugin) Run(eval *sdk.ScalingCheckEvaluation, count int64) (*sdk.ScalingCheckEvaluation, error) {
	targetCount, err := s.calculateTargetCount(eval.Check.Strategy.Config, time.Now)
	if err != nil {
		return eval, err
	}

	// Identify the direction of scaling, if any.
	eval.Action.Direction = s.calculateDirection(count, targetCount)
	if eval.Action.Direction == sdk.ScaleDirectionNone {
		return eval, nil
	}

	// Log at trace level the details of the strategy calculation. This is
	// helpful in ultra-debugging situations when there is a need to understand
	// all the calculations made.
	s.logger.Trace("calculated scaling strategy results",
		"check_name", eval.Check.Name, "current_count", count, "new_count", targetCount,
		"direction", eval.Action.Direction)

	eval.Action.Count = targetCount
	eval.Action.Reason = fmt.Sprintf("scaling %s because cron value is %d", eval.Action.Direction, targetCount)

	return eval, nil
}

// calculateDirection is used to calculate the direction of scaling that should
// occur, if any at all.
func (s *StrategyPlugin) calculateDirection(count, target int64) sdk.ScaleDirection {
	if count == target {
		return sdk.ScaleDirectionNone
	} else if count < target {
		return sdk.ScaleDirectionUp
	}
	return sdk.ScaleDirectionDown
}

func (s *StrategyPlugin) calculateTargetCount(config map[string]string, timer func() time.Time) (int64, error) {
	now := timer()

	var value int64 = 1
	var rules []*Rule

	for k, element := range config {
		if k == runConfigKeyCount {
			v, err := strconv.ParseInt(element, 10, 64)
			if err != nil {
				return -1, fmt.Errorf("invalid value for `%s`: %v (%T)", runConfigKeyCount, element, element)
			}
			value = v
		}

		if strings.HasPrefix(k, runConfigKeyPeriodPrefix) {
			rule, err := parsePeriodRule(k, element, s.separator)
			if err != nil {
				return -1, err
			}

			inPeriod := rule.InPeriod(now)

			s.logger.Trace("checking period", "period", rule.period, "in_period", inPeriod, "priority", rule.priority)

			if inPeriod {
				rules = append(rules, rule)
			}
		}
	}

	if len(rules) == 0 {
		return value, nil
	} else if len(rules) == 1 {
		s.logger.Trace("selected period", "period", rules[0].period, "priority", rules[0].priority, "count", rules[0].count)
		return rules[0].count, nil
	} else {
		sort.Sort(RuleSorter(rules))
		s.logger.Trace("selected period", "period", rules[0].period, "priority", rules[0].priority, "count", rules[0].count)
		return rules[0].count, nil
	}
}
