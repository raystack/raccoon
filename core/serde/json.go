package serde

import "encoding/json"

func DeserializeJSON(b []byte, i interface{}) error {
	return json.Unmarshal(b, i)
}

func SerializeJSON(m interface{}) ([]byte, error) {
	return json.Marshal(m)
}
