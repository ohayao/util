package file

import (
	"fmt"
	"testing"
)

func TestGetJSON(t *testing.T) {
	data := struct {
		Age   int
		Name  string
		Hobby []string
	}{}
	err := GetJSON("./temp.json", &data)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(data)
}
