package config

type KafkaConfig struct {
	brokerList                string
	topic					  string
}

func (kc KafkaConfig) BrokerList() string {
	return kc.brokerList
}

func (kc KafkaConfig) Topic() string {
	return kc.topic
}

func NewKafkaConfig() KafkaConfig {
	kc := KafkaConfig{
		brokerList:                mustGetString("KAFKA_BROKER_LIST"),
		topic:					   mustGetString("KAFKA_TOPIC"),
	}
	return kc
}
