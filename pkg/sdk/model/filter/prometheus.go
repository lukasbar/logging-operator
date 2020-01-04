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

// +docName:"Fluentd prometheus filter"
// Prometheus Filter Plugin to count Incoming Records
// More information at https://github.com/fluent/fluent-plugin-prometheus#prometheus-outputfilter-plugin
//
// #### Example record configurations
// ```
// spec:
//  filters:
//    - prometheus:
//        metric:
//          type: counter
//          name: message_foo_counter
//          desc: "The total number of foo in message."
//          key: foo
//          labels:
//            foo: bar
// ```
type _docPrometheus interface{}

// +kubebuilder:object:generate=true
type PrometheusConfig struct {
	// +docLink:"Metrics Section,#Parse-Section"
	Metric *MetricSection `json:"metric,omitempty"`
}

// +kubebuilder:object:generate=true
// +docName:"Metric Section"
type MetricSection struct {
	// Metric name
	Name string `json:"name"`
	//Metric type
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

func (r Label) merge(input Label) {
	for k, v := range input {
		r[k] = v
	}
}

func (m *MetricSection) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	meta := types.PluginMeta{
		Directive: "metric",
	}
	metric := m.DeepCopy()

	if m.Labels == nil {
		m.Labels = Label{}
	}
	m.Labels.merge(Label{
		"tag":  `$.kubernetes.tag`,
		"host": `$.kubernetes.host`,
	})

	return types.NewFlatDirective(meta, metric, secretLoader)
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
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(prometheusConfig); err != nil {
		return nil, err
	} else {
		prometheus.Params = params
	}
	if p.Metric != nil {
		if format, err := p.Metric.ToDirective(secretLoader, ""); err != nil {
			return nil, err
		} else {
			prometheus.SubDirectives = append(prometheus.SubDirectives, format)
		}
	}

	return prometheus, nil
}
