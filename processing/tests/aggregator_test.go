package processing

import (
	"testing"
	"reflect"
	"../../processing"
)

type TestData struct {
	data []int
	expected[]int
}

func TestRemoveDuplicates(t *testing.T) {
	ShouldRemoveDuplicates(t)
}

func ShouldRemoveDuplicates(t *testing.T) {

}

func TestSort(t *testing.T) {
	ShouldSortCorrectly(t)
}

func ShouldSortCorrectly(t *testing.T) {

}

func TestAggregate(t *testing.T) {
	ShouldAggregateCorrectly(t)
}

func ShouldAggregateCorrectly(t *testing.T) {

	tests := map[string]*TestData {
		"Should return an empty array.": &TestData {
			data: []int{},
			expected: []int{},
		},
		"Should return an array without duplicats.": &TestData {
			data: []int{1, 2, 1, 2, 3},
			expected: []int{1, 2, 3},
		},
		"Should return a sorted array": &TestData {
			data: []int{5, 1, 2, 3},
			expected: []int{1, 2, 3, 5},
		},
		"Should return a sorted array without duplicates": &TestData {
			data: []int{5, 1, 2, 3, 5 , 1, 7, 1},
			expected: []int{1, 2, 3, 5, 7},
		},
	}

	aggregator := processing.NewAggregator()

	for caseName, testCase := range tests {

		got := aggregator.Process(testCase.data)

		if(!reflect.DeepEqual(testCase.expected, got)) {
			t.Error(
				"For", caseName,
				"expected", testCase.expected,
				"got", got,
			)
		}
	}
}