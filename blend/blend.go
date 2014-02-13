package blend

const (
	blendFunctionAdd = "+"
)

var functionMap = map[string]bool{
	blendFunctionAdd: true,
}

func keyIsFunction(key string) bool {
	return functionMap[key]
}

func Blend(source, dest map[string]interface{}) {

	for key, value := range source {
		if _, exists := dest[key]; !exists {
			dest[key] = make([]interface{}, 0)
		}

		for itemKey, itemValue := range value.(map[string]interface{}) {
			if keyIsFunction(itemKey) {
				executeFunction(itemKey, key, itemValue, dest)
			}
		}
	}

}

func executeFunction(function, key string, value interface{}, dest map[string]interface{}) {
	switch function {
	case blendFunctionAdd:
		dest[key] = append(dest[key].([]interface{}), value)
	default:
	}
}
