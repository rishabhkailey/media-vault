package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

// todo add check that incomfing value is a pointer similar to json.unmarshall
// todo tests with different types pointer, without pointer, pointer to pointer
func UnmarshalInterface[T comparable](input interface{}, output *T) error {
	rawJson, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("[InterfaceToUrlValues]: json marshal failed %w", err)
	}
	// update ptr -> updated ptr -> value
	err = json.Unmarshal(rawJson, output)
	if err != nil {
		return fmt.Errorf("[InterfaceToUrlValues]: json unmarshal failed %w", err)
	}
	return nil
}
