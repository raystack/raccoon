package deserialization

type Deserializer interface {
	Deserialize(b []byte, i interface{}) error
}

type DeserializeFunc func(b []byte, i interface{}) error

func (f DeserializeFunc) Deserialize(b []byte, i interface{}) error {
	return f(b, i)
}
