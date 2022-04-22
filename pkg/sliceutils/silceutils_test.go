package sliceutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MapFunc_SimpleSituation(t *testing.T) {
	str := []string{"hello", "to", "you"}
	lens := MapFunc(str, func(a string) int {
		return len(a)
	})

	a := assert.New(t)
	a.Equal(len(lens), 3, "Length is not 3")
}

func Test_ReduceFunc_SimpleReduce(t *testing.T) {
	str := []string{"hello", "world", "hello", "again"}
	a := assert.New(t)
	hellos := ReduceFunc(str, func(a string) bool {
		return a == "hello"
	})
	a.Equal(len(hellos), 2)
	a.Equal(hellos[0], "hello")
	a.Equal(hellos[1], "hello")
}

func Test_AnyFunc_ReturnsTrue(t *testing.T) {
	// Arrange
	str := []string{"hello", "my", "name", "is"}
	a := assert.New(t)

	// Act
	result := AnyFunc(str, func(item string) bool {
		return item == "name"
	})

	// Assert
	a.True(result)
}

func Test_AnyFunc_ReturnsFalse(t *testing.T) {
	// Arrange
	str := []string{"hello", "my", "name", "is"}
	a := assert.New(t)

	// Act
	result := AnyFunc(str, func(item string) bool {
		return item == "not found"
	})

	// Assert
	a.False(result)
}

func Test_AllFunc_ReturnsTrue(t *testing.T) {
	// Arrange
	nums := []int{1, 2, 3, 4, 5}
	a := assert.New(t)

	// Act
	result := AllFunc(nums, func(v int) bool {
		return v < 6
	})

	// Assert
	a.True(result)
}
func Test_AllFunc_ReturnsFalse(t *testing.T) {
	// Arrange
	nums := []int{6, 7, 8, 9}
	a := assert.New(t)

	// Act
	result := AllFunc(nums, func(v int) bool {
		return v > 8
	})

	// Assert
	a.False(result)
}

func Test_Sort_Numbers(t *testing.T) {
	nums := []int{5, 1, 4, 3, 2, 0}
	Sort(nums)
	a := assert.New(t)
	for idx, val := range nums {
		a.Equal(idx, val)
	}
}

func Test_Chunk(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	chunks := Chunk(nums, 2)
	if len(chunks) != 4 {
		t.Fatal("Expected 4 slices.")
	}
	for i := 0; i < len(nums); i++ {
		if chunks[i/2][i%2] != nums[i] {
			t.Error("Expected ", nums[i])
		}
	}
}
