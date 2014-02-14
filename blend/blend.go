package blend

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
			if item == value {
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

}
func BlendFuncMergeDirect(source, dest map[string]interface{}) {
	for key, value := range source {
		dest[key] = value
	}
}
func BlendFuncMergeShallow(source, dest map[string]interface{}) {

}
func BlendFuncMergeDeep(source, dest map[string]interface{}) {

}
