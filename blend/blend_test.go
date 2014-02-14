package blend

import (
	"encoding/json"
	"fmt"
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
	// direct blending
	{
		name:        "Direct blending",
		sources:     []string{`{"^":{"name":"Mat"}}`},
		destination: `{"age":31}`,
		expected:    `{"name":"Mat","age":31}`,
	},
	{
		name:        "Direct blending with nested data",
		sources:     []string{`{"^":{"grandpa":{"parent":{"another":"Tyler"}}}}`},
		destination: `{"grandpa":{"parent":{"child":"Mat"}}}`,
		expected:    `{"grandpa":{"parent":{"another":"Tyler"}}}`,
	},
	// shallow blending
	{
		name:        "Shallow blending",
		sources:     []string{`{"<": {"contact": {"address2": "Unit 1D"}}}`},
		destination: `{"contact": {"address1": "123 Fake Street"}}`,
		expected:    `{"contact": {"address1": "123 Fake Street", "address2": "Unit 1D"}}`,
	},
	// deep blending
	{
		name:        "Deep blending",
		sources:     []string{`{"<<":{"grandpa":{"parent":{"another":"Tyler"}}}}`, `{"<<":{"grandpa":{"parent":{"athird":"Ryan"}}}}`},
		destination: `{"grandpa":{"parent":{"child":"Mat"}}}`,
		expected:    `{"grandpa":{"parent":{"child":"Mat","another":"Tyler","athird":"Ryan"}}}`,
	},
	{
		name:        "Double Deep blending",
		sources:     []string{`{"<<":{"grandpa":{"parent":{"child":{"another":"Tyler"}}}}}`, `{"<<":{"grandpa":{"parent":{"child":{"athird":"Ryan"}}}}}`},
		destination: `{"grandpa":{"parent":{"child":{"first":"Mat"}}}}`,
		expected:    `{"grandpa":{"parent":{"child":{"first":"Mat","another":"Tyler","athird":"Ryan"}}}}`,
	},
	// + - adding to arrays
	{
		name:        "Create array",
		sources:     []string{`{"+":{"name":"Mat"}}`},
		destination: `{}`,
		expected:    `{"name":["Mat"]}`,
	},
	{
		name:        "Add to existing array",
		sources:     []string{`{"+":{"name":"Tyler"}}`},
		destination: `{"name":["Mat"]}`,
		expected:    `{"name":["Mat","Tyler"]}`,
	},
	{
		name:        "Add to deep array",
		sources:     []string{`{"<<":{"grandpa":{"parent":{"child":{"+":{"names":"Tyler"}}}}}}`},
		destination: `{"grandpa":{"parent":{"child":{"names":["Mat"]}}}}`,
		expected:    `{"grandpa":{"parent":{"child":{"names":["Mat","Tyler"]}}}}`,
	},
	// +? - ensure in array
	{
		name:        "Add if not there to existing array",
		sources:     []string{`{"+?":{"name":"Tyler"}}`, `{"+?":{"name":"Mat"}}`, `{"+?":{"name":"Tyler"}}`},
		destination: `{"name":["Mat"]}`,
		expected:    `{"name":["Mat","Tyler"]}`,
	},
	{
		name:        "Add if not there to new existing array",
		sources:     []string{`{"+?":{"name":"Tyler"}}`, `{"+?":{"name":"Mat"}}`, `{"+?":{"name":"Tyler"}}`},
		destination: `{}`,
		expected:    `{"name":["Tyler","Mat"]}`,
	},
	{
		name:        "Add if not there to deep array",
		sources:     []string{`{"<<":{"grandpa":{"parent":{"child":{"+?":{"names":"Tyler"}}}}}}`, `{"<<":{"grandpa":{"parent":{"child":{"+?":{"names":"Mat"}}}}}}`},
		destination: `{"grandpa":{"parent":{"child":{"names":["Mat"]}}}}`,
		expected:    `{"grandpa":{"parent":{"child":{"names":["Mat","Tyler"]}}}}`,
	},
}

func jsonToMSI(jsonString string) (msi map[string]interface{}) {
	if len(jsonString) == 0 {
		jsonString = "{}"
	}
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
		fmt.Printf("Blending %s\n", test.name)
		for _, sourceStr := range test.sources {
			source := jsonToMSI(sourceStr)
			Blend(source, current)
		}

		assert.True(t, reflect.DeepEqual(current, expected), "%s failed - Actual: %#v is not equal to Expected: %#v", test.name, MSIToJson(current), MSIToJson(expected))

	}

}
