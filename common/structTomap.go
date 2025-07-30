package common

import "encoding/json"

func StructToStringMap(obj any) (map[string]string, error) {
	var data = make(map[string]string)
	b, err := json.Marshal(&obj)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	return data, err
}
