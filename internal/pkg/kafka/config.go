package kafka

type KafkaConfig struct {
	Brokers           []string `mapstructure:"brokers"`
	GroupID           string   `mapstructure:"groupId"`
	ClientID          string   `mapstructure:"clientId"`
	Topic             string   `mapstructure:"topic"`
	PartitionCount    int      `mapstructure:"partitionCount"`
	ReplicationFactor int      `mapstructure:"replicationFactor"`
	RetentionHours    int      `mapstructure:"retentionHours"`
}

type KafkaClientsConfig map[string]KafkaConfig
