package sliceutils

type comp interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string
}

// MapFunc runs a function over each item in a slice and projects them to a new slice
func MapFunc[T any, S any](src []T, f func(item T) S) []S {
	result := make([]S, 0, len(src))
	for _, item := range src {
		result = append(result, f(item))
	}
	return result
}

// ReduceFunc reduces a slice by creating a new slice only with items that match the given function
func ReduceFunc[T any](src []T, f func(item T) bool) []T {
	result := make([]T, 0)
	for _, item := range src {
		if f(item) {
			result = append(result, item)
		}
	}
	return result
}

// Contains returns true if the slice contains the value
func Contains[T comparable](src []T, value T) bool {
	return AnyFunc(src, func(item T) bool {
		return item == value
	})
}

// AllFunc returns true if the function evaluates to true for all members of the slice
func AllFunc[T any](src []T, f func(item T) bool) bool {
	for _, item := range src {
		if !f(item) {
			return false
		}
	}
	return true
}

// AnyFunc returns true if the function evaluates to true for any member of the slice
func AnyFunc[T any](src []T, f func(item T) bool) bool {
	for _, item := range src {
		if f(item) {
			return true
		}
	}
	return false
}

// Sort sorts a set of comparable numbers
func Sort[T comp](src []T) {
	SortFunc(src, sortComps[T])
}

// SortFunc sorts a set of items.  This currently uses a bubble sort which has O(n^2) performance.
func SortFunc[T any](src []T, f func(item1 T, item2 T) int) {
	// TODO: implement a more efficient stort. Bubble sort for now, should implement a more efficient sort.
	for i := 0; i < len(src); i++ {
		for j := i + 1; j < len(src); j++ {
			item0 := src[i]
			item1 := src[j]
			if f(item1, item0) < 0 {
				src[i] = item1
				src[j] = item0
			}
		}
	}
}

func sortComps[T comp](a T, b T) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

// ToMapFunc creates a map from a slice
func ToMapFunc[T any, K comparable](slice []T, keySelector func(origin T) K) map[K]T {
	ret := make(map[K]T, len(slice))
	for _, value := range slice {
		key := keySelector(value)
		ret[key] = value
	}
	return ret
}

func Chunk[T any](src []T, maxCount int) [][]T {
	chunks := make([][]T, 0, len(src)/maxCount+1)
	for i := 0; i < len(src); i += maxCount {
		end := i + maxCount
		if end <= len(src) {
			chunks = append(chunks, src[i:end])
		} else {
			chunks = append(chunks, src[i:])
		}
	}
	return chunks
}
