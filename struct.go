package pkgep

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func StructField(fields []string) reflect.Value {
	// Create a slice of reflect.StructField
	var structFields []reflect.StructField
	for _, q := range fields {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}
		fieldName := FieldName(q)
		fmt.Printf("fieldName: %v\n", fieldName)
		structFields = append(structFields, reflect.StructField{
			Name: fieldName,
			Type: reflect.TypeOf(""), // Assuming all fields are of type string for simplicity
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, fieldName)),
		})

	}

	fmt.Printf("structFields: %v\n", structFields)

	// Create a temporary struct type
	tempStructType := reflect.StructOf(structFields)

	// Create a slice of the temporary struct type
	tempStructSlice := reflect.MakeSlice(reflect.SliceOf(tempStructType), 0, 0)

	// Create a pointer to the slice
	tempStructSlicePtr := reflect.New(tempStructSlice.Type())

	return tempStructSlicePtr

}

func FieldName(field string) string {
	// Replace invalid characters with underscores
	field = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, field)

	// Ensure the first letter is uppercase
	if len(field) > 0 {
		field = strings.ToUpper(string(field[0])) + field[1:]
	}

	return field
}
