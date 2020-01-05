package argmap

import "fmt"

// IsPresent just tells if an argument is present in the map
func IsPresent(aMap map[string]interface{}, key string) bool {
	_, ok := aMap[key]
	return ok
}

// GetSFArray searches the map and possibly returns the list of argument values of a StringFlag.
// An error is returned if the key is not in the map or the identifier does not indicate a
// StringFlag output.
func GetSFArray(aMap map[string]interface{}, key string) ([]string, error) {
	if argList, ok := aMap[key]; ok {
		if valuesList, ok := argList.([]string); ok {
			return valuesList, nil
		}
		return nil, fmt.Errorf("Error: argument is not a list")
	}
	return nil, fmt.Errorf("Error: key not found in map")
}

// GetSFValue searches the map and the list of StringFlag output values in order to return
// the one at the specified index. An error is returned if the index exceeds the slice bounds.
func GetSFValue(aMap map[string]interface{}, key string, index int) (string, error) {
	valuesList, err := GetSFArray(aMap, key)
	if err != nil {
		return "", err
	} else if index >= len(valuesList) || index < 0 {
		return "", fmt.Errorf("Error: index out of bound")
	}
	return valuesList[index], nil
}

// GetBool searches the map for the boolean value of a BoolFlag. If not present, returns false.
func GetBool(aMap map[string]interface{}, key string) bool {
	if boolValue, ok := aMap[key]; ok {
		if b, ok := boolValue.(bool); ok {
			return b
		}
	}
	return false
}

// GetPositional returns the string value (if present) of the indicated positional argument.
// Returns an error if it isn't a positional or the key isn't to be found
func GetPositional(aMap map[string]interface{}, key string) (string, error) {
	if posArg, ok := aMap[key]; ok {
		if s, ok := posArg.(string); ok {
			return s, nil
		}
		return "", fmt.Errorf("Error: argument is not a string")
	}
	return "", fmt.Errorf("Error: key not found in map")
}

// GetCommandMap returns the name of the inserted command in the map and the corresponding argument
// map for that command. Returns an error if no command has been invoked by the user
func GetCommandMap(aMap map[string]interface{}) (string, map[string]interface{}, error) {
	for key, value := range aMap {
		if cmdMap, ok := value.(map[string]interface{}); ok {
			return key, cmdMap, nil
		}
	}
	return "", nil, fmt.Errorf("Error: no command found in map")
}
