package maps

import (
	"bytes"
	"encoding/gob"
)

func RevertStringMap(m map[string]string) map[string]string {
	new := make(map[string]string, len(m))

	for k, v := range m {
		new[v] = k
	}

	return new
}

func DeepCopyStringMap(old map[string]string) map[string]string {
	new := make(map[string]string, len(old))

	for k, v := range old {
		new[k] = v
	}

	return new
}

func DeepCopyStringMapExceptKeys(old map[string]string, keys ...string) map[string]string {
	new := DeepCopyStringMap(old)

	for _, key := range keys {
		delete(new, key)
	}

	return new
}

func DeepCopyStringMapRemoveEmptyValue(old map[string]string) map[string]string {
	new := DeepCopyStringMap(old)

	for k, v := range new {
		if v == "" {
			delete(new, k)
		}
	}

	return new
}

func DeepCopyMap(in, out interface{}) {
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(in)
	gob.NewDecoder(buf).Decode(out)
}
