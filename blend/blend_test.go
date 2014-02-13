package blend

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var tests = []struct {
	destination string
	source      string
	expected    string
}{
	{
		source:      `{"name":{"+":"Mat"}}`,
		destination: `{}`,
		expected:    `{"name":["Mat"]}`,
	},
	{
		source:      `{"name":{"+":"Tyler"}}`,
		destination: `{"name":{"+":"Mat"}}`,
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

		source := jsonToMSI(test.source)
		destination := jsonToMSI(test.destination)
		expected := jsonToMSI(test.expected)

		actual := Blend(source, destination)

		assert.True(t, reflect.DeepEqual(actual, expected), "Actual: %#v is not equal to Expected: %#v", MSIToJson(actual), MSIToJson(expected))

	}

}
