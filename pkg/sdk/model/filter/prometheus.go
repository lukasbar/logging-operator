// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
)

// +kubebuilder:object:generate=true
// +docName:"[Prometheus Filter](https://github.com/fluent/fluent-plugin-prometheus#prometheus-outputfilter-plugin)"
// Prometheus Filter Plugin to count Incoming Records
type _docPrometheus interface{}

// +kubebuilder:object:generate=true
type PrometheusConfig struct {
	// +docLink:"Metrics Section,#Metrics-Section"
	Metrics []MetricSection `json:"metrics,omitempty"`
}

// +kubebuilder:object:generate=true
// +docName:"Metrics Section"
type MetricSection struct {
	// Metrics name
	Name string `json:"name"`
	//Metrics type
	Type string `json:"type"`
	//Description of metric
	Desc string `json:"desc"`
	//Key name of record for instrumentation.
	Key string `json:"key"`
	//Buckets of record for instrumentation
	Buckets string `json:"buckets,omitempty"`
	//Additional labels for this metric
	Labels Label `json:"labels,omitempty"`
}

// #### Example `Regexp` filter configurations
// ```
//apiVersion: logging.banzaicloud.io/v1beta1
//kind: Flow
//metadata:
//  name: demo-flow
//spec:
//  filters:
//    - prometheus:
//        metrics:
//        - type: counter
//          name: message_foo_counter
//          desc: "The total number of foo in message."
//          key: foo
//          labels:
//            foo: bar
//  selectors: {}
//  outputRefs:
//    - demo-output
// ```
//
// #### Fluentd Config Result
// ```
//<filter **>
//  @type prometheus
//  @id test_prometheus
//  <metric>
//    desc The total number of foo in message.
//    key foo
//    name message_foo_counter
//    type counter
//    <label>
//      foo bar
//    </label>
//  </metric>
//</filter>
// ```
type _expPrometheus interface{}

type Label map[string]string

func (r Label) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	recordSet := types.PluginMeta{
		Directive: "label",
	}
	directive := &types.GenericDirective{
		PluginMeta: recordSet,
		Params:     r,
	}
	return directive, nil
}

func (m *MetricSection) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	metricSection := &types.GenericDirective{
		PluginMeta: types.PluginMeta{
			Directive: "metric",
		},
	}
	metric := m.DeepCopy()
	// Render Labels as subdirective
	if metric.Labels != nil {
		if format, err := metric.Labels.ToDirective(secretLoader, ""); err != nil {
			return nil, err
		} else {
			metricSection.SubDirectives = append(metricSection.SubDirectives, format)
		}
	}
	// Remove Labels from parameter rendering
	metric.Labels = nil
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(metric); err != nil {
		return nil, err
	} else {
		metricSection.Params = params
	}
	return metricSection, nil
}

func (p *PrometheusConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "prometheus"
	prometheus := &types.GenericDirective{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "filter",
			Tag:       "**",
			Id:        id + "_" + pluginType,
		},
	}

	prometheusConfig := p.DeepCopy()

	if len(p.Metrics) > 0 {
		for _, metrics := range p.Metrics {
			if meta, err := metrics.ToDirective(secretLoader, ""); err != nil {
				return nil, err
			} else {
				prometheus.SubDirectives = append(prometheus.SubDirectives, meta)
			}
		}
	}

	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(prometheusConfig); err != nil {
		return nil, err
	} else {
		prometheus.Params = params
	}

	return prometheus, nil
}
