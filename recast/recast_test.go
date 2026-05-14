package recast

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

/*
TestGetNumberAsInt is a test which allows you to verify that numeric values of various types are correctly converted to
an integer.

Example:

	TestGetNumberAsInt(t)
	Inputs: float32(-3.312), float64(-3.312), int(-3), int8(-3), int16(-3), int32(-3), int64(-3), uint(3), uint8(3), uint16(3), uint32(3), uint64(3)
	Outputs: int(-3), int(-3), int(-3), int(-3), int(-3), int(-3), int(-3), int(3), int(3), int(3), int(3), int(3)
*/
func TestGetNumberAsInt(test *testing.T) {
	var floatValues = GetArrayOfInterfaces(float32(-3.312), float64(-3.312))
	var intValues = GetArrayOfInterfaces(int(-3), int8(-3), int16(-3), int32(-3), int64(-3))
	var uIntValues = GetArrayOfInterfaces(uint(3), uint8(3), uint16(3), uint32(3), uint64(3))
	intermediateExpectedResult := -3.312
	expectedResultFloat := int(intermediateExpectedResult)
	for _, currentValue := range floatValues {
		obtainedResult := GetNumberAsInt(currentValue)
		assert.Equalf(test, expectedResultFloat, obtainedResult, "The variable of type '"+reflect.TypeOf(currentValue).String()+"' did not equal the expected amount.")
	}
	expectedResult := int(-3)
	for _, currentValue := range intValues {
		obtainedResult := GetNumberAsInt(currentValue)
		assert.Equalf(test, expectedResult, obtainedResult, "The variable of type '"+reflect.TypeOf(currentValue).String()+"' did not equal the expected amount.")
	}
	expectedResult = int(3)
	for _, currentValue := range uIntValues {
		obtainedResult := GetNumberAsInt(currentValue)
		assert.Equalf(test, expectedResult, obtainedResult, "The variable of type '"+reflect.TypeOf(currentValue).String()+"' did not equal the expected amount.")
	}
}

/*
TestGetNumberAsInt64 is a test which allows you to verify that numeric values of various types are correctly converted
to an int64.

Example:

	TestGetNumberAsInt64(t)
	Inputs: float32(-3.312), float64(-3.312), int(-3), int8(-3), int16(-3), int32(-3), int64(-3), uint(3), uint8(3), uint16(3), uint32(3), uint64(3)
	Outputs: int64(-3), int64(-3), int64(-3), int64(-3), int64(-3), int64(-3), int64(-3), int64(3), int64(3), int64(3), int64(3), int64(3)
*/
func TestGetNumberAsInt64(test *testing.T) {
	var floatValues = GetArrayOfInterfaces(float32(-3.312), float64(-3.312))
	var intValues = GetArrayOfInterfaces(int(-3), int8(-3), int16(-3), int32(-3), int64(-3))
	var uIntValues = GetArrayOfInterfaces(uint(3), uint8(3), uint16(3), uint32(3), uint64(3))
	intermediateExpectedResult := -3.312
	expectedResultFloat := int64(intermediateExpectedResult)
	for _, currentValue := range floatValues {
		obtainedResult := GetNumberAsInt64(currentValue)
		assert.Equalf(test, expectedResultFloat, obtainedResult, "The variable of type '"+reflect.TypeOf(currentValue).String()+"' did not equal the expected amount.")
	}
	expectedResult := int64(-3)
	for _, currentValue := range intValues {
		obtainedResult := GetNumberAsInt64(currentValue)
		assert.Equalf(test, expectedResult, obtainedResult, "The variable of type '"+reflect.TypeOf(currentValue).String()+"' did not equal the expected amount.")
	}
	expectedResult = int64(3)
	for _, currentValue := range uIntValues {
		obtainedResult := GetNumberAsInt64(currentValue)
		assert.Equalf(test, expectedResult, obtainedResult, "The variable of type '"+reflect.TypeOf(currentValue).String()+"' did not equal the expected amount.")
	}
}

/*
TestGetNumberAsFloat64 is a test which allows you to verify that numeric values of various types are correctly converted
to a float64.

Example:

	TestGetNumberAsFloat64(t)
	Inputs: float32(-3.3), float64(-3.3)
	Outputs: float64(-3.3), float64(-3.3)
*/
func TestGetNumberAsFloat64(test *testing.T) {
	var floatValues = GetArrayOfInterfaces(float32(-3.3), float64(-3.3))
	expectedResult := float64(-3.3)
	for _, currentValue := range floatValues {
		obtainedResult := GetNumberAsFloat64(currentValue)
		assert.InDelta(test, expectedResult, obtainedResult, 0.0000001, "The variable of type '"+reflect.TypeOf(currentValue).String()+"' did not equal the expected amount.")
	}
}
