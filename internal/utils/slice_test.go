package utils_test

import (
	"testing"

	"github.com/rishabhkailey/media-vault/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	stringTest := []struct {
		name     string
		slice    []string
		element  string
		contains bool
	}{
		{
			name:     "empty string slice",
			slice:    []string{},
			element:  "abc",
			contains: false,
		},
		{
			name:     "slice containing element",
			slice:    []string{"abc", "def", "ghi"},
			element:  "abc",
			contains: true,
		},
		{
			name:     "slice not containing element",
			slice:    []string{"abc1", "def", "ghi"},
			element:  "abc",
			contains: false,
		},
	}
	intTest := []struct {
		name     string
		slice    []int
		element  int
		contains bool
	}{
		{
			name:     "empty string slice",
			slice:    []int{},
			element:  1,
			contains: false,
		},
		{
			name:     "slice containing element",
			slice:    []int{1, 2, 3},
			element:  1,
			contains: true,
		},
		{
			name:     "slice not containing element",
			slice:    []int{2, 3, 4},
			element:  1,
			contains: false,
		},
	}

	for _, test := range stringTest {
		assert.EqualValues(t, test.contains, utils.Contains(test.slice, test.element))
	}
	for _, test := range intTest {
		assert.EqualValues(t, test.contains, utils.Contains(test.slice, test.element))
	}
}

func TestContainsSlice(t *testing.T) {
	stringTest := []struct {
		name     string
		slice    []string
		subSlice []string
		contains bool
	}{
		{
			name:     "empty string slice",
			slice:    []string{},
			subSlice: []string{"abc"},
			contains: false,
		},
		{
			name:     "both empty string slice",
			slice:    []string{},
			subSlice: []string{},
			contains: true,
		},
		{
			name:     "empty sub string slice",
			slice:    []string{"abc", "def"},
			subSlice: []string{},
			contains: true,
		},
		{
			name:     "slice containing subSlice",
			slice:    []string{"abc", "def", "ghi"},
			subSlice: []string{"abc"},
			contains: true,
		},
		{
			name:     "2 slice containing subSlice",
			slice:    []string{"abc", "def", "ghi"},
			subSlice: []string{"abc", "def", "ghi"},
			contains: true,
		},
		{
			name:     "slice not containing subSlice",
			slice:    []string{"abc1", "def", "ghi"},
			subSlice: []string{"abc"},
			contains: false,
		},
	}
	intTest := []struct {
		name     string
		slice    []int
		subSlice []int
		contains bool
	}{
		{
			name:     "empty string slice",
			slice:    []int{},
			subSlice: []int{1},
			contains: false,
		},
		{
			name:     "both empty string slice",
			slice:    []int{},
			subSlice: []int{},
			contains: true,
		},
		{
			name:     "empty sub string slice",
			slice:    []int{1, 2},
			subSlice: []int{},
			contains: true,
		},
		{
			name:     "slice containing subSlice",
			slice:    []int{1, 2, 3},
			subSlice: []int{1},
			contains: true,
		},
		{
			name:     "2 slice containing subSlice",
			slice:    []int{1, 2, 3},
			subSlice: []int{1, 2, 3},
			contains: true,
		},
		{
			name:     "slice not containing subSlice",
			slice:    []int{2, 3, 4},
			subSlice: []int{1},
			contains: false,
		},
	}

	for _, test := range stringTest {
		assert.EqualValues(t, test.contains, utils.ContainsSlice(test.slice, test.subSlice), test.name)
	}
	for _, test := range intTest {
		assert.EqualValues(t, test.contains, utils.ContainsSlice(test.slice, test.subSlice), test.name)
	}
}
