package math

import (
	"fmt"
	"github.com/supercom32/consolizer/recast"
	"math"
	"strconv"
)

/*
GetAbsoluteValueAsFloat64 is a method which allows you to get the absolute value of a number as a float64.

:param number: The number to get the absolute value for.

:return: The absolute value of the number as a float64.

Example:

	absoluteValue := GetAbsoluteValueAsFloat64(-3.3)
*/
func GetAbsoluteValueAsFloat64(number interface{}) float64 {
	numberAsFloat64 := recast.GetNumberAsFloat64(number)
	if numberAsFloat64 < 0 {
		return -numberAsFloat64
	}
	return numberAsFloat64
}

/*
GetAbsoluteValueAsInt is a method which allows you to get the absolute value of a number as an int.

:param number: The number to get the absolute value for.

:return: The absolute value of the number as an int.

Example:

	absoluteValue := GetAbsoluteValueAsInt(-3)
*/
func GetAbsoluteValueAsInt(number interface{}) int {
	numberAsInt := recast.GetNumberAsInt(number)
	if numberAsInt < 0 {
		return -numberAsInt
	}
	return numberAsInt
}

/*
RoundToWholeNumber is a method which allows you to round a number to the nearest whole number.

:param number: The number to round.

:return: The rounded whole number as a float64.

Example:

	roundedNumber := RoundToWholeNumber(5.5)
*/
func RoundToWholeNumber(number interface{}) float64 {
	numberAsFloat64 := recast.GetNumberAsFloat64(number)
	return math.Round(numberAsFloat64)
}

/*
RoundToDecimal is a method which allows you to round a number to a specific number of decimal places.

:param number: The number to round.
:param numberOfPlaces: The number of decimal places to round to.

:return: The rounded number as a float64.

Example:

	roundedNumber := RoundToDecimal(5.1234567, 4)
*/
func RoundToDecimal(number interface{}, numberOfPlaces int) float64 {
	numberAsFloat64 := recast.GetNumberAsFloat64(number)
	numberFormat := "%." + strconv.Itoa(numberOfPlaces) + "f"
	numberAsString := fmt.Sprintf(numberFormat, numberAsFloat64)
	roundedNumber, _ := strconv.ParseFloat(numberAsString, 64)
	return roundedNumber
}

/*
IsNumberEven is a method which allows you to check if a number is even.

:param number: The number to check.

:return: True if the number is even, false otherwise.

Example:

	isEven := IsNumberEven(12)
*/
func IsNumberEven(number interface{}) bool {
	numberAsInt64 := recast.GetNumberAsInt64(number)
	remainder := numberAsInt64 % 2
	if remainder == 0 {
		return true
	}
	return false
}

/*
IsFloatEffectivelyEqual is a method which allows you to check if two floating point numbers are effectively equal to
each other. Since floating point operations perform approximate arithmetic, it is normal that there will be an
accumulation of rounding errors in floating-point operations. By using this method, you can check if your numbers are
for most practical purposes, equal or not by automatically rounding numbers down to 7 places.

:param firstNumber: The first number to compare.
:param secondNumber: The second number to compare.

:return: True if the numbers are effectively equal, false otherwise.

Example:

	isEqual := IsFloatEffectivelyEqual(1.00000001, 1.00000002)
*/
func IsFloatEffectivelyEqual(firstNumber, secondNumber float64) bool {
	firstNumberRounded := RoundToDecimal(firstNumber, 7)
	secondNumberRounded := RoundToDecimal(secondNumber, 7)
	return firstNumberRounded == secondNumberRounded
}
