module github.com/jsiebens/nomad-autoscaler-plugin-strategy-cron

go 1.17

require (
	github.com/hashicorp/go-hclog v1.0.0
	github.com/hashicorp/nomad-autoscaler v0.3.3
	github.com/stretchr/testify v1.7.0
)

require github.com/robfig/cron/v3 v3.0.1
