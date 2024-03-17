package utils

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// InterfaceToInt attempts to convert an interface{} type to an int.
// This is commonly needed when dealing with JSON input, as JSON numbers are typically decoded as float64.
// It supports converting from float64 (handling JSON numbers) and int types.
// Returns the converted int value or an error if the type cannot be directly converted to int.
func InterfaceToInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil // JSON numbers are decoded as float64
	case int:
		return v, nil
	default:
		errMsg := "type assertion to int failed"
		log.Error().Str("type", fmt.Sprintf("%T", val)).Msg(errMsg)
		return 0, fmt.Errorf(errMsg)
	}
}

// Contains checks if a slice of int contains a specific int element.
// It iterates through the slice and returns true if the element is found.
// This is a generic utility function to simplify checking for element existence in a slice of ints.
func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ContainsInterfaceAsInt checks if a slice of interface{} contains a specific int element.
// It attempts to convert each interface{} element to an int using InterfaceToInt function.
// If the conversion is successful and the element matches the target int, it returns true.
// This function is particularly useful when dealing with slices of mixed types or when the exact type is not known in advance.
// Logs an error and continues checking the next element if any conversion fails.
func ContainsInterfaceAsInt(s []interface{}, e int) bool {
	for _, a := range s {
		if aInt, err := InterfaceToInt(a); err == nil && aInt == e {
			return true
		} else if err != nil {
			log.Error().Err(err).Interface("element", a).Msg("Failed to convert element to int in ContainsInterfaceAsInt")
		}
	}
	return false
}
