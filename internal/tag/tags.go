package tag

import (
	"encoding/json"
	"go.opentelemetry.io/otel/attribute"
)

type Tags []attribute.KeyValue

func (t Tags) Len() int {
	return len(t)
}

func (t Tags) Less(i, j int) bool {
	return t[i].Key < t[j].Key
}

func (t Tags) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t Tags) String() string {
	tagMap := make(map[string]interface{})
	for _, tag := range t {
		tagMap[string(tag.Key)] = tag.Value.AsInterface()
	}
	data, _ := json.Marshal(tagMap)
	return string(data)
}
