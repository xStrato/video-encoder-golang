package utils

import "encoding/json"

func IsJson(j string) error {
	var js struct{}
	if err := json.Unmarshal([]byte(j), &js); err != nil {
		return err
	}
	return nil
}
