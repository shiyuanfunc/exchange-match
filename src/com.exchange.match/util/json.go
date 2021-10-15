package util

import "encoding/json"

// obj 转 json
func ToJsonString(obj interface{}) string {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jsonStr)
}
