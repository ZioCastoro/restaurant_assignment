package utils

// SafeDereference returns the value of the pointer if it is not nil, otherwise it returns the default value.
// If no default value is provided, it returns the zero value of the type.
func SafeDereference[T any](ptr *T, defaultValue ...T) T {
	if ptr == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return *new(T)
	}

	return *ptr
}
