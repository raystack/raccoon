package serialization

type SerializeFunc func(m interface{}) ([]byte, error)
