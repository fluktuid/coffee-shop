package util

import (
	"encoding/json"
)

func Marshal(i interface{}) string {
	b, _ := json.Marshal(i)
	return string(b)
}

func Unmarshal(s string, i interface{}) {
	json.Unmarshal([]byte(s), i)
}
