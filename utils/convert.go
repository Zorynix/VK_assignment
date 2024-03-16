package utils

import "fmt"

func InterfaceToInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil // JSON numbers are decoded as float64
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("type assertion to int failed")
	}
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsInterfaceAsInt(s []interface{}, e int) bool {
	for _, a := range s {
		if aInt, err := InterfaceToInt(a); err == nil && aInt == e {
			return true
		}
	}
	return false
}
