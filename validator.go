package fuselage

import (
	"reflect"
	"strconv"
	"strings"
)

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStruct validates struct fields
func ValidateStruct(v interface{}) []ValidationError {
	var errors []ValidationError
	val := reflect.ValueOf(v)
	
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return errors
	}
	
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}
		
		fieldName := getJSONFieldName(&fieldType)
		
		if strings.Contains(tag, "required") {
			if isEmpty(field) {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Message: "Field is required",
				})
			}
		}
		
		if strings.Contains(tag, "min=") {
			if minVal := extractMinValue(tag); minVal > 0 {
				if field.Kind() == reflect.String && len(field.String()) < minVal {
					errors = append(errors, ValidationError{
						Field:   fieldName,
						Message: "Field is too short",
					})
				}
			}
		}
	}
	
	return errors
}

func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

func getJSONFieldName(field *reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			return jsonTag[:idx]
		}
		return jsonTag
	}
	return field.Name
}

func extractMinValue(tag string) int {
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		if strings.HasPrefix(part, "min=") {
			if val, err := strconv.Atoi(strings.TrimPrefix(part, "min=")); err == nil {
				return val
			}
		}
	}
	return 0
}