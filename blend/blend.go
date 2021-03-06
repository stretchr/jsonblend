package blend

import (
	"encoding/json"
	"errors"
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

type blendFunc func(source, dest map[string]interface{}) error

var functionMap = map[string]blendFunc{
	blendFunctionAdd:             Add,
	blendFunctionAddIfNotPresent: AddIfNotPresent,
	blendFunctionRemove:          Remove,
	blendFunctionMergeDirect:     MergeDirect,
	blendFunctionMergeShallow:    MergeShallow,
}

// defaultBlendFunc is the default behaviour when non-keys are encountered.
var defaultBlendFunc blendFunc = MergeDirect

// keyIsFunction gets whether the specified key is a special
// JSON blend function or not.
func keyIsFunction(key string) bool {
	return functionMap[key] != nil
}

// Blend blends the source string into the destination map using the
// blending functions present in the maps.
func BlendJSON(source string, dest map[string]interface{}) error {
	sourceMap, err := JsonToMSI(source)
	if err != nil {
		return err
	}
	return Blend(sourceMap, dest)
}

// Blend blends the source into the destination using the
// blending functions present in the maps.
func Blend(source, dest map[string]interface{}) error {

	if _, ok := functionMap[blendFunctionMergeDeep]; !ok {
		functionMap[blendFunctionMergeDeep] = MergeDeep
	}

	for key, value := range source {
		if keyIsFunction(key) {
			if err := functionMap[key](value.(map[string]interface{}), dest); err != nil {
				return err
			}
		} else {
			// do default behaviour
			if err := defaultBlendFunc(source, dest); err != nil {
				return err
			}
		}
	}
	return nil
}

func Add(source, dest map[string]interface{}) error {
	for key, value := range source {
		if _, exists := dest[key]; !exists {
			dest[key] = make([]interface{}, 0)
		}
		dest[key] = append(dest[key].([]interface{}), value)
	}
	return nil
}
func AddIfNotPresent(source, dest map[string]interface{}) error {
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
	return nil
}
func Remove(source, dest map[string]interface{}) error {
	for key, value := range source {
		if _, exists := dest[key]; !exists {
			return nil
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
	return nil
}
func MergeDirect(source, dest map[string]interface{}) error {
	for key, value := range source {
		dest[key] = value
	}
	return nil
}
func MergeShallow(source, dest map[string]interface{}) error {
	for key, _ := range source {
		if _, exists := dest[key]; exists {
			for sourceKey, sourceValue := range source[key].(map[string]interface{}) {
				dest[key].(map[string]interface{})[sourceKey] = sourceValue
			}
		} else {
			dest[key] = source[key]
		}
	}
	return nil
}
func MergeDeep(source, dest map[string]interface{}) error {
	for sKey, sValue := range source {
		if keyIsFunction(sKey) {
			Blend(source, dest)
			continue
		}
		var dValue interface{}
		var exists bool
		if dValue, exists = dest[sKey]; !exists {
			if isMap(sValue) {
				if mapContainsCommand(sValue.(map[string]interface{})) {
					dest[sKey] = make(map[string]interface{})
					dValue = dest[sKey]
				} else {
					dest[sKey] = sValue
					continue
				}
			} else {
				dest[sKey] = sValue
				continue
			}
		}
		if isMap(sValue) && dValue == nil {
			dest[sKey] = make(map[string]interface{})
			dValue = dest[sKey]
		}
		if sValue == nil {
			continue
		}
		if isMap(sValue) && isMap(dValue) {
			// Both values are maps, we can recurse
			MergeDeep(sValue.(map[string]interface{}), dValue.(map[string]interface{}))
		} else {
			// One of them is not a map, cannot proceed
			// TODO: improve this to merge intelligently when keys have different values/types
			// TODO: unknown when this will be the case. Needs tests. Will this ever happen?
			return errors.New(fmt.Sprintf("Cannot recurse. Both maps contain key \"%s\" but values are not both maps. Both values must be maps in order to merge", sKey))
		}

	}
	return nil
}
func mapContainsCommand(source map[string]interface{}) bool {
	retval := false
	for k, v := range source {
		if keyIsFunction(k) {
			retval = true
			break
		}
		if isMap(v) {
			retval = mapContainsCommand(v.(map[string]interface{}))
		}
	}
	return retval
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
