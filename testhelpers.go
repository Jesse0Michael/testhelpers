package testhelpers

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

// LoadJSONFile will read a json test data file and unmarshal the data into the provided interface
// This is a test helper function that can be used in two ways without having to handle errors yourself
// Either use the variable who's address was passed in or cast the result to your type
func LoadJSONFile(t *testing.T, file string, i interface{}) interface{} {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("failed to load JSON file: %s", file)
	}
	err = json.Unmarshal(b, &i)
	if err != nil {
		t.Errorf("failed to unmarshal JSON file: %s", string(b))
	}

	return i
}
