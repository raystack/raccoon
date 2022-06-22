package serialization

import "encoding/json"

func SerializeJSON(m interface{}) ([]byte, error) {
	return json.Marshal(m)
}
