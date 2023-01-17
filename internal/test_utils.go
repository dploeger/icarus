package internal

// BoolAddr creates a pointer to a bool value
func BoolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}

// StringAddr creates a pointer to a string value
func StringAddr(s string) *string {
	StringVar := s
	return &StringVar
}

func IntAddr(i int) *int {
	IntVar := i
	return &IntVar
}

func IgnoreError(val interface{}, _ error) interface{} {
	return val
}
