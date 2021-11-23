package serialization

import "encoding/json"

func JSONSerializer() Serializer {
	return SerializeFunc(func(m interface{}) ([]byte, error) {
		return json.Marshal(m)
	})
}
