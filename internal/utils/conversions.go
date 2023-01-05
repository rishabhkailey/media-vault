package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

func UnmarshalInterface[T comparable](input interface{}, output T) error {
	rawJson, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("[InterfaceToUrlValues]: json marshal failed %w", err)
	}
	err = json.Unmarshal(rawJson, &output)
	if err != nil {
		return fmt.Errorf("[InterfaceToUrlValues]: json unmarshal failed %w", err)
	}
	return nil
}
