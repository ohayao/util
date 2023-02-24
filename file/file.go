package file

import (
	"encoding/json"
	"os"
)

func GetBytes(path string) []byte {
	bs, _ := os.ReadFile(path)
	return bs
}

func GetString(path string) string {
	bs, _ := os.ReadFile(path)
	return string(bs)
}

func GetJSON(path string, value interface{}) error {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return json.Unmarshal(bs, value)
}
