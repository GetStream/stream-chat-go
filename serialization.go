package stream_chat

type marshalMap interface {
	marshalMap() map[string]interface{}
}

type unmarshalMap interface {
	unmarshalMap(map[string]interface{})
}

type unmarshalSlice interface {
	unmarshalSlice([]interface{})
}
