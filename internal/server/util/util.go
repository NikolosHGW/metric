package util

import (
	"sort"
	"strings"
)

func SortMetrics(metrics []string) []string {
	sort.Slice(metrics, func(i, j int) bool {
		k1 := strings.Split(metrics[i], ":")[0]
		k2 := strings.Split(metrics[j], ":")[0]

		return k1 < k2
	})

	return metrics
}
