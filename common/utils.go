package common

import "reflect"

// Check if value is in enum struct
func IsValidEnumMember[E EnumStruct](enumStruct E, value interface{}) bool {
	enumStructInfo := reflect.ValueOf(&enumStruct).Elem()
	for i := 0; i < enumStructInfo.NumField(); i++ {
		enumStructValue := enumStructInfo.Field(i).Interface()
		if enumStructValue == value {
			return true
		}
	}
	return false
}

// Check if value is in array
func IsValueInArray[V comparable](array []V, value V) bool {
	for _, val := range array {
		if val == value {
			return true
		}
	}
	return false
}
