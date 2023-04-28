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

type Values map[string]any

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

func (value Values) Get(key string) string {
	v, ok := value[key]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func (value Values) Keys() []string {
	return maps.Keys(value)
}

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

func Input(s string) Values {
	return Values{DefaultKey: s}
}
