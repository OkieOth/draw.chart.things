package boxes

import (
	"fmt"
	y "gopkg.in/yaml.v3"
	"os"
)

func LoadInputFromFile[T any](filePath string) (*T, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading input file: %v", err)
	}

	var input T
	if err := y.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("error while unmarshal input: %v", err)
	}

	return &input, nil
}
