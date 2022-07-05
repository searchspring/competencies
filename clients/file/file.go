package file

import (
	"fmt"
	"io/ioutil"
)

func Read(path string, name string) (string, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, name))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
