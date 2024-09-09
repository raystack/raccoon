package deserialization

type DeserializeFunc func(b []byte, i interface{}) error
