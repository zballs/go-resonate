package coala

import (
	"github.com/kazarena/json-gold/ld"
	. "github.com/zballs/go_resonate/util"
)

// json to map[string]interface{} for json-ld interpretation
func MapJSON(json []byte) map[string]interface{} {
	data := make(map[string]interface{})
	UnmarshalJSON(json, &data)
	return data
}

// json-ld methods
func CompactJSON(json []byte) (map[string]interface{}, error) {
	proc := ld.NewJsonLdProcessor()
	data := MapJSON(json)
	output, err := proc.Compact(data, nil, nil)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func ExpandJSON(json []byte) ([]interface{}, error) {
	proc := ld.NewJsonLdProcessor()
	data := MapJSON(json)
	output, err := proc.Expand(data, nil)
	if err != nil {
		return nil, err
	}
	return output, nil
}
