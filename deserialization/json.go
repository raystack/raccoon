package deserialization

import "encoding/json"

type JSONDeserializer struct{}

func (j *JSONDeserializer) Deserialize(b []byte, i interface{}) error {
	return json.Unmarshal(b, i)
}
