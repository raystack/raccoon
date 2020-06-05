package config

type KafkaConfig struct {
	brokerList                string
}

func (kc KafkaConfig) BrokerList() string {
	return kc.brokerList
}

func NewKafkaConfig() KafkaConfig {
	kc := KafkaConfig{
		brokerList:                mustGetString("KAFKA_BROKER_LIST"),
	}
	return kc
}
