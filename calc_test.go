package rpccalc

import (
	"testing"
)

type TestCase struct {
	input1   int
	input2   int
	expected int
}

var AddTestCases = []TestCase{
	{input1: 10, input2: 1, expected: 11},
	{input1: 10, input2: 0, expected: 10},
	{input1: 0, input2: 0, expected: 0},
	{input1: -10, input2: 0, expected: -10},
	{input1: -10, input2: 0, expected: -10},
	{input1: -10, input2: -1, expected: -11},
}

func TestAdd(t *testing.T) {
	c := new(CalcService)
	for i, test := range AddTestCases {
		if result := c.Add(test.input1, test.input2); result != test.expected {
			t.Fatalf("Test case %d: expected: %v, got: %v", i, test.expected, result)
		}
	}
}

var SubtractTestCases = []TestCase{
	{input1: 10, input2: 1, expected: 9},
	{input1: 1, input2: 10, expected: -9},
	{input1: 10, input2: 0, expected: 10},
	{input1: 0, input2: 0, expected: 0},
	{input1: -10, input2: 0, expected: -10},
	{input1: -10, input2: -1, expected: -9},
	{input1: -1, input2: -10, expected: 9},
	{input1: -1, input2: 10, expected: -11},
	{input1: 1, input2: -10, expected: 11},
}

func TestSubtract(t *testing.T) {
	c := new(CalcService)
	for i, test := range SubtractTestCases {
		if result := c.Subtract(test.input1, test.input2); result != test.expected {
			t.Fatalf("Test case %d: expected: %v, actual: %v", i, test.expected, result)
		}
	}
}
