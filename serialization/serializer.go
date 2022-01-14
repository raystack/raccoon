package serialization

type Serializer interface {
	Serialize(m interface{}) ([]byte, error)
}
