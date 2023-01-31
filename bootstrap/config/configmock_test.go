//
// Copyright (c) 2020 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package config

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v3/config"
)

// WritableInfo is used to hold configuration information that is considered "live" or can be changed on the fly without a restart of the service.
type WritableInfo struct {
	// Set level of logging to report
	//
	// example: TRACE
	// required: true
	// enum: TRACE,DEBUG,INFO,WARN,ERROR
	LogLevel        string
	Pipeline        PipelineInfo
	StoreAndForward StoreAndForwardInfo
	InsecureSecrets config.InsecureSecrets
	Telemetry       config.TelemetryInfo
}

// ConfigurationMockStruct
// swagger:model ConfigurationMockStruct
type ConfigurationMockStruct struct {
	// Writable contains the configuration that change be change on the fly
	Writable WritableInfo
	// Registry contains the configuration for connecting the Registry service
	Registry config.RegistryInfo
	// Service contains the standard 'service' configuration for the Application service
	Service config.ServiceInfo
	// HttpServer contains the configuration for the HTTP Server
	HttpServer HttpConfig
	// MessageBus contains the configuration for connecting to the EdgeX MessageBus
	MessageBus config.MessageBusInfo
	// Trigger contains the configuration for the Function Pipeline Trigger
	Trigger TriggerInfo
	// ApplicationSettings contains the custom configuration for the Application service
	ApplicationSettings map[string]string
	// Clients contains the configuration for connecting to the dependent Edgex clients
	Clients map[string]config.ClientInfo
}

// TriggerInfo contains Metadata associated with each Trigger
type TriggerInfo struct {
	// Type of trigger to start pipeline
	// enum: http, edgex-messagebus, or external-mqtt
	Type string
	// Used when Type=external-mqtt
	ExternalMqtt ExternalMqttConfig
}

// HttpConfig contains the addition configuration for HTTP Server
type HttpConfig struct {
	// Protocol is the for the HTTP Server to use HTTP or HTTPS
	// If HTTPS then the Secret store must contain the HTTPS cert and key
	Protocol string
	// SecretName is the name in the secret store for the HTTPS cert and key
	SecretName string
	// HTTPSCertName is name of the HTTPS cert in the secret store
	HTTPSCertName string
	// HTTPSKeyName is name of the HTTPS key in the secret store
	HTTPSKeyName string
}

// ExternalMqttConfig contains the MQTT broker configuration for MQTT Trigger
type ExternalMqttConfig struct {
	// Url contains the fully qualified URL to connect to the MQTT broker
	Url string
	// SubscribeTopics is a comma separated list of topics in which to subscribe
	SubscribeTopics string
	// PublishTopic is the topic to publish pipeline output (if any)
	PublishTopic string
	// ClientId to connect to the broker with.
	ClientId string
	// ConnectTimeout is a time duration indicating how long to wait timing out on the broker connection
	ConnectTimeout string
	// AutoReconnect indicated whether or not to retry connection if disconnected
	AutoReconnect bool
	// KeepAlive is seconds between client ping when no active data flowing to avoid client being disconnected
	KeepAlive int64
	// QoS for MQTT Connection
	QoS byte
	// Retain setting for MQTT Connection
	Retain bool
	// SkipCertVerify indicates if the certificate verification should be skipped
	SkipCertVerify bool
	// SecretPath is the name of the path in secret provider to retrieve your secrets
	SecretPath string
	// AuthMode indicates what to use when connecting to the broker. Options are "none", "cacert" , "usernamepassword", "clientcert".
	// If a CA Cert exists in the SecretPath then it will be used for all modes except "none".
	AuthMode string
	// RetryDuration indicates how long (in seconds) to wait timing out on the MQTT client creation
	RetryDuration int
	// RetryInterval indicates the time (in seconds) that will be waited between attempts to create MQTT client
	RetryInterval int
}

// PipelineInfo defines the top level data for configurable pipelines
type PipelineInfo struct {
	// ExecutionOrder is a list of functions, in execution order, for the default configurable pipeline
	ExecutionOrder string
	// PerTopicPipelines is a collection of pipelines that only execute if the incoming topic matched the pipelines configured topic
	PerTopicPipelines map[string]TopicPipeline
	// UseTargetTypeOfByteArray indicates if raw []byte type is to be used for the TargetType
	UseTargetTypeOfByteArray bool
	// TODO: for 3.0 rework this to be a single TargetType with empty, "raw" and "metric" as to allowed values
	// UseTargetTypeOfMetric indicates the Dtos.Metric is to be used for the TargetType
	UseTargetTypeOfMetric bool
	// Functions is a collection of pipeline functions with configured parameters to be used in the ExecutionOrder of one
	// of the configured pipelines (default or pre topic)
	Functions map[string]PipelineFunction
}

// TopicPipeline define the data to a Per Topics functions pipeline
type TopicPipeline struct {
	// Id is the unique ID of the pipeline instance
	Id string
	// Topics is the set of comma separated topics matched against the incoming to determine if pipeline should execute
	Topics string
	// ExecutionOrder is a list of functions, in execution order, for the pipeline instance
	ExecutionOrder string
}

// PipelineFunction is a collection of built-in pipeline functions configurations.
// The map key must be unique start with the name of one of the built-in configurable functions
type PipelineFunction struct {
	// Parameters is the collection of configurable parameters specific to the built-in configurable function specified by the map key.
	Parameters map[string]string
}

type StoreAndForwardInfo struct {
	Enabled       bool
	RetryInterval string
	MaxRetryCount int
}

// Credentials encapsulates username-password attributes.
type Credentials struct {
	Username string
	Password string
}

// UpdateFromRaw converts configuration received from the registry to a service-specific configuration struct which is
// then used to overwrite the service's existing configuration struct.
func (c *ConfigurationMockStruct) UpdateFromRaw(rawConfig interface{}) bool {
	configuration, ok := rawConfig.(*ConfigurationMockStruct)
	if ok {
		*c = *configuration
	}
	return ok
}

// EmptyWritablePtr returns a pointer to an empty WritableInfo struct.  It is used by the bootstrap to
// provide the appropriate structure for Config Client's WatchForChanges().
func (c *ConfigurationMockStruct) EmptyWritablePtr() interface{} {
	return &WritableInfo{}
}

// UpdateWritableFromRaw updates the Writeable section of configuration from raw update received from Configuration Provider.
func (c *ConfigurationMockStruct) UpdateWritableFromRaw(rawWritable interface{}) bool {
	writable, ok := rawWritable.(*WritableInfo)
	if ok {
		c.Writable = *writable
	}
	return ok
}

// GetBootstrap returns the configuration elements required by the bootstrap.
func (c *ConfigurationMockStruct) GetBootstrap() config.BootstrapConfiguration {
	return config.BootstrapConfiguration{
		Clients:    c.Clients,
		Service:    c.transformToBootstrapServiceInfo(),
		Registry:   c.Registry,
		MessageBus: c.MessageBus,
	}
}

// GetLogLevel returns log level from the configuration
func (c *ConfigurationMockStruct) GetLogLevel() string {
	return c.Writable.LogLevel
}

// GetRegistryInfo returns the RegistryInfo section from the configuration
func (c *ConfigurationMockStruct) GetRegistryInfo() config.RegistryInfo {
	return c.Registry
}

// GetInsecureSecrets returns the service's InsecureSecrets.
func (c *ConfigurationMockStruct) GetInsecureSecrets() config.InsecureSecrets {
	return c.Writable.InsecureSecrets
}

// GetTelemetryInfo returns the service's Telemetry settings.
func (c *ConfigurationMockStruct) GetTelemetryInfo() *config.TelemetryInfo {
	return &c.Writable.Telemetry
}

// transformToBootstrapServiceInfo transforms the SDK's ServiceInfo to the bootstrap's version of ServiceInfo
func (c *ConfigurationMockStruct) transformToBootstrapServiceInfo() config.ServiceInfo {
	return c.Service
}
