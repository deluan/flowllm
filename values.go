package pipelm

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/maps"
)

const (
	DefaultKey     = "text"
	DefaultChatKey = "_chat_messages"
)

// Values is a map of string to any value. This is the type used to pass values between handlers.
type Values map[string]any

// Merge merges multiple Values into one.
func (value Values) Merge(values ...Values) Values {
	res := Values{}
	for k, v := range value {
		res[k] = v
	}
	for _, v := range values {
		for k, vv := range v {
			res[k] = vv
		}
	}
	return res
}

// Get returns the value for a given key as a string. If the key does not exist, it returns an empty string.
func (value Values) Get(key string) string {
	v, ok := value[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// Keys returns the keys of the Values object.
func (value Values) Keys() []string {
	return maps.Keys(value)
}

// String returns a string representation of the Values object. If the Values object has only one key,
// it returns the value of that key. If the Values object has multiple keys, it returns a JSON representation.
func (value Values) String() string {
	if len(value) == 0 {
		return ""
	}
	if v, ok := value[DefaultKey]; ok {
		return v.(string)
	}
	if len(value) == 1 {
		for _, v := range value {
			return v.(string)
		}
	}
	j, _ := json.MarshalIndent(value, "", "  ")
	return string(j)
}
