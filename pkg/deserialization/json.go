package deserialization

import "encoding/json"

func JSONDeserializer() Deserializer {
	return DeserializeFunc(func(b []byte, i interface{}) error {
		return json.Unmarshal(b, i)
	})
}
