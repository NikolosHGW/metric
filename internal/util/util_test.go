package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckGaugeType(t *testing.T) {
	testCases := []struct {
		name        string
		metricType  string
		metricValue string
		expected    bool
	}{
		{"корректный случай #1", GaugeType, "3.14", true},
		{"корректный случай #2", GaugeType, "42", true},
		{"некорректное значение", GaugeType, "foo", false},
		{"некорректный тип", CounterType, "42", false},
		{"пустые строки", "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CheckGaugeType(tc.metricType, tc.metricValue)

			if actual != tc.expected {
				t.Errorf("CheckGaugeType(%q, %q) = %v; want %v", tc.metricType, tc.metricValue, actual, tc.expected)
			}
		})
	}
}

func TestCheckCounterType(t *testing.T) {
	testCases := []struct {
		name        string
		metricType  string
		metricValue string
		expected    bool
	}{
		{"корректный случай #1", CounterType, "3", true},
		{"корректный случай #2", CounterType, "42", true},
		{"некорректное значение", CounterType, "foo", false},
		{"некорректный тип", GaugeType, "42", false},
		{"пустые строки", "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CheckCounterType(tc.metricType, tc.metricValue)

			if actual != tc.expected {
				t.Errorf("CheckCounterType(%q, %q) = %v; want %v", tc.metricType, tc.metricValue, actual, tc.expected)
			}
		})
	}
}

func TestSliceStrings(t *testing.T) {
	testCases := []struct {
		name     string
		strings  []string
		i        int
		expected []string
	}{
		{"корректный случай #1", []string{"a", "b", "c"}, 1, []string{"a", "c"}},
		{"корректный случай #2", []string{"a", "b", "c"}, 0, []string{"b", "c"}},
		{"корректный случай #3", []string{"a", "b", "c"}, 2, []string{"a", "b"}},
		{"индекс за пределами слайса", []string{"a", "b", "c"}, 3, []string{"a", "b", "c"}},
		{"отрицательный индекс", []string{"a", "b", "c"}, -1, []string{"a", "b", "c"}},
		{"пустой слайс", []string{}, 0, []string{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := SliceStrings(tc.strings, tc.i)
			if !assert.ElementsMatch(t, actual, tc.expected) {
				t.Errorf("SliceStrings(%v, %d) = %v; want %v", tc.strings, tc.i, actual, tc.expected)
			}
		})
	}
}
