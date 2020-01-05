# [Prometheus Filter](https://github.com/fluent/fluent-plugin-prometheus#prometheus-outputfilter-plugin)
## Overview
 Prometheus Filter Plugin to count Incoming Records

## Configuration
### PrometheusConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| metrics | []MetricSection | No | - | [Metrics Section](#Metrics-Section)<br> |
### Metrics Section
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| name | string | Yes | - | Metrics name<br> |
| type | string | Yes | - | Metrics type [counter](https://github.com/fluent/fluent-plugin-prometheus#counter-type), [gauge](https://github.com/fluent/fluent-plugin-prometheus#gauge-type), [summary](https://github.com/fluent/fluent-plugin-prometheus#summary-type), [histogram](https://github.com/fluent/fluent-plugin-prometheus#histogram-type)<br> |
| desc | string | Yes | - | Description of metric<br> |
| key | string | Yes | - | Key name of record for instrumentation.<br> |
| buckets | string | No | - | Buckets of record for instrumentation<br> |
| labels | Label | No | - | Additional labels for this metric<br> |
 #### Example `Regexp` filter configurations
 ```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: demo-flow
spec:
  filters:
    - prometheus:
        metrics:
        - type: counter
          name: message_foo_counter
          desc: "The total number of foo in message."
          key: foo
          labels:
            foo: bar
  selectors: {}
  outputRefs:
    - demo-output
 ```

 #### Fluentd Config Result
 ```
<filter **>
  @type prometheus
  @id test_prometheus
  <metric>
    desc The total number of foo in message.
    key foo
    name message_foo_counter
    type counter
    <label>
      foo bar
    </label>
  </metric>
</filter>
 ```

---
