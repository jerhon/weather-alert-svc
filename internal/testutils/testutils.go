package testutils

import (
	"io/ioutil"
	"path"
)

func ReadTestFile(fileName string) (string, error) {

	filePath := GetTestFilepath(fileName)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	text := string(content)
	return text, nil
}

func GetTestFilepath(fileName string) string {
	filepath := path.Join(path.Dir(fileName), "testdata", fileName)
	return filepath
}
