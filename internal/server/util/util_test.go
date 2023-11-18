package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortMetrics(t *testing.T) {
	metrics := []string{"foo: 10", "bar: 20", "baz: 30"}
	expected := []string{"bar: 20", "baz: 30", "foo: 10"}

	actual := SortMetrics(metrics)

	assert.Equal(t, expected, actual)
}
