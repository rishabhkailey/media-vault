package utils

import (
	"encoding/json"
	"fmt"
)

type StringOrStringSlice struct {
	internalList []string
}

func (s *StringOrStringSlice) List() []string {
	return s.internalList
}

func (s *StringOrStringSlice) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data")
	}
	if data[0] == '[' && data[len(data)-1] == ']' {
		return json.Unmarshal(data, &s.internalList)
	}
	if data[0] == '"' && data[len(data)-1] == '"' {
		s.internalList = []string{
			string(data[1 : len(data)-1]),
		}
		return nil
	}
	s.internalList = []string{}
	return fmt.Errorf("invalid first and last characters")
}

// type SingleOrList[T any] struct {
// 	values []T
// }

// func (s *SingleOrList[T]) UnmarshalJSON(data []byte) error {
// 	if len(data) == 0 {
// 		return fmt.Errorf("empty data")
// 	}
// 	if data[0] == '[' && data[len(data)-1] == ']' {
// 		return json.Unmarshal(data, &s.values)
// 	}
// 	json.Unmarshal()
// 	return nil
// }
