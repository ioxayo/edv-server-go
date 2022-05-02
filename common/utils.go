package common

import "reflect"

// Check if value is in enum struct
func ValueInEnumStruct[E EnumStruct](value interface{}, enumStruct E) bool {
	enumStructInfo := reflect.ValueOf(&enumStruct).Elem()
	for i := 0; i < enumStructInfo.NumField(); i++ {
		enumStructValue := enumStructInfo.Field(i).Interface()
		if enumStructValue == value {
			return true
		}
	}
	return false
}
