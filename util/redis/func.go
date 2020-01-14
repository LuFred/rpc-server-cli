package redis

import (
	"encoding/json"
	"fmt"
)

func convertInterfaceToString(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Errorf("Cache: Convert Interface To String:%s", err.Error())
	}
	return string(b), nil
}
