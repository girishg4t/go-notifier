package libs

import (
	jsoniter "github.com/json-iterator/go"
)

var j = jsoniter.ConfigCompatibleWithStandardLibrary

func DoMarshal(o map[string]interface{}) ([]byte, error) {
	body, err := j.Marshal(&o)
	return body, err
}
