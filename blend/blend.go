package blend

const (
	blendFunctionAdd       = "+"
	blendFunctionDeepMerge = "<"
	blendFunctionSetDirect = ""
)

var functionMap = map[string]bool{
	blendFunctionAdd:       true,
	blendFunctionDeepMerge: true,
}

func keyIsFunction(key string) bool {
	return functionMap[key]
}

func Blend(source, dest map[string]interface{}) {

	for key, value := range source {

		if _, isMSI := value.(map[string]interface{}); isMSI {
			for itemKey, itemValue := range value.(map[string]interface{}) {
				if keyIsFunction(itemKey) {
					executeFunction(itemKey, key, itemValue, dest)
				} else {
					executeFunction(blendFunctionSetDirect, key, value, dest)
				}
			}
		} else {
			executeFunction(blendFunctionSetDirect, key, value, dest)
		}
	}

}

func executeFunction(function, key string, value interface{}, dest map[string]interface{}) {
	switch function {
	case blendFunctionAdd:
		if _, exists := dest[key]; !exists {
			dest[key] = make([]interface{}, 0)
		}
		dest[key] = append(dest[key].([]interface{}), value)
	case blendFunctionDeepMerge:
		deepMerge(dest[key].(map[string]interface{}), value.(map[string]interface{}))
	case blendFunctionSetDirect:
		dest[key] = value
	}
}

func deepMerge(current map[string]interface{}, value map[string]interface{}) {

	// TODO: this needs to be redone entirely to support arbitrary levels of keys on BOTH the source and dest

	for _, currentValue := range current {
		for leftKey, _ := range current {
			if _, ok := value[leftKey]; ok {
				for valueKey, valueData := range value[leftKey].(map[string]interface{}) {
					current[leftKey].(map[string]interface{})[valueKey] = valueData
				}
			}
		}
		if msiValue, ok := currentValue.(map[string]interface{}); ok {
			deepMerge(msiValue, value)
		}
	}

}
