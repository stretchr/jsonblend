package blend

import (
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

}
