package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// if range end not provided
const defaultRangeSize int64 = 1000000 // 1mb

// todo move to models directory?

// range is inclusive
type Range struct {
	Start int64 // start will always be provided
	End   int64 // end = -1 if not provided
}
type RangeHeader struct {
	Unit   string
	Ranges []Range
}

func ParseRangeHeader(value string) (*RangeHeader, error) {
	s := strings.Split(value, "=")
	if len(s) != 2 {
		return nil, fmt.Errorf("invalid range header")
	}
	unit := s[0]
	var ranges []Range
	rangesStr := strings.Split(s[1], ",")
	for _, r := range rangesStr {
		r, err := parseRange(r)
		if err != nil {
			return nil, fmt.Errorf("invalid range header: %w", err)
		}
		ranges = append(ranges, *r)
	}
	return &RangeHeader{
		Unit:   unit,
		Ranges: ranges,
	}, nil
}

func parseRange(r string) (*Range, error) {
	r = strings.TrimSpace(r)
	separatorIndex := strings.Index(r, "-")
	// validations
	if separatorIndex == -1 || separatorIndex == 0 {
		return nil, fmt.Errorf("invalid range: %v", r)
	}
	rangeStartEndArr := strings.Split(r, "-")
	startStr := rangeStartEndArr[0]
	endStr := rangeStartEndArr[1]
	if len(rangeStartEndArr) > 2 {
		return nil, fmt.Errorf("invalid range: %v", r)
	}
	// range start
	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid range: %w", err)
	}

	// range end
	end, err := strconv.ParseInt(rangeStartEndArr[1], 10, 64)
	if err != nil {
		if len(endStr) != 0 {
			return nil, fmt.Errorf("invalid range: %w", err)
		}
		end = start + defaultRangeSize
	}
	return &Range{
		Start: int64(start),
		End:   int64(end),
	}, nil
}
