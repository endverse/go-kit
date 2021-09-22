package codec

import (
	"encoding/json"
	"errors"
	"reflect"
)

func Scan(data, value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("unmarshal json failed")
	}

	return json.Unmarshal(bytes, data)
}

func Value(data interface{}) ([]byte, error) {
	vi := reflect.ValueOf(data)
	if vi.IsZero() {
		return nil, nil
	}

	return json.Marshal(data)
}
