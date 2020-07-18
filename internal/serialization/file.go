package serialization

import (
	"encoding/json"
	"fmt"
	"os"
)

func SerializeFile(filename string, target interface{}) error {
	wordsfile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Could not create file %s: %v", filename, err)
	}

	encoder := json.NewEncoder(wordsfile)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(target)
	if err != nil {
		return fmt.Errorf("Could not write JSON to file %s: %v", filename, err)
	}

	return nil
}

func DeserializeFile(filename string, target interface{}) error {
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return SerializeFile(filename, target)
	}

	if err != nil {
		return fmt.Errorf("Could not open file %s: %v", filename, err)
	}

	err = json.NewDecoder(file).Decode(&target)
	if err != nil {
		return fmt.Errorf("Could not deserialize JSON %s: %v", filename, err)
	}

	return nil
}
