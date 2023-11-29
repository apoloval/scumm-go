package collections

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// VisitMap visits a map in order of keys.
func VisitMap[K constraints.Integer, V any](m map[K]V, fn func(K, V)) {
	keys := make([]int, 0, len(m))
	for key := range m {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)
	for _, key := range keys {
		fn(K(key), m[K(key)])
	}
}
