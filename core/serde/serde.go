package serde

type DeserializeFunc func(b []byte, i interface{}) error

type SerializeFunc func(m interface{}) ([]byte, error)
