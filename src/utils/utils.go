package utils

import (
	"encoding/json"
)

func Copy(src interface{}, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
