package serialization

import "encoding/json"

type JSONSerializer struct{}

func (s *JSONSerializer) Serialize(m interface{}) ([]byte, error) {
	return json.Marshal(m)
}
