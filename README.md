# Nomad Autoscaler Cron Strategy

A cron-like strategy plugin, where task groups are scaled based on a predefined scheduled.

```hcl
job "webapp" {
  ...
  group "demo" {
    ...
    scaling {
      ...
      policy {
        check "business hours" {
          source = "prometheus"

          strategy "cron" {
            count = 2
            period_business = "* * 9-17 * * mon-fri * -> 5"
            period_weekend = "* * * * * sat,sun * -> 1"
          }
        }
      }
    }
  }
}
```

In the example above, every weekday between 9 and 16:59, the number of instances is increased to 5 to handle the large traffic during operating hours. 
On the weekend, the number drops to 1. 
The rest of the time, the value is taken from the default `count = 2`