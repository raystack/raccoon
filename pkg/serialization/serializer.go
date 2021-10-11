package serialization

type Serializer interface {
	Serialize(m interface{}) ([]byte, error)
}

type SerializeFunc func(m interface{}) ([]byte, error)

func (f SerializeFunc) Serialize(m interface{}) ([]byte, error) {
	return f(m)
}
