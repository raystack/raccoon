package deserialization

import "encoding/json"

func DeserializeJSON(b []byte, i interface{}) error {
	return json.Unmarshal(b, i)
}
