package reader

import (
	"io/ioutil"
)

//GetFiles func
func GetFiles(location string) ([]string, error) {

	var allFiles []string

	files, err := ioutil.ReadDir(location)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		allFiles = append(allFiles, f.Name())
	}

	return allFiles, nil
}
