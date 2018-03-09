package MsgPack

import msgpack "gopkg.in/vmihailenco/msgpack.v2"

// Marshal data in format of msgpack
func Marshal(data interface{}) ([]byte, error) {
	return msgpack.Marshal(data)
}

// Unmarshal unmarshal data in msgpack format
func Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}
