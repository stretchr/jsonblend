package blend

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var tests = []struct {
	name        string
	destination string
	sources     []string
	expected    string
}{
	// + - adding to arrays
	{
		name:        "Create array",
		sources:     []string{`{"name":{"+":"Mat"}}`},
		destination: `{}`,
		expected:    `{"name":["Mat"]}`,
	},
	{
		name:        "Add to existing array",
		sources:     []string{`{"name":{"+":"Tyler"}}`},
		destination: `{"name":["Mat"]}`,
		expected:    `{"name":["Mat","Tyler"]}`,
	},
	// +? - ensure in array
	{
		name:        "Add if not there to existing array",
		sources:     []string{`{"name":{"+?":"Tyler"}}`, `{"name":{"+?":"Mat"}}`, `{"name":{"+?":"Tyler"}}`},
		destination: `{"name":["Mat"]}`,
		expected:    `{"name":["Mat","Tyler"]}`,
	},
	{
		name:        "Add if not there to new existing array",
		sources:     []string{`{"name":{"+?":"Tyler"}}`, `{"name":{"+?":"Mat"}}`, `{"name":{"+?":"Tyler"}}`},
		destination: `{}`,
		expected:    `{"name":["Mat","Tyler"]}`,
	},
}

func jsonToMSI(jsonString string) (msi map[string]interface{}) {
	err := json.Unmarshal([]byte(jsonString), &msi)
	if err != nil {
		panic(err)
	}
	return
}
func MSIToJson(msi map[string]interface{}) (jsonString string) {
	jsonBytes, err := json.Marshal(msi)
	if err != nil {
		panic(err)
	}
	jsonString = string(jsonBytes)
	return
}

func TestAll(t *testing.T) {

	for _, test := range tests {

		destination := jsonToMSI(test.destination)
		expected := jsonToMSI(test.expected)

		current := destination
		for _, sourceStr := range test.sources {
			source := jsonToMSI(sourceStr)
			current = Blend(source, current)
		}

		assert.True(t, reflect.DeepEqual(current, expected), "%s failed - Actual: %#v is not equal to Expected: %#v", test.name, MSIToJson(current), MSIToJson(expected))

	}

}
