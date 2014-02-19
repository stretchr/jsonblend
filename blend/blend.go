package blend

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	blendFunctionAdd             = "+"
	blendFunctionAddIfNotPresent = "+?"
	blendFunctionRemove          = "-"
	blendFunctionMergeDirect     = "^"
	blendFunctionMergeShallow    = "<"
	blendFunctionMergeDeep       = "<<"
)

type blendFunc func(source, dest map[string]interface{})

var functionMap = map[string]blendFunc{
	blendFunctionAdd:             BlendFuncAdd,
	blendFunctionAddIfNotPresent: BlendFuncAddIfNotPresent,
	blendFunctionRemove:          BlendFuncRemove,
	blendFunctionMergeDirect:     BlendFuncMergeDirect,
	blendFunctionMergeShallow:    BlendFuncMergeShallow,
	blendFunctionMergeDeep:       BlendFuncMergeDeep,
}

func keyIsFunction(key string) bool {
	return functionMap[key] != nil
}

// Blend blends the source into the destination using the
// blending functions present in the maps.
func Blend(source, dest map[string]interface{}) {

	for key, value := range source {
		if keyIsFunction(key) {
			functionMap[key](value.(map[string]interface{}), dest)
		}
	}
}

func BlendFuncAdd(source, dest map[string]interface{}) {
	for key, value := range source {
		if _, exists := dest[key]; !exists {
			dest[key] = make([]interface{}, 0)
		}
		dest[key] = append(dest[key].([]interface{}), value)
	}
}
func BlendFuncAddIfNotPresent(source, dest map[string]interface{}) {
	for key, value := range source {
		if _, exists := dest[key]; !exists {
			dest[key] = make([]interface{}, 0)
		}
		found := false
		for _, item := range dest[key].([]interface{}) {
			if reflect.DeepEqual(item, value) {
				found = true
				break
			}
		}
		if !found {
			dest[key] = append(dest[key].([]interface{}), value)
		}
	}
}
func BlendFuncRemove(source, dest map[string]interface{}) {
	for key, value := range source {
		if _, exists := dest[key]; !exists {
			return
		}
		location := -1
		for index, item := range dest[key].([]interface{}) {
			if reflect.DeepEqual(item, value) {
				location = index
				break
			}
		}
		if location != -1 {
			dest[key] = append(dest[key].([]interface{})[:location], dest[key].([]interface{})[location+1:]...)
		}
	}
}
func BlendFuncMergeDirect(source, dest map[string]interface{}) {
	for key, value := range source {
		dest[key] = value
	}
}
func BlendFuncMergeShallow(source, dest map[string]interface{}) {
	for key, _ := range dest {
		if _, exists := source[key]; exists {
			for sourceKey, sourceValue := range source[key].(map[string]interface{}) {
				dest[key].(map[string]interface{})[sourceKey] = sourceValue
			}
		}
	}
}
func BlendFuncMergeDeep(source, dest map[string]interface{}) {
	recurseDeepMerge(source, dest, nil, nil)
}
func recurseDeepMerge(source, dest, sourceParent, destParent map[string]interface{}) {
	for sKey, sValue := range source {
		if dValue, exists := dest[sKey]; exists {
			if isMap(sValue) && isMap(dValue) {
				// Both values are maps, we can recurse
				recurseDeepMerge(sValue.(map[string]interface{}), dValue.(map[string]interface{}), source, dest)
			} else {
				// One of them is not a map, cannot proceed
				// TODO: improve this to merge intelligently when keys are different
				panic(fmt.Sprintf("Cannot recurse. Both maps contain key but both are not maps: %#v,%#v\n", sValue, dValue))
			}
		} else {
			// The destination map does not contain the key from the source map. Merge.
			dest[sKey] = sValue
		}
	}
}

// Utility functions
func JsonToMSI(jsonString string) (msi map[string]interface{}, err error) {
	if len(jsonString) == 0 {
		jsonString = "{}"
	}
	err = json.Unmarshal([]byte(jsonString), &msi)
	return
}
func MSIToJson(msi map[string]interface{}) (jsonString string, err error) {
	jsonBytes, err := json.Marshal(msi)
	if err == nil {
		jsonString = string(jsonBytes)
	}
	return
}

func isMap(value interface{}) bool {
	_, isMap := value.(map[string]interface{})
	return isMap
}
