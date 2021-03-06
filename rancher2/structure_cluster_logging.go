package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterLogging(d *schema.ResourceData, in *managementClient.ClusterLogging) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("cluster_id", in.ClusterID)
	if err != nil {
		return err
	}

	err = d.Set("name", in.Name)
	if err != nil {
		return err
	}

	kind := d.Get("kind").(string)
	if kind == "" {
		if in.ElasticsearchConfig != nil {
			kind = loggingElasticsearchKind
		}
		if in.FluentForwarderConfig != nil {
			kind = loggingFluentdKind
		}
		if in.KafkaConfig != nil {
			kind = loggingKafkaKind
		}
		if in.SplunkConfig != nil {
			kind = loggingSplunkKind
		}
		if in.SyslogConfig != nil {
			kind = loggingSyslogKind
		}

		err = d.Set("kind", kind)
		if err != nil {
			return err
		}
	}

	switch kind {
	case loggingElasticsearchKind:
		elkConfig, err := flattenLoggingElasticsearchConfig(in.ElasticsearchConfig)
		if err != nil {
			return err
		}
		err = d.Set("elasticsearch_config", elkConfig)
		if err != nil {
			return err
		}
	case loggingFluentdKind:
		fluentdConfig, err := flattenLoggingFluentdConfig(in.FluentForwarderConfig)
		if err != nil {
			return err
		}
		err = d.Set("fluentd_config", fluentdConfig)
		if err != nil {
			return err
		}
	case loggingKafkaKind:
		kafkaConfig, err := flattenLoggingKafkaConfig(in.KafkaConfig)
		if err != nil {
			return err
		}
		err = d.Set("kafka_config", kafkaConfig)
		if err != nil {
			return err
		}
	case loggingSplunkKind:
		splunkConfig, err := flattenLoggingSplunkConfig(in.SplunkConfig)
		if err != nil {
			return err
		}
		err = d.Set("splunk_config", splunkConfig)
		if err != nil {
			return err
		}
	case loggingSyslogKind:
		syslogConfig, err := flattenLoggingSyslogConfig(in.SyslogConfig)
		if err != nil {
			return err
		}
		err = d.Set("syslog_config", syslogConfig)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("[ERROR] Flattening Cluster Logging resource data: Kind %s not supported", kind)
	}

	err = d.Set("namespace_id", in.NamespaceId)
	if err != nil {
		return err
	}

	err = d.Set("output_flush_interval", int(in.OutputFlushInterval))
	if err != nil {
		return err
	}

	err = d.Set("output_tags", toMapInterface(in.OutputTags))
	if err != nil {
		return err
	}

	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil

}

// Expanders

func expandClusterLogging(in *schema.ResourceData) (*managementClient.ClusterLogging, error) {
	obj := &managementClient.ClusterLogging{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] Expanding cluster config: Schema Resource data is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.ClusterID = in.Get("cluster_id").(string)
	obj.Name = in.Get("name").(string)

	switch kind := in.Get("kind").(string); kind {
	case loggingElasticsearchKind:
		elkConfig, err := expandLoggingElasticsearchConfig(in.Get("elasticsearch_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.ElasticsearchConfig = elkConfig
	case loggingFluentdKind:
		fluentdConfig, err := expandLoggingFluentdConfig(in.Get("fluentd_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.FluentForwarderConfig = fluentdConfig
	case loggingKafkaKind:
		kafkaConfig, err := expandLoggingKafkaConfig(in.Get("kafka_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.KafkaConfig = kafkaConfig
	case loggingSplunkKind:
		splunkConfig, err := expandLoggingSplunkConfig(in.Get("splunk_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.SplunkConfig = splunkConfig
	case loggingSyslogKind:
		syslogConfig, err := expandLoggingSyslogConfig(in.Get("syslog_config").([]interface{}))
		if err != nil {
			return obj, err
		}
		obj.SyslogConfig = syslogConfig
	}

	if v, ok := in.Get("namespace_id").(string); ok && len(v) > 0 {
		obj.NamespaceId = v
	}

	if v, ok := in.Get("output_flush_interval").(int); ok && v > 0 {
		obj.OutputFlushInterval = int64(v)
	}

	if v, ok := in.Get("output_tags").(map[string]interface{}); ok && len(v) > 0 {
		obj.OutputTags = toMapString(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}
