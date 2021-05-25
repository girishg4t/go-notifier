package libs

import (
	"log"

	"gopkg.in/yaml.v2"
)

func GetConfigObject(ys string) ([]byte, error) {
	y := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(ys), &y)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return DoMarshal(y)
}
