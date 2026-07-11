package json

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Marshal returns the JSON encoding of the given value.
// Unlike the default json.Marshal function, it doesn't escape HTML tags.
func Marshal(value any) ([]byte, error) {
	buffer := new(bytes.Buffer)

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(value); err != nil {
		return nil, fmt.Errorf("encoder.Encode(): %w", err)
	}

	// Encode(), unlike Marshal() adds a newline at the end of the buffer.
	// We return the buffer without the last element in order to preserve the behavior.
	res := append([]byte{}, buffer.Bytes()[:buffer.Len()-1]...)

	return res, nil
}

// Unmarshal wraps the standard json.Unmarshal function.
func Unmarshal[T any](data []byte) (T, error) {
	var res, zero T

	if err := json.Unmarshal(data, &res); err != nil {
		return zero, fmt.Errorf("json.Unmarshal(): %w", err)
	}

	return res, nil
}

// MustMarshal wraps Marshal but panics if there is an error.
// Should be used when we are sure that the value can be marshalled.
// If the intent is to panic on error, prefer using Marshal to clarify intent.
func MustMarshal(value any) []byte {
	data, err := Marshal(value)
	if err != nil {
		panic(err)
	}

	return data
}

// MustUnmarshal wraps Unmarshal but panics if there is an error.
// Should be used when the data was already validated.
// If the intent is to panic on error, prefer using Unmarshal to clarify intent.
func MustUnmarshal[T any](data []byte) T {
	res, err := Unmarshal[T](data)
	if err != nil {
		panic(err)
	}

	return res
}
